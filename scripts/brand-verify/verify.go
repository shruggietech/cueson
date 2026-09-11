package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

const (
	defaultManifestPath       = "brand/cueson/1.0.0/import-manifest.json"
	officialSourceURL         = "https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip"
	maximumArchiveBytes int64 = 512 << 20
	maximumEntryBytes   int64 = 64 << 20
	maximumPayloadBytes int64 = 512 << 20
	maximumEntries            = 10_000
)

type importManifest struct {
	SchemaVersion int                   `json:"schema_version"`
	Brand         string                `json:"brand"`
	KitVersion    string                `json:"kit_version"`
	SourceURL     string                `json:"source_url"`
	EffectiveURL  string                `json:"effective_url"`
	AcquiredAt    string                `json:"acquired_at"`
	Archive       archiveIdentity       `json:"archive"`
	PayloadRoot   string                `json:"payload_root"`
	Entries       []manifestEntry       `json:"entries"`
	References    []repositoryReference `json:"references"`
}

type archiveIdentity struct {
	Path              string `json:"path"`
	Bytes             int64  `json:"bytes"`
	SHA256            string `json:"sha256"`
	EntryCount        int    `json:"entry_count"`
	UncompressedBytes int64  `json:"uncompressed_bytes"`
}

type manifestEntry struct {
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
}

type repositoryReference struct {
	Document  string `json:"document"`
	Asset     string `json:"asset"`
	Reference string `json:"reference"`
	Purpose   string `json:"purpose"`
}

type verificationResult struct {
	entries    int
	references int
	violations []string
}

type inspectedEntry struct {
	path   string
	bytes  int64
	sha256 string
}

func verifyRepository(repoRoot, manifestRelative string) (verificationResult, error) {
	var result verificationResult
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return result, fmt.Errorf("resolve repository root: %w", err)
	}
	rootInfo, err := os.Stat(root)
	if err != nil {
		return result, fmt.Errorf("inspect repository root: %w", err)
	}
	if !rootInfo.IsDir() {
		return result, fmt.Errorf("repository root is not a directory")
	}

	manifestPath, err := resolveRegularRepositoryFile(root, manifestRelative)
	if err != nil {
		return result, fmt.Errorf("manifest: %w", err)
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return result, fmt.Errorf("read manifest: %w", err)
	}
	manifest, err := decodeManifest(data)
	if err != nil {
		result.violations = []string{"manifest: " + err.Error()}
		return result, nil
	}

	violations := validateManifest(manifest, manifestRelative)
	if len(violations) != 0 {
		result.violations = sortedUnique(violations)
		return result, nil
	}

	archivePath, err := resolveRegularRepositoryFile(root, manifest.Archive.Path)
	if err != nil {
		violations = append(violations, "archive: "+err.Error())
	} else {
		archiveEntries, archiveViolations := inspectArchive(archivePath, manifest.Archive)
		violations = append(violations, archiveViolations...)
		violations = append(violations, compareCatalog(manifest.Entries, archiveEntries)...)
	}

	payloadEntries, payloadViolations := inspectPayload(root, manifest.PayloadRoot)
	violations = append(violations, payloadViolations...)
	violations = append(violations, comparePayload(manifest.Entries, payloadEntries)...)
	violations = append(violations, verifyReferences(root, manifest)...)

	result.entries = len(manifest.Entries)
	result.references = len(manifest.References)
	result.violations = sortedUnique(violations)
	return result, nil
}

func decodeManifest(data []byte) (importManifest, error) {
	var manifest importManifest
	if !utf8.Valid(data) {
		return manifest, errors.New("not valid UTF-8")
	}
	if err := rejectDuplicateJSONFields(data); err != nil {
		return manifest, err
	}
	if err := validateManifestJSONShape(data); err != nil {
		return manifest, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return manifest, fmt.Errorf("invalid JSON: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return manifest, err
	}
	return manifest, nil
}

func validateManifestJSONShape(data []byte) error {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if top == nil {
		return errors.New("invalid JSON: manifest must be an object")
	}
	topFields := []string{"schema_version", "brand", "kit_version", "source_url", "effective_url", "acquired_at", "archive", "payload_root", "entries", "references"}
	if err := requireExactFields("manifest", top, topFields); err != nil {
		return err
	}

	var archive map[string]json.RawMessage
	if err := json.Unmarshal(top["archive"], &archive); err != nil || archive == nil {
		return errors.New("invalid JSON: archive must be an object")
	}
	if err := requireExactFields("archive", archive, []string{"path", "bytes", "sha256", "entry_count", "uncompressed_bytes"}); err != nil {
		return err
	}

	var entries []map[string]json.RawMessage
	if err := json.Unmarshal(top["entries"], &entries); err != nil || entries == nil {
		return errors.New("invalid JSON: entries must be an array")
	}
	for index, entry := range entries {
		if entry == nil {
			return fmt.Errorf("invalid JSON: entries[%d] must be an object", index)
		}
		if err := requireExactFields(fmt.Sprintf("entries[%d]", index), entry, []string{"path", "bytes", "sha256"}); err != nil {
			return err
		}
	}

	var references []map[string]json.RawMessage
	if err := json.Unmarshal(top["references"], &references); err != nil || references == nil {
		return errors.New("invalid JSON: references must be an array")
	}
	for index, reference := range references {
		if reference == nil {
			return fmt.Errorf("invalid JSON: references[%d] must be an object", index)
		}
		if err := requireExactFields(fmt.Sprintf("references[%d]", index), reference, []string{"document", "asset", "reference", "purpose"}); err != nil {
			return err
		}
	}
	return nil
}

func requireExactFields(location string, object map[string]json.RawMessage, required []string) error {
	allowed := make(map[string]struct{}, len(required))
	for _, field := range required {
		allowed[field] = struct{}{}
		if _, exists := object[field]; !exists {
			return fmt.Errorf("invalid JSON: %s is missing required field %q", location, field)
		}
	}
	var unknown []string
	for field := range object {
		if _, exists := allowed[field]; !exists {
			unknown = append(unknown, field)
		}
	}
	if len(unknown) != 0 {
		sort.Strings(unknown)
		return fmt.Errorf("invalid JSON: %s contains unknown field %q", location, unknown[0])
	}
	return nil
}

func rejectDuplicateJSONFields(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := inspectJSONValue(decoder, "$", make(map[string]struct{})); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return requireJSONEOF(decoder)
}

func inspectJSONValue(decoder *json.Decoder, location string, scratch map[string]struct{}) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return nil
	}
	switch delimiter {
	case '{':
		keys := scratch
		clear(keys)
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errors.New("object key is not a string")
			}
			if _, exists := keys[key]; exists {
				return fmt.Errorf("duplicate field %q at %s", key, location)
			}
			keys[key] = struct{}{}
			if err := inspectJSONValue(decoder, location+"."+key, make(map[string]struct{})); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim('}') {
			return errors.New("unterminated object")
		}
	case '[':
		index := 0
		for decoder.More() {
			if err := inspectJSONValue(decoder, fmt.Sprintf("%s[%d]", location, index), make(map[string]struct{})); err != nil {
				return err
			}
			index++
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim(']') {
			return errors.New("unterminated array")
		}
	default:
		return fmt.Errorf("unexpected delimiter %q", delimiter)
	}
	return nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("invalid JSON: multiple top-level values")
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func validateManifest(manifest importManifest, manifestRelative string) []string {
	var violations []string
	if manifest.SchemaVersion != 1 {
		violations = append(violations, "manifest: schema_version must be 1")
	}
	if manifest.Brand != "cueson" {
		violations = append(violations, "manifest: brand must be cueson")
	}
	if manifest.KitVersion != "1.0.0" {
		violations = append(violations, "manifest: kit_version must be 1.0.0")
	}
	if manifest.SourceURL != officialSourceURL {
		violations = append(violations, "manifest: source_url must name the operator-designated Cueson archive")
	}
	if !validHTTPSURL(manifest.EffectiveURL) {
		violations = append(violations, "manifest: effective_url must be an absolute HTTPS URL")
	}
	parsedAcquiredAt, acquiredAtError := time.Parse(time.RFC3339, manifest.AcquiredAt)
	_, acquiredAtOffset := parsedAcquiredAt.Zone()
	if acquiredAtError != nil || acquiredAtOffset != 0 {
		violations = append(violations, "manifest: acquired_at must be a UTC RFC 3339 timestamp")
	}

	for label, value := range map[string]string{
		"archive.path": manifest.Archive.Path,
		"payload_root": manifest.PayloadRoot,
	} {
		if err := validateRepositoryPath(value); err != nil {
			violations = append(violations, fmt.Sprintf("manifest: %s: %v", label, err))
		}
	}
	manifestDirectory := path.Dir(filepath.ToSlash(manifestRelative))
	if manifest.Archive.Path != path.Join(manifestDirectory, "archive", "cueson-brand-1.0.0.zip") {
		violations = append(violations, "manifest: archive.path must name the versioned Cueson archive")
	}
	if manifest.PayloadRoot != path.Join(manifestDirectory, "kit") {
		violations = append(violations, "manifest: payload_root must name the versioned kit directory")
	}
	if manifest.Archive.Bytes <= 0 || manifest.Archive.Bytes > maximumArchiveBytes {
		violations = append(violations, fmt.Sprintf("manifest: archive.bytes must be between 1 and %d", maximumArchiveBytes))
	}
	if !validDigest(manifest.Archive.SHA256) {
		violations = append(violations, "manifest: archive.sha256 must be 64 lowercase hexadecimal characters")
	}
	if manifest.Archive.EntryCount <= 0 || manifest.Archive.EntryCount > maximumEntries {
		violations = append(violations, fmt.Sprintf("manifest: archive.entry_count must be between 1 and %d", maximumEntries))
	}
	if manifest.Archive.UncompressedBytes <= 0 || manifest.Archive.UncompressedBytes > maximumPayloadBytes {
		violations = append(violations, fmt.Sprintf("manifest: archive.uncompressed_bytes must be between 1 and %d", maximumPayloadBytes))
	}
	if len(manifest.Entries) == 0 || len(manifest.Entries) > maximumEntries {
		violations = append(violations, fmt.Sprintf("manifest: entries must contain between 1 and %d records", maximumEntries))
	}
	if manifest.Archive.EntryCount != len(manifest.Entries) {
		violations = append(violations, "manifest: archive.entry_count does not equal entries length")
	}

	var total int64
	var prior string
	seenEntries := make(map[string]struct{})
	for index, entry := range manifest.Entries {
		prefix := fmt.Sprintf("manifest: entries[%d]", index)
		if err := validatePortablePath(entry.Path); err != nil {
			violations = append(violations, prefix+": path: "+err.Error())
		}
		if index > 0 && prior >= entry.Path {
			violations = append(violations, "manifest: entries are not strictly path-sorted")
		}
		if _, exists := seenEntries[entry.Path]; exists {
			violations = append(violations, prefix+": duplicate path")
		}
		seenEntries[entry.Path] = struct{}{}
		prior = entry.Path
		if entry.Bytes < 0 || entry.Bytes > maximumEntryBytes {
			violations = append(violations, fmt.Sprintf("%s: bytes must be between 0 and %d", prefix, maximumEntryBytes))
		} else if total > maximumPayloadBytes-entry.Bytes {
			violations = append(violations, "manifest: entry byte total exceeds the allowed limit")
		} else {
			total += entry.Bytes
		}
		if !validDigest(entry.SHA256) {
			violations = append(violations, prefix+": sha256 must be 64 lowercase hexadecimal characters")
		}
	}
	violations = append(violations, caseCollisionViolations("manifest entries", entryPaths(manifest.Entries))...)
	if total != manifest.Archive.UncompressedBytes {
		violations = append(violations, "manifest: archive.uncompressed_bytes does not equal entry byte total")
	}

	if len(manifest.References) == 0 {
		violations = append(violations, "manifest: references must contain at least one record")
	}
	var priorReference string
	seenReferences := make(map[string]struct{})
	referenceAssets := make(map[string]string)
	for index, reference := range manifest.References {
		prefix := fmt.Sprintf("manifest: references[%d]", index)
		documentError := validateRepositoryPath(reference.Document)
		if documentError != nil {
			violations = append(violations, prefix+": document: "+documentError.Error())
		}
		assetError := validatePortablePath(reference.Asset)
		if assetError != nil {
			violations = append(violations, prefix+": asset: "+assetError.Error())
		}
		if reference.Reference == "" {
			violations = append(violations, prefix+": reference must not be empty")
		} else if documentError == nil && assetError == nil {
			expectedPath, err := documentRelativeAssetPath(reference.Document, manifest.PayloadRoot, reference.Asset)
			if err != nil {
				violations = append(violations, prefix+": derive asset reference: "+err.Error())
			} else if !strings.Contains(reference.Reference, expectedPath) {
				violations = append(violations, fmt.Sprintf("%s: reference must contain document-relative asset path %q", prefix, expectedPath))
			}
		}
		if reference.Purpose == "" {
			violations = append(violations, prefix+": purpose must not be empty")
		}
		key := reference.Document + "\x00" + reference.Asset + "\x00" + reference.Reference
		if index > 0 && priorReference >= key {
			violations = append(violations, "manifest: references are not strictly sorted by document, asset, and reference")
		}
		priorReference = key
		if _, exists := seenReferences[key]; exists {
			violations = append(violations, prefix+": duplicate reference record")
		}
		seenReferences[key] = struct{}{}
		groupKey := reference.Document + "\x00" + reference.Reference
		if priorAsset, exists := referenceAssets[groupKey]; exists && priorAsset != reference.Asset {
			violations = append(violations, prefix+": identical document reference text cannot identify different assets")
		} else {
			referenceAssets[groupKey] = reference.Asset
		}
	}
	return violations
}

func documentRelativeAssetPath(document, payloadRoot, asset string) (string, error) {
	documentDirectory := filepath.FromSlash(path.Dir(document))
	assetPath := filepath.FromSlash(path.Join(payloadRoot, asset))
	relative, err := filepath.Rel(documentDirectory, assetPath)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relative), nil
}

func validHTTPSURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.IsAbs() && parsed.Scheme == "https" && parsed.Host != "" && parsed.User == nil
}

func validDigest(value string) bool {
	if len(value) != sha256.Size*2 || strings.ToLower(value) != value {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func validateRepositoryPath(value string) error {
	if err := validatePortablePath(value); err != nil {
		return err
	}
	return nil
}

func validatePortablePath(value string) error {
	if value == "" {
		return errors.New("must not be empty")
	}
	if !utf8.ValidString(value) {
		return errors.New("must be valid UTF-8")
	}
	if strings.Contains(value, "\\") {
		return errors.New("backslashes are not allowed")
	}
	if strings.HasPrefix(value, "/") || path.IsAbs(value) || filepath.VolumeName(value) != "" {
		return errors.New("absolute or drive-qualified paths are not allowed")
	}
	if strings.HasSuffix(value, "/") || path.Clean(value) != value {
		return errors.New("path must be clean and must not end with a slash")
	}
	for _, character := range value {
		if character == 0 || character < 0x20 || character == 0x7f {
			return errors.New("control characters are not allowed")
		}
	}
	for _, component := range strings.Split(value, "/") {
		if component == "" || component == "." || component == ".." {
			return errors.New("empty, current, and parent components are not allowed")
		}
		if strings.HasSuffix(component, ".") || strings.HasSuffix(component, " ") {
			return errors.New("components must not end in a dot or space")
		}
		if strings.ContainsAny(component, `<>:"|?*`) {
			return errors.New("windows-reserved path characters are not allowed")
		}
		base := strings.ToUpper(strings.SplitN(component, ".", 2)[0])
		if windowsReservedName(base) {
			return fmt.Errorf("windows-reserved name %q is not allowed", component)
		}
	}
	return nil
}

func windowsReservedName(value string) bool {
	switch value {
	case "CON", "PRN", "AUX", "NUL", "CLOCK$", "CONIN$", "CONOUT$", "COM¹", "COM²", "COM³", "LPT¹", "LPT²", "LPT³":
		return true
	}
	if len(value) == 4 && (strings.HasPrefix(value, "COM") || strings.HasPrefix(value, "LPT")) {
		return value[3] >= '1' && value[3] <= '9'
	}
	return false
}

func resolveRegularRepositoryFile(root, relative string) (string, error) {
	resolved, info, err := resolveRepositoryNode(root, relative)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", filepath.ToSlash(relative))
	}
	return resolved, nil
}

func resolveRepositoryDirectory(root, relative string) (string, error) {
	resolved, info, err := resolveRepositoryNode(root, relative)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", filepath.ToSlash(relative))
	}
	return resolved, nil
}

func resolveRepositoryNode(root, relative string) (string, fs.FileInfo, error) {
	if err := validateRepositoryPath(filepath.ToSlash(relative)); err != nil {
		return "", nil, err
	}
	resolved := filepath.Join(root, filepath.FromSlash(relative))
	within, err := filepath.Rel(root, resolved)
	if err != nil || within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator)) {
		return "", nil, errors.New("path escapes repository root")
	}
	current := root
	for _, component := range strings.Split(filepath.FromSlash(relative), string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return "", nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", nil, fmt.Errorf("%s is a symbolic link", filepath.ToSlash(relative))
		}
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", nil, err
	}
	return resolved, info, nil
}

func inspectArchive(archivePath string, expected archiveIdentity) (map[string]inspectedEntry, []string) {
	entries := make(map[string]inspectedEntry)
	var violations []string
	info, err := os.Stat(archivePath)
	if err != nil {
		return entries, []string{"archive: " + err.Error()}
	}
	if info.Size() != expected.Bytes {
		violations = append(violations, fmt.Sprintf("archive: byte size is %d, manifest records %d", info.Size(), expected.Bytes))
	}
	if info.Size() > maximumArchiveBytes {
		violations = append(violations, fmt.Sprintf("archive: size exceeds %d bytes", maximumArchiveBytes))
	}
	if digest, err := digestFile(archivePath); err != nil {
		violations = append(violations, "archive: hash: "+err.Error())
	} else if digest != expected.SHA256 {
		violations = append(violations, fmt.Sprintf("archive: sha256 is %s, manifest records %s", digest, expected.SHA256))
	}

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return entries, append(violations, "archive: open ZIP: "+err.Error())
	}
	defer reader.Close()
	if len(reader.File) > maximumEntries {
		violations = append(violations, fmt.Sprintf("archive: contains %d entries, limit is %d", len(reader.File), maximumEntries))
	}
	var total int64
	var names []string
	for index, file := range reader.File {
		label := fmt.Sprintf("archive entry[%d] %q", index, file.Name)
		if file.NonUTF8 || !utf8.ValidString(file.Name) {
			violations = append(violations, label+": name is not portable UTF-8")
			continue
		}
		if err := validatePortablePath(file.Name); err != nil {
			violations = append(violations, label+": "+err.Error())
			continue
		}
		if !file.Mode().IsRegular() {
			violations = append(violations, label+": only regular files are allowed")
			continue
		}
		if file.UncompressedSize64 > uint64(maximumEntryBytes) {
			violations = append(violations, fmt.Sprintf("%s: uncompressed size exceeds %d bytes", label, maximumEntryBytes))
			continue
		}
		if _, exists := entries[file.Name]; exists {
			violations = append(violations, label+": duplicate path")
			continue
		}
		if total > maximumPayloadBytes-int64(file.UncompressedSize64) {
			violations = append(violations, label+": cumulative uncompressed size exceeds limit")
			continue
		}
		opened, err := file.Open()
		if err != nil {
			violations = append(violations, label+": open: "+err.Error())
			continue
		}
		hasher := sha256.New()
		count, readErr := io.Copy(hasher, io.LimitReader(opened, maximumEntryBytes+1))
		closeErr := opened.Close()
		if readErr != nil {
			violations = append(violations, label+": read: "+readErr.Error())
			continue
		}
		if closeErr != nil {
			violations = append(violations, label+": close: "+closeErr.Error())
			continue
		}
		if count != int64(file.UncompressedSize64) {
			violations = append(violations, fmt.Sprintf("%s: read %d bytes, header records %d", label, count, file.UncompressedSize64))
			continue
		}
		entry := inspectedEntry{path: file.Name, bytes: count, sha256: hex.EncodeToString(hasher.Sum(nil))}
		entries[file.Name] = entry
		names = append(names, file.Name)
		total += count
	}
	violations = append(violations, caseCollisionViolations("archive entries", names)...)
	if len(entries) != expected.EntryCount {
		violations = append(violations, fmt.Sprintf("archive: accepted %d regular entries, manifest records %d", len(entries), expected.EntryCount))
	}
	if total != expected.UncompressedBytes {
		violations = append(violations, fmt.Sprintf("archive: uncompressed byte total is %d, manifest records %d", total, expected.UncompressedBytes))
	}
	return entries, violations
}

func inspectPayload(root, payloadRoot string) (map[string]inspectedEntry, []string) {
	entries := make(map[string]inspectedEntry)
	var violations []string
	if err := validateRepositoryPath(payloadRoot); err != nil {
		return entries, []string{"payload: " + err.Error()}
	}
	payloadPath, err := resolveRepositoryDirectory(root, payloadRoot)
	if err != nil {
		return entries, []string{"payload: " + err.Error()}
	}
	var names []string
	err = filepath.WalkDir(payloadPath, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			violations = append(violations, "payload: walk: "+walkErr.Error())
			return nil
		}
		if filePath == payloadPath {
			return nil
		}
		relative, relErr := filepath.Rel(payloadPath, filePath)
		if relErr != nil {
			violations = append(violations, "payload: relative path: "+relErr.Error())
			return nil
		}
		name := filepath.ToSlash(relative)
		info, infoErr := os.Lstat(filePath)
		if infoErr != nil {
			violations = append(violations, fmt.Sprintf("payload %q: inspect: %v", name, infoErr))
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			violations = append(violations, fmt.Sprintf("payload %q: symbolic links are not allowed", name))
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			violations = append(violations, fmt.Sprintf("payload %q: only regular files are allowed", name))
			return nil
		}
		if err := validatePortablePath(name); err != nil {
			violations = append(violations, fmt.Sprintf("payload %q: %v", name, err))
			return nil
		}
		if info.Size() > maximumEntryBytes {
			violations = append(violations, fmt.Sprintf("payload %q: size exceeds %d bytes", name, maximumEntryBytes))
			return nil
		}
		digest, digestErr := digestFile(filePath)
		if digestErr != nil {
			violations = append(violations, fmt.Sprintf("payload %q: hash: %v", name, digestErr))
			return nil
		}
		entries[name] = inspectedEntry{path: name, bytes: info.Size(), sha256: digest}
		names = append(names, name)
		return nil
	})
	if err != nil {
		violations = append(violations, "payload: walk: "+err.Error())
	}
	violations = append(violations, caseCollisionViolations("payload entries", names)...)
	return entries, violations
}

func compareCatalog(manifest []manifestEntry, actual map[string]inspectedEntry) []string {
	return compareEntries("archive", manifest, actual)
}

func comparePayload(manifest []manifestEntry, actual map[string]inspectedEntry) []string {
	return compareEntries("payload", manifest, actual)
}

func compareEntries(source string, manifest []manifestEntry, actual map[string]inspectedEntry) []string {
	var violations []string
	expected := make(map[string]manifestEntry, len(manifest))
	for _, entry := range manifest {
		expected[entry.Path] = entry
		found, exists := actual[entry.Path]
		if !exists {
			violations = append(violations, fmt.Sprintf("%s: missing manifest entry %q", source, entry.Path))
			continue
		}
		if found.bytes != entry.Bytes {
			violations = append(violations, fmt.Sprintf("%s entry %q: byte size is %d, manifest records %d", source, entry.Path, found.bytes, entry.Bytes))
		}
		if found.sha256 != entry.SHA256 {
			violations = append(violations, fmt.Sprintf("%s entry %q: sha256 is %s, manifest records %s", source, entry.Path, found.sha256, entry.SHA256))
		}
	}
	for name := range actual {
		if _, exists := expected[name]; !exists {
			violations = append(violations, fmt.Sprintf("%s: unrecorded entry %q", source, name))
		}
	}
	return violations
}

func verifyReferences(root string, manifest importManifest) []string {
	var violations []string
	entrySet := make(map[string]struct{}, len(manifest.Entries))
	for _, entry := range manifest.Entries {
		entrySet[entry.Path] = struct{}{}
	}
	type referenceSpan struct {
		start int
		end   int
		label string
	}
	documents := make(map[string][]byte)
	assigned := make(map[string][]referenceSpan)
	for _, reference := range manifest.References {
		label := fmt.Sprintf("reference %q -> %q", reference.Document, reference.Asset)
		if _, exists := entrySet[reference.Asset]; !exists {
			violations = append(violations, label+": asset is not a manifest entry")
			continue
		}
		assetPath := path.Join(manifest.PayloadRoot, reference.Asset)
		if _, err := resolveRegularRepositoryFile(root, assetPath); err != nil {
			violations = append(violations, label+": asset: "+err.Error())
			continue
		}
		content, exists := documents[reference.Document]
		if !exists {
			documentPath, err := resolveRegularRepositoryFile(root, reference.Document)
			if err != nil {
				violations = append(violations, label+": document: "+err.Error())
				continue
			}
			content, err = os.ReadFile(documentPath)
			if err != nil {
				violations = append(violations, label+": read document: "+err.Error())
				continue
			}
			documents[reference.Document] = content
		}
		positions := referencePositions(content, []byte(reference.Reference))
		if len(positions) != 1 {
			violations = append(violations, fmt.Sprintf("%s: exact reference text occurs %d times; expected exactly 1 distinct use", label, len(positions)))
			continue
		}
		position := positions[0]
		span := referenceSpan{start: position, end: position + len(reference.Reference), label: label}
		for _, prior := range assigned[reference.Document] {
			if span.start < prior.end && prior.start < span.end {
				violations = append(violations, fmt.Sprintf("%s: exact reference overlaps the distinct use assigned to %s", label, prior.label))
			}
		}
		assigned[reference.Document] = append(assigned[reference.Document], span)
	}
	return violations
}

func referencePositions(content, reference []byte) []int {
	if len(reference) == 0 {
		return nil
	}
	var positions []int
	for offset := 0; offset <= len(content)-len(reference); {
		index := bytes.Index(content[offset:], reference)
		if index < 0 {
			break
		}
		position := offset + index
		positions = append(positions, position)
		offset = position + len(reference)
	}
	return positions
}

func caseCollisionViolations(label string, names []string) []string {
	var violations []string
	seen := make(map[string]string, len(names))
	for _, name := range names {
		key := repositoryEquivalentPathKey(name)
		if prior, exists := seen[key]; exists && prior != name {
			violations = append(violations, fmt.Sprintf("%s: repository-equivalent collision between %q and %q", label, prior, name))
			continue
		}
		seen[key] = name
	}
	return violations
}

func repositoryEquivalentPathKey(value string) string {
	decomposed := norm.NFD.String(value)
	return norm.NFD.String(cases.Fold().String(decomposed))
}

func entryPaths(entries []manifestEntry) []string {
	paths := make([]string, len(entries))
	for index, entry := range entries {
		paths[index] = entry.Path
	}
	return paths
}

func digestFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func sortedUnique(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	result := values[:0]
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
