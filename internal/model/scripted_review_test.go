package model

import (
	"encoding/base64"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestScriptedInertScalarGrammarPreservesUnicodeAndPunctuation(t *testing.T) {
	grammar := regexp.MustCompile(`^[\p{L}\p{N} ._+&=#%()\[\]{}!;-]*$`)
	for _, c := range []struct {
		value string
		valid bool
	}{
		{"", true}, {"opaque native content", true}, {"字幕 École Ελληνικά ١٢ Ⅷ ²", true}, {" ._+&=#%()[]{}!;-", true},
		{"opaque/path", false}, {`opaque\path`, false}, {"opaque:value", false}, {"opaque~value", false}, {"opaque,value", false}, {"opaque>value", false}, {"opaque\u00a0value", false}, {"opaque\tvalue", false}, {"opaque\nvalue", false}, {"opaque😀value", false}, {"opaque\xffvalue", false},
	} {
		if got := isInertNativeScalar(c.value); got != c.valid || got != grammar.MatchString(c.value) {
			t.Fatalf("inert native grammar for %q: got %v, want %v", c.value, got, c.valid)
		}
		field := ScriptedField{FieldName: "X-Inert", RawValue: c.value}
		if err := validateScriptedScalar(field, "ass", 0, "style"); (err == nil) != c.valid {
			t.Fatalf("native extension acceptance for %q: %v", c.value, err)
		}
	}
	for r := rune(0); r < 128; r++ {
		value := string(r)
		if isInertNativeScalar(value) != grammar.MatchString(value) {
			t.Fatalf("ASCII native grammar changed at %U", r)
		}
	}
}

func TestScriptedMetadataPrefilterPreservesIdentityMatcher(t *testing.T) {
	for _, value := range []string{
		"ordinary native content", "字幕 Ελληνικά ١٢ Ⅷ", "opaque scalar", "LOWERCASE", "", "\xff",
		`C:\private\source.ass`, "D:/private/source.ass", `\\machine\share`, "FILE:source", "https://private.invalid", "git+ssh://private.invalid",
		"/workspace/source.ass", "note>/opt/source.ass", "note\u00a0~/source.ass", "~someone/source.ass", `~someone\source.ass`, "LOCALHOST", "localhoſt",
		"hostname:private", "machineid=private", "username\u00a0:private", "userid=private", "machineNAME=private",
	} {
		if containsScriptedMetadataIdentity(value) != (strings.ContainsAny(value, `/\`) || metadataIdentity.MatchString(value)) {
			t.Fatalf("metadata prefilter changed identity match for %q", value)
		}
	}
}

func TestScriptedMalformedSourceRecordsClassifyBeforeContentExemptions(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, prefix := range []string{"Dialogue", "Comment", "Style"} {
			for _, value := range []string{"/home/alice/private.ass", `relative\private.ass`, "relative/private.ass", "code once", "template syl", "unclassifiable😀"} {
				d := scriptedExample(t, format)
				b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
				section := ""
				if prefix == "Style" {
					section = "[V4+ Styles]\n"
					if format == "ssa" {
						section = "[V4 Styles]\n"
					}
				}
				b = append(b, []byte(section+prefix+": "+value+"\n")...)
				setScriptedSourceBytes(&d, b)
				if err := d.Validate(); err == nil {
					t.Fatalf("%s malformed %s with %q accepted", format, prefix, value)
				}
			}
		}
		for _, prefix := range []string{"Comment", "Style"} {
			d := scriptedExample(t, format)
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			if prefix == "Style" {
				section := "[V4+ Styles]\n"
				if format == "ssa" {
					section = "[V4 Styles]\n"
				}
				b = append(b, []byte(section)...)
			}
			b = append(b, []byte(prefix+": bounded,未知 ١٢\n")...)
			setScriptedSourceBytes(&d, b)
			if err := scriptedSourcePrivacy(d.Source, format); err != nil {
				t.Fatalf("bounded inert malformed %s cannot be retained: %v", prefix, err)
			}
		}
		for _, declaration := range []string{"", "Format: Text\n", "Format: Text, Start, End\n", "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text, Text\n"} {
			d := scriptedExample(t, format)
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = append(b, []byte("[Events]\n"+declaration+"Dialogue: harmless\n")...)
			setScriptedSourceBytes(&d, b)
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "malformed_native_record") {
				t.Fatalf("malformed Dialogue declaration %q accepted: %v", declaration, err)
			}
		}
		d := scriptedExample(t, format)
		b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		lines := strings.Split(string(b), "\n")
		lines[4] += ",relative/private.ass"
		setScriptedSourceBytes(&d, []byte(strings.Join(lines, "\n")))
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
			t.Fatalf("%s excess Style value obtained comma-suffix exemption: %v", format, err)
		}
		for _, role := range []string{"text", "fontname"} {
			d = scriptedExample(t, format)
			b, _ = base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			if role == "text" {
				b = []byte(strings.Replace(string(b), "Hello", "content discusses relative/private.ass, and /root/private.ass", 1))
			} else {
				b = []byte(strings.Replace(string(b), "Arial", "font content relative/private.ass", 1))
			}
			setScriptedSourceBytes(&d, b)
			if err := scriptedSourcePrivacy(d.Source, format); err != nil {
				t.Fatalf("%s interpretable original %s lost content role: %v", format, role, err)
			}
		}
	}
}

func TestScriptedSourceFramingRequiresCompleteCanonicalDeclaration(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		names := assEventNames
		if format == "ssa" {
			names = ssaEventNames
		}
		for _, change := range []string{"missing", "duplicate", "nonfinal", "alias-duplicate"} {
			declaration := append([]string{}, names...)
			switch change {
			case "missing":
				declaration = declaration[1:]
			case "duplicate":
				declaration = append(declaration, "Text")
			case "nonfinal":
				declaration[0], declaration[len(declaration)-1] = declaration[len(declaration)-1], declaration[0]
			case "alias-duplicate":
				declaration = append(declaration[:len(declaration)-1], "Actor", "Text")
			}
			value := strings.Repeat("0,", len(declaration)-1) + "content"
			if _, valid := scriptedSourceFields(value, declaration, format+"-event"); valid {
				t.Fatalf("%s %s declaration obtained content roles with exact value count", format, change)
			}
		}
		if values, valid := scriptedSourceFields(strings.Repeat("0,", len(names)-1)+"content, keeps, commas", names, format+"-event"); !valid || values[len(values)-1] != "content, keeps, commas" {
			t.Fatalf("%s final Text comma suffix not preserved: %v %v", format, values, valid)
		}
	}
}

func TestScriptedRetainedMalformedOwnersRemainSafeAndInspectable(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, kind := range []string{"event", "style"} {
			d := scriptedExample(t, format)
			n := d.FormatData.ASS
			if n == nil {
				n = d.FormatData.SSA
			}
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			order := 8
			line := "Comment: bounded,未知 ١٢"
			sectionID := "section-events"
			ownerID := "event-malformed"
			if kind == "style" {
				order = 9
				sectionID = "section-malformed-styles"
				name := "V4+ Styles"
				if format == "ssa" {
					name = "V4 Styles"
				}
				header := "[" + name + "]"
				n.Sections = append(n.Sections, ScriptedSection{SectionID: sectionID, SourceOrder: 8, Name: name, RawHeader: &header})
				b = append(b, []byte(header+"\n")...)
				line = "Style: bounded,未知 ١٢"
				ownerID = "style-malformed"
			}
			b = append(b, []byte(line+"\n")...)
			setScriptedSourceBytes(&d, b)
			record := ScriptedRecord{RecordID: "record-malformed", SourceOrder: order, SectionID: sectionID, Kind: kind, RawLine: &line}
			if kind == "style" {
				record.StyleID = &ownerID
				n.Styles = append(n.Styles, ScriptedStyle{StyleID: ownerID, RecordID: record.RecordID, Fields: []ScriptedField{}})
			} else {
				record.EventID = &ownerID
				text := "bounded,未知 ١٢"
				facts, err := ProjectScriptedText(text, 0, 0, 1)
				if err != nil {
					t.Fatal(err)
				}
				n.Events = append(n.Events, ScriptedEvent{EventID: ownerID, RecordID: record.RecordID, EventType: "comment", Text: text, Fields: []ScriptedField{}, Spans: facts.Spans, Tags: facts.Tags, Karaoke: facts.Karaoke})
			}
			n.Records = append(n.Records, record)
			d.Diagnostics = []Diagnostic{{Severity: "warning", Code: "malformed_native_record", Message: "Bounded uninterpretable native occurrence retained.", SourceOrder: &order}}
			d.Stats.DiagnosticCount, d.Stats.WarningCount = 1, 1
			if err := d.Validate(); err != nil {
				t.Fatalf("%s safe malformed %s cannot be retained: %v", format, kind, err)
			}
			value := "relative/private.ass"
			field := ScriptedField{FieldName: "Fontname", RawValue: value, TypedValue: &ScriptedValue{Kind: "string", String: &value}}
			if kind == "style" {
				n.Styles[len(n.Styles)-1].Fields = []ScriptedField{field}
			} else {
				field.FieldName = "Text"
				n.Events[len(n.Events)-1].Fields = []ScriptedField{field}
			}
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
				t.Fatalf("%s malformed %s borrowed typed content role: %v", format, kind, err)
			}
		}
	}
}

func TestScriptedAegisubExecutionEffectsRejectOriginalAndEditedValues(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, effect := range []string{"code once", "code line", "code syl", " CODE once noblank", "\tCoDe syl all", "code", "code custom_modifier", "template line", "!code once"} {
			d := scriptedExample(t, format)
			n := d.FormatData.ASS
			if n == nil {
				n = d.FormatData.SSA
			}
			for i := range n.Events[0].Fields {
				if nativeName(n.Events[0].Fields[i].FieldName) == "effect" {
					n.Events[0].Fields[i].RawValue = effect
				}
			}
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_active_content") {
				t.Fatalf("edited %s Effect %q accepted: %v", format, effect, err)
			}
			d = scriptedExample(t, format)
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = []byte(strings.Replace(string(b), ",,Hello", ","+effect+",Hello", 1))
			setScriptedSourceBytes(&d, b)
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_active_content") {
				t.Fatalf("original %s Effect %q accepted: %v", format, effect, err)
			}
		}
		for _, effect := range []string{"", "Banner;100;0", "Scroll up;0;100;20", "barcode once", "decode syl"} {
			if err := scriptedContentFieldPrivacy("Effect", effect, 0, format+"-event"); err != nil {
				t.Fatalf("inert native Effect %q rejected: %v", effect, err)
			}
		}
	}
}

func TestScriptedNativeDialogueIntervalIndependentOfCommonEdits(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, captured := range []bool{true, false} {
			for _, start := range []string{"0:00:02.00", "0:00:03.00"} {
				d := scriptedExample(t, format)
				n := d.FormatData.ASS
				cue := d.Cues[0].FormatData.ASS
				if n == nil {
					n = d.FormatData.SSA
					cue = d.Cues[0].FormatData.SSA
				}
				for i := range n.Events[0].Fields {
					if nativeName(n.Events[0].Fields[i].FieldName) == "start" {
						n.Events[0].Fields[i].RawValue = start
					}
				}
				if captured {
					b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
					b = []byte(strings.Replace(string(b), "0:00:01.00", start, 1))
					setScriptedSourceBytes(&d, b)
					line := strings.Replace(*n.Records[4].RawLine, "0:00:01.00", start, 1)
					n.Records[4].RawLine = &line
					cue.StartTimestampRaw = &start
				} else {
					n.Records[4].RawLine = nil
					cue.StartTimestampRaw, cue.EndTimestampRaw = nil, nil
				}
				if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "invalid_native_interval") {
					t.Fatalf("%s captured=%t native Start %s with End 2s accepted under positive common interval: %v", format, captured, start, err)
				}
			}
		}
		for _, deletedCues := range []bool{false, true} {
			d := scriptedExample(t, format)
			n := d.FormatData.ASS
			cue := d.Cues[0].FormatData.ASS
			if n == nil {
				n = d.FormatData.SSA
				cue = d.Cues[0].FormatData.SSA
			}
			n.Records[4].RawLine = nil
			cue.StartTimestampRaw, cue.EndTimestampRaw = nil, nil
			if deletedCues {
				n.Events = []ScriptedEvent{}
				n.Records = n.Records[:4]
				d.Cues = []Cue{}
				d.Document = DocumentSummary{}
				d.Stats = Stats{}
			}
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = []byte(strings.Replace(string(b), "0:00:01.00", "0:00:03.00", 1))
			setScriptedSourceBytes(&d, b)
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "invalid_native_interval") {
				t.Fatalf("%s original invalid interval escaped with omitted captures/deleted=%t: %v", format, deletedCues, err)
			}
		}
	}
}

func TestScriptedUnknownFieldsHaveContextualInertRoles(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, c := range []struct{ context, name, value string }{
			{"style", "SourceFile", "/home/alice/project.ass"}, {"style", "SourceFile", "relative.ass"}, {"style", "X-Native", "/workspace/private.ass"}, {"style", "ScriptType", "/root/private.ass"}, {"style", "Title", "/opt/private.ass"}, {"style", "Text", "/mnt/private.ass"}, {"event", "Fontname", "/workspace/private.ass"}, {"event", "SourceFile", "~/private.ass"}, {"style", "X-Native", "opaque/relative/path"},
		} {
			t.Run(format+"/"+c.context+"/"+c.name+"/"+c.value, func(t *testing.T) {
				d := scriptedExample(t, format)
				n := d.FormatData.ASS
				if n == nil {
					n = d.FormatData.SSA
				}
				field := ScriptedField{FieldName: c.name, RawValue: c.value}
				if c.context == "style" {
					n.Styles[0].Fields = append(n.Styles[0].Fields, field)
					n.Records[1].DeclarationFields = append(n.Records[1].DeclarationFields, ScriptedDeclarationField{FieldName: c.name})
				} else {
					event := &n.Events[0]
					event.Fields = append(event.Fields[:len(event.Fields)-1], field, event.Fields[len(event.Fields)-1])
					decl := &n.Records[3]
					decl.DeclarationFields = append(decl.DeclarationFields[:len(decl.DeclarationFields)-1], ScriptedDeclarationField{FieldName: c.name}, decl.DeclarationFields[len(decl.DeclarationFields)-1])
				}
				if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
					t.Fatalf("unsafe structured extension accepted: %v", err)
				}
				// The original grammar uses the same actual style/event context, including
				// names that happen to be canonical only in another native context.
				original := scriptedExample(t, format)
				b, _ := base64.StdEncoding.DecodeString(original.Source.Assets[0].DataBase64)
				lines := strings.Split(string(b), "\n")
				if c.context == "style" {
					lines[3] += "," + c.name
					lines[4] += "," + c.value
				} else {
					lines[6] = strings.Replace(lines[6], ", Text", ", "+c.name+", Text", 1)
					lines[7] = strings.Replace(lines[7], ",,Hello", ",,"+c.value+",Hello", 1)
				}
				setScriptedSourceBytes(&original, []byte(strings.Join(lines, "\n")))
				if err := original.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
					t.Fatalf("unsafe original extension accepted: %v", err)
				}
			})
		}
	}
	for _, profile := range []string{"ass-style", "ssa-style", "ass-event", "ssa-event"} {
		name := "Text"
		if strings.HasSuffix(profile, "-style") {
			name = "Fontname"
		}
		if err := scriptedContentFieldPrivacy(name, "content discusses /workspace/captions.ass and ~/notes", 0, profile); err != nil {
			t.Fatalf("explicit native content role %s: %v", profile, err)
		}
	}
}

func TestScriptedStructuralNamesCannotCarryResourceIdentity(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, name := range []string{"/home/alice/private.ass", "relative/private.ass", `relative\private.ass`, "localhost", "hostname=private"} {
			d := scriptedExample(t, format)
			n := d.FormatData.ASS
			if n == nil {
				n = d.FormatData.SSA
			}
			// Even an unused declaration is structural metadata, and an empty value
			// cannot make a resource-bearing field name inert.
			n.Records[1].DeclarationFields = append(n.Records[1].DeclarationFields, ScriptedDeclarationField{FieldName: name})
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
				t.Fatalf("%s edited declaration name %q accepted: %v", format, name, err)
			}
			if err := validateScriptedScalar(ScriptedField{FieldName: name}, format, 4, "style"); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
				t.Fatalf("%s native field name %q accepted: %v", format, name, err)
			}
			d = scriptedExample(t, format)
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = append(b, []byte("Format: "+name+"\n")...)
			setScriptedSourceBytes(&d, b)
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
				t.Fatalf("%s uninterpretable original declaration name %q accepted: %v", format, name, err)
			}
		}
		if err := validateScriptedScalar(ScriptedField{FieldName: "X-未知", RawValue: "inert extension"}, format, 4, "style"); err != nil {
			t.Fatalf("safe unknown Unicode field name rejected: %v", err)
		}
	}
}

func TestScriptedActorAliasBelongsOnlyToEventName(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		d := scriptedExample(t, format)
		n := d.FormatData.ASS
		if n == nil {
			n = d.FormatData.SSA
		}
		n.Styles[0].Fields[0].FieldName = "Actor"
		n.Records[1].DeclarationFields[0].FieldName = "Actor"
		if err := d.Validate(); err == nil {
			t.Fatalf("%s style Actor replaced mandatory Name", format)
		}
		styleNames := append([]string{}, assStyleNames...)
		if format == "ssa" {
			styleNames = append([]string{}, ssaStyleNames...)
		}
		styleNames[0] = "Actor"
		if _, valid := scriptedSourceFields(strings.Repeat("0,", len(styleNames)-1)+"0", styleNames, format+"-style"); valid {
			t.Fatalf("%s original style Actor replaced mandatory Name", format)
		}
		d = scriptedExample(t, format)
		n = d.FormatData.ASS
		if n == nil {
			n = d.FormatData.SSA
		}
		n.Styles[0].Fields = append(n.Styles[0].Fields, ScriptedField{FieldName: "Actor", RawValue: "inert extension"})
		n.Records[1].DeclarationFields = append(n.Records[1].DeclarationFields, ScriptedDeclarationField{FieldName: "Actor"})
		if err := d.Validate(); err != nil {
			t.Fatalf("%s safe extra style Actor was not retained as unknown: %v", format, err)
		}
		n.Styles[0].Fields[len(n.Styles[0].Fields)-1].RawValue = "relative/private.ass"
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
			t.Fatalf("%s extra style Actor borrowed event content exemption: %v", format, err)
		}
		d = scriptedExample(t, format)
		n = d.FormatData.ASS
		cue := d.Cues[0].FormatData.ASS
		if n == nil {
			n = d.FormatData.SSA
			cue = d.Cues[0].FormatData.SSA
		}
		n.Events[0].Fields[4].FieldName = "Actor"
		n.Records[3].DeclarationFields[4].FieldName = "Actor"
		actor := "Actor"
		cue.Projection.SpeakerFieldName = &actor
		b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		b = []byte(strings.Replace(string(b), ", Style, Name,", ", Style, Actor,", 1))
		setScriptedSourceBytes(&d, b)
		line := strings.Replace(*n.Records[3].RawLine, ", Style, Name,", ", Style, Actor,", 1)
		n.Records[3].RawLine = &line
		if err := d.Validate(); err != nil {
			t.Fatalf("%s genuine captured event Actor alias rejected: %v", format, err)
		}
		if n.Events[0].Fields[4].FieldName != "Actor" || n.Records[3].DeclarationFields[4].FieldName != "Actor" {
			t.Fatal("event Actor raw field names changed")
		}
	}
}

func TestScriptedSemicolonCommentsHaveOnlyAnInertRole(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, comment := range []string{"; generated from /home/alice/private.ass", `; content discussion C:\examples\captions.ass`, "; generated from relative/private.ass", "; hostname=private", "; code once", "; unclassifiable😀"} {
			for _, section := range []string{"Events", "V4 Styles", "V4+ Styles", "Script Info"} {
				if err := scriptedRecordPrivacy(ScriptedRecord{Kind: "comment", SourceOrder: 8, RawLine: &comment}, section); err == nil {
					t.Fatalf("%s structured comment in %s with %q accepted", format, section, comment)
				}
			}
			d := scriptedExample(t, format)
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = append(b, []byte(comment+"\n")...)
			setScriptedSourceBytes(&d, b)
			if err := d.Validate(); err == nil {
				t.Fatalf("%s original comment with %q accepted", format, comment)
			}
		}
		d := scriptedExample(t, format)
		n := d.FormatData.ASS
		if n == nil {
			n = d.FormatData.SSA
		}
		comment := "; safe 字幕 ١٢, ordered native notes."
		b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		b = append(b, []byte(comment+"\n")...)
		setScriptedSourceBytes(&d, b)
		n.Records = append(n.Records, ScriptedRecord{RecordID: "record-inert-comment", SourceOrder: 8, SectionID: "section-events", Kind: "comment", RawLine: &comment})
		if err := d.Validate(); err != nil {
			t.Fatalf("%s safe captured comment cannot be retained in order: %v", format, err)
		}
	}
}

func TestScriptedMetadataDetectsEveryAbsoluteAndHomePath(t *testing.T) {
	for _, path := range []string{"/workspace/private.ass", "/opt/private.ass", "/root/private.ass", "/mnt/private.ass", "/unusual-prefix/private.ass", "~/private.ass", "~someone/private.ass", "note=/workspace/private.ass", "note[/opt/private.ass", "note>/workspace/private.ass", "note\u00a0~/private.ass", "(/root/private.ass)", "ftp://private.invalid/media", "relative/private.ass", `relative\private.ass`, "./private.ass", `..\private.ass`} {
		d := scriptedExample(t, "ass")
		title := path
		d.Metadata.Title = &title
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
			t.Fatalf("document metadata path %q: %v", path, err)
		}
		d = scriptedExample(t, "ass")
		order := 7
		d.Diagnostics = []Diagnostic{{Severity: "warning", Code: "source_variant", Message: "Observed " + path, SourceOrder: &order}}
		d.Stats.DiagnosticCount = 1
		d.Stats.WarningCount = 1
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
			t.Fatalf("diagnostic path %q: %v", path, err)
		}
		d = scriptedExample(t, "ass")
		b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		b = []byte(strings.Replace(string(b), "ScriptType: v4.00+", "ScriptType: v4.00+\nTitle: "+path, 1))
		setScriptedSourceBytes(&d, b)
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "unsafe_source_metadata") {
			t.Fatalf("original metadata path %q: %v", path, err)
		}
	}
}

func TestScriptedCaptureObservationsRequireVerifiedPhysicalEvidence(t *testing.T) {
	for _, c := range []struct {
		name   string
		mutate func(*Document)
	}{
		{"raw recognized record", func(d *Document) { s := "Style: unrelated safe content"; d.FormatData.ASS.Records[2].RawLine = &s }},
		{"raw header", func(d *Document) { s := "[Unrelated Safe Section]"; d.FormatData.ASS.Sections[2].RawHeader = &s }},
		{"cue timestamp", func(d *Document) { s := "0:00:01.10"; d.Cues[0].FormatData.ASS.StartTimestampRaw = &s }},
		{"native timestamp occurrence", func(d *Document) { d.FormatData.ASS.Events[0].Fields[1].RawValue = "0:00:01.10" }},
	} {
		t.Run(c.name, func(t *testing.T) {
			d := scriptedExample(t, "ass")
			c.mutate(&d)
			if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "inconsistent_source_capture") {
				t.Fatalf("forged capture: %v", err)
			}
		})
	}
	for _, part := range []string{"hash", "length"} {
		d := scriptedExample(t, "ass")
		if part == "hash" {
			d.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64)
		} else {
			d.Source.Assets[0].Size.Bytes++
		}
		if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "source integrity") {
			t.Fatalf("unverified source %s: %v", part, err)
		}
	}
	for _, ending := range []string{"no_final_lf", "bom_crlf"} {
		d := scriptedExample(t, "ass")
		b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		if ending == "no_final_lf" {
			b = []byte(strings.TrimSuffix(string(b), "\n"))
		} else {
			b = []byte("\ufeff" + strings.ReplaceAll(string(b), "\n", "\r\n"))
			d.Source.Assets[0].Encoding = nil
		}
		setScriptedSourceBytes(&d, b)
		if err := d.Validate(); err != nil {
			t.Fatalf("truthful %s decoded observations: %v", ending, err)
		}
	}
	// Structured edits may diverge from original raw capture. Original timing
	// remains observed while common timing owns an independently editable value.
	d := scriptedExample(t, "ass")
	v := 22.0
	d.FormatData.ASS.Styles[0].Fields[2].RawValue = "22"
	d.FormatData.ASS.Styles[0].Fields[2].TypedValue.Decimal = &v
	start, span := int64(1011), int64(989)
	d.Cues[0].Timing.StartMilliseconds = start
	d.Cues[0].Timing.DurationMilliseconds = span
	d.Document.MediaStartMilliseconds = &start
	d.Document.MediaSpanMilliseconds = &span
	d.Stats.MediaSpanMilliseconds = &span
	if err := d.Validate(); err != nil {
		t.Fatalf("editable owners with truthful old captures: %v", err)
	}
	n := d.FormatData.ASS
	for i := range n.Sections {
		n.Sections[i].RawHeader = nil
	}
	for i := range n.Records {
		n.Records[i].RawLine = nil
	}
	n.Styles[0].DeclarationID = nil
	n.Events[0].DeclarationID = nil
	d.Cues[0].FormatData.ASS.StartTimestampRaw = nil
	d.Cues[0].FormatData.ASS.EndTimestampRaw = nil
	n.Events[0].Fields[1].RawValue = "0:00:01.10"
	if err := d.Validate(); err != nil {
		t.Fatalf("constructed owners without invented captures: %v", err)
	}
}

func TestScriptedCaptureOriginalReorderedStartFirstDeclaration(t *testing.T) {
	d := scriptedExample(t, "ass")
	n := d.FormatData.ASS
	event := &n.Events[0]
	event.Fields[0], event.Fields[1] = event.Fields[1], event.Fields[0]
	decl := &n.Records[3]
	decl.DeclarationFields[0], decl.DeclarationFields[1] = decl.DeclarationFields[1], decl.DeclarationFields[0]
	b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	lines := strings.Split(string(b), "\n")
	lines[6] = strings.Replace(lines[6], "Layer, Start", "Start, Layer", 1)
	lines[7] = strings.Replace(lines[7], "Dialogue: 0,0:00:01.00", "Dialogue:   0:00:01.00,0", 1)
	decl.RawLine = &lines[6]
	n.Records[4].RawLine = &lines[7]
	setScriptedSourceBytes(&d, []byte(strings.Join(lines, "\n")))
	if err := d.Validate(); err != nil {
		t.Fatalf("original Start-first framing: %v", err)
	}
}

func TestScriptedUnknownPrefixOverridesRemainExactAndDiagnosed(t *testing.T) {
	for _, text := range []string{`{\random1}text`, `{\fnonsense1}text`, `{\rdefault}text`, `{\fnarial}text`, `{\rUndeclared}text`, `{\fnUnknown Font}text`} {
		facts, err := ProjectScriptedText(text, 0, 1000, 2000)
		if err != nil || len(facts.Tags) != 1 || !slices.Contains(facts.DiagnosticCodes, "unsupported_override") {
			t.Fatalf("opaque prefix %q: %#v %v", text, facts, err)
		}
		if facts.Tags[0].Name == "r" || facts.Tags[0].Name == "fn" || facts.Tags[0].Raw != text[1:len(text)-5] {
			t.Fatalf("native tag identity/lexeme lost: %#v", facts.Tags)
		}
	}
	names := ScriptedTextNames{ResetStyles: map[string]bool{"Default": true, "Style1": true}, FontFamilies: map[string]bool{"Arial": true, "Noto Sans": true}}
	for _, c := range []struct{ text, name, parameter string }{{`{\rDefault}text`, "r", "Default"}, {`{\rStyle1}text`, "r", "Style1"}, {`{\fnArial}text`, "fn", "Arial"}, {`{\fnNoto Sans}text`, "fn", "Noto Sans"}} {
		facts, err := ProjectScriptedText(c.text, 0, 1000, 2000, names)
		if err != nil || len(facts.Tags) != 1 || facts.Tags[0].Name != c.name || facts.Tags[0].Parameter != c.parameter || len(facts.DiagnosticCodes) != 0 {
			t.Fatalf("declared prefix identity %q: %#v %v", c.text, facts, err)
		}
	}
}

func TestScriptedCapturedCommentAndRepeatedEventsEvidence(t *testing.T) {
	d := scriptedExample(t, "ass")
	n := d.FormatData.ASS
	event := &n.Events[0]
	event.EventType = "comment"
	event.CueID = nil
	d.Cues = []Cue{}
	d.Document = DocumentSummary{}
	d.Stats = Stats{}
	b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	lines := strings.Split(string(b), "\n")
	lines[7] = strings.Replace(lines[7], "Dialogue:", "Comment:", 1)
	n.Records[4].RawLine = &lines[7]
	setScriptedSourceBytes(&d, []byte(strings.Join(lines, "\n")))
	if err := d.Validate(); err != nil {
		t.Fatalf("valid captured Comment: %v", err)
	}
	event.Fields[1].RawValue = "0:00:01.10"
	if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "inconsistent_source_capture") {
		t.Fatalf("forged Comment timestamp: %v", err)
	}
	event.Fields[1].RawValue = "0:00:01.00"
	second := *event
	second.EventID = "event-second"
	second.RecordID = "record-comment-second"
	second.Fields = append([]ScriptedField{}, event.Fields...)
	second.Fields[1].RawValue = "0:00:03.00"
	second.Fields[2].RawValue = "0:00:04.00"
	second.Fields[0], second.Fields[1] = second.Fields[1], second.Fields[0]
	declarationID := "record-format-second"
	second.DeclarationID = &declarationID
	declaration := append([]ScriptedDeclarationField{}, n.Records[3].DeclarationFields...)
	declaration[0], declaration[1] = declaration[1], declaration[0]
	header := "[Events]"
	formatLine := strings.Replace(lines[6], "Layer, Start", "Start, Layer", 1)
	values := []string{}
	for _, field := range second.Fields {
		values = append(values, field.RawValue)
	}
	commentLine := "Comment:   " + strings.Join(values, ",")
	n.Sections = append(n.Sections, ScriptedSection{SectionID: "section-events-second", SourceOrder: 8, Name: "Events", RawHeader: &header})
	n.Records = append(n.Records, ScriptedRecord{RecordID: declarationID, SourceOrder: 9, SectionID: "section-events-second", Kind: "format_declaration", RawLine: &formatLine, DeclarationFields: declaration}, ScriptedRecord{RecordID: second.RecordID, SourceOrder: 10, SectionID: "section-events-second", Kind: "event", RawLine: &commentLine, EventID: &second.EventID})
	n.Events = append(n.Events, second)
	b, _ = base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	setScriptedSourceBytes(&d, append(b, []byte(header+"\n"+formatLine+"\n"+commentLine+"\n")...))
	if err := d.Validate(); err != nil {
		t.Fatalf("repeated Events with distinct active Start-first Format: %v", err)
	}
	n.Events[1].Fields[0].RawValue = "0:00:01.00"
	if err := d.Validate(); err == nil || !strings.Contains(err.Error(), "inconsistent_source_capture") {
		t.Fatalf("capture borrowed from prior section: %v", err)
	}
}
