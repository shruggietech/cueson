package conformance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/testutil"
)

type conformanceMatrix struct {
	Version int            `json:"version"`
	Formats []matrixFormat `json:"formats"`
}

type matrixFormat struct {
	Format string      `json:"format"`
	Rows   []matrixRow `json:"rows"`
}

type matrixRow struct {
	RowID       string           `json:"row_id"`
	Description string           `json:"description"`
	Evidence    []matrixEvidence `json:"evidence"`
}

type matrixEvidence struct {
	Kind      string `json:"kind"`
	Reference string `json:"reference"`
	Reason    string `json:"reason,omitempty"`
}

var matrixIDPattern = regexp.MustCompile(`^(srt|vtt)-[a-z0-9]+(?:-[a-z0-9]+)*$`)

var requiredMatrixEvidenceKinds = []string{
	"accepted_fixture",
	"malformed_fixture",
	"render_test",
	"conversion_test",
	"fuzz_target",
	"platform_test",
}

func TestConformanceMatrixResolvesDocumentationFixturesAndTests(t *testing.T) {
	t.Parallel()

	testdataRoot := fixtureRoot(t)
	repositoryRoot := filepath.Dir(testdataRoot)
	matrix := loadConformanceMatrix(t, filepath.Join(testdataRoot, "conformance-matrix.json"))
	manifest := verifiedManifest(t, testdataRoot)
	if matrix.Version != 1 || len(matrix.Formats) != 2 {
		t.Fatalf("matrix identity = version %d, formats %d", matrix.Version, len(matrix.Formats))
	}

	wantFormats := []string{"srt", "vtt"}
	rowIDs := make(map[string]struct{})
	seenKinds := make(map[string]bool)
	for formatIndex, format := range matrix.Formats {
		if format.Format != wantFormats[formatIndex] || len(format.Rows) == 0 {
			t.Fatalf("format record %d = %q with %d rows", formatIndex, format.Format, len(format.Rows))
		}
		documentName := format.Format + ".md"
		if format.Format == "vtt" {
			documentName = "webvtt.md"
		}
		document := readRepositoryFile(t, filepath.Join(repositoryRoot, "docs", "formats", documentName))
		for _, row := range format.Rows {
			if !matrixIDPattern.MatchString(row.RowID) || !strings.HasPrefix(row.RowID, format.Format+"-") {
				t.Errorf("row_id %q is not canonical for %s", row.RowID, format.Format)
			}
			if _, exists := rowIDs[row.RowID]; exists {
				t.Errorf("row_id %q is duplicated", row.RowID)
			}
			rowIDs[row.RowID] = struct{}{}
			if !strings.Contains(document, "`"+row.RowID+"`") {
				t.Errorf("row_id %q is absent from docs/formats/%s", row.RowID, documentName)
			}
			if strings.TrimSpace(row.Description) == "" || len(row.Evidence) == 0 {
				t.Errorf("row %q has no description or evidence", row.RowID)
			}
			rowKinds := make(map[string]bool)
			inapplicableKinds := make(map[string]bool)
			evidenceKeys := make(map[string]bool)
			for _, evidence := range row.Evidence {
				key := evidence.Kind + "\x00" + evidence.Reference
				if evidenceKeys[key] {
					t.Errorf("row %q duplicates %s evidence %q", row.RowID, evidence.Kind, evidence.Reference)
				}
				evidenceKeys[key] = true
				switch evidence.Kind {
				case "accepted_fixture", "malformed_fixture":
					rowKinds[evidence.Kind] = true
					seenKinds[evidence.Kind] = true
					fixture, exists := manifest.FixtureByID(evidence.Reference)
					if !exists {
						t.Errorf("row %q references missing fixture %q", row.RowID, evidence.Reference)
						continue
					}
					want := "accepted"
					if evidence.Kind == "malformed_fixture" {
						want = "rejected"
					}
					if fixture.Expectation.Result != want {
						t.Errorf("row %q fixture %q result = %q, want %q", row.RowID, evidence.Reference, fixture.Expectation.Result, want)
					}
				case "render_test", "conversion_test", "fuzz_target", "platform_test":
					rowKinds[evidence.Kind] = true
					seenKinds[evidence.Kind] = true
					if !repositoryFunctionExists(t, repositoryRoot, evidence.Reference) {
						t.Errorf("row %q references missing function %q", row.RowID, evidence.Reference)
					}
				case "inapplicable":
					if !isRequiredMatrixEvidenceKind(evidence.Reference) || strings.TrimSpace(evidence.Reason) == "" {
						t.Errorf("row %q has incomplete inapplicability evidence", row.RowID)
						continue
					}
					inapplicableKinds[evidence.Reference] = true
				default:
					t.Errorf("row %q has unknown evidence kind %q", row.RowID, evidence.Kind)
				}
			}
			for _, kind := range requiredMatrixEvidenceKinds {
				switch {
				case rowKinds[kind] && inapplicableKinds[kind]:
					t.Errorf("row %q marks %s both evidenced and inapplicable", row.RowID, kind)
				case !rowKinds[kind] && !inapplicableKinds[kind]:
					t.Errorf("row %q omits %s evidence without an inapplicability reason", row.RowID, kind)
				}
			}
		}
	}
	for _, kind := range []string{"accepted_fixture", "malformed_fixture", "render_test", "conversion_test", "fuzz_target", "platform_test"} {
		if !seenKinds[kind] {
			t.Errorf("matrix has no %s evidence", kind)
		}
	}
}

func isRequiredMatrixEvidenceKind(kind string) bool {
	for _, required := range requiredMatrixEvidenceKinds {
		if kind == required {
			return true
		}
	}
	return false
}

func loadConformanceMatrix(t *testing.T, path string) conformanceMatrix {
	t.Helper()
	data := []byte(readRepositoryFile(t, path))
	if !utf8.Valid(data) || bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		t.Fatal("conformance matrix must be UTF-8 without BOM")
	}
	var matrix conformanceMatrix
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&matrix); err != nil {
		t.Fatalf("decode conformance matrix: %v", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err == nil {
		t.Fatal("conformance matrix contains multiple JSON values")
	}
	if err := testutil.CheckNoForbidden("conformance-matrix", "metadata", data); err != nil {
		t.Fatal(err)
	}
	return matrix
}

func repositoryFunctionExists(t *testing.T, repositoryRoot, reference string) bool {
	t.Helper()
	directory, name, found := strings.Cut(reference, ":")
	if !found || directory == "" || !regexp.MustCompile(`^(Test|Fuzz)[A-Za-z0-9_]+$`).MatchString(name) {
		return false
	}
	matches, err := filepath.Glob(filepath.Join(repositoryRoot, filepath.FromSlash(directory), "*_test.go"))
	if err != nil {
		t.Fatalf("resolve test reference %q: %v", reference, err)
	}
	needle := []byte("func " + name + "(")
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read test reference %q: %v", reference, err)
		}
		if bytes.Contains(data, needle) {
			return true
		}
	}
	return false
}

func readRepositoryFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(fmt.Errorf("read %s: %w", filepath.Base(path), err))
	}
	return string(data)
}
