package conformance_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/testutil"
)

var requiredScriptedMatrixRowIDs = []string{
	"scripted-dialect-detection", "scripted-utf8-source", "scripted-declared-fields", "scripted-native-order", "scripted-styles", "scripted-event-semantics", "scripted-centisecond-timer", "scripted-text-projection", "scripted-overrides", "scripted-karaoke-speakers", "scripted-attachments", "scripted-privacy-active-content", "scripted-edit-ownership", "scripted-render-cycle", "scripted-complete-conversion", "scripted-hostile-bounds", "scripted-stable-gate",
}

func validateScriptedMatrixRows(rows []testutil.MatrixRow) error {
	if len(rows) != len(requiredScriptedMatrixRowIDs) {
		return fmt.Errorf("scripted matrix has %d rows, want all seventeen selected-profile rows", len(rows))
	}
	required := map[string]bool{}
	for _, id := range requiredScriptedMatrixRowIDs {
		required[id] = true
	}
	seen := map[string]bool{}
	for _, row := range rows {
		if !required[row.RowID] || seen[row.RowID] {
			return fmt.Errorf("scripted matrix has unapproved or duplicate row %q", row.RowID)
		}
		seen[row.RowID] = true
	}
	return nil
}

var matrixIssueReference = regexp.MustCompile(`#[0-9]+\b`)

func allowedStableMatrixDeferral(format, rowID string, evidence testutil.MatrixEvidence) bool {
	if format != "scripted" || rowID != "scripted-stable-gate" || evidence.Kind != "deferred" || (evidence.Reference != "render_test" && evidence.Reference != "platform_test") || strings.TrimSpace(evidence.Reason) == "" {
		return false
	}
	issues := matrixIssueReference.FindAllString(evidence.Reason, -1)
	if len(issues) == 0 {
		return false
	}
	for _, issue := range issues {
		if issue != "#64" && issue != "#65" {
			return false
		}
	}
	return true
}

func TestScriptedMatrixRejectsSubstitutedRowsAndDevelopmentDeferrals(t *testing.T) {
	t.Parallel()
	rows := make([]testutil.MatrixRow, len(requiredScriptedMatrixRowIDs))
	for index, id := range requiredScriptedMatrixRowIDs {
		rows[index].RowID = id
	}
	if err := validateScriptedMatrixRows(rows); err != nil {
		t.Fatal(err)
	}
	rows[0].RowID = "scripted-invented-but-seventeen"
	if err := validateScriptedMatrixRows(rows); err == nil {
		t.Fatal("substituted required row accepted")
	}
	rows[0].RowID = rows[1].RowID
	if err := validateScriptedMatrixRows(rows); err == nil {
		t.Fatal("duplicate required row accepted")
	}
	for _, test := range []struct {
		row, reference, reason string
		want                   bool
	}{
		{"scripted-stable-gate", "render_test", "Stable freeze and reviewed candidate proof belong to #64 and #65.", true},
		{"scripted-stable-gate", "platform_test", "Native immutable candidate release proof belongs to #65.", true},
		{"scripted-stable-gate", "platform_test", "Hardening remains #63; stable release #64 and #65.", false},
		{"scripted-hostile-bounds", "platform_test", "Hardening #64 and #65.", false},
		{"scripted-stable-gate", "conversion_test", "Stable #64 and #65.", false},
		{"scripted-stable-gate", "render_test", "Stable release #66.", false},
		{"scripted-stable-gate", "render_test", "Hashtag #junk.", false},
	} {
		t.Run(test.row+"/"+test.reason, func(t *testing.T) {
			evidence := testutil.MatrixEvidence{Kind: "deferred", Reference: test.reference, Reason: test.reason}
			if got := allowedStableMatrixDeferral("scripted", test.row, evidence); got != test.want {
				t.Fatalf("deferral allowed=%t want=%t", got, test.want)
			}
		})
	}
}
