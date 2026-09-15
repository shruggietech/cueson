package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

type scriptedOwners struct {
	native       *model.ScriptedDocumentData
	events       map[string]int
	styles       map[string]int
	styleNames   map[string][]int
	recordOrders map[string]int
	emphasis     map[int]scriptedEmphasis
	represented  map[int]scriptedEmphasis
}

var scriptedWebVTTTextEscape = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func indexedScriptedOwners(document model.Document) (scriptedOwners, error) {
	native := document.FormatData.ASS
	if document.Format == "ssa" {
		native = document.FormatData.SSA
	}
	if native == nil {
		return scriptedOwners{}, fmt.Errorf("missing_scripted_native")
	}
	owners := scriptedOwners{native: native, events: map[string]int{}, styles: map[string]int{}, styleNames: map[string][]int{}, recordOrders: map[string]int{}, emphasis: map[int]scriptedEmphasis{}, represented: map[int]scriptedEmphasis{}}
	for index, event := range native.Events {
		owners.events[event.EventID] = index
	}
	for index, style := range native.Styles {
		owners.styles[style.StyleID] = index
		owners.styleNames[style.Name] = append(owners.styleNames[style.Name], index)
		state, err := scriptedStyleEmphasis(style)
		if err != nil {
			return scriptedOwners{}, err
		}
		owners.emphasis[index] = state
	}
	for _, record := range native.Records {
		owners.recordOrders[record.RecordID] = record.SourceOrder
	}
	return owners, nil
}

type scriptedEmphasis struct{ bold, italic, underline bool }

func scriptedStyleEmphasis(style model.ScriptedStyle) (scriptedEmphasis, error) {
	state := scriptedEmphasis{}
	for _, field := range style.Fields {
		name := strings.ToLower(strings.TrimSpace(field.FieldName))
		if name != "bold" && name != "italic" && name != "underline" {
			continue
		}
		value, err := strconv.ParseInt(strings.TrimSpace(field.RawValue), 10, 64)
		if err != nil || (value != -1 && value != 0 && value != 1) {
			return state, fmt.Errorf("invalid_scripted_emphasis")
		}
		switch name {
		case "bold":
			state.bold = value != 0
		case "italic":
			state.italic = value != 0
		case "underline":
			state.underline = value != 0
		}
	}
	return state, nil
}

// translateScriptedPayload reads validated editable native owners, never source bytes.
func translateScriptedPayload(document model.Document, cueIndex int, targetFormat string) (payloadTranslation, error) {
	owners, err := indexedScriptedOwners(document)
	if err != nil {
		return payloadTranslation{}, err
	}
	translation, _, err := translateIndexedScriptedPayload(document, cueIndex, targetFormat, owners)
	return translation, err
}

type scriptedTextLoss struct {
	code string
	kind Kind
	path string
}

func translateIndexedScriptedPayload(document model.Document, cueIndex int, targetFormat string, owners scriptedOwners) (payloadTranslation, []scriptedTextLoss, error) {
	if targetFormat != "subrip" && targetFormat != "webvtt" {
		return payloadTranslation{}, nil, fmt.Errorf("unsupported_scripted_text_target")
	}
	if cueIndex < 0 || cueIndex >= len(document.Cues) {
		return payloadTranslation{}, nil, fmt.Errorf("invalid_scripted_cue")
	}
	cue := document.Cues[cueIndex]
	nativeCue := cue.FormatData.ASS
	if document.Format == "ssa" {
		nativeCue = cue.FormatData.SSA
	}
	if nativeCue == nil {
		return payloadTranslation{}, nil, fmt.Errorf("missing_scripted_cue_owner")
	}
	eventIndex, found := owners.events[nativeCue.EventID]
	if !found {
		return payloadTranslation{}, nil, fmt.Errorf("missing_scripted_event_owner")
	}
	styleIndex, found := owners.styles[nativeCue.StyleID]
	if !found {
		return payloadTranslation{}, nil, fmt.Errorf("missing_scripted_style_owner")
	}
	initialStyle := styleIndex
	state := owners.emphasis[styleIndex]
	var err error
	event := owners.native.Events[eventIndex]
	if strings.TrimSpace(cue.Payload.PlainText) == "" {
		return payloadTranslation{}, nil, fmt.Errorf("unreadable_scripted_dialogue")
	}
	base := fmt.Sprintf("/format_data/%s/events/%d", document.Format, eventIndex)
	losses := []scriptedTextLoss{}
	appendLoss := func(code string, kind Kind, path string) error {
		if len(losses) >= MaxLosses {
			return fmt.Errorf("scripted_payload_loss_limit")
		}
		losses = append(losses, scriptedTextLoss{code: code, kind: kind, path: path})
		return nil
	}
	var output, readable strings.Builder
	active := scriptedEmphasis{}
	closeActive := func() {
		if active.underline {
			output.WriteString("</u>")
		}
		if active.italic {
			output.WriteString("</i>")
		}
		if active.bold {
			output.WriteString("</b>")
		}
		active = scriptedEmphasis{}
	}
	writeRun := func(value string) {
		if value == "" {
			return
		}
		if strings.TrimSpace(value) != "" {
			// Each style dimension needs a matching readable run, including false
			// defaults. Resets without text and drawing spans cannot preserve it.
			defaults := owners.emphasis[styleIndex]
			matched := owners.represented[styleIndex]
			matched.bold = matched.bold || state.bold == defaults.bold
			matched.italic = matched.italic || state.italic == defaults.italic
			matched.underline = matched.underline || state.underline == defaults.underline
			owners.represented[styleIndex] = matched
		}
		if state != active {
			closeActive()
			if state.bold {
				output.WriteString("<b>")
			}
			if state.italic {
				output.WriteString("<i>")
			}
			if state.underline {
				output.WriteString("<u>")
			}
			active = state
		}
		readable.WriteString(value)
		if targetFormat == "webvtt" {
			value = scriptedWebVTTTextEscape.Replace(value)
		}
		output.WriteString(value)
	}
	wrap := nativeCue.Projection.WrapStyle
	tagCursor := 0
	runes := []rune(event.Text)
	for spanIndex, span := range event.Spans {
		if span.Kind == "drawing" {
			if err := appendLoss(LossCodeScriptedDrawingOmitted, KindOmitted, fmt.Sprintf("%s/spans/%d", base, spanIndex)); err != nil {
				return payloadTranslation{}, nil, err
			}
			continue
		}
		if span.Kind != "override" {
			var literal strings.Builder
			characters := []rune(span.Raw)
			for index := 0; index < len(characters); {
				if characters[index] == '\\' && index+1 < len(characters) {
					switch characters[index+1] {
					case 'N':
						literal.WriteByte('\n')
						index += 2
						continue
					case 'n':
						if wrap == 2 {
							literal.WriteByte('\n')
						} else {
							literal.WriteByte(' ')
						}
						index += 2
						continue
					case 'h':
						literal.WriteRune('\u00a0')
						index += 2
						continue
					}
				}
				literal.WriteRune(characters[index])
				index++
			}
			writeRun(literal.String())
			continue
		}
		residueStart := span.StartScalar + 1
		for tagCursor < len(event.Tags) && event.Tags[tagCursor].StartScalar < span.EndScalar {
			tagIndex := tagCursor
			tag := event.Tags[tagCursor]
			tagCursor++
			if tag.StartScalar < span.StartScalar {
				return payloadTranslation{}, nil, fmt.Errorf("inconsistent_scripted_tag_order")
			}
			if strings.TrimSpace(string(runes[residueStart:tag.StartScalar])) != "" {
				if err := appendLoss(LossCodeScriptedOverrideCommentOmitted, KindOmitted, fmt.Sprintf("%s/spans/%d", base, spanIndex)); err != nil {
					return payloadTranslation{}, nil, err
				}
			}
			residueStart = tag.EndScalar
			path := fmt.Sprintf("%s/tags/%d", base, tagIndex)
			switch tag.Name {
			case "b", "i", "u":
				defaults := owners.emphasis[styleIndex]
				value := false
				parameter := strings.TrimSpace(tag.Parameter)
				if parameter == "" {
					switch tag.Name {
					case "b":
						value = defaults.bold
					case "i":
						value = defaults.italic
					case "u":
						value = defaults.underline
					}
				} else {
					integer, parseErr := strconv.ParseInt(parameter, 10, 64)
					if parseErr != nil {
						return payloadTranslation{}, nil, fmt.Errorf("invalid_scripted_emphasis")
					}
					if integer != -1 && integer != 0 && integer != 1 {
						if tag.Name != "b" || integer <= 1 {
							return payloadTranslation{}, nil, fmt.Errorf("invalid_scripted_emphasis")
						}
						value = integer >= 600
						if err := appendLoss(LossCodeScriptedOverrideDegraded, KindDegraded, path); err != nil {
							return payloadTranslation{}, nil, err
						}
					} else {
						value = integer != 0
					}
				}
				switch tag.Name {
				case "b":
					state.bold = value
				case "i":
					state.italic = value
				case "u":
					state.underline = value
				}
			case "r":
				styleIndex = initialStyle
				if tag.Parameter != "" {
					matches := owners.styleNames[tag.Parameter]
					if len(matches) != 1 {
						return payloadTranslation{}, nil, fmt.Errorf("ambiguous_scripted_reset")
					}
					styleIndex = matches[0]
				}
				state = owners.emphasis[styleIndex]
				if err := appendLoss(LossCodeScriptedOverrideDegraded, KindDegraded, path); err != nil {
					return payloadTranslation{}, nil, err
				}
			case "q":
				wrap, err = strconv.Atoi(strings.TrimSpace(tag.Parameter))
				if err != nil || wrap < 0 || wrap > 3 {
					return payloadTranslation{}, nil, fmt.Errorf("invalid_scripted_wrap")
				}
				if err := appendLoss(LossCodeScriptedOverrideOmitted, KindOmitted, path); err != nil {
					return payloadTranslation{}, nil, err
				}
			default:
				if err := appendLoss(LossCodeScriptedOverrideOmitted, KindOmitted, path); err != nil {
					return payloadTranslation{}, nil, err
				}
			}
		}
		if strings.TrimSpace(string(runes[residueStart:span.EndScalar-1])) != "" {
			if err := appendLoss(LossCodeScriptedOverrideCommentOmitted, KindOmitted, fmt.Sprintf("%s/spans/%d", base, spanIndex)); err != nil {
				return payloadTranslation{}, nil, err
			}
		}
	}
	closeActive()
	if readable.String() != cue.Payload.PlainText {
		return payloadTranslation{}, nil, fmt.Errorf("scripted_readable_projection_disagreement")
	}
	collector := payloadIssueCollector{issues: []payloadIssue{}}
	text := replaceEmptyLines(output.String(), readable.String(), "<i></i>", &collector)
	if collector.overflow {
		return payloadTranslation{}, nil, fmt.Errorf("scripted_payload_loss_limit")
	}
	if targetFormat == "subrip" {
		if subRipPlainText(text) != readable.String() || subRipPayloadAmbiguous(text) {
			return payloadTranslation{}, nil, fmt.Errorf("scripted_payload_subrip_reinterpretation")
		}
	}
	return payloadTranslation{Text: text, Issues: collector.issues}, losses, nil
}
