package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// NativeProof records only governed identities and counts, never temporary paths.
type NativeProof struct {
	ConsumerAnnotations           *ConsumerAnnotationProof `json:"consumer_annotations,omitempty"`
	PublishedConsumers            []PublishedConsumerProof `json:"published_consumers,omitempty"`
	Formats                       []string                 `json:"formats"`
	HistoricalInputCount          int                      `json:"historical_input_count"`
	HistoricalStrictRefusalCount  int                      `json:"historical_strict_refusal_count,omitempty"`
	OldConsumerVersion            string                   `json:"old_consumer_version,omitempty"`
	OldConsumerRevision           string                   `json:"old_consumer_revision,omitempty"`
	OldConsumerArchive            string                   `json:"old_consumer_archive,omitempty"`
	OldConsumerArchiveSHA256      string                   `json:"old_consumer_archive_sha256,omitempty"`
	OldConsumerRefusalCount       int                      `json:"old_consumer_refusal_count,omitempty"`
	OldConsumerIdentityProbeCount int                      `json:"old_consumer_identity_probe_count,omitempty"`
}

type nativeEnvelope struct {
	Schema       string             `json:"$schema"`
	Version      string             `json:"schema_version"`
	Format       string             `json:"format"`
	Capabilities nativeCapabilities `json:"format_support"`
	Producer     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"producer"`
	Source struct {
		Assets []struct {
			Data string `json:"data_base64"`
			Size struct {
				Bytes int `json:"bytes"`
			} `json:"size"`
			Hashes struct {
				SHA256 string `json:"sha256"`
			} `json:"hashes"`
		} `json:"assets"`
	} `json:"source"`
}

type nativeCapabilities struct {
	Status  string `json:"status"`
	Ingest  bool   `json:"ingest_supported"`
	Render  bool   `json:"render_supported"`
	Restore bool   `json:"restore_supported"`
	OCR     bool   `json:"ocr_required_for_semantic_output"`
}

func (capabilities nativeCapabilities) stableComplete() bool {
	return capabilities.Status == "stable" && capabilities.Ingest && capabilities.Render && capabilities.Restore && !capabilities.OCR
}

func verifyStableHostCLI(ctx context.Context, binary []byte, name, repository, version string) (NativeProof, error) {
	return verifyHostWorkflows(ctx, binary, name, repository, version, true)
}

func verifyDevelopmentHostCLI(ctx context.Context, binary []byte, name, repository, version string) (NativeProof, error) {
	return verifyHostWorkflows(ctx, binary, name, repository, version, false)
}

func verifyHostWorkflows(ctx context.Context, binary []byte, name, repository, version string, publishedConsumer bool) (NativeProof, error) {
	var proof NativeProof
	directory, err := os.MkdirTemp("", "cueson-stable-native-")
	if err != nil {
		return proof, err
	}
	defer os.RemoveAll(directory)
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, binary, 0700); err != nil {
		return proof, err
	}
	formats, err := verifyNativeWorkflows(ctx, path, repository, directory, version)
	if err != nil {
		return proof, err
	}
	historical, err := verifyHistoricalNative(ctx, path, repository, directory)
	if err != nil {
		return proof, err
	}
	if publishedConsumer {
		proof, err = verifyPublishedOldConsumer(ctx, repository, directory)
		if err != nil {
			return proof, err
		}
	}
	proof.Formats, proof.HistoricalInputCount = formats, historical
	if version == "1.2.0" {
		proof.HistoricalStrictRefusalCount = historical * 2
		annotations, err := verifyConsumerNativeWorkflows(ctx, path, repository, directory)
		if err != nil {
			return proof, err
		}
		proof.ConsumerAnnotations = &annotations
		if publishedConsumer {
			legacy := PublishedConsumerProof{Version: proof.OldConsumerVersion, Revision: proof.OldConsumerRevision, Archive: proof.OldConsumerArchive, ArchiveSHA256: proof.OldConsumerArchiveSHA256, RefusalCount: proof.OldConsumerRefusalCount, IdentityProbeCount: proof.OldConsumerIdentityProbeCount, SelectorIdentityProbeCount: proof.OldConsumerIdentityProbeCount}
			for _, format := range []string{"subrip", "webvtt"} {
				input := filepath.Join(repository, "internal", "schema", "testdata", "historical-v1.0.0-"+format+".json")
				if _, err := nativeSuccess(ctx, filepath.Join(directory, "old-"+name), "validate", input); err != nil {
					return proof, err
				}
				legacy.PositiveControlCount++
			}
			consumer, err := verifyPublishedV110Consumer(ctx, repository, directory)
			if err != nil {
				return proof, err
			}
			proof.PublishedConsumers = []PublishedConsumerProof{legacy, consumer}
		}
	}
	return proof, nil
}

func invokeNative(ctx context.Context, binary string, args ...string) (int, []byte, []byte, error) {
	command := exec.CommandContext(ctx, binary, args...)
	configureProcess(command)
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if err == nil {
		return 0, stdout.Bytes(), stderr.Bytes(), nil
	}
	var failure *exec.ExitError
	if errors.As(err, &failure) {
		return failure.ExitCode(), stdout.Bytes(), stderr.Bytes(), nil
	}
	return -1, stdout.Bytes(), stderr.Bytes(), err
}

func nativeSuccess(ctx context.Context, binary string, args ...string) ([]byte, error) {
	status, stdout, stderr, err := invokeNative(ctx, binary, args...)
	if err != nil {
		return nil, fmt.Errorf("native %s could not execute: %w", args[0], err)
	}
	if status != 0 {
		return nil, fmt.Errorf("native %s returned %d", args[0], status)
	}
	if args[0] == "validate" || args[0] == "inspect" {
		if err := verifyConsumerReadoutPrivacy(stdout, stderr); err != nil {
			return nil, err
		}
	}
	if args[0] == "validate" && len(stdout) != 0 {
		return nil, fmt.Errorf("successful validation unexpectedly produced a payload")
	}
	// Successful codec/conversion diagnostics are allowed by the CLI contract.
	if err := verifyDiagnosticPrivacy(stderr, append(deriveForbidden(filepath.Dir(binary), nil), proofSpeakerID)); err != nil {
		return nil, err
	}
	return stdout, nil
}

func verifyNativePayload(payload, original []byte, version, format string, forbidden []string) error {
	if err := scanForbidden("native Cue JSON", payload, forbidden); err != nil {
		return err
	}
	var envelope nativeEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return err
	}
	if envelope.Version != version || envelope.Schema != "https://cueson.io/schema/v"+version+"/cueson.schema.json" || envelope.Producer.Version != version || envelope.Producer.Name != "cueson" || envelope.Format != format || !envelope.Capabilities.stableComplete() {
		return fmt.Errorf("native encode identity differs from exact current contract")
	}
	if len(envelope.Source.Assets) != 1 {
		return fmt.Errorf("native encode source asset count differs")
	}
	asset := envelope.Source.Assets[0]
	decoded, err := base64.StdEncoding.DecodeString(asset.Data)
	if err != nil || !bytes.Equal(decoded, original) || asset.Size.Bytes != len(original) || asset.Hashes.SHA256 != digestHex(original) {
		return fmt.Errorf("native encode source bytes, length or digest differ")
	}
	return nil
}

func digestHex(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func verifyRefusal(status int, stdout, stderr []byte, destination string, forbidden []string) error {
	if status != 1 || len(stdout) != 0 || len(stderr) == 0 {
		return fmt.Errorf("refusal must return runtime failure, diagnostics and no payload")
	}
	if err := verifyDiagnosticPrivacy(stderr, forbidden); err != nil {
		return err
	}
	if destination != "" {
		data, err := os.ReadFile(destination)
		if err != nil || string(data) != "sentinel" {
			return fmt.Errorf("refusal replaced or removed existing destination")
		}
	}
	return nil
}

// Diagnostics legitimately contain absolute JSON pointers such as /cues/0.
// Check concrete local identifiers here; structural path scanning belongs to
// public JSON payloads, where a JSON pointer diagnostic is not source metadata.
func verifyDiagnosticPrivacy(data []byte, forbidden []string) error {
	text := strings.ToLower(string(data))
	for _, needle := range forbidden {
		if strings.Contains(text, needle) {
			return fmt.Errorf("diagnostics contain a forbidden local identifier")
		}
	}
	return nil
}

func verifyNativeWorkflows(ctx context.Context, binary, repository, directory, version string) ([]string, error) {
	formats := []string{"srt", "vtt", "ass", "ssa"}
	forbidden := deriveForbidden(repository, []string{directory})
	for _, format := range formats {
		original := []byte("1\n00:00:01,250 --> 00:00:04,200\n<i>Hello, world.</i>\n")
		if format == "vtt" {
			original = []byte("WEBVTT\n\n00:00:01.250 --> 00:00:04.200\nHello, world.\n")
		}
		if format == "ass" || format == "ssa" {
			var err error
			original, err = readRegularFile(filepath.Join(repository, "testdata", "fixtures", "scripted", format+"-bom-crlf", "source", "captions."+format))
			if err != nil {
				return nil, err
			}
		}
		input := filepath.Join(directory, "source."+format)
		jsonPath := input + ".cueson.json"
		if err := os.WriteFile(input, original, 0600); err != nil {
			return nil, err
		}
		payload, err := nativeSuccess(ctx, binary, "encode", input, "--stdout")
		if err != nil {
			return nil, err
		}
		wantFormat := format
		if format == "srt" {
			wantFormat = "subrip"
		}
		if format == "vtt" {
			wantFormat = "webvtt"
		}
		if err := verifyNativePayload(payload, original, version, wantFormat, forbidden); err != nil {
			return nil, err
		}
		if err := os.WriteFile(jsonPath, payload, 0600); err != nil {
			return nil, err
		}
		for _, path := range []string{input, jsonPath} {
			validated, err := nativeSuccess(ctx, binary, "validate", path)
			if err != nil || len(validated) != 0 {
				return nil, fmt.Errorf("native validation %s: %v", format, err)
			}
			inspected, err := nativeSuccess(ctx, binary, "inspect", path, "--json")
			if err != nil {
				return nil, err
			}
			if !json.Valid(inspected) {
				return nil, fmt.Errorf("native inspection is not JSON")
			}
			var report struct {
				Format string `json:"format"`
				Schema struct {
					Version string `json:"version"`
				} `json:"schema"`
				Integrity struct {
					Status string `json:"status"`
				} `json:"integrity"`
				Capabilities struct {
					Declared  nativeCapabilities `json:"declared"`
					Installed struct {
						Ingest   bool `json:"ingest"`
						Render   bool `json:"render"`
						Restore  bool `json:"restore"`
						Validate bool `json:"validate"`
						Inspect  bool `json:"inspect"`
					} `json:"installed"`
				} `json:"capabilities"`
			}
			if err := json.Unmarshal(inspected, &report); err != nil {
				return nil, err
			}
			wantFormat := format
			if format == "srt" {
				wantFormat = "subrip"
			}
			if format == "vtt" {
				wantFormat = "webvtt"
			}
			installed := report.Capabilities.Installed
			if report.Format != wantFormat || report.Schema.Version != version || report.Integrity.Status != "verified" || !report.Capabilities.Declared.stableComplete() || !installed.Ingest || !installed.Render || !installed.Restore || !installed.Validate || !installed.Inspect {
				return nil, fmt.Errorf("native inspection format, identity or integrity differs")
			}
			if err := scanForbidden("native inspection", inspected, forbidden); err != nil {
				return nil, err
			}
		}
		rendered, err := nativeSuccess(ctx, binary, "render", jsonPath, "--to", format, "--output", "-")
		if err != nil || len(rendered) == 0 {
			return nil, fmt.Errorf("native render %s: %v", format, err)
		}
		renderPath := filepath.Join(directory, "rendered."+format)
		if err := os.WriteFile(renderPath, rendered, 0600); err != nil {
			return nil, err
		}
		if _, err := nativeSuccess(ctx, binary, "validate", renderPath); err != nil {
			return nil, err
		}
		output := filepath.Join(directory, "restored."+format)
		restoredOutput, err := nativeSuccess(ctx, binary, "restore", jsonPath, "--output", output, "--no-metadata")
		if err != nil || len(restoredOutput) != 0 {
			return nil, fmt.Errorf("native restore %s: %v", format, err)
		}
		restored, err := os.ReadFile(output)
		if err != nil || !bytes.Equal(original, restored) {
			return nil, fmt.Errorf("native restore changed source")
		}
		target := "vtt"
		if format == "vtt" {
			target = "srt"
		}
		converted, err := nativeSuccess(ctx, binary, "convert", jsonPath, "--to", target)
		if err != nil || len(converted) == 0 {
			return nil, fmt.Errorf("native convert %s: %v", format, err)
		}
		convertedPath := filepath.Join(directory, "converted-"+format+"."+target)
		if err := os.WriteFile(convertedPath, converted, 0600); err != nil {
			return nil, err
		}
		if _, err := nativeSuccess(ctx, binary, "validate", convertedPath); err != nil {
			return nil, err
		}
		if format == "ass" || format == "ssa" {
			if err := verifyForcedRefusal(ctx, binary, directory, jsonPath, "convert", "vtt", false, "--strict"); err != nil {
				return nil, err
			}
		}
		corrupt := bytes.Replace(payload, []byte(digestHex(original)), []byte(strings.Repeat("0", 64)), 1)
		corruptPath := filepath.Join(directory, "corrupt-"+format+".json")
		if err := os.WriteFile(corruptPath, corrupt, 0600); err != nil {
			return nil, err
		}
		for _, action := range []string{"restore", "render", "convert"} {
			if err := verifyForcedRefusal(ctx, binary, directory, corruptPath, action, target, false); err != nil {
				return nil, err
			}
		}
		after, err := os.ReadFile(jsonPath)
		if err != nil || !bytes.Equal(payload, after) {
			return nil, fmt.Errorf("native workflow mutated Cue JSON")
		}
		inputAfter, err := os.ReadFile(input)
		if err != nil || !bytes.Equal(original, inputAfter) {
			return nil, fmt.Errorf("native workflow mutated original source bytes")
		}
	}
	return formats, nil
}

func verifyForcedRefusal(ctx context.Context, binary, directory, input, action, target string, identity bool, extra ...string) error {
	payloadBefore, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	if identity {
		if err := verifyOldConsumerIdentityGate(ctx, binary, input, directory); err != nil {
			return err
		}
	}
	destination := filepath.Join(directory, "refusal-destination")
	if err := os.WriteFile(destination, []byte("sentinel"), 0600); err != nil {
		return err
	}
	args := []string{action, input, "--output", destination, "--force"}
	if action == "restore" {
		args = append(args, "--no-metadata")
	} else {
		args = append(args, "--to", target)
	}
	args = append(args, extra...)
	before, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	status, stdout, stderr, err := invokeNative(ctx, binary, args...)
	if err != nil {
		return err
	}
	if err := verifyRefusal(status, stdout, stderr, destination, deriveForbidden(directory, nil)); err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	after, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(before, after) {
		return fmt.Errorf("refusal created staging or publication artifacts")
	}
	payloadAfter, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(payloadBefore, payloadAfter) {
		return fmt.Errorf("refused operation mutated candidate payload")
	}
	return nil
}

func verifyHistoricalNative(ctx context.Context, binary, repository, directory string) (int, error) {
	inputs := []struct{ version, name, native, target string }{
		{"1.0.0", "subrip", "srt", "vtt"}, {"1.0.0", "webvtt", "vtt", "srt"},
		{"1.1.0", "subrip", "srt", "vtt"}, {"1.1.0", "webvtt", "vtt", "srt"},
		{"1.1.0", "ass", "ass", "srt"}, {"1.1.0", "ssa", "ssa", "vtt"},
	}
	for _, format := range inputs {
		payload, err := readRegularFile(filepath.Join(repository, "internal", "schema", "testdata", "historical-v"+format.version+"-"+format.name+".json"))
		if err != nil {
			return 0, err
		}
		var envelope nativeEnvelope
		if err := json.Unmarshal(payload, &envelope); err != nil {
			return 0, err
		}
		if envelope.Version != format.version || envelope.Schema != "https://cueson.io/schema/v"+format.version+"/cueson.schema.json" || len(envelope.Source.Assets) != 1 {
			return 0, fmt.Errorf("frozen historical fixture identity differs")
		}
		original, err := base64.StdEncoding.DecodeString(envelope.Source.Assets[0].Data)
		if err != nil {
			return 0, err
		}
		input := filepath.Join(directory, "historical-"+format.version+"-"+format.name+".json")
		output := filepath.Join(directory, "historical-restored-"+format.version+"."+format.native)
		if err := os.WriteFile(input, payload, 0600); err != nil {
			return 0, err
		}
		for _, args := range [][]string{{"validate", input}, {"inspect", input, "--json"}, {"render", input, "--to", format.native}, {"convert", input, "--to", format.target}, {"restore", input, "--output", output, "--no-metadata"}} {
			result, err := nativeSuccess(ctx, binary, args...)
			if err != nil {
				return 0, err
			}
			if args[0] == "inspect" && !bytes.Contains(result, []byte(`"version":"`+format.version+`"`)) {
				return 0, fmt.Errorf("historical inspection relabeled identity")
			}
		}
		restored, err := os.ReadFile(output)
		if err != nil || !bytes.Equal(original, restored) {
			return 0, fmt.Errorf("historical source restoration differs")
		}
		for _, occupied := range []bool{true, false} {
			if err := verifyNativeDestinationRefusal(ctx, binary, directory, input, "convert", format.target, occupied, "--strict"); err != nil {
				return 0, fmt.Errorf("historical strict conversion %s/%s: %w", format.version, format.name, err)
			}
		}
		corrupt := bytes.Replace(payload, []byte(envelope.Source.Assets[0].Hashes.SHA256), []byte(strings.Repeat("0", 64)), 1)
		corruptPath := input + ".corrupt.json"
		if err := os.WriteFile(corruptPath, corrupt, 0600); err != nil {
			return 0, err
		}
		for _, action := range []string{"restore", "render", "convert"} {
			if err := verifyForcedRefusal(ctx, binary, directory, corruptPath, action, format.target, false); err != nil {
				return 0, err
			}
		}
		after, err := os.ReadFile(input)
		if err != nil || !bytes.Equal(payload, after) {
			return 0, fmt.Errorf("historical payload changed")
		}
	}
	return len(inputs), nil
}

func verifyPublishedOldConsumer(ctx context.Context, repository, directory string) (NativeProof, error) {
	var proof NativeProof
	contractBytes, err := readRegularFile(filepath.Join(repository, "specs", "S020-publish-v1-release", "contracts", "release-publication-contract.json"))
	if err != nil {
		return proof, err
	}
	var contract struct {
		Version  string `json:"version"`
		Revision string `json:"source_revision"`
		Assets   []struct {
			Name   string `json:"name"`
			Size   int64  `json:"size"`
			SHA256 string `json:"sha256"`
			GOOS   string `json:"goos"`
			GOARCH string `json:"goarch"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(contractBytes, &contract); err != nil {
		return proof, err
	}
	if contract.Version != "1.0.0" || contract.Revision != "2cad4c816340404289b4d1d87179a4071713bb46" {
		return proof, fmt.Errorf("published old-consumer contract binding changed")
	}
	targets := expectedTargets("1.0.0")
	for _, target := range targets {
		if target.GOOS != runtime.GOOS || target.GOARCH != runtime.GOARCH {
			continue
		}
		for _, asset := range contract.Assets {
			if asset.Name != target.Archive {
				continue
			}
			if asset.GOOS != target.GOOS || asset.GOARCH != target.GOARCH || asset.Size <= 0 || asset.Size > 16<<20 {
				return proof, fmt.Errorf("invalid published archive binding")
			}
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://github.com/shruggietech/cueson/releases/download/v1.0.0/"+asset.Name, nil)
			if err != nil {
				return proof, err
			}
			response, err := (&http.Client{Timeout: 90 * time.Second}).Do(request)
			if err != nil {
				return proof, fmt.Errorf("download published old consumer: %w", err)
			}
			archive, readErr := io.ReadAll(io.LimitReader(response.Body, asset.Size+1))
			closeErr := response.Body.Close()
			if response.StatusCode != http.StatusOK || readErr != nil || closeErr != nil || int64(len(archive)) != asset.Size || digestHex(archive) != asset.SHA256 {
				return proof, fmt.Errorf("published archive HTTP status, length or digest differs")
			}
			members, err := readArchive(asset.Name, archive)
			if err != nil {
				return proof, err
			}
			schema, err := readRegularFile(filepath.Join(repository, "schema", "releases", "v1.0.0", releaseSchemaFilename))
			if err != nil {
				return proof, err
			}
			license, err := readRegularFile(filepath.Join(repository, "LICENSE"))
			if err != nil {
				return proof, err
			}
			notice, err := readRegularFile(filepath.Join(repository, "NOTICE"))
			if err != nil {
				return proof, err
			}
			binary, err := verifyMembers(target, members, schema, license, notice, "1.0.0")
			if err != nil {
				return proof, err
			}
			if err := verifyBuildInfo(binary, target, contract.Revision, deriveForbidden(repository, nil)); err != nil {
				return proof, err
			}
			path := filepath.Join(directory, "old-"+target.Binary)
			if err := os.WriteFile(path, binary, 0700); err != nil {
				return proof, err
			}
			version, err := nativeSuccess(ctx, path, "version")
			if err != nil || string(version) != "1.0.0\n" {
				return proof, fmt.Errorf("published consumer version differs")
			}
			for _, format := range []string{"srt", "vtt", "ass", "ssa"} {
				input := filepath.Join(directory, "source."+format+".cueson.json")
				for _, action := range []string{"validate", "inspect"} {
					payloadBefore, err := os.ReadFile(input)
					if err != nil {
						return proof, err
					}
					if err := verifyOldConsumerIdentityGate(ctx, path, input, directory); err != nil {
						return proof, err
					}
					args := []string{action, input}
					if action == "inspect" {
						args = append(args, "--json")
					}
					status, stdout, stderr, err := invokeNative(ctx, path, args...)
					if err != nil {
						return proof, err
					}
					if err := verifyRefusal(status, stdout, stderr, "", deriveForbidden(directory, nil)); err != nil {
						return proof, err
					}
					proof.OldConsumerRefusalCount++
					proof.OldConsumerIdentityProbeCount++
					payloadAfter, err := os.ReadFile(input)
					if err != nil || !bytes.Equal(payloadBefore, payloadAfter) {
						return proof, fmt.Errorf("old consumer mutated candidate payload")
					}
				}
				for _, action := range []string{"restore", "render", "convert"} {
					// Old software only accepts text target selectors. Identity refusal is
					// tested before dispatch using the same valid selector for all sources.
					if err := verifyForcedRefusal(ctx, path, directory, input, action, "vtt", true); err != nil {
						return proof, err
					}
					proof.OldConsumerRefusalCount++
					proof.OldConsumerIdentityProbeCount++
					if err := verifyAbsentRefusal(ctx, path, directory, input, action, "vtt"); err != nil {
						return proof, err
					}
					proof.OldConsumerRefusalCount++
					proof.OldConsumerIdentityProbeCount++
				}
			}
			proof.OldConsumerVersion, proof.OldConsumerRevision, proof.OldConsumerArchive, proof.OldConsumerArchiveSHA256 = "1.0.0", contract.Revision, asset.Name, asset.SHA256
			return proof, nil
		}
	}
	return proof, fmt.Errorf("no frozen published archive matches host")
}

func verifyAbsentRefusal(ctx context.Context, binary, directory, input, action, target string) error {
	payloadBefore, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	if err := verifyOldConsumerIdentityGate(ctx, binary, input, directory); err != nil {
		return err
	}
	destination := filepath.Join(directory, "absent-destination")
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return fmt.Errorf("absent destination precondition failed")
	}
	args := []string{action, input, "--output", destination, "--force"}
	if action == "restore" {
		args = append(args, "--no-metadata")
	} else {
		args = append(args, "--to", target)
	}
	before, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	status, stdout, stderr, err := invokeNative(ctx, binary, args...)
	if err != nil {
		return err
	}
	if err := verifyRefusal(status, stdout, stderr, "", deriveForbidden(directory, nil)); err != nil {
		return err
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return fmt.Errorf("refusal published absent destination")
	}
	after, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(before, after) {
		return fmt.Errorf("refusal created staging artifacts")
	}
	payloadAfter, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(payloadBefore, payloadAfter) {
		return fmt.Errorf("old consumer mutated candidate payload")
	}
	return nil
}

// Released commands such as inspect intentionally suppress validation detail.
// Prove the exact identity rejection through its validate command immediately
// before each owning command, against the same unchanged payload. Generic
// runtime refusal alone never establishes this compatibility gate.
func verifyOldConsumerIdentityGate(ctx context.Context, binary, input, directory string) error {
	before, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	if err := verifyMinimalIdentityGate(ctx, binary, input, directory, "1.0.0"); err != nil {
		return err
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	status, stdout, stderr, err := invokeNative(ctx, binary, "validate", input)
	if err != nil {
		return err
	}
	if err := verifyRefusal(status, stdout, stderr, "", deriveForbidden(directory, nil)); err != nil {
		return err
	}
	if err := verifyOldIdentityDiagnostic(stderr); err != nil {
		return err
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(before, after) {
		return fmt.Errorf("old identity probe mutated candidate input")
	}
	afterEntries, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(entries, afterEntries) {
		return fmt.Errorf("full-document identity probe created publication or staging artifacts")
	}
	return nil
}

func verifyOldIdentityDiagnostic(stderr []byte) error {
	text := string(stderr)
	if !strings.Contains(text, "1.0.0") || (!strings.Contains(text, "$schema") && !strings.Contains(text, "schema_version")) {
		return fmt.Errorf("published consumer did not reject the exact schema identity: %q", text)
	}
	return nil
}
