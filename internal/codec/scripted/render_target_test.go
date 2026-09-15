package scripted

import (
	"context"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func targetFixture(t *testing.T, original model.Document) model.Document {
	t.Helper()
	target := renderFixture(t, "ass")
	target.Source = original.Source
	native := target.FormatData.ASS
	for index := range native.Sections {
		native.Sections[index].RawHeader = nil
	}
	for index := range native.Records {
		native.Records[index].RawLine = nil
	}
	for index := range target.Cues {
		target.Cues[index].FormatData.ASS.StartTimestampRaw, target.Cues[index].FormatData.ASS.EndTimestampRaw = nil, nil
	}
	return target
}

func TestRenderTargetPreservesOriginalBoundary(t *testing.T) {
	original := renderFixture(t, "ssa")
	target := targetFixture(t, original)
	if result, err := Render(context.Background(), target, false); err == nil || len(result.Bytes) != 0 {
		t.Fatal("ordinary renderer accepted mismatched source")
	}
	result, err := RenderTarget(context.Background(), original, target)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := Parse(context.Background(), result.Bytes, "ass")
	if err != nil || parsed.Cues[0].Payload.PlainText != target.Cues[0].Payload.PlainText {
		t.Fatalf("target reparse: %v", err)
	}
	for _, ctx := range []context.Context{nil, cancelledTargetContext()} {
		if result, err := RenderTarget(ctx, original, target); err == nil || len(result.Bytes) != 0 {
			t.Fatal("invalid context published payload")
		}
	}
}

func cancelledTargetContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestRenderTargetChecksAllSourceIntegrity(t *testing.T) {
	original := renderFixture(t, "ssa")
	extra := original.Source.Assets[0]
	extra.ID, extra.Role, extra.FileName = "asset-extra", "related", "related.ssa"
	extra.Hashes.SHA256 = strings.Repeat("0", 64)
	original.Source.Assets = append(original.Source.Assets, extra)
	target := targetFixture(t, original)
	result, err := RenderTarget(context.Background(), original, target)
	if err == nil || len(result.Bytes) != 0 || !strings.Contains(err.Error(), "invalid_source_integrity") {
		t.Fatalf("secondary integrity: %v", err)
	}
}
