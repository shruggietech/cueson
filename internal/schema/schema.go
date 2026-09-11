// Package schema owns the embedded Cue JSON schema and its validation boundary.
package schema

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode/utf8"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/shruggietech/cueson/internal/model"
)

const (
	schemaID      = "https://cueson.io/schema/v1.0.0/cueson.schema.json"
	schemaVersion = "1.0.0"
)

var (
	//go:embed cueson.schema.json
	canonicalBytes []byte

	//go:embed testdata/representative.cueson.json
	representativeBytes []byte

	compileOnce   sync.Once
	compiled      *jsonschema.Schema
	compiledError error
)

// ID returns the current canonical schema identifier.
func ID() string {
	return schemaID
}

// Version returns the embedded Cue JSON contract version.
func Version() string {
	return schemaVersion
}

// Bytes returns a copy of the canonical embedded schema bytes.
func Bytes() []byte {
	return bytes.Clone(canonicalBytes)
}

// Representative returns a copy of the canonical representative document.
func Representative() []byte {
	return bytes.Clone(representativeBytes)
}

// Compiled returns the process-wide compiled Draft 2020-12 schema.
func Compiled() (*jsonschema.Schema, error) {
	compileOnce.Do(func() {
		compiled, compiledError = compileCanonical()
	})
	return compiled, compiledError
}

// Validate applies JSON parsing, structural validation, and model semantics.
func Validate(data []byte) error {
	_, err := Decode(data)
	return err
}

// Decode validates Cue JSON and returns its typed document.
func Decode(data []byte) (model.Document, error) {
	var document model.Document
	if !utf8.Valid(data) {
		return document, fmt.Errorf("parse Cue JSON: input is not valid UTF-8")
	}
	if err := validateUnicodeEscapes(data); err != nil {
		return document, fmt.Errorf("parse Cue JSON: %w", err)
	}
	instance, err := decodeOne(data)
	if err != nil {
		return document, fmt.Errorf("parse Cue JSON: %w", err)
	}
	contract, err := Compiled()
	if err != nil {
		return document, fmt.Errorf("compile embedded schema: %w", err)
	}
	if err := contract.Validate(instance); err != nil {
		return document, fmt.Errorf("validate Cue JSON structure: %w", err)
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return document, fmt.Errorf("decode Cue JSON model: %w", err)
	}
	if err := document.Validate(); err != nil {
		return document, fmt.Errorf("validate Cue JSON semantics: %w", err)
	}
	return document, nil
}

func validateUnicodeEscapes(data []byte) error {
	inString := false
	for index := 0; index < len(data); index++ {
		switch data[index] {
		case '"':
			inString = !inString
		case '\\':
			if !inString || index+1 >= len(data) {
				continue
			}
			if data[index+1] != 'u' {
				index++
				continue
			}
			value, valid := decodeHexEscape(data, index+2)
			if !valid {
				continue
			}
			switch {
			case value >= 0xd800 && value <= 0xdbff:
				if index+12 > len(data) || data[index+6] != '\\' || data[index+7] != 'u' {
					return fmt.Errorf("unpaired high-surrogate Unicode escape at byte %d", index)
				}
				low, lowValid := decodeHexEscape(data, index+8)
				if !lowValid || low < 0xdc00 || low > 0xdfff {
					return fmt.Errorf("unpaired high-surrogate Unicode escape at byte %d", index)
				}
				index += 11
			case value >= 0xdc00 && value <= 0xdfff:
				return fmt.Errorf("unpaired low-surrogate Unicode escape at byte %d", index)
			default:
				index += 5
			}
		}
	}
	return nil
}

func decodeHexEscape(data []byte, start int) (uint16, bool) {
	if start+4 > len(data) {
		return 0, false
	}
	var value uint16
	for _, character := range data[start : start+4] {
		value <<= 4
		switch {
		case character >= '0' && character <= '9':
			value += uint16(character - '0')
		case character >= 'a' && character <= 'f':
			value += uint16(character-'a') + 10
		case character >= 'A' && character <= 'F':
			value += uint16(character-'A') + 10
		default:
			return 0, false
		}
	}
	return value, true
}

// CheckLockstep verifies that official software and schema versions match.
func CheckLockstep(softwareVersion string) error {
	if softwareVersion != schemaVersion {
		return fmt.Errorf("software version %q does not match schema version %q", softwareVersion, schemaVersion)
	}
	if _, err := Compiled(); err != nil {
		return fmt.Errorf("schema identity check failed: %w", err)
	}
	return nil
}

func compileCanonical() (*jsonschema.Schema, error) {
	artifact, err := decodeOne(canonicalBytes)
	if err != nil {
		return nil, fmt.Errorf("parse canonical schema: %w", err)
	}
	header, ok := artifact.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("canonical schema root is not an object")
	}
	if err := validateArtifactIdentity(header); err != nil {
		return nil, err
	}

	compiler := jsonschema.NewCompiler()
	compiler.AssertContent()
	compiler.AssertFormat()
	if err := compiler.AddResource(schemaID, artifact); err != nil {
		return nil, fmt.Errorf("register canonical schema: %w", err)
	}
	contract, err := compiler.Compile(schemaID)
	if err != nil {
		return nil, fmt.Errorf("compile canonical schema: %w", err)
	}
	if contract.DraftVersion != 2020 {
		return nil, fmt.Errorf("canonical schema draft is %d, want 2020", contract.DraftVersion)
	}
	return contract, nil
}

func validateArtifactIdentity(header map[string]any) error {
	if !strings.HasSuffix(schemaID, "/v"+schemaVersion+"/cueson.schema.json") {
		return fmt.Errorf("schema identifier %q does not contain schema version %q", schemaID, schemaVersion)
	}
	if got := stringValue(header["$id"]); got != schemaID {
		return fmt.Errorf("canonical schema $id is %q, want %q", got, schemaID)
	}
	properties, ok := header["properties"].(map[string]any)
	if !ok {
		return fmt.Errorf("canonical schema properties are missing")
	}
	instanceSchemaProperty, ok := properties["$schema"].(map[string]any)
	if !ok || stringValue(instanceSchemaProperty["const"]) != schemaID {
		return fmt.Errorf("canonical instance $schema const must be %q", schemaID)
	}
	versionProperty, ok := properties["schema_version"].(map[string]any)
	if !ok || stringValue(versionProperty["const"]) != schemaVersion {
		return fmt.Errorf("canonical schema schema_version const must be %q", schemaVersion)
	}
	return nil
}

func decodeOne(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values")
		}
		return nil, fmt.Errorf("trailing data: %w", err)
	}
	return value, nil
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}
