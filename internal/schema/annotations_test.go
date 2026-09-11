package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/shruggietech/cueson/internal/testutil"
)

const (
	normativeSchemaSHA256 = "33ef838ab89ab491d26bdf8e97ddd6a0828542e173e011be3ce1293305323902"
	releasedSchemaSHA256  = "d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975"
)

var identifyingExample = regexp.MustCompile(`(?i)(?:^|["'\s])(?:[a-z]:[\\/]|/(?:home|users)/|\\\\[^\\]+\\|file://|localhost)`)

func TestAnnotationCoverage(t *testing.T) {
	t.Parallel()

	artifact := annotationArtifact(t)
	reachable := reachableDefinitions(t, artifact)
	if got, want := len(reachable), 44; got != want {
		t.Fatalf("reachable definition count = %d, want %d", got, want)
	}

	checkPropertyDescriptions(t, "$", artifact["properties"])
	definitions := artifact["$defs"].(map[string]any)
	for _, name := range reachable {
		definition := definitions[name].(map[string]any)
		requireAnnotationText(t, "$/$defs/"+name, definition, "title")
		requireAnnotationText(t, "$/$defs/"+name, definition, "description")
		if properties, ok := definition["properties"]; ok {
			checkPropertyDescriptions(t, "$/$defs/"+name, properties)
		}
	}

	for _, pointer := range enumPointers(artifact, reachable) {
		node := schemaNodeAt(t, artifact, pointer)
		examples, ok := node["examples"].([]any)
		if !ok || len(examples) == 0 {
			t.Errorf("%s enum has no examples", pointer)
		}
	}

	requiredExamples := []string{
		"$/properties/$schema",
		"$/properties/schema_version",
		"$/properties/format",
		"$/properties/format_support",
		"$/properties/producer",
		"$/properties/source",
		"$/properties/metadata",
		"$/properties/document",
		"$/properties/cues",
		"$/properties/format_data",
		"$/properties/diagnostics",
		"$/properties/stats",
		"$/$defs/identifier",
		"$/$defs/safe_basename",
		"$/$defs/asset_hashes/properties/sha256",
		"$/$defs/source_asset/properties/data_base64",
		"$/$defs/timestamp",
		"$/$defs/timing",
		"$/$defs/diagnostic",
		"$/$defs/subrip_cue_data",
		"$/$defs/subrip_document_data",
		"$/$defs/webvtt_cue_data",
		"$/$defs/webvtt_document_data",
		"$/$defs/webvtt_block",
		"$/$defs/webvtt_region_data",
	}
	for _, pointer := range requiredExamples {
		node := schemaNodeAt(t, artifact, pointer)
		if examples, ok := node["examples"].([]any); !ok || len(examples) == 0 {
			t.Errorf("%s has no required examples", pointer)
		}
	}
}

func TestAnnotationExamplesValidate(t *testing.T) {
	t.Parallel()

	contract, err := Compiled()
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[*jsonschema.Schema]bool)
	validated := 0
	walkCompiledSchema(contract, seen, func(fragment *jsonschema.Schema) {
		for index, example := range fragment.Examples {
			validated++
			if err := fragment.Validate(example); err != nil {
				t.Errorf("example %d at %s is invalid: %v", index, fragment.Location, err)
			}
			checkExamplePrivacy(t, fragment.Location, example)
		}
	})
	if validated == 0 {
		t.Fatal("compiled schema exposed no examples")
	}

	for index, example := range contract.Examples {
		encoded, err := json.Marshal(example)
		if err != nil {
			t.Fatalf("marshal root example %d: %v", index, err)
		}
		if err := Validate(encoded); err != nil {
			t.Errorf("root example %d fails structural or semantic validation: %v", index, err)
		}
	}
}

func TestAnnotationIdentityExamplesAreCurrent(t *testing.T) {
	t.Parallel()

	artifact := annotationArtifact(t)
	assertExamplesEqual(t, schemaNodeAt(t, artifact, "$/properties/$schema"), ID())
	assertExamplesEqual(t, schemaNodeAt(t, artifact, "$/properties/schema_version"), Version())
	for index, raw := range artifact["examples"].([]any) {
		example := raw.(map[string]any)
		if example["$schema"] != ID() || example["schema_version"] != Version() {
			t.Errorf("root example %d has stale schema identity", index)
		}
	}
}

func TestAnnotationExamplesPreserveNativeFidelity(t *testing.T) {
	t.Parallel()

	artifact := annotationArtifact(t)
	for index, raw := range artifact["examples"].([]any) {
		example := raw.(map[string]any)
		for cueIndex, rawCue := range example["cues"].([]any) {
			cue := rawCue.(map[string]any)
			payload := cue["payload"].(map[string]any)
			lines := payload["lines"].([]any)
			textLines := make([]string, len(lines))
			for lineIndex, line := range lines {
				textLines[lineIndex] = line.(string)
			}
			if got, want := strings.Join(textLines, "\n"), payload["raw_text"]; got != want {
				t.Errorf("root example %d cue %d payload lines join to %q, want raw_text %q", index, cueIndex, got, want)
			}
		}
	}

	regionExamples := schemaNodeAt(t, artifact, "$/$defs/webvtt_region_data")["examples"].([]any)
	for index, raw := range regionExamples {
		region := raw.(map[string]any)
		rawSettings := strings.Fields(region["settings_raw"].(string))
		occurrences := region["setting_occurrences"].([]any)
		if len(occurrences) != len(rawSettings) {
			t.Errorf("REGION example %d has %d raw settings and %d occurrences", index, len(rawSettings), len(occurrences))
			continue
		}
		effective := make(map[string]string)
		for occurrenceIndex, rawOccurrence := range occurrences {
			occurrence := rawOccurrence.(map[string]any)
			if occurrence["raw"] != rawSettings[occurrenceIndex] {
				t.Errorf("REGION example %d occurrence %d raw = %v, want %q", index, occurrenceIndex, occurrence["raw"], rawSettings[occurrenceIndex])
			}
			name, value, found := strings.Cut(rawSettings[occurrenceIndex], ":")
			if !found || occurrence["name"] != name || occurrence["value"] != value {
				t.Errorf("REGION example %d occurrence %d does not parse %q faithfully", index, occurrenceIndex, rawSettings[occurrenceIndex])
			}
			if occurrence["recognized"] == true && occurrence["valid"] == true {
				effective[name] = value
			}
		}
		settings := region["settings"].(map[string]any)
		if len(settings) != len(effective) {
			t.Errorf("REGION example %d has %d effective settings and %d valid recognized occurrences", index, len(settings), len(effective))
		}
		for name, value := range effective {
			if settings[name] != value {
				t.Errorf("REGION example %d effective setting %q = %v, want %q", index, name, settings[name], value)
			}
		}
	}
}

func TestAnnotationsDoNotChangeNormativeSchema(t *testing.T) {
	t.Parallel()

	artifact := annotationArtifact(t)
	stripAnnotations(artifact)
	encoded, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	got := sha256.Sum256(encoded)
	if hex.EncodeToString(got[:]) != normativeSchemaSHA256 {
		t.Fatalf("normative schema digest = %s, want %s", hex.EncodeToString(got[:]), normativeSchemaSHA256)
	}
}

func TestReleasedSchemaRemainsImmutable(t *testing.T) {
	t.Parallel()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate annotation test")
	}
	path := filepath.Join(filepath.Dir(currentFile), "..", "..", "schema", "releases", "v0.0.0", "cueson.schema.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := sha256.Sum256(data)
	if hex.EncodeToString(got[:]) != releasedSchemaSHA256 {
		t.Fatalf("released v0.0.0 schema digest = %s, want %s", hex.EncodeToString(got[:]), releasedSchemaSHA256)
	}
}

func annotationArtifact(t *testing.T) map[string]any {
	t.Helper()
	value, err := decodeOne(Bytes())
	if err != nil {
		t.Fatal(err)
	}
	artifact, ok := value.(map[string]any)
	if !ok {
		t.Fatal("canonical schema is not an object")
	}
	return artifact
}

func reachableDefinitions(t *testing.T, artifact map[string]any) []string {
	t.Helper()
	definitions := artifact["$defs"].(map[string]any)
	root := make(map[string]any, len(artifact)-1)
	for key, value := range artifact {
		if key != "$defs" {
			root[key] = value
		}
	}
	queue := localDefinitionReferences(root)
	seen := make(map[string]bool)
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		if seen[name] {
			continue
		}
		definition, ok := definitions[name]
		if !ok {
			t.Fatalf("reference to missing local definition %q", name)
		}
		seen[name] = true
		queue = append(queue, localDefinitionReferences(definition)...)
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func localDefinitionReferences(value any) []string {
	var result []string
	var walk func(any)
	walk = func(current any) {
		switch current := current.(type) {
		case map[string]any:
			if reference, ok := current["$ref"].(string); ok && strings.HasPrefix(reference, "#/$defs/") {
				name := strings.Split(strings.TrimPrefix(reference, "#/$defs/"), "/")[0]
				result = append(result, name)
			}
			for _, child := range current {
				walk(child)
			}
		case []any:
			for _, child := range current {
				walk(child)
			}
		}
	}
	walk(value)
	return result
}

func checkPropertyDescriptions(t *testing.T, owner string, raw any) {
	t.Helper()
	properties, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("%s properties are not an object", owner)
	}
	for name, rawProperty := range properties {
		property := rawProperty.(map[string]any)
		requireAnnotationText(t, owner+"/properties/"+name, property, "description")
	}
}

func requireAnnotationText(t *testing.T, pointer string, node map[string]any, keyword string) {
	t.Helper()
	value, ok := node[keyword].(string)
	minimum := 12
	if keyword == "title" {
		minimum = 3
	}
	if !ok || len(strings.TrimSpace(value)) < minimum {
		t.Errorf("%s has missing or non-specific %s", pointer, keyword)
	}
}

func enumPointers(artifact map[string]any, reachable []string) []string {
	var result []string
	var walk func(any, string)
	walk = func(value any, pointer string) {
		switch value := value.(type) {
		case map[string]any:
			if _, ok := value["enum"]; ok {
				result = append(result, pointer)
			}
			for key, child := range value {
				if key != "$defs" {
					walk(child, pointer+"/"+key)
				}
			}
		case []any:
			for index, child := range value {
				walk(child, fmt.Sprintf("%s/%d", pointer, index))
			}
		}
	}
	root := make(map[string]any, len(artifact)-1)
	for key, value := range artifact {
		if key != "$defs" {
			root[key] = value
		}
	}
	walk(root, "$")
	definitions := artifact["$defs"].(map[string]any)
	for _, name := range reachable {
		walk(definitions[name], "$/$defs/"+name)
	}
	sort.Strings(result)
	return result
}

func schemaNodeAt(t *testing.T, artifact map[string]any, pointer string) map[string]any {
	t.Helper()
	var current any = artifact
	for _, token := range strings.Split(strings.TrimPrefix(pointer, "$/"), "/") {
		switch value := current.(type) {
		case map[string]any:
			current = value[token]
		case []any:
			var index int
			if _, err := fmt.Sscanf(token, "%d", &index); err != nil {
				t.Fatalf("invalid array token %q in %s", token, pointer)
			}
			current = value[index]
		default:
			t.Fatalf("%s does not resolve", pointer)
		}
	}
	node, ok := current.(map[string]any)
	if !ok {
		t.Fatalf("%s is not a schema object", pointer)
	}
	return node
}

func walkCompiledSchema(schema *jsonschema.Schema, seen map[*jsonschema.Schema]bool, visit func(*jsonschema.Schema)) {
	if schema == nil || seen[schema] {
		return
	}
	seen[schema] = true
	visit(schema)
	children := []*jsonschema.Schema{schema.Ref, schema.RecursiveRef, schema.Not, schema.If, schema.Then, schema.Else, schema.PropertyNames, schema.UnevaluatedProperties, schema.Contains, schema.UnevaluatedItems, schema.ContentSchema}
	children = append(children, schema.AllOf...)
	children = append(children, schema.AnyOf...)
	children = append(children, schema.OneOf...)
	children = append(children, schema.PrefixItems...)
	for _, child := range schema.Properties {
		children = append(children, child)
	}
	for _, child := range schema.PatternProperties {
		children = append(children, child)
	}
	for _, child := range schema.DependentSchemas {
		children = append(children, child)
	}
	if child, ok := schema.AdditionalProperties.(*jsonschema.Schema); ok {
		children = append(children, child)
	}
	if child, ok := schema.Items.(*jsonschema.Schema); ok {
		children = append(children, child)
	}
	if childrenList, ok := schema.Items.([]*jsonschema.Schema); ok {
		children = append(children, childrenList...)
	}
	children = append(children, schema.Items2020)
	for _, child := range children {
		walkCompiledSchema(child, seen, visit)
	}
}

func checkExamplePrivacy(t *testing.T, location string, example any) {
	t.Helper()
	encoded, err := json.Marshal(example)
	if err != nil {
		t.Fatalf("marshal example at %s: %v", location, err)
	}
	if identifyingExample.Match(encoded) {
		t.Errorf("example at %s contains an absolute path or machine identifier: %s", location, encoded)
	}
	home, _ := os.UserHomeDir()
	host, _ := os.Hostname()
	if err := testutil.CheckNoForbidden("schema-annotations", location, encoded, home, host, os.Getenv("USERNAME"), os.Getenv("USER")); err != nil {
		t.Error(err)
	}
}

func assertExamplesEqual(t *testing.T, node map[string]any, want string) {
	t.Helper()
	examples, ok := node["examples"].([]any)
	if !ok || len(examples) == 0 {
		t.Fatal("identity property has no examples")
	}
	for _, example := range examples {
		if example != want {
			t.Errorf("identity example = %v, want %q", example, want)
		}
	}
}

func stripAnnotations(value any) {
	switch value := value.(type) {
	case map[string]any:
		delete(value, "title")
		delete(value, "description")
		delete(value, "examples")
		for keyword, child := range value {
			if keyword == "properties" || keyword == "$defs" {
				for _, schema := range child.(map[string]any) {
					stripAnnotations(schema)
				}
				continue
			}
			stripAnnotations(child)
		}
	case []any:
		for _, child := range value {
			stripAnnotations(child)
		}
	}
}
