package testutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ConformanceMatrix is the domain-neutral shape of the governed evidence inventory.
type ConformanceMatrix struct {
	Version int            `json:"version"`
	Formats []MatrixFormat `json:"formats"`
}

type MatrixFormat struct {
	Format string      `json:"format"`
	Rows   []MatrixRow `json:"rows"`
}

type MatrixRow struct {
	RowID       string           `json:"row_id"`
	Description string           `json:"description"`
	Evidence    []MatrixEvidence `json:"evidence"`
}

type MatrixEvidence struct {
	Kind      string `json:"kind"`
	Reference string `json:"reference"`
	Reason    string `json:"reason,omitempty"`
}

// DecodeConformanceMatrix accepts exactly one UTF-8 JSON object, without BOM or unknown fields.
func DecodeConformanceMatrix(data []byte) (ConformanceMatrix, error) {
	if !utf8.Valid(data) || bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		return ConformanceMatrix{}, fmt.Errorf("conformance matrix must be UTF-8 without BOM")
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return ConformanceMatrix{}, fmt.Errorf("conformance matrix must be a JSON object")
	}
	var matrix ConformanceMatrix
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&matrix); err != nil {
		return ConformanceMatrix{}, fmt.Errorf("decode conformance matrix: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return ConformanceMatrix{}, fmt.Errorf("conformance matrix contains trailing data")
	}
	return matrix, nil
}

// EvidenceFunctionExists resolves a real, runnable Go test/fuzz declaration without escaping the repository.
func EvidenceFunctionExists(repositoryRoot, reference, kind string) (bool, error) {
	directory, name, found := strings.Cut(reference, ":")
	if !found || directory == "." || !fs.ValidPath(directory) || strings.ContainsAny(directory, `\:`) {
		return false, nil
	}
	wantType := "T"
	switch kind {
	case "fuzz_target":
		if !validEvidenceName(name, "Fuzz") {
			return false, nil
		}
		wantType = "F"
	case "render_test", "conversion_test", "platform_test":
		if !validEvidenceName(name, "Test") {
			return false, nil
		}
	default:
		return false, nil
	}
	resolved := repositoryRoot
	for _, component := range strings.Split(directory, "/") {
		resolved = filepath.Join(resolved, component)
		info, err := os.Lstat(resolved)
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("inspect evidence directory: %w", err)
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return false, nil
		}
	}
	entries, err := os.ReadDir(resolved)
	if err != nil {
		return false, fmt.Errorf("read evidence directory: %w", err)
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return false, fmt.Errorf("inspect evidence source: %w", err)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(resolved, entry.Name()), nil, 0)
		if err != nil {
			return false, fmt.Errorf("parse evidence source: %w", err)
		}
		testingName := ""
		for _, item := range file.Imports {
			path, err := strconv.Unquote(item.Path.Value)
			if err == nil && path == "testing" {
				testingName = "testing"
				if item.Name != nil {
					testingName = item.Name.Name
				}
			}
		}
		for _, declaration := range file.Decls {
			fn, ok := declaration.(*ast.FuncDecl)
			if !ok || fn.Name.Name != name || fn.Recv != nil || fn.Body == nil || fn.Type.TypeParams != nil || fn.Type.Params == nil || len(fn.Type.Params.List) != 1 || fn.Type.Results != nil && len(fn.Type.Results.List) != 0 {
				continue
			}
			parameter := fn.Type.Params.List[0]
			if len(parameter.Names) > 1 {
				continue
			}
			pointer, ok := parameter.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			selector, ok := pointer.X.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != wantType {
				continue
			}
			owner, ok := selector.X.(*ast.Ident)
			if ok && testingName != "" && testingName != "_" && owner.Name == testingName {
				return true, nil
			}
		}
	}
	return false, nil
}

func validEvidenceName(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) || len(name) == len(prefix) {
		return false
	}
	for _, character := range name {
		if character != '_' && !(character >= 'A' && character <= 'Z') && !(character >= 'a' && character <= 'z') && !(character >= '0' && character <= '9') {
			return false
		}
	}
	first := name[len(prefix)]
	return !(first >= 'a' && first <= 'z')
}
