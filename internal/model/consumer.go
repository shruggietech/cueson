package model

import (
	"fmt"
	"math"
	"unicode"
	"unicode/utf8"
)

// SpeakerAttribution stores opaque consumer-assigned identity independently
// from native subtitle observations. Consumers determine identifier scope.
type SpeakerAttribution struct {
	SpeakerID         string `json:"speaker_id"`
	StartMilliseconds *int64 `json:"start_milliseconds,omitempty"`
	EndMilliseconds   *int64 `json:"end_milliseconds,omitempty"`
}

// MediaTiming is a consumer declaration, never a duration inferred from cues
// or an independently verified measurement of source media.
type MediaTiming struct {
	DurationMilliseconds      int64  `json:"duration_milliseconds"`
	TimelineStartMilliseconds *int64 `json:"timeline_start_milliseconds,omitempty"`
}

// ConsumerAnnotationSummary contains counts and states, never speaker values.
type ConsumerAnnotationSummary struct {
	AttributionCount        int    `json:"attribution_count"`
	TimedAttributionCount   int    `json:"timed_attribution_count"`
	UntimedAttributionCount int    `json:"untimed_attribution_count"`
	MediaTimingPresent      bool   `json:"media_timing_present"`
	MediaBoundaryCheck      string `json:"media_boundary_check"`
	CueMediaConflictCount   int    `json:"cue_media_conflict_count"`
}

func hasConsumerContract(document Document) bool {
	return document.SchemaVersion == "1.2.0-dev" && document.Schema == "https://cueson.io/schema/v1.2.0-dev/cueson.schema.json"
}

func hasScriptedContract(document Document) bool {
	return hasConsumerContract(document) || document.SchemaVersion == "1.1.0" && document.Schema == "https://cueson.io/schema/v1.1.0/cueson.schema.json"
}

// ValidateConsumerAnnotations revalidates consumer additions for native APIs
// that accept partial native representations rather than full source envelopes.
func (document Document) ValidateConsumerAnnotations() error {
	return validateConsumerAnnotations(document)
}

func validateConsumerAnnotations(document Document) error {
	current := hasConsumerContract(document)
	if document.MediaTiming != nil && !current {
		return fmt.Errorf("media_timing requires the current consumer contract")
	}
	var mediaStart, mediaEnd int64
	if document.MediaTiming != nil {
		var err error
		mediaStart, mediaEnd, err = consumerMediaBounds(document.MediaTiming)
		if err != nil {
			return err
		}
	}
	// Bounds precede item traversal, including direct typed callers.
	if err := boundedLength("cues", len(document.Cues), MaxDocumentItems); err != nil {
		return err
	}
	for cueIndex, cue := range document.Cues {
		prefix := fmt.Sprintf("cues[%d].speaker_attributions", cueIndex)
		if cue.SpeakerAttributions != nil && !current {
			return fmt.Errorf("%s requires the current consumer contract", prefix)
		}
		if err := boundedLength(prefix, len(cue.SpeakerAttributions), MaxItemOccurrences); err != nil {
			return err
		}
		for index, attribution := range cue.SpeakerAttributions {
			path := fmt.Sprintf("%s[%d]", prefix, index)
			if !validConsumerSpeakerID(attribution.SpeakerID) {
				return fmt.Errorf("%s.speaker_id must contain 1..256 Unicode scalar values without excluded controls or boundary whitespace", path)
			}
			if (attribution.StartMilliseconds == nil) != (attribution.EndMilliseconds == nil) {
				return fmt.Errorf("%s requires both start_milliseconds and end_milliseconds or neither", path)
			}
			if attribution.StartMilliseconds == nil {
				continue
			}
			start, end := *attribution.StartMilliseconds, *attribution.EndMilliseconds
			if start < 0 || end <= start {
				return fmt.Errorf("%s requires a positive nonnegative integer millisecond interval", path)
			}
			if start < cue.Timing.StartMilliseconds || end > cue.Timing.EndMilliseconds {
				return fmt.Errorf("%s timing falls outside its cue", path)
			}
			if document.MediaTiming != nil && (start < mediaStart || end > mediaEnd) {
				return fmt.Errorf("%s timing falls outside declared media_timing", path)
			}
		}
	}
	return nil
}

func validConsumerSpeakerID(value string) bool {
	if value == "" || len(value) > 256*utf8.UTFMax || !utf8.ValidString(value) || utf8.RuneCountInString(value) > 256 {
		return false
	}
	first, _ := utf8.DecodeRuneInString(value)
	last, _ := utf8.DecodeLastRuneInString(value)
	if unicode.IsSpace(first) || unicode.IsSpace(last) {
		return false
	}
	for _, r := range value {
		if r <= 0x1f || r >= 0x7f && r <= 0x9f || r == 0x2028 || r == 0x2029 || r == 0x061c || r == 0x200e || r == 0x200f || r >= 0x202a && r <= 0x202e || r >= 0x2066 && r <= 0x2069 || r == 0xfeff {
			return false
		}
	}
	return true
}

func consumerMediaBounds(timing *MediaTiming) (int64, int64, error) {
	start := int64(0)
	if timing.TimelineStartMilliseconds != nil {
		start = *timing.TimelineStartMilliseconds
	}
	if timing.DurationMilliseconds < 0 {
		return 0, 0, fmt.Errorf("media_timing.duration_milliseconds must be nonnegative")
	}
	if start > math.MaxInt64-timing.DurationMilliseconds {
		return 0, 0, fmt.Errorf("media_timing endpoint exceeds the integer millisecond range")
	}
	return start, start + timing.DurationMilliseconds, nil
}

// ConsumerAnnotationsSummary assesses already-validated consumer values without
// mutating source diagnostics or claiming identity/media verification.
func (document Document) ConsumerAnnotationsSummary() *ConsumerAnnotationSummary {
	summary := ConsumerAnnotationSummary{MediaTimingPresent: document.MediaTiming != nil, MediaBoundaryCheck: "unavailable"}
	for _, cue := range document.Cues {
		for _, attribution := range cue.SpeakerAttributions {
			summary.AttributionCount++
			if attribution.StartMilliseconds != nil && attribution.EndMilliseconds != nil {
				summary.TimedAttributionCount++
			} else {
				summary.UntimedAttributionCount++
			}
		}
	}
	if !summary.MediaTimingPresent && summary.AttributionCount == 0 {
		return nil
	}
	if summary.MediaTimingPresent {
		summary.MediaBoundaryCheck = "not_evaluated"
		if summary.TimedAttributionCount > 0 {
			summary.MediaBoundaryCheck = "checked"
		}
		start, end, err := consumerMediaBounds(document.MediaTiming)
		if err == nil {
			for _, cue := range document.Cues {
				if cue.Timing.StartMilliseconds < start || cue.Timing.EndMilliseconds > end {
					summary.CueMediaConflictCount++
				}
			}
		}
	}
	return &summary
}

// ConsumerTimingDiagnostics reports source cue conflicts as one bounded runtime
// warning. These conflicts never invalidate a valid preserved source envelope.
func (document Document) ConsumerTimingDiagnostics() ([]Diagnostic, error) {
	if document.MediaTiming == nil {
		return []Diagnostic{}, nil
	}
	if _, _, err := consumerMediaBounds(document.MediaTiming); err != nil {
		return nil, err
	}
	summary := document.ConsumerAnnotationsSummary()
	if summary.CueMediaConflictCount == 0 {
		return []Diagnostic{}, nil
	}
	return []Diagnostic{{Severity: "warning", Code: "consumer_cue_media_conflict", Message: fmt.Sprintf("media_timing conflicts with %d cue intervals; source cues remain unchanged", summary.CueMediaConflictCount)}}, nil
}

// ConsumerAnnotationLossDiagnostics reports each omitted consumer field before
// native output publication. The cap is checked before allocating observations.
func (document Document) ConsumerAnnotationLossDiagnostics() ([]Diagnostic, error) {
	count := 0
	if document.MediaTiming != nil {
		count++
	}
	for _, cue := range document.Cues {
		if len(cue.SpeakerAttributions) > MaxDiagnostics-count {
			return nil, fmt.Errorf("consumer annotation omissions exceed maximum diagnostics %d", MaxDiagnostics)
		}
		count += len(cue.SpeakerAttributions)
	}
	result := make([]Diagnostic, 0, count)
	for cueIndex, cue := range document.Cues {
		for index := range cue.SpeakerAttributions {
			result = append(result, Diagnostic{Severity: "warning", Code: "consumer_speaker_attribution_omitted", Message: fmt.Sprintf("cues[%d].speaker_attributions[%d] is omitted from native output", cueIndex, index)})
		}
	}
	if document.MediaTiming != nil {
		result = append(result, Diagnostic{Severity: "warning", Code: "consumer_media_timing_omitted", Message: "media_timing is omitted from native output"})
	}
	return result, nil
}
