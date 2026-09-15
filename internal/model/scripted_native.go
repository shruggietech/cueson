package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ScriptedCanonicalFields returns a copy of the ratified dialect field order.
func ScriptedCanonicalFields(format, context string) []string {
	var names []string
	switch format + "-" + context {
	case "ass-style":
		names = assStyleNames
	case "ssa-style":
		names = ssaStyleNames
	case "ass-event":
		names = assEventNames
	case "ssa-event":
		names = ssaEventNames
	}
	return append([]string(nil), names...)
}

// ScriptedFieldIdentity recognizes ASCII field names in their owning context.
func ScriptedFieldIdentity(name, format, context string) (string, bool) {
	profile := format + "-" + context
	if context == "metadata" {
		profile = context
	}
	n := nativeFieldName(name, profile)
	return n, scriptedFieldProfiles[profile][n]
}

// ScriptedSplitFields applies the same declaration framing as source validation.
func ScriptedSplitFields(value string, declaration []string, format, context string) ([]string, bool) {
	return scriptedSourceFields(value, declaration, format+"-"+context)
}

// ScriptedMilliseconds parses the exact checked native timestamp profile.
func ScriptedMilliseconds(raw string) (int64, error) { return scriptedMilliseconds(raw) }

// ValidateScriptedNativeField shares scalar and privacy authority with codecs.
func ValidateScriptedNativeField(field ScriptedField, format, context string, order int) error {
	return validateScriptedScalar(field, format, order, context)
}

// ValidateScriptedNativeRecord shares closed record-context privacy authority.
func ValidateScriptedNativeRecord(record ScriptedRecord, section string) error {
	return scriptedRecordPrivacy(record, section)
}

// ValidateScriptedNativeSection applies the model's section name restrictions.
func ValidateScriptedNativeSection(name string, order int) error {
	if !scriptedPhysical(name) || strings.ContainsAny(name, "[]") || strings.TrimSpace(name) == "" {
		return scriptedError("invalid_native_section", order)
	}
	if containsScriptedMetadataIdentity(name) {
		return scriptedError("unsafe_source_metadata", order)
	}
	return nil
}

// ValidateScriptedMalformedValue grants no native content role exemptions.
func ValidateScriptedMalformedValue(value string, order int) error {
	return scriptedMalformedValuePrivacy(value, order)
}

// ScriptedNativeSafeName applies the portable embedded-attachment name profile.
func ScriptedNativeSafeName(name string) bool { return nativeSafeName(name) }

// ScriptedScalarValue derives a checked lexical scalar without altering it.
// Validation remains the authority, including dialect/context and privacy.
func ScriptedScalarValue(field ScriptedField, format, context string, order int) (ScriptedValue, error) {
	if err := ValidateScriptedNativeField(field, format, context, order); err != nil {
		return ScriptedValue{}, err
	}
	n, known := ScriptedFieldIdentity(field.FieldName, format, context)
	raw := strings.TrimSpace(field.RawValue)
	v := ScriptedValue{Kind: "string", String: &field.RawValue}
	if !known {
		return v, nil
	}
	switch n {
	case "bold", "italic", "underline", "strikeout":
		i, _ := strconv.ParseInt(raw, 10, 64)
		b := i != 0
		v = ScriptedValue{Kind: "boolean", Boolean: &b}
	case "fontsize", "scalex", "scaley", "spacing", "angle", "outline", "shadow":
		d, _ := strconv.ParseFloat(raw, 64)
		v = ScriptedValue{Kind: "decimal", Decimal: &d}
	case "layer", "marked", "borderstyle", "alignment", "marginl", "marginr", "marginv", "encoding", "alphalevel":
		if n == "marked" {
			raw = strings.TrimPrefix(raw, "Marked=")
		}
		i, _ := strconv.ParseInt(raw, 10, 64)
		v = ScriptedValue{Kind: "integer", Integer: &i}
	case "primarycolour", "secondarycolour", "outlinecolour", "tertiarycolour", "backcolour":
		var bits uint64
		if strings.HasPrefix(strings.ToUpper(raw), "&H") {
			bits, _ = strconv.ParseUint(strings.TrimSuffix(raw[2:], "&"), 16, 32)
		} else if strings.HasPrefix(raw, "-") {
			i, err := strconv.ParseInt(raw, 10, 32)
			if err != nil {
				return ScriptedValue{}, scriptedError("invalid_native_field_value", order)
			}
			bits = uint64(uint32(int32(i)))
		} else {
			bits, _ = strconv.ParseUint(strings.TrimPrefix(raw, "+"), 10, 32)
		}
		c := ScriptedColor{Alpha: uint8(bits >> 24), Blue: uint8(bits >> 16), Green: uint8(bits >> 8), Red: uint8(bits)}
		v = ScriptedValue{Kind: "color", Color: &c}
	}
	// Recheck the derived union against the owning authority. This protects the
	// adapter if the accepted scalar profile changes independently later.
	check := field
	check.TypedValue = &v
	if err := ValidateScriptedNativeField(check, format, context, order); err != nil {
		return ScriptedValue{}, err
	}
	return v, nil
}

// ScriptedAttachmentFacts checks retained 6-bit content without decoding assets.
func ScriptedAttachmentFacts(lines []string) (decodedBytes int, malformed bool, err error) {
	count := 0
	var last byte
	for _, line := range lines {
		if len(line) > 80 || count > math.MaxInt-len(line) {
			return 0, false, fmt.Errorf("complexity_limit")
		}
		for i := 0; i < len(line); i++ {
			if line[i] < 33 || line[i] > 96 {
				malformed = true
			}
			last = line[i]
		}
		count += len(line)
	}
	decodedBytes = count / 4 * 3
	switch count % 4 {
	case 1:
		malformed = true
	case 2:
		decodedBytes++
		malformed = malformed || ((last-33)&15 != 0)
	case 3:
		decodedBytes += 2
		malformed = malformed || ((last-33)&3 != 0)
	}
	if decodedBytes > 16<<20 {
		return 0, false, fmt.Errorf("complexity_limit")
	}
	return decodedBytes, malformed, nil
}
