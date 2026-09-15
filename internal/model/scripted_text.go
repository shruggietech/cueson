package model

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// ScriptedTextFacts is a deterministic derived view, never replacement source.
// The bounded scanner is shared model validation, not a document codec.
type ScriptedTextFacts struct {
	Lines           []string
	Spans           []ScriptedSpan
	Tags            []ScriptedTag
	Karaoke         []ScriptedKaraoke
	Tokens          []Token
	DrawingExcluded bool
	DiagnosticCodes []string
}
type ScriptedTextNames struct {
	ResetStyles  map[string]bool
	FontFamilies map[string]bool
}

// ProjectScriptedText computes the ratified explicit-break/drawing/karaoke view.
// It neither loads fonts nor simulates layout, wrapping, or pixels.
func ProjectScriptedText(text string, wrap int, start, end int64, names ...ScriptedTextNames) (ScriptedTextFacts, error) {
	f := ScriptedTextFacts{Lines: []string{}, Spans: []ScriptedSpan{}, Tags: []ScriptedTag{}, Karaoke: []ScriptedKaraoke{}, Tokens: []Token{}, DiagnosticCodes: []string{}}
	if !scriptedPhysical(text) || wrap < 0 || wrap > 3 || start < 0 || end <= start {
		return f, fmt.Errorf("inconsistent_projection")
	}
	r := []rune(text)
	var readable strings.Builder
	drawing := false
	karaokeBad := false
	elapsed := int64(0)
	tokenStart, tokenEnd := int64(0), int64(0)
	var syllable strings.Builder
	activeKaraoke := false
	diag := func(code string) {
		if !slices.Contains(f.DiagnosticCodes, code) {
			f.DiagnosticCodes = append(f.DiagnosticCodes, code)
		}
	}
	flushToken := func() {
		if activeKaraoke {
			native := strconv.FormatInt((tokenEnd-tokenStart)/10, 10)
			f.Tokens = append(f.Tokens, Token{Text: syllable.String(), StartMilliseconds: tokenStart, EndMilliseconds: tokenEnd, NativeTiming: &native})
			syllable.Reset()
		}
	}
	literal := func(a, b int, ambiguous bool) {
		kind := "literal"
		if drawing && !ambiguous {
			kind = "drawing"
			f.DrawingExcluded = true
		}
		f.Spans = append(f.Spans, ScriptedSpan{Kind: kind, StartScalar: a, EndScalar: b, Raw: string(r[a:b])})
		if kind == "drawing" {
			return
		}
		for i := a; i < b; i++ {
			if !ambiguous && r[i] == '\\' && i+1 < b {
				switch r[i+1] {
				case 'N':
					readable.WriteByte('\n')
					if activeKaraoke {
						syllable.WriteByte('\n')
					}
					i++
					continue
				case 'n':
					c := ' '
					if wrap == 2 {
						c = '\n'
					}
					readable.WriteRune(c)
					if activeKaraoke {
						syllable.WriteRune(c)
					}
					i++
					continue
				case 'h':
					readable.WriteRune('\u00a0')
					if activeKaraoke {
						syllable.WriteRune('\u00a0')
					}
					i++
					continue
				}
			}
			readable.WriteRune(r[i])
			if activeKaraoke {
				syllable.WriteRune(r[i])
			}
		}
	}
	for i := 0; i < len(r); {
		if r[i] != '{' {
			j := i + 1
			for j < len(r) && r[j] != '{' {
				j++
			}
			literal(i, j, false)
			i = j
			continue
		}
		j := i + 1
		nested := false
		for j < len(r) && r[j] != '}' {
			if r[j] == '{' {
				nested = true
			}
			j++
		}
		if j == len(r) || nested {
			diag("malformed_override")
			if j < len(r) {
				j++
			}
			literal(i, j, true)
			i = j
			continue
		}
		f.Spans = append(f.Spans, ScriptedSpan{Kind: "override", StartScalar: i, EndScalar: j + 1, Raw: string(r[i : j+1])})
		for p := i + 1; p < j; {
			if r[p] != '\\' {
				p++
				continue
			}
			a := p
			p++
			nameStart := p
			for p < j && unicode.IsDigit(r[p]) {
				p++
			}
			for p < j && unicode.IsLetter(r[p]) {
				p++
			}
			name := string(r[nameStart:p])
			paramStart := p
			depth := 0
			for p < j {
				if r[p] == '(' {
					depth++
					if depth > MaxScriptedDepth {
						return f, fmt.Errorf("complexity_limit")
					}
				}
				if r[p] == ')' {
					depth--
					if depth < 0 {
						diag("malformed_override")
						depth = 0
					}
				}
				if r[p] == '\\' && depth == 0 {
					break
				}
				p++
			}
			if depth != 0 {
				diag("malformed_override")
			}
			// Prefix-like unknown alphabetic identities are ambiguous. Split only
			// when the complete suffix identifies declared native content exactly.
			if len(names) > 0 {
				full := string(r[nameStart:p])
				for _, prefix := range []string{"fn", "r"} {
					if !strings.HasPrefix(name, prefix) || len(name) <= len(prefix) {
						continue
					}
					declared := names[0].ResetStyles
					if prefix == "fn" {
						declared = names[0].FontFamilies
					}
					suffix := full[len(prefix):]
					if declared[suffix] {
						name = prefix
						paramStart = nameStart + len(prefix)
						break
					}
				}
			}
			parameter := string(r[paramStart:p])
			tag := ScriptedTag{Name: name, Parameter: parameter, StartScalar: a, EndScalar: p, Raw: string(r[a:p])}
			f.Tags = append(f.Tags, tag)
			switch name {
			case "q":
				v, e := strconv.Atoi(strings.TrimSpace(parameter))
				if e != nil || v < 0 || v > 3 {
					return f, fmt.Errorf("inconsistent_projection")
				}
				wrap = v
			case "p":
				v, e := strconv.ParseInt(strings.TrimSpace(parameter), 10, 64)
				if e != nil || v < 0 {
					return f, fmt.Errorf("inconsistent_projection")
				}
				drawing = v > 0
			case "k", "K", "kf", "ko", "kt":
				flushToken()
				activeKaraoke = false
				v, e := strconv.ParseInt(strings.TrimSpace(parameter), 10, 64)
				supported := e == nil && v >= 0 && name != "kt"
				var duration *int64
				if e == nil {
					duration = &v
				}
				if !supported || v > (end-start-elapsed)/10 {
					karaokeBad = true
					diag("unsupported_karaoke_timing")
					supported = false
				}
				f.Karaoke = append(f.Karaoke, ScriptedKaraoke{Variant: name, StartScalar: a, EndScalar: p, DurationCentiseconds: duration, Supported: supported})
				if supported {
					tokenStart = start + elapsed
					elapsed += v * 10
					tokenEnd = start + elapsed
					activeKaraoke = true
				}
			default:
				known := []string{"b", "i", "u", "s", "r", "fn", "fs", "fscx", "fscy", "fsp", "fr", "frx", "fry", "frz", "fax", "fay", "fe", "c", "1c", "2c", "3c", "4c", "alpha", "1a", "2a", "3a", "4a", "a", "an", "pos", "move", "org", "clip", "iclip", "t", "fad", "fade", "bord", "xbord", "ybord", "shad", "xshad", "yshad", "blur", "be", "pbo"}
				if !slices.Contains(known, name) {
					diag("unsupported_override")
				}
			}
			if len(f.Tags) > MaxItemOccurrences || len(f.Karaoke) > MaxItemOccurrences {
				return f, fmt.Errorf("complexity_limit")
			}
		}
		i = j + 1
		if len(f.Spans) > MaxItemOccurrences {
			return f, fmt.Errorf("complexity_limit")
		}
	}
	flushToken()
	if karaokeBad {
		f.Tokens = []Token{}
	}
	if len(f.Spans) > MaxItemOccurrences {
		return f, fmt.Errorf("complexity_limit")
	}
	f.Lines = strings.Split(readable.String(), "\n")
	if len(f.Lines) > MaxItemOccurrences {
		return f, fmt.Errorf("complexity_limit")
	}
	return f, nil
}

func validateScriptedProjection(c *Cue, e ScriptedEvent, n *ScriptedCueData, wrap int, diagnosticIndex map[string]bool, textNames ScriptedTextNames) error {
	p := n.Projection
	if p.Origin != "native_derived" || p.SourceEventID != e.EventID || nativeName(p.TextFieldName) != "text" || p.WrapStyle < 0 || p.WrapStyle > 3 {
		return scriptedError("inconsistent_projection", c.SourceOrder)
	}
	if p.WrapStyle != wrap {
		return scriptedError("inconsistent_projection", c.SourceOrder)
	}
	f, err := ProjectScriptedText(e.Text, wrap, c.Timing.StartMilliseconds, c.Timing.EndMilliseconds, textNames)
	if err != nil {
		return scriptedError(err.Error(), c.SourceOrder)
	}
	if c.Payload.RawText != e.Text || !slices.Equal(c.Payload.Lines, f.Lines) || c.Payload.PlainText != strings.Join(f.Lines, "\n") || p.DrawingExcluded != f.DrawingExcluded || !slices.Equal(e.Spans, f.Spans) || !slices.Equal(e.Tags, f.Tags) || !equalScriptedKaraoke(e.Karaoke, f.Karaoke) || !equalScriptedTokens(c.Tokens, f.Tokens) {
		return scriptedError("inconsistent_projection", c.SourceOrder)
	}
	for _, code := range f.DiagnosticCodes {
		if !diagnosticIndex[code+":"+strconv.Itoa(c.SourceOrder)] {
			return scriptedError("missing_native_diagnostic", c.SourceOrder)
		}
	}
	actor := ""
	speakerField := ""
	for _, field := range e.Fields {
		if nativeFieldName(field.FieldName, "event") == "name" {
			actor = field.RawValue
			speakerField = field.FieldName
		}
	}
	if actor == "" {
		if len(c.Speakers) != 0 || p.SpeakerFieldName != nil {
			return scriptedError("inconsistent_speaker_provenance", c.SourceOrder)
		}
	} else {
		if len(c.Speakers) != 1 || c.Speakers[0].Name != actor || c.Speakers[0].Origin != "native" || p.SpeakerFieldName == nil || *p.SpeakerFieldName != speakerField {
			return scriptedError("inconsistent_speaker_provenance", c.SourceOrder)
		}
	}
	return nil
}
func equalScriptedKaraoke(a, b []ScriptedKaraoke) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Variant != b[i].Variant || a[i].StartScalar != b[i].StartScalar || a[i].EndScalar != b[i].EndScalar || a[i].Supported != b[i].Supported || !equalOptionalInt64(a[i].DurationCentiseconds, b[i].DurationCentiseconds) {
			return false
		}
	}
	return true
}
func equalOptionalInt64(a, b *int64) bool {
	return a == nil && b == nil || a != nil && b != nil && *a == *b
}
func equalScriptedTokens(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Text != b[i].Text || a[i].StartMilliseconds != b[i].StartMilliseconds || a[i].EndMilliseconds != b[i].EndMilliseconds || !equalOptionalString(a[i].NativeTiming, b[i].NativeTiming) {
			return false
		}
	}
	return true
}
