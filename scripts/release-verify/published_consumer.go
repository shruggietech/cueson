package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type PublishedConsumerProof struct {
	Version                    string `json:"version"`
	Revision                   string `json:"source_revision"`
	Archive                    string `json:"archive"`
	ArchiveSHA256              string `json:"archive_sha256"`
	PositiveControlCount       int    `json:"positive_control_count"`
	RefusalCount               int    `json:"refusal_count"`
	IdentityProbeCount         int    `json:"identity_probe_count"`
	SelectorIdentityProbeCount int    `json:"selector_identity_probe_count"`
}

type publishedConsumerAsset struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	GOOS   string `json:"goos,omitempty"`
	GOARCH string `json:"goarch,omitempty"`
}

type publishedConsumerContract struct {
	Version       string                   `json:"version"`
	Revision      string                   `json:"source_revision"`
	URL           string                   `json:"release_url"`
	Manifest      publishedConsumerAsset   `json:"checksum_manifest"`
	SchemaSHA256  string                   `json:"release_schema_sha256"`
	LicenseSHA256 string                   `json:"license_sha256"`
	NoticeSHA256  string                   `json:"notice_sha256"`
	Assets        []publishedConsumerAsset `json:"assets"`
}

func loadPublishedV110Contract(repository string) (publishedConsumerContract, error) {
	var contract publishedConsumerContract
	data, err := readRegularFile(filepath.Join(repository, "specs", "S034-prepare-v1-2-release-candidate", "contracts", "published-v1.1.0-consumer.json"))
	if err != nil {
		return contract, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&contract); err != nil {
		return contract, err
	}
	if contract.Version != "1.1.0" || contract.Revision != "7ff45c1d8cd8df377e1fb568b9785286b649fd7c" || contract.URL != "https://github.com/shruggietech/cueson/releases/tag/v1.1.0" || contract.Manifest.Name != "cueson_1.1.0_checksums.txt" || contract.Manifest.Size != 588 || contract.Manifest.SHA256 != "57089835e30184f9625d5b9b493ebb41d58ca5ac7ec76af84d0c72aa8ca5e908" || contract.SchemaSHA256 != "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7" || contract.LicenseSHA256 != "4e2221f7a8060092f01e563dd2f399fa573a17823e0d25d34a5c28df75a6a654" || contract.NoticeSHA256 != "9f7a0dfad459dd180b4d0ebd8c616afd8c710d3cd8322c5ec59d294eeabae2af" {
		return contract, fmt.Errorf("published 1.1 consumer identity binding changed")
	}
	want := []publishedConsumerAsset{
		{"cueson_1.1.0_linux_amd64.tar.gz", 2244709, "689738ae27098f67bb9308ac690e9e3304ba732464c98835fff34fe7473406a1", "linux", "amd64"},
		{"cueson_1.1.0_darwin_amd64.tar.gz", 2260001, "c48225a74eff45956fbdafc4523b75bc546ba22334aa133aecda3abfe767bcca", "darwin", "amd64"},
		{"cueson_1.1.0_windows_amd64.zip", 2306834, "2aff230ca03127eeca49d4be447b5a49563026caab96b5c545189a4181e8e672", "windows", "amd64"},
	}
	if len(contract.Assets) != len(want) {
		return contract, fmt.Errorf("published 1.1 host archive count differs")
	}
	for i, asset := range contract.Assets {
		if asset != want[i] {
			return contract, fmt.Errorf("published 1.1 host archive binding changed")
		}
	}
	return contract, nil
}

func downloadPublishedConsumerAsset(ctx context.Context, asset publishedConsumerAsset) ([]byte, error) {
	if asset.Size <= 0 || asset.Size > 16<<20 {
		return nil, fmt.Errorf("invalid bounded published asset length")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://github.com/shruggietech/cueson/releases/download/v1.1.0/"+asset.Name, nil)
	if err != nil {
		return nil, err
	}
	response, err := (&http.Client{Timeout: 90 * time.Second}).Do(request)
	if err != nil {
		return nil, fmt.Errorf("download published 1.1 consumer: %w", err)
	}
	data, readErr := io.ReadAll(io.LimitReader(response.Body, asset.Size+1))
	closeErr := response.Body.Close()
	if response.StatusCode != http.StatusOK || readErr != nil || closeErr != nil || int64(len(data)) != asset.Size || digestHex(data) != asset.SHA256 {
		return nil, fmt.Errorf("published 1.1 asset HTTP status, length or digest differs")
	}
	return data, nil
}

func verifyVersionedIdentityDiagnostic(stderr []byte, version string) error {
	if version == "1.0.0" {
		return verifyOldIdentityDiagnostic(stderr)
	}
	if version != "1.1.0" || !strings.Contains(string(stderr), "select exact contract: unsupported or mismatched Cue JSON identity ($schema, schema_version)") {
		return fmt.Errorf("published consumer did not reject exact contract selection")
	}
	return nil
}

// Probe only the exact envelope selectors so source integrity and native data
// cannot explain the old consumer's rejection of a current identity.
func minimalIdentityPayload(payload []byte) ([]byte, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		return nil, err
	}
	if len(fields["$schema"]) == 0 || len(fields["schema_version"]) == 0 {
		return nil, fmt.Errorf("identity probe requires both exact selectors")
	}
	return json.Marshal(map[string]json.RawMessage{"$schema": fields["$schema"], "schema_version": fields["schema_version"]})
}

func verifyMinimalIdentityGate(ctx context.Context, binary, input, directory, version string) error {
	before, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	minimal, err := minimalIdentityPayload(before)
	if err != nil {
		return err
	}
	path := filepath.Join(directory, "selector-identity-probe.json")
	if err := os.WriteFile(path, minimal, 0600); err != nil {
		return err
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	status, stdout, stderr, err := invokeNative(ctx, binary, "validate", path)
	if err != nil {
		return err
	}
	if err := verifyRefusal(status, stdout, stderr, "", append(deriveForbidden(directory, nil), proofSpeakerID)); err != nil {
		return err
	}
	if err := verifyVersionedIdentityDiagnostic(stderr, version); err != nil {
		return err
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(minimal, after) {
		return fmt.Errorf("minimal selector probe mutated its payload")
	}
	after, err = os.ReadFile(input)
	if err != nil || !bytes.Equal(before, after) {
		return fmt.Errorf("minimal selector probe mutated the owning document")
	}
	afterEntries, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(entries, afterEntries) {
		return fmt.Errorf("minimal selector probe created publication or staging artifacts")
	}
	return nil
}

func verifyVersionedIdentityGate(ctx context.Context, binary, input, directory, version string) error {
	before, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	if err := verifyMinimalIdentityGate(ctx, binary, input, directory, version); err != nil {
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
	if err := verifyRefusal(status, stdout, stderr, "", append(deriveForbidden(directory, nil), proofSpeakerID)); err != nil {
		return err
	}
	if err := verifyVersionedIdentityDiagnostic(stderr, version); err != nil {
		return err
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(before, after) {
		return fmt.Errorf("old consumer identity probe mutated input")
	}
	afterEntries, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(entries, afterEntries) {
		return fmt.Errorf("full-document identity probe created publication or staging artifacts")
	}
	return nil
}

func verifyPublishedV110Consumer(ctx context.Context, repository, directory string) (PublishedConsumerProof, error) {
	var proof PublishedConsumerProof
	contract, err := loadPublishedV110Contract(repository)
	if err != nil {
		return proof, err
	}
	manifest, err := downloadPublishedConsumerAsset(ctx, contract.Manifest)
	if err != nil {
		return proof, err
	}
	checksums, err := parseChecksums(manifest)
	if err != nil {
		return proof, err
	}
	if len(checksums) != 6 {
		return proof, fmt.Errorf("published 1.1 checksum manifest archive count differs")
	}
	for _, asset := range contract.Assets {
		if checksums[asset.Name] != asset.SHA256 {
			return proof, fmt.Errorf("published archive differs from authenticated checksum manifest")
		}
	}
	for _, asset := range contract.Assets {
		if asset.GOOS != runtime.GOOS || asset.GOARCH != runtime.GOARCH {
			continue
		}
		archive, err := downloadPublishedConsumerAsset(ctx, asset)
		if err != nil {
			return proof, err
		}
		members, err := readArchive(asset.Name, archive)
		if err != nil {
			return proof, err
		}
		target := Target{GOOS: asset.GOOS, GOARCH: asset.GOARCH, Archive: asset.Name, Binary: "cueson"}
		if asset.GOOS == "windows" {
			target.Binary += ".exe"
		}
		schema, err := readRegularFile(filepath.Join(repository, "schema", "releases", "v1.1.0", releaseSchemaFilename))
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
		if digestHex(schema) != contract.SchemaSHA256 || digestHex(license) != contract.LicenseSHA256 || digestHex(notice) != contract.NoticeSHA256 {
			return proof, fmt.Errorf("published 1.1 packaged schema or legal binding changed")
		}
		binary, err := verifyMembers(target, members, schema, license, notice, contract.Version)
		if err != nil {
			return proof, err
		}
		if err := verifyBuildInfo(binary, target, contract.Revision, deriveForbidden(repository, nil)); err != nil {
			return proof, err
		}
		path := filepath.Join(directory, "old-1.1.0-"+target.Binary)
		if err := os.WriteFile(path, binary, 0700); err != nil {
			return proof, err
		}
		version, err := nativeSuccess(ctx, path, "version")
		if err != nil || string(version) != "1.1.0\n" {
			return proof, fmt.Errorf("published 1.1 executable version differs")
		}
		// Positive controls prove this authentic executable accepts its own contract.
		for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
			input := filepath.Join(repository, "internal", "schema", "testdata", "historical-v1.1.0-"+format+".json")
			if _, err := nativeSuccess(ctx, path, "validate", input); err != nil {
				return proof, err
			}
			proof.PositiveControlCount++
		}
		for _, format := range []string{"srt", "vtt", "ass", "ssa"} {
			for _, input := range []string{filepath.Join(directory, "source."+format+".cueson.json"), filepath.Join(directory, "annotated-"+format+".json")} {
				for _, action := range []string{"validate", "inspect"} {
					if err := verifyVersionedIdentityGate(ctx, path, input, directory, contract.Version); err != nil {
						return proof, err
					}
					before, err := os.ReadFile(input)
					if err != nil {
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
					if err := verifyRefusal(status, stdout, stderr, "", append(deriveForbidden(directory, nil), proofSpeakerID)); err != nil {
						return proof, err
					}
					after, err := os.ReadFile(input)
					if err != nil || !bytes.Equal(before, after) {
						return proof, fmt.Errorf("published 1.1 consumer mutated candidate payload")
					}
					proof.IdentityProbeCount++
					proof.SelectorIdentityProbeCount++
					proof.RefusalCount++
				}
				for _, action := range []string{"restore", "render", "convert"} {
					for _, occupied := range []bool{true, false} {
						if err := verifyVersionedIdentityGate(ctx, path, input, directory, contract.Version); err != nil {
							return proof, err
						}
						if err := verifyNativeDestinationRefusal(ctx, path, directory, input, action, "vtt", occupied); err != nil {
							return proof, err
						}
						proof.IdentityProbeCount++
						proof.SelectorIdentityProbeCount++
						proof.RefusalCount++
					}
				}
			}
		}
		proof.Version, proof.Revision, proof.Archive, proof.ArchiveSHA256 = contract.Version, contract.Revision, asset.Name, asset.SHA256
		return proof, nil
	}
	return proof, fmt.Errorf("no authenticated published 1.1 archive matches native host")
}
