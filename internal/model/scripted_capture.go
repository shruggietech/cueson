package model

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

type scriptedSourceCapture struct {
	lines       []string
	eventFields map[int]map[string]string
}

func (capture scriptedSourceCapture) matches(order int, line string) bool {
	return order >= 0 && order < len(capture.lines) && capture.lines[order] == line
}

// inspectScriptedCaptureSource verifies the source independently of editable
// native owners, then records bounded physical capture evidence. Original bytes
// retain BOM/terminators; decoded observations intentionally exclude both.
func inspectScriptedCaptureSource(source SourceEnvelope) (scriptedSourceCapture, error) {
	capture := scriptedSourceCapture{lines: []string{}, eventFields: map[int]map[string]string{}}
	for _, asset := range source.Assets {
		if asset.ID != source.PrimaryAssetID {
			continue
		}
		// This boundary also validates private conversion targets, which skip
		// the source grammar's earlier size preflight.
		if len(asset.DataBase64) > ((64<<20)+2)/3*4 {
			return capture, fmt.Errorf("complexity_limit: original capture source")
		}
		if strings.ContainsAny(asset.DataBase64, "\r\n") {
			return capture, fmt.Errorf("invalid scripted source integrity")
		}
		raw, err := base64.StdEncoding.Strict().DecodeString(asset.DataBase64)
		if err != nil || len(raw) > 64<<20 || int64(len(raw)) != asset.Size.Bytes || fmt.Sprintf("%x", sha256.Sum256(raw)) != asset.Hashes.SHA256 {
			return capture, fmt.Errorf("invalid scripted source integrity")
		}
		text := strings.TrimPrefix(string(raw), "\ufeff")
		text = strings.ReplaceAll(text, "\r\n", "\n")
		section := ""
		var declaration []string
		for text != "" {
			if len(capture.lines) >= MaxDocumentItems {
				return capture, fmt.Errorf("complexity_limit: original capture lines")
			}
			line, remaining, _ := strings.Cut(text, "\n")
			text = remaining
			order := len(capture.lines)
			capture.lines = append(capture.lines, line)
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
				section = trim[1 : len(trim)-1]
				declaration = nil
				continue
			}
			if !strings.EqualFold(section, "Events") {
				continue
			}
			prefix, value, ok := strings.Cut(trim, ":")
			if !ok {
				continue
			}
			if strings.EqualFold(prefix, "Format") {
				declaration = strings.SplitN(value, ",", MaxScriptedDeclarationFields+1)
				if len(declaration) > MaxScriptedDeclarationFields {
					return capture, fmt.Errorf("complexity_limit: original declaration")
				}
				continue
			}
			if (!strings.EqualFold(prefix, "Dialogue") && !strings.EqualFold(prefix, "Comment")) || len(declaration) == 0 {
				continue
			}
			// Prefix whitespace is native framing, not part of the first event field.
			value = strings.TrimLeft(value, " ")
			values := strings.SplitN(value, ",", len(declaration))
			if len(values) != len(declaration) {
				continue
			}
			fields := map[string]string{}
			for i, name := range declaration {
				fields[nativeFieldName(name, "event")] = values[i]
			}
			capture.eventFields[order] = fields
		}
		return capture, nil
	}
	return capture, fmt.Errorf("missing scripted primary capture asset")
}
