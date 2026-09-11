package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/version"
)

func TestCanonicalSchemaAndRepresentative(t *testing.T) {
	t.Parallel()

	if _, err := Compiled(); err != nil {
		t.Fatalf("Compiled() error = %v", err)
	}
	if err := Validate(Representative()); err != nil {
		t.Fatalf("Validate(Representative()) error = %v", err)
	}
	if got, want := ID(), "https://cueson.io/schema/v0.1.0/cueson.schema.json"; got != want {
		t.Errorf("ID() = %q, want %q", got, want)
	}
	if got, want := Version(), "0.1.0"; got != want {
		t.Errorf("Version() = %q, want %q", got, want)
	}
	if bytes.HasPrefix(Bytes(), []byte{0xef, 0xbb, 0xbf}) {
		t.Error("Bytes() has UTF-8 BOM")
	}
	if !bytes.HasSuffix(Bytes(), []byte("\n")) {
		t.Error("Bytes() does not end with LF")
	}
}

func TestBytesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	first := Bytes()
	first[0] = 'x'
	if bytes.Equal(first, Bytes()) {
		t.Error("Bytes() exposes mutable embedded storage")
	}
}

func TestValidateRejectsStructuralViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(map[string]any)
		want   string
	}{
		{name: "schema version", mutate: func(doc map[string]any) { doc["schema_version"] = "0.0.0" }, want: "schema_version"},
		{name: "subrip envelope capability", mutate: func(doc map[string]any) {
			doc["format_support"] = envelopeOnlySupport()
		}, want: "format_support"},
		{name: "unsafe path", mutate: func(doc map[string]any) { firstAsset(doc)["file_name"] = "../captions.srt" }, want: "file_name"},
		{name: "reserved device", mutate: func(doc map[string]any) { firstAsset(doc)["file_name"] = "CoN.txt" }, want: "file_name"},
		{name: "superscript reserved device", mutate: func(doc map[string]any) { firstAsset(doc)["file_name"] = "COM¹.txt" }, want: "file_name"},
		{name: "trailing period", mutate: func(doc map[string]any) { firstAsset(doc)["file_name"] = "captions." }, want: "file_name"},
		{name: "uppercase digest", mutate: func(doc map[string]any) {
			firstAsset(doc)["hashes"].(map[string]any)["sha256"] = strings.Repeat("A", 64)
		}, want: "sha256"},
		{name: "bad base64", mutate: func(doc map[string]any) { firstAsset(doc)["data_base64"] = "%%%" }, want: "data_base64"},
		{name: "missing payload lines", mutate: func(doc map[string]any) { delete(firstCue(doc)["payload"].(map[string]any), "lines") }, want: "lines"},
		{name: "extra root property", mutate: func(doc map[string]any) { doc["source_path"] = "secret" }, want: "source_path"},
		{name: "mismatched cue format", mutate: func(doc map[string]any) {
			firstCue(doc)["format_data"] = map[string]any{"webvtt": map[string]any{"identifier_raw": nil, "timing_line_raw": "x", "settings_raw": "", "settings": map[string]any{}, "raw_payload": "x"}}
		}, want: "format_data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := representativeMap(t)
			tt.mutate(doc)
			encoded, err := json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			err = Validate(encoded)
			if err == nil {
				t.Fatal("Validate() error = nil, want rejection")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Validate() error = %q, want mention of %q", err, tt.want)
			}
		})
	}
}

func TestValidateRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	for _, input := range [][]byte{[]byte("{"), append(append([]byte{}, Representative()...), []byte("{}")...)} {
		if err := Validate(input); err == nil {
			t.Errorf("Validate(%q) error = nil, want rejection", input)
		}
	}
}

func TestDecodeReturnsTypedDocumentAndRejectsInvalidUTF8(t *testing.T) {
	t.Parallel()

	document, err := Decode(Representative())
	if err != nil {
		t.Fatalf("Decode(Representative()) error = %v", err)
	}
	if document.SchemaVersion != Version() || !document.FormatSupport.RestoreSupported {
		t.Fatalf("Decode() returned unexpected document: version = %q, restore = %t", document.SchemaVersion, document.FormatSupport.RestoreSupported)
	}
	invalid := append([]byte{}, Representative()...)
	invalid[len(invalid)-2] = 0xff
	if _, err := Decode(invalid); err == nil || !strings.Contains(err.Error(), "UTF-8") {
		t.Fatalf("Decode(invalid UTF-8) error = %v, want UTF-8 rejection", err)
	}
}

func TestDecodeRejectsUnpairedSurrogateEscapes(t *testing.T) {
	t.Parallel()

	for _, escapedName := range []string{`bad\ud800.srt`, `bad\udc00.srt`, `bad\ud800\u0041.srt`} {
		input := bytes.Replace(schemaRepresentativeForEscapeTest(), []byte("captions.srt"), []byte(escapedName), 1)
		if _, err := Decode(input); err == nil || !strings.Contains(err.Error(), "surrogate") {
			t.Errorf("Decode(%q) error = %v, want surrogate rejection", escapedName, err)
		}
	}
	paired := bytes.Replace(schemaRepresentativeForEscapeTest(), []byte("captions.srt"), []byte(`face\ud83d\ude00.srt`), 1)
	if _, err := Decode(paired); err != nil {
		t.Fatalf("Decode(paired surrogate) error = %v", err)
	}
	if err := validateUnicodeEscapes([]byte(`{"value":"literal\\ud800"}`)); err != nil {
		t.Fatalf("validateUnicodeEscapes(escaped literal) error = %v", err)
	}
}

func schemaRepresentativeForEscapeTest() []byte {
	return append([]byte(nil), Representative()...)
}

func TestValidateRejectsSemanticViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(map[string]any)
		want   string
	}{
		{name: "unresolved primary", mutate: func(doc map[string]any) { doc["source"].(map[string]any)["primary_asset_id"] = "missing" }, want: "primary_asset_id"},
		{name: "cue duration", mutate: func(doc map[string]any) {
			firstCue(doc)["timing"].(map[string]any)["duration_milliseconds"] = float64(1)
		}, want: "duration_milliseconds"},
		{name: "timestamp instant", mutate: func(doc map[string]any) {
			firstAsset(doc)["timestamps"].(map[string]any)["modified"].(map[string]any)["unix_ns"] = float64(0)
		}, want: "timestamps.modified"},
		{name: "summary count", mutate: func(doc map[string]any) { doc["stats"].(map[string]any)["cue_count"] = float64(2) }, want: "cue_count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := representativeMap(t)
			tt.mutate(doc)
			encoded, err := json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			err = Validate(encoded)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Validate() error = %v, want semantic rejection containing %q", err, tt.want)
			}
		})
	}
}

func TestLockstep(t *testing.T) {
	t.Parallel()

	if err := CheckLockstep(version.String()); err != nil {
		t.Fatalf("CheckLockstep(version.String()) error = %v", err)
	}
	if err := CheckLockstep("0.0.0"); err == nil || !strings.Contains(err.Error(), "0.0.0") {
		t.Fatalf("CheckLockstep(0.0.0) error = %v, want version mismatch", err)
	}
}

func TestCanonicalIdentityRejectsDrift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(map[string]any)
		want   string
	}{
		{name: "artifact id", mutate: func(artifact map[string]any) { artifact["$id"] = "https://cueson.io/schema/v0.0.0/cueson.schema.json" }, want: "$id"},
		{name: "instance schema", mutate: func(artifact map[string]any) {
			artifact["properties"].(map[string]any)["$schema"].(map[string]any)["const"] = "https://cueson.io/schema/v0.0.0/cueson.schema.json"
		}, want: "$schema"},
		{name: "schema version", mutate: func(artifact map[string]any) {
			artifact["properties"].(map[string]any)["schema_version"].(map[string]any)["const"] = "0.0.0"
		}, want: "schema_version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var artifact map[string]any
			if err := json.Unmarshal(Bytes(), &artifact); err != nil {
				t.Fatal(err)
			}
			tt.mutate(artifact)
			err := validateArtifactIdentity(artifact)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("validateArtifactIdentity() error = %v, want drift rejection containing %q", err, tt.want)
			}
		})
	}
}

func TestValidateRecognizedFormatsAndMultipleAssets(t *testing.T) {
	t.Parallel()

	t.Run("multiple assets", func(t *testing.T) {
		doc := representativeMap(t)
		asset := firstAsset(doc)
		companion := make(map[string]any, len(asset))
		for key, value := range asset {
			companion[key] = value
		}
		companion["id"] = "asset-1"
		companion["role"] = "companion"
		companion["file_name"] = "captions.bin"
		doc["source"].(map[string]any)["assets"] = append(doc["source"].(map[string]any)["assets"].([]any), companion)
		assertMapValid(t, doc)
	})

	t.Run("webvtt", func(t *testing.T) {
		doc := representativeMap(t)
		doc["format"] = "webvtt"
		doc["format_support"] = envelopeOnlySupport()
		firstAsset(doc)["file_name"] = "captions.vtt"
		firstAsset(doc)["media_type"] = "text/vtt"
		doc["format_data"] = map[string]any{"webvtt": map[string]any{"signature": "WEBVTT", "description": nil, "metadata_lines": []any{}, "blocks": []any{}}}
		firstCue(doc)["format_data"] = map[string]any{"webvtt": map[string]any{"identifier_raw": nil, "timing_line_raw": "00:00:01.250 --> 00:00:04.200", "settings_raw": "", "settings": map[string]any{}, "raw_payload": "Hello, world."}}
		assertMapValid(t, doc)
	})

	t.Run("ocr observation", func(t *testing.T) {
		doc := representativeMap(t)
		firstCue(doc)["ocr_observations"] = []any{map[string]any{"id": "ocr-0", "derived": true, "engine": "test-engine", "text": "Hello, world.", "lines": []any{"Hello, world."}, "source_asset_id": "asset-0"}}
		assertMapValid(t, doc)
	})
}

func envelopeOnlySupport() map[string]any {
	return map[string]any{
		"status": "envelope_only", "ingest_supported": false, "render_supported": false,
		"restore_supported": true, "ocr_required_for_semantic_output": false,
	}
}

func TestCanonicalProjectNamesAreLowerSnakeCase(t *testing.T) {
	t.Parallel()

	var artifact any
	if err := json.Unmarshal(Bytes(), &artifact); err != nil {
		t.Fatal(err)
	}
	if err := checkProjectNames(artifact, "$"); err != nil {
		t.Fatal(err)
	}
}

func checkProjectNames(value any, path string) error {
	object, ok := value.(map[string]any)
	if !ok {
		if values, ok := value.([]any); ok {
			for _, item := range values {
				if err := checkProjectNames(item, path); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if properties, ok := object["properties"].(map[string]any); ok {
		for name := range properties {
			if name != "$schema" && !isLowerSnakeCase(name) {
				return &nameError{path: path, name: name}
			}
		}
	}
	if values, ok := object["enum"].([]any); ok && !strings.Contains(path, "/safe_basename/") {
		for _, raw := range values {
			value, ok := raw.(string)
			if ok && !isLowerSnakeCase(value) {
				return &nameError{path: path + "/enum", name: value}
			}
		}
	}
	for key, child := range object {
		if err := checkProjectNames(child, path+"/"+key); err != nil {
			return err
		}
	}
	return nil
}

func TestProjectNameCheckerRejectsDrift(t *testing.T) {
	t.Parallel()

	for _, artifact := range []any{
		map[string]any{"properties": map[string]any{"camelCase": map[string]any{}}},
		map[string]any{"enum": []any{"Not_Snake"}},
	} {
		if err := checkProjectNames(artifact, "$"); err == nil {
			t.Errorf("checkProjectNames(%v) error = nil, want rejection", artifact)
		}
	}
}

func isLowerSnakeCase(value string) bool {
	if value == "" || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	previousUnderscore := false
	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z', character >= '0' && character <= '9':
			previousUnderscore = false
		case character == '_' && !previousUnderscore:
			previousUnderscore = true
		default:
			return false
		}
	}
	return !previousUnderscore
}

type nameError struct {
	path string
	name string
}

func (err *nameError) Error() string {
	return "non-snake-case project property " + err.name + " at " + err.path
}

func representativeMap(t *testing.T) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(Representative(), &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func firstAsset(doc map[string]any) map[string]any {
	return doc["source"].(map[string]any)["assets"].([]any)[0].(map[string]any)
}

func firstCue(doc map[string]any) map[string]any {
	return doc["cues"].([]any)[0].(map[string]any)
}

func assertMapValid(t *testing.T, doc map[string]any) {
	t.Helper()
	encoded, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(encoded); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
