package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func externalScript(format, text string) string {
	dialect, section := "v4.00+", "V4+ Styles"
	style := "Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,-1,0,0,0,100,100,0,0,1,2,1,2,10,10,10,1"
	first := "0"
	if format == "ssa" {
		dialect = "v4.00"
		section = "V4 Styles"
		style = "Default,Arial,20,16777215,255,0,-2147483648,-1,0,1,2,1,2,10,10,10,0,1"
		first = "Marked=0"
	}
	return "[Script Info]\nScriptType: " + dialect + "\n[" + section + "]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "style"), ",") + "\nStyle: " + style + "\n[Events]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "event"), ",") + "\nDialogue: " + first + ",0:00:01.00,0:00:03.00,Default,Actor,0,0,0,," + text + "\n"
}

func TestVerifyScriptedContentCyclesAndPreservationOnly(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	inputs := map[string]string{
		"nested/valid.bin":        "\ufeff" + strings.ReplaceAll(externalScript("ass", `{\k10}Hello\Nworld`), "\n", "\r\n"),
		"valid.ssa":               externalScript("ssa", "Hello, SSA"),
		"unknown.ass":             externalScript("ass", `{\future9}Hello`),
		"preserve-style.ass":      strings.Replace(externalScript("ass", "Hello"), "[Events]", "Style: inert\n[Events]", 1),
		"preserve-attachment.ssa": externalScript("ssa", "Hello") + "[Fonts]\nfontname: example.ttf\n!\"\n",
		"preserve-header.ass":     externalScript("ass", "Hello") + "[Fonts]\nfontname: empty.ttf\n",
	}
	for identity, input := range inputs {
		writeTestFile(t, filepath.Join(root, filepath.FromSlash(identity)), input)
	}
	before := externalByteSnapshot(t, root)
	report, err := Verify(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Files != 6 || report.Accepted != 6 || report.Rejected != 0 || report.Unsupported != 0 || report.RenderRejected != 3 || report.Formats.ASS != 4 || report.Formats.SSA != 2 {
		t.Fatalf("wrong native classifications: %+v", report)
	}
	if len(report.Failures) != 3 {
		t.Fatalf("preservation failure inventory: %+v", report.Failures)
	}
	for _, f := range report.Failures {
		if f.Result != "render_rejected" {
			t.Fatalf("wrong preservation disposition: %+v", f)
		}
	}
	if after := externalByteSnapshot(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("verifier changed corpus bytes or entries")
	}
	assertExternalReportPrivate(t, report, root)
}

func TestVerifyRejectedScriptedCandidatesIgnoreExtension(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	unsafe := strings.Replace(externalScript("ass", "Hello"), "ScriptType: v4.00+", "ScriptType: v4.00+\nTitle: secret/thing", 1)
	writeTestFile(t, filepath.Join(root, "unsafe.bin"), unsafe)
	writeTestFile(t, filepath.Join(root, "mixed.bin"), externalScript("ass", "Hello")+"[V4 Styles]\n")
	writeTestFile(t, filepath.Join(root, "unsupported.bin"), "[Script Info]\nScriptType: v4.00++\n[V4+ Styles]\n")
	writeTestFile(t, filepath.Join(root, "plain.bin"), "not subtitles")
	before := externalByteSnapshot(t, root)
	report, err := Verify(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Files != 4 || report.Rejected != 3 || report.Unsupported != 1 || report.Accepted != 0 || len(report.Failures) != 3 {
		t.Fatalf("extension hid rejected native content: %+v", report)
	}
	for _, f := range report.Failures {
		if f.Result != "rejected" {
			t.Fatalf("unexpected failure: %+v", f)
		}
	}
	if after := externalByteSnapshot(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("rejection changed corpus bytes or entries")
	}
	assertExternalReportPrivate(t, report, root, "secret/thing", "v4.00++")
}

func TestExpectedScriptedRenderRejectionDoesNotHideOtherFailures(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "preserve.ass")
	writeTestFile(t, path, externalScript("ass", "Hello")+"[Fonts]\nfontname: empty.ttf\n")
	var stdout, stderr bytes.Buffer
	if status := cli.Run(context.Background(), []string{"encode", "--stdout", path}, strings.NewReader(""), &stdout, &stderr); status != cli.ExitSuccess {
		t.Fatalf("encode %d: %s", status, stderr.String())
	}
	document, err := schema.Decode(stdout.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	_, renderErr := scripted.Render(context.Background(), document, false)
	if !expectedScriptedRenderRejection(document, renderErr) {
		t.Fatalf("real preservation-only refusal not classified: %v", renderErr)
	}
	for _, other := range []error{context.Canceled, errors.New("render scripted: invalid_source_integrity"), errors.New("render scripted: unrepresentable_centiseconds"), errors.New("render scripted: malformed_attachment"), errors.New("unrelated renderer failure"), nil} {
		if expectedScriptedRenderRejection(document, other) {
			t.Fatalf("diagnosis hid unexpected renderer failure: %v", other)
		}
	}
	clean := document
	clean.Diagnostics = nil
	if expectedScriptedRenderRejection(clean, renderErr) {
		t.Fatal("missing diagnosis accepted as expected refusal")
	}
	clean = document
	native := *document.FormatData.ASS
	native.Records = append([]model.ScriptedRecord(nil), native.Records...)
	for i := range native.Records {
		if native.Records[i].Kind == "malformed" {
			native.Records[i].Kind = "unknown"
		}
	}
	clean.FormatData.ASS = &native
	if expectedScriptedRenderRejection(clean, renderErr) {
		t.Fatal("diagnosis without malformed owner hid renderer failure")
	}
}

func TestVerifyScriptedEntryCancellationIsNotRejection(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "valid.ass")
	writeTestFile(t, path, externalScript("ass", "Hello"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	report := Report{Failures: []Failure{}}
	err := verifyEntry(ctx, path, "valid.ass", &report)
	if !errors.Is(err, context.Canceled) || report.Rejected != 0 || report.Unsupported != 0 || report.Accepted != 0 {
		t.Fatalf("cancellation became corpus classification: %+v %v", report, err)
	}
	if strings.Contains(err.Error(), root) {
		t.Fatal("cancellation leaked corpus root")
	}
}

func externalByteSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		payload, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		identity, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(identity)] = string(payload)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func assertExternalReportPrivate(t *testing.T, report Report, root string, private ...string) {
	t.Helper()
	payload, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range append(private, root, filepath.ToSlash(root)) {
		encoded, _ := json.Marshal(value)
		if strings.Contains(string(payload), value) || strings.Contains(string(payload), string(encoded[1:len(encoded)-1])) {
			t.Fatalf("report leaked private value: %s", payload)
		}
	}
}
