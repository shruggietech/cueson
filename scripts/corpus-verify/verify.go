package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
)

const MaxExternalFiles = 10_000

type Report struct {
	Files          int       `json:"files"`
	Accepted       int       `json:"accepted"`
	Rejected       int       `json:"rejected"`
	Unsupported    int       `json:"unsupported"`
	RenderRejected int       `json:"render_rejected"`
	Formats        Formats   `json:"formats"`
	Failures       []Failure `json:"failures"`
}

type Formats struct {
	SRT int `json:"srt"`
	VTT int `json:"vtt"`
	ASS int `json:"ass"`
	SSA int `json:"ssa"`
}

type Failure struct {
	RelativeIdentity string `json:"relative_identity"`
	Result           string `json:"result"`
}

func Verify(ctx context.Context, root string) (Report, error) {
	return verifyWithLimits(ctx, root, MaxExternalFiles, source.MaxCaptureBytes)
}

func verifyWithLimits(ctx context.Context, root string, maxFiles int, maxBytes int64) (Report, error) {
	report := Report{Failures: []Failure{}}
	if ctx == nil {
		return report, fmt.Errorf("nil context")
	}
	if maxFiles < 1 || maxBytes < 0 {
		return report, fmt.Errorf("invalid verifier limits")
	}
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return report, fmt.Errorf("cannot inspect corpus root")
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
		return report, fmt.Errorf("corpus root must be a directory and not a link")
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("verification canceled: %w", err)
		}
		if walkErr != nil {
			return fmt.Errorf("cannot inspect corpus entry")
		}
		if path == root {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil || relative == "." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
			return fmt.Errorf("corpus entry has an unsafe relative identity")
		}
		identity := filepath.ToSlash(relative)
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("entry %q is a symbolic link", identity)
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("entry %q is not a regular file", identity)
		}
		report.Files++
		if report.Files > maxFiles {
			return fmt.Errorf("corpus contains more than %d files", maxFiles)
		}
		if info.Size() > maxBytes {
			return fmt.Errorf("entry %q exceeds the %d-byte limit", identity, maxBytes)
		}
		return verifyEntry(ctx, path, identity, &report)
	})
	if err != nil {
		return Report{}, err
	}
	return report, nil
}

func verifyEntry(ctx context.Context, path, identity string, report *Report) error {
	payload, err := source.ReadFileContext(ctx, path)
	if err != nil {
		if canceled := ctx.Err(); canceled != nil {
			return fmt.Errorf("verification canceled: %w", canceled)
		}
		return fmt.Errorf("entry %q cannot be read safely", identity)
	}
	var stdout, stderr bytes.Buffer
	status := cli.Run(ctx, []string{"encode", "--stdout", path}, strings.NewReader(""), &stdout, &stderr)
	if canceled := ctx.Err(); canceled != nil {
		return fmt.Errorf("verification canceled: %w", canceled)
	}
	if status != cli.ExitSuccess {
		extension := strings.ToLower(filepath.Ext(path))
		if extension != ".srt" && extension != ".vtt" && extension != ".ass" && extension != ".ssa" && !scripted.Detect(payload).Candidate {
			report.Unsupported++
			return nil
		}
		report.Rejected++
		report.Failures = append(report.Failures, Failure{RelativeIdentity: identity, Result: "rejected"})
		return nil
	}
	document, err := schema.Decode(stdout.Bytes())
	if err != nil {
		return fmt.Errorf("entry %q produced invalid Cue JSON", identity)
	}
	switch document.Format {
	case "subrip":
		rendered, renderErr := subrip.Render(document.Cues)
		if renderErr != nil {
			return fmt.Errorf("entry %q failed the render cycle", identity)
		}
		if _, parseErr := subrip.Parse(string(rendered), subrip.Options{}); parseErr != nil {
			return fmt.Errorf("entry %q failed the parser-renderer cycle", identity)
		}
		report.Formats.SRT++
	case "webvtt":
		rendered, renderErr := webvtt.Render(document, webvtt.RenderOptions{})
		if renderErr != nil {
			return fmt.Errorf("entry %q failed the render cycle", identity)
		}
		if _, parseErr := webvtt.Parse(string(rendered.Bytes)); parseErr != nil {
			return fmt.Errorf("entry %q failed the parser-renderer cycle", identity)
		}
		report.Formats.VTT++
	case "ass", "ssa":
		if document.Format == "ass" {
			report.Formats.ASS++
		} else {
			report.Formats.SSA++
		}
		rendered, renderErr := scripted.Render(ctx, document, false)
		if canceled := ctx.Err(); canceled != nil {
			return fmt.Errorf("verification canceled: %w", canceled)
		}
		if renderErr != nil {
			if !expectedScriptedRenderRejection(document, renderErr) {
				return fmt.Errorf("entry %q failed the scripted render cycle", identity)
			}
			report.RenderRejected++
			report.Failures = append(report.Failures, Failure{RelativeIdentity: identity, Result: "render_rejected"})
		} else if _, parseErr := scripted.Parse(ctx, rendered.Bytes, document.Format); parseErr != nil {
			return fmt.Errorf("entry %q failed the scripted parser-renderer cycle", identity)
		}
	default:
		return fmt.Errorf("entry %q produced an unsupported document format", identity)
	}
	if canceled := ctx.Err(); canceled != nil {
		return fmt.Errorf("verification canceled: %w", canceled)
	}
	report.Accepted++
	return nil
}

// Only the parser's two preservation-only classifications may account for a
// refused native render. A diagnosis must not hide cancellation, integrity,
// ownership, precision, or other unexpected renderer failures.
func expectedScriptedRenderRejection(document model.Document, renderErr error) bool {
	if renderErr == nil {
		return false
	}
	code := ""
	switch renderErr.Error() {
	case "render scripted: malformed_native_record":
		code = "malformed_native_record"
	case "render scripted: malformed_attachment":
		code = "malformed_attachment"
	default:
		return false
	}
	diagnosed := false
	for _, d := range document.Diagnostics {
		if d.Code == code {
			diagnosed = true
			break
		}
	}
	if !diagnosed {
		return false
	}
	native := document.FormatData.ASS
	if document.Format == "ssa" {
		native = document.FormatData.SSA
	}
	if native == nil {
		return false
	}
	if code == "malformed_native_record" {
		for _, r := range native.Records {
			if r.Kind == "malformed" {
				return true
			}
		}
		for _, s := range native.Styles {
			if !s.Valid {
				return true
			}
		}
		for _, e := range native.Events {
			if !e.Valid {
				return true
			}
		}
		return false
	}
	byID := make(map[string]model.ScriptedRecord, len(native.Records))
	byOrder := make(map[int]model.ScriptedRecord, len(native.Records))
	for _, r := range native.Records {
		byID[r.RecordID] = r
		byOrder[r.SourceOrder] = r
	}
	for _, a := range native.Attachments {
		start, found := byID[a.DataStartRecordID]
		if !found || a.DataRecordCount < 1 || a.DataRecordCount > len(native.Records) {
			continue
		}
		lines := make([]string, 0, a.DataRecordCount)
		for offset := 0; offset < a.DataRecordCount; offset++ {
			r, found := byOrder[start.SourceOrder+offset]
			if !found || r.RawLine == nil || r.Kind != "attachment_data" {
				break
			}
			lines = append(lines, *r.RawLine)
		}
		if len(lines) != a.DataRecordCount {
			continue
		}
		_, malformed, err := model.ScriptedAttachmentFacts(lines)
		if err == nil && malformed {
			return true
		}
	}
	return false
}
