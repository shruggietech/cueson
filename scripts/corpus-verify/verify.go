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
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
)

const MaxExternalFiles = 10_000

type Report struct {
	Files       int       `json:"files"`
	Accepted    int       `json:"accepted"`
	Rejected    int       `json:"rejected"`
	Unsupported int       `json:"unsupported"`
	Formats     Formats   `json:"formats"`
	Failures    []Failure `json:"failures"`
}

type Formats struct {
	SRT int `json:"srt"`
	VTT int `json:"vtt"`
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
	if _, err := source.ReadFileContext(ctx, path); err != nil {
		return fmt.Errorf("entry %q cannot be read safely", identity)
	}
	var stdout, stderr bytes.Buffer
	status := cli.Run(ctx, []string{"encode", "--stdout", path}, strings.NewReader(""), &stdout, &stderr)
	if status != cli.ExitSuccess {
		extension := strings.ToLower(filepath.Ext(path))
		if extension != ".srt" && extension != ".vtt" {
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
	}
	report.Accepted++
	return nil
}
