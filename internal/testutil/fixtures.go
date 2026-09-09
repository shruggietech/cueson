package testutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	fixtureIDPattern       = regexp.MustCompile(`^[a-z0-9]+(?:[/-][a-z0-9]+)*$`)
	artifactIDPattern      = regexp.MustCompile(`^[a-z0-9]+(?:[_-][a-z0-9]+)*$`)
	sha256Pattern          = regexp.MustCompile(`^[0-9a-f]{64}$`)
	localDrivePathPattern  = regexp.MustCompile(`(?i)^[a-z]:[\\/]`)
	localPathInTextPattern = regexp.MustCompile(`(?i)(^|[[:space:]'"(=])(?:[a-z]:[\\/]|file:(?:/{0,2})|\\\\|/(?:[^/[:space:]]+/)+)`)
)

// Manifest is the versioned, ordered root fixture inventory.
type Manifest struct {
	ManifestVersion int       `json:"manifest_version"`
	Fixtures        []Fixture `json:"fixtures"`
}

// Fixture groups related payloads under one provenance and expectation record.
type Fixture struct {
	ID             string         `json:"id"`
	Purpose        string         `json:"purpose"`
	Class          string         `json:"class"`
	Origin         Origin         `json:"origin"`
	Redistribution Redistribution `json:"redistribution"`
	Artifacts      []Artifact     `json:"artifacts"`
	Expectation    Expectation    `json:"expectation"`
}

// Origin describes where fixture bytes came from without recording local paths.
type Origin struct {
	Kind        string  `json:"kind"`
	Source      string  `json:"source"`
	Recipe      *string `json:"recipe,omitempty"`
	RetrievedOn *string `json:"retrieved_on,omitempty"`
}

// Redistribution records the explicit decision permitting committed bytes.
type Redistribution struct {
	Status         string `json:"status"`
	License        string `json:"license"`
	Attribution    string `json:"attribution"`
	NoticeRequired *bool  `json:"notice_required"`
}

// Artifact governs one exact payload beneath the root fixture inventory.
type Artifact struct {
	ID           string       `json:"id"`
	Role         string       `json:"role"`
	Path         string       `json:"path"`
	MediaType    string       `json:"media_type"`
	SizeBytes    *int64       `json:"size_bytes"`
	SHA256       string       `json:"sha256"`
	ByteContract ByteContract `json:"byte_contract"`
}

// ByteContract records characteristics that checkout and validation preserve.
type ByteContract struct {
	Encoding     string `json:"encoding"`
	BOM          string `json:"bom"`
	LineEndings  string `json:"line_endings"`
	FinalNewline string `json:"final_newline"`
}

// Expectation declares acceptance or one stable rejection boundary.
type Expectation struct {
	Result             string  `json:"result"`
	Stage              *string `json:"stage,omitempty"`
	DiagnosticContains *string `json:"diagnostic_contains,omitempty"`
}

type declaredArtifact struct {
	fixtureID  string
	artifactID string
	artifact   Artifact
}

// VerifyFixtures loads and verifies a fixture root without rewriting payloads.
func VerifyFixtures(root string) (Manifest, error) {
	var manifest Manifest
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		return manifest, fmt.Errorf("fixture manifest: cannot read manifest")
	}
	if err := decodeStrict(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("fixture manifest: decode manifest: %v", err)
	}
	declared, err := validateManifest(manifest)
	if err != nil {
		return Manifest{}, err
	}
	for _, fixture := range manifest.Fixtures {
		for _, artifact := range fixture.Artifacts {
			if err := verifyArtifact(root, fixture.ID, artifact); err != nil {
				return Manifest{}, err
			}
		}
	}
	if err := verifyInventory(root, declared); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// ArtifactByID returns one logical artifact from a fixture.
func (fixture Fixture) ArtifactByID(id string) (Artifact, bool) {
	for _, artifact := range fixture.Artifacts {
		if artifact.ID == id {
			return artifact, true
		}
	}
	return Artifact{}, false
}

// FixtureByID returns one fixture from a manifest.
func (manifest Manifest) FixtureByID(id string) (Fixture, bool) {
	for _, fixture := range manifest.Fixtures {
		if fixture.ID == id {
			return fixture, true
		}
	}
	return Fixture{}, false
}

// ArtifactFile resolves a previously validated artifact path beneath root.
func ArtifactFile(root string, artifact Artifact) string {
	return filepath.Join(root, filepath.FromSlash(artifact.Path))
}

func decodeStrict(data []byte, destination any) error {
	if !utf8.Valid(data) || bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		return fmt.Errorf("manifest must be UTF-8 without BOM")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return fmt.Errorf("trailing data: %v", err)
	}
	return nil
}

func validateManifest(manifest Manifest) (map[string]declaredArtifact, error) {
	if manifest.ManifestVersion != 1 {
		return nil, fmt.Errorf("fixture manifest: manifest_version must be 1")
	}
	if len(manifest.Fixtures) == 0 {
		return nil, fmt.Errorf("fixture manifest: fixtures must not be empty")
	}
	fixtureKeys := make(map[string]struct{}, len(manifest.Fixtures))
	pathOwners := make(map[string]declaredArtifact)
	for fixtureIndex, fixture := range manifest.Fixtures {
		label := fmt.Sprintf("fixture record %d", fixtureIndex)
		if !fixtureIDPattern.MatchString(fixture.ID) {
			return nil, fmt.Errorf("%s: fixture identifier is not portable", label)
		}
		fixtureKey := strings.ToLower(fixture.ID)
		if _, exists := fixtureKeys[fixtureKey]; exists {
			return nil, fmt.Errorf("fixture %s: fixture identifier is duplicated", fixture.ID)
		}
		fixtureKeys[fixtureKey] = struct{}{}
		if strings.TrimSpace(fixture.Purpose) == "" {
			return nil, fmt.Errorf("fixture %s: purpose is required", fixture.ID)
		}
		if fixture.Class != "accepted" && fixture.Class != "malformed" && fixture.Class != "fuzz_regression" {
			return nil, fmt.Errorf("fixture %s: class is invalid", fixture.ID)
		}
		if err := validateOrigin(fixture.ID, fixture.Origin); err != nil {
			return nil, err
		}
		if err := validateRedistribution(fixture.ID, fixture.Redistribution); err != nil {
			return nil, err
		}
		if err := validateExpectation(fixture.ID, fixture.Class, fixture.Expectation); err != nil {
			return nil, err
		}
		if len(fixture.Artifacts) == 0 {
			return nil, fmt.Errorf("fixture %s: artifacts must not be empty", fixture.ID)
		}
		artifactIDs := make(map[string]struct{}, len(fixture.Artifacts))
		for _, artifact := range fixture.Artifacts {
			if !artifactIDPattern.MatchString(artifact.ID) {
				return nil, fmt.Errorf("fixture %s: artifact identifier is not portable", fixture.ID)
			}
			if _, exists := artifactIDs[artifact.ID]; exists {
				return nil, fmt.Errorf("fixture %s artifact %s: artifact identifier is duplicated", fixture.ID, artifact.ID)
			}
			artifactIDs[artifact.ID] = struct{}{}
			if err := validateArtifactRecord(fixture.ID, artifact); err != nil {
				return nil, err
			}
			pathKey, err := PortablePathKey(artifact.Path)
			if err != nil || (!strings.HasPrefix(artifact.Path, "fixtures/") && !strings.HasPrefix(artifact.Path, "malformed/") && !strings.HasPrefix(artifact.Path, "fuzz/")) {
				return nil, fmt.Errorf("fixture %s artifact %s path is not portable", fixture.ID, artifact.ID)
			}
			if owner, exists := pathOwners[pathKey]; exists {
				return nil, fmt.Errorf("fixture %s artifact %s: artifact path is duplicated by fixture %s artifact %s", fixture.ID, artifact.ID, owner.fixtureID, owner.artifactID)
			}
			pathOwners[pathKey] = declaredArtifact{fixtureID: fixture.ID, artifactID: artifact.ID, artifact: artifact}
		}
	}
	return pathOwners, nil
}

func validateOrigin(fixtureID string, origin Origin) error {
	switch origin.Kind {
	case "project_authored":
		if origin.Recipe != nil || origin.RetrievedOn != nil {
			return fmt.Errorf("fixture %s: project-authored origin has inapplicable fields", fixtureID)
		}
	case "synthetic", "derived":
		if origin.Recipe == nil || strings.TrimSpace(*origin.Recipe) == "" {
			return fmt.Errorf("fixture %s: synthetic or derived origin recipe is required", fixtureID)
		}
		if origin.RetrievedOn != nil {
			return fmt.Errorf("fixture %s: synthetic or derived origin has retrieved_on", fixtureID)
		}
	case "third_party":
		if origin.RetrievedOn == nil {
			return fmt.Errorf("fixture %s: third-party retrieved_on is required", fixtureID)
		}
		if _, err := time.Parse(time.DateOnly, *origin.RetrievedOn); err != nil {
			return fmt.Errorf("fixture %s: third-party retrieved_on must be an ISO date", fixtureID)
		}
	default:
		return fmt.Errorf("fixture %s: origin kind is invalid", fixtureID)
	}
	if strings.TrimSpace(origin.Source) == "" {
		return fmt.Errorf("fixture %s: origin source is required", fixtureID)
	}
	lowerSource := strings.ToLower(strings.TrimSpace(origin.Source))
	if strings.HasPrefix(lowerSource, "file:") || strings.HasPrefix(lowerSource, "/") || strings.HasPrefix(lowerSource, `\\`) || localDrivePathPattern.MatchString(lowerSource) || localPathInTextPattern.MatchString(origin.Source) {
		return fmt.Errorf("fixture %s: origin source must not be a local path", fixtureID)
	}
	if origin.Recipe != nil && localPathInTextPattern.MatchString(*origin.Recipe) {
		return fmt.Errorf("fixture %s: origin recipe must not contain a local path", fixtureID)
	}
	return nil
}

func validateRedistribution(fixtureID string, redistribution Redistribution) error {
	if redistribution.Status != "approved" {
		return fmt.Errorf("fixture %s: redistribution status must be approved", fixtureID)
	}
	if strings.TrimSpace(redistribution.License) == "" {
		return fmt.Errorf("fixture %s: redistribution license is required", fixtureID)
	}
	if strings.TrimSpace(redistribution.Attribution) == "" {
		return fmt.Errorf("fixture %s: redistribution attribution is required", fixtureID)
	}
	if redistribution.NoticeRequired == nil {
		return fmt.Errorf("fixture %s: redistribution NOTICE decision is required", fixtureID)
	}
	return nil
}

func validateExpectation(fixtureID, class string, expectation Expectation) error {
	switch expectation.Result {
	case "accepted":
		if class != "accepted" || expectation.Stage != nil || expectation.DiagnosticContains != nil {
			return fmt.Errorf("fixture %s: accepted expectation is inconsistent", fixtureID)
		}
	case "rejected":
		if class == "accepted" || expectation.Stage == nil || expectation.DiagnosticContains == nil || strings.TrimSpace(*expectation.DiagnosticContains) == "" {
			return fmt.Errorf("fixture %s: rejected expectation is incomplete", fixtureID)
		}
		switch *expectation.Stage {
		case "parse", "structure", "semantics", "integrity":
		default:
			return fmt.Errorf("fixture %s: rejection stage is invalid", fixtureID)
		}
	default:
		return fmt.Errorf("fixture %s: expectation result is invalid", fixtureID)
	}
	return nil
}

func validateArtifactRecord(fixtureID string, artifact Artifact) error {
	switch artifact.Role {
	case "source", "expected_model", "expected_diagnostics", "expected_bytes", "malformed_input":
	default:
		return fmt.Errorf("fixture %s artifact %s: role is invalid", fixtureID, artifact.ID)
	}
	if strings.TrimSpace(artifact.MediaType) == "" {
		return fmt.Errorf("fixture %s artifact %s: media_type is required", fixtureID, artifact.ID)
	}
	if artifact.SizeBytes == nil {
		return fmt.Errorf("fixture %s artifact %s: size_bytes is required", fixtureID, artifact.ID)
	}
	if *artifact.SizeBytes < 0 {
		return fmt.Errorf("fixture %s artifact %s: size_bytes must be non-negative", fixtureID, artifact.ID)
	}
	if !sha256Pattern.MatchString(artifact.SHA256) {
		return fmt.Errorf("fixture %s artifact %s: sha256 must be 64 lowercase hexadecimal digits", fixtureID, artifact.ID)
	}
	return validateByteContractShape(fixtureID, artifact.ID, artifact.ByteContract)
}

func validateByteContractShape(fixtureID, artifactID string, contract ByteContract) error {
	switch contract.Encoding {
	case "utf-8", "ascii", "invalid_utf8":
		if contract.BOM != "present" && contract.BOM != "absent" {
			return fmt.Errorf("fixture %s artifact %s: text BOM declaration is invalid", fixtureID, artifactID)
		}
		if contract.LineEndings != "lf" && contract.LineEndings != "crlf" && contract.LineEndings != "cr" && contract.LineEndings != "mixed" && contract.LineEndings != "none" {
			return fmt.Errorf("fixture %s artifact %s: text line_endings declaration is invalid", fixtureID, artifactID)
		}
		if contract.FinalNewline != "present" && contract.FinalNewline != "absent" {
			return fmt.Errorf("fixture %s artifact %s: text final_newline declaration is invalid", fixtureID, artifactID)
		}
	case "binary":
		if contract.BOM != "not_applicable" || contract.LineEndings != "not_applicable" || contract.FinalNewline != "not_applicable" {
			return fmt.Errorf("fixture %s artifact %s: binary byte contract must use not_applicable", fixtureID, artifactID)
		}
	default:
		return fmt.Errorf("fixture %s artifact %s: encoding declaration is invalid", fixtureID, artifactID)
	}
	return nil
}

func verifyArtifact(root, fixtureID string, artifact Artifact) error {
	label := fmt.Sprintf("fixture %s artifact %s", fixtureID, artifact.ID)
	fullPath := filepath.Join(root, filepath.FromSlash(artifact.Path))
	current := root
	for _, component := range strings.Split(artifact.Path, "/") {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%s: declared payload is missing", label)
			}
			return fmt.Errorf("%s: cannot inspect declared payload", label)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s: symbolic links are not allowed", label)
		}
	}
	info, err := os.Lstat(fullPath)
	if err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("%s: declared payload is not a regular file", label)
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("%s: cannot read declared payload", label)
	}
	if int64(len(data)) != *artifact.SizeBytes {
		return fmt.Errorf("%s: byte length does not match", label)
	}
	sum := sha256.Sum256(data)
	if hex.EncodeToString(sum[:]) != artifact.SHA256 {
		return fmt.Errorf("%s: SHA-256 does not match", label)
	}
	if err := verifyByteContract(data, artifact.ByteContract); err != nil {
		return fmt.Errorf("%s: %v", label, err)
	}
	return nil
}

func verifyByteContract(data []byte, contract ByteContract) error {
	hasBOM := bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf})
	if contract.Encoding != "binary" {
		if (contract.BOM == "present") != hasBOM {
			return fmt.Errorf("BOM does not match declaration")
		}
		switch contract.Encoding {
		case "utf-8":
			if !utf8.Valid(data) {
				return fmt.Errorf("encoding does not match utf-8 declaration")
			}
		case "ascii":
			for _, value := range data {
				if value > 0x7f {
					return fmt.Errorf("encoding does not match ascii declaration")
				}
			}
		case "invalid_utf8":
			if utf8.Valid(data) {
				return fmt.Errorf("encoding does not match invalid_utf8 declaration")
			}
		}
		if classifyLineEndings(data) != contract.LineEndings {
			return fmt.Errorf("line endings do not match declaration")
		}
		finalNewline := "absent"
		if len(data) > 0 && (data[len(data)-1] == '\n' || data[len(data)-1] == '\r') {
			finalNewline = "present"
		}
		if finalNewline != contract.FinalNewline {
			return fmt.Errorf("final newline does not match declaration")
		}
	}
	return nil
}

func classifyLineEndings(data []byte) string {
	lf := 0
	crlf := 0
	cr := 0
	for index := 0; index < len(data); index++ {
		switch data[index] {
		case '\r':
			if index+1 < len(data) && data[index+1] == '\n' {
				crlf++
				index++
			} else {
				cr++
			}
		case '\n':
			lf++
		}
	}
	kinds := 0
	if lf > 0 {
		kinds++
	}
	if crlf > 0 {
		kinds++
	}
	if cr > 0 {
		kinds++
	}
	switch {
	case kinds == 0:
		return "none"
	case kinds > 1:
		return "mixed"
	case lf > 0:
		return "lf"
	case crlf > 0:
		return "crlf"
	default:
		return "cr"
	}
}

func verifyInventory(root string, declared map[string]declaredArtifact) error {
	seen := make(map[string]struct{}, len(declared))
	for _, subtree := range []string{"fixtures", "malformed", "fuzz"} {
		base := filepath.Join(root, subtree)
		if _, err := os.Lstat(base); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("fixture inventory: cannot inspect payload tree")
		}
		err := filepath.WalkDir(base, func(current string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("cannot inspect payload")
			}
			relative, err := filepath.Rel(root, current)
			if err != nil {
				return fmt.Errorf("cannot derive payload identity")
			}
			portable := filepath.ToSlash(relative)
			if portable == subtree {
				return nil
			}
			key, err := PortablePathKey(portable)
			if err != nil {
				return fmt.Errorf("payload path is not portable")
			}
			owner, declaredPath := declared[key]
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("cannot inspect payload")
			}
			if info.Mode()&os.ModeSymlink != 0 {
				if declaredPath {
					return fmt.Errorf("fixture %s artifact %s: symbolic links are not allowed", owner.fixtureID, owner.artifactID)
				}
				return fmt.Errorf("fixture inventory: symbolic links are not allowed")
			}
			if entry.IsDir() {
				return nil
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("fixture inventory: non-regular payloads are not allowed")
			}
			if !declaredPath {
				return fmt.Errorf("fixture inventory: payload is not declared in manifest")
			}
			if portable != owner.artifact.Path {
				return fmt.Errorf("fixture %s artifact %s: undeclared payload collides with the declared portable path", owner.fixtureID, owner.artifactID)
			}
			seen[key] = struct{}{}
			return nil
		})
		if err != nil {
			return err
		}
	}
	missing := make([]string, 0)
	for key, owner := range declared {
		if _, exists := seen[key]; !exists {
			missing = append(missing, fmt.Sprintf("fixture %s artifact %s", owner.fixtureID, owner.artifactID))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("%s: declared payload is missing from inventory", missing[0])
	}
	return nil
}
