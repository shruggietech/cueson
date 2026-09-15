package conformance_test

import (
	"fmt"
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
		for _, evidence := range row.Evidence {
			if evidence.Kind == "deferred" {
				return fmt.Errorf("stable scripted row %q has deferred evidence", row.RowID)
			}
		}
	}
	return nil
}

func TestScriptedMatrixRejectsSubstitutedRowsAndAllDeferrals(t *testing.T) {
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
	rows[0].RowID = requiredScriptedMatrixRowIDs[0]
	for index, id := range requiredScriptedMatrixRowIDs {
		for _, kind := range requiredMatrixEvidenceKinds {
			rows[index].Evidence = []testutil.MatrixEvidence{{Kind: "deferred", Reference: kind, Reason: "Stable freeze #64 and candidate proof #65."}}
			if err := validateScriptedMatrixRows(rows); err == nil {
				t.Fatalf("stable row %s accepted deferred %s evidence", id, kind)
			}
			rows[index].Evidence = nil
		}
	}
}
