package scripted

import (
	"context"
	"github.com/shruggietech/cueson/internal/model"
	"strings"
)

func (p *parser) project(ctx context.Context) error {
	names := model.ScriptedTextNames{ResetStyles: map[string]bool{}, FontFamilies: map[string]bool{}}
	styles := map[string][]string{}
	styleCounts := map[string]int{}
	for _, s := range p.result.Native.Styles {
		styleCounts[s.Name]++
		if !s.Valid {
			continue
		}
		styles[s.Name] = append(styles[s.Name], s.StyleID)
		names.ResetStyles[s.Name] = true
		for _, f := range s.Fields {
			identity, _ := model.ScriptedFieldIdentity(f.FieldName, p.format, "style")
			if identity == "fontname" {
				names.FontFamilies[f.RawValue] = true
			}
		}
	}
	orders := map[string]int{}
	for _, r := range p.result.Native.Records {
		orders[r.RecordID] = r.SourceOrder
	}
	previous := int64(-1)
	for i := range p.result.Native.Events {
		if err := ctx.Err(); err != nil {
			return err
		}
		e := &p.result.Native.Events[i]
		if !e.Valid {
			continue
		}
		order := orders[e.RecordID]
		start, _ := model.ScriptedMilliseconds(field(e.Fields, "start"))
		end, _ := model.ScriptedMilliseconds(field(e.Fields, "end"))
		facts, err := model.ProjectScriptedText(e.Text, p.wrap, start, end, names)
		if err != nil {
			return failure(err.Error(), order)
		}
		e.Spans = facts.Spans
		e.Tags = facts.Tags
		e.Karaoke = facts.Karaoke
		p.spans += len(facts.Spans) + len(facts.Tags)
		if p.spans > model.MaxScriptedSpans {
			return failure("complexity_limit", order)
		}
		for _, code := range facts.DiagnosticCodes {
			if err = p.warn(code, order); err != nil {
				return err
			}
		}
		if e.EventType != "dialogue" {
			continue
		}
		styleIDs := styles[field(e.Fields, "style")]
		if len(styleIDs) != 1 || styleCounts[field(e.Fields, "style")] != 1 {
			return failure("unresolved_style", order)
		}
		ordinal := len(p.result.Cues)
		id := occurrence("cue", ordinal)
		e.CueID = ptr(id)
		native := &model.ScriptedCueData{EventID: e.EventID, StyleID: styleIDs[0], StartTimestampRaw: ptr(field(e.Fields, "start")), EndTimestampRaw: ptr(field(e.Fields, "end")), Projection: model.ScriptedProjection{Origin: "native_derived", SourceEventID: e.EventID, TextFieldName: "Text", WrapStyle: p.wrap, DrawingExcluded: facts.DrawingExcluded}}
		cue := model.Cue{ID: id, Ordinal: ordinal, SourceOrder: order, Timing: model.Timing{StartMilliseconds: start, EndMilliseconds: end, DurationMilliseconds: end - start}, Payload: model.Payload{RawText: e.Text, PlainText: strings.Join(facts.Lines, "\n"), Lines: facts.Lines}, Speakers: []model.Speaker{}, Tokens: facts.Tokens, OCRObservations: []model.OCRObservation{}}
		for _, f := range e.Fields {
			identity, _ := model.ScriptedFieldIdentity(f.FieldName, p.format, "event")
			if identity == "text" {
				native.Projection.TextFieldName = f.FieldName
			}
			if identity == "name" && f.RawValue != "" {
				native.Projection.SpeakerFieldName = ptr(f.FieldName)
				cue.Speakers = append(cue.Speakers, model.Speaker{Name: f.RawValue, Origin: "native"})
			}
		}
		if p.format == "ass" {
			cue.FormatData.ASS = native
		} else {
			cue.FormatData.SSA = native
		}
		if previous > start {
			if err = p.warn("decreasing_start", order); err != nil {
				return err
			}
		}
		previous = start
		p.result.Cues = append(p.result.Cues, cue)
	}
	return nil
}
