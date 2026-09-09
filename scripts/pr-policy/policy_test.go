package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	testHead = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testOld  = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

var (
	testStart = time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	codex     = Actor{Login: CodexRESTLogin, ID: CodexUserID, Type: "Bot"}
	operator  = Actor{Login: OperatorLogin, ID: 46768484, Type: "User"}
)

type fixtureCase struct {
	Name          string      `json:"name"`
	Scenario      string      `json:"scenario"`
	ExpectedState StatusState `json:"expected_state"`
	RequestSecond bool        `json:"request_second"`
}

func TestFixtureDocumentsAreValidAndUnique(t *testing.T) {
	t.Parallel()

	for _, group := range []string{"links", "reviews", "overrides"} {
		cases := loadFixtureCases(t, group)
		if len(cases) == 0 {
			t.Fatalf("%s fixture corpus is empty", group)
		}
		seen := make(map[string]bool, len(cases))
		for _, fixture := range cases {
			if fixture.Name == "" || fixture.Scenario == "" || fixture.ExpectedState == "" {
				t.Fatalf("%s has an incomplete fixture: %+v", group, fixture)
			}
			if seen[fixture.Name] {
				t.Fatalf("%s repeats fixture name %q", group, fixture.Name)
			}
			seen[fixture.Name] = true
		}
	}
}

func TestActorIdentityIsExact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		actor      Actor
		codex      bool
		dependabot bool
	}{
		{name: "codex REST", actor: codex, codex: true},
		{name: "codex GraphQL spelling", actor: Actor{Login: "chatgpt-codex-connector", ID: CodexUserID, Type: "Bot"}, codex: true},
		{name: "codex wrong id", actor: Actor{Login: CodexRESTLogin, ID: 1, Type: "Bot"}},
		{name: "codex wrong type", actor: Actor{Login: CodexRESTLogin, ID: CodexUserID, Type: "User"}},
		{name: "lookalike bot", actor: Actor{Login: "chatgpt-codex-connector-helper[bot]", ID: CodexUserID, Type: "Bot"}},
		{name: "dependabot", actor: Actor{Login: DependabotLogin, ID: DependabotUserID, Type: "Bot"}, dependabot: true},
		{name: "dependabot wrong id", actor: Actor{Login: DependabotLogin, ID: 1, Type: "Bot"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsCodexActor(test.actor); got != test.codex {
				t.Fatalf("IsCodexActor() = %v, want %v", got, test.codex)
			}
			if got := IsDependabotActor(test.actor); got != test.dependabot {
				t.Fatalf("IsDependabotActor() = %v, want %v", got, test.dependabot)
			}
		})
	}
}

func TestDirectiveParsingIsExact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, author, body, policy, reason string
		ok                                 bool
	}{
		{name: "issue link", author: OperatorLogin, body: "skip: issue-link - bootstrap tracking is external", policy: "issue-link", reason: "bootstrap tracking is external", ok: true},
		{name: "review", author: OperatorLogin, body: "context\nskip: codex-review - infrastructure failure\nmore", policy: "codex-review", reason: "infrastructure failure", ok: true},
		{name: "wrong author", author: "another-maintainer", body: "skip: codex-review - no", ok: false},
		{name: "missing reason", author: OperatorLogin, body: "skip: codex-review -   ", ok: false},
		{name: "wrong case", author: OperatorLogin, body: "Skip: codex-review - no", ok: false},
		{name: "not exact line", author: OperatorLogin, body: "> skip: codex-review - quoted", ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy, reason, ok := ParseDirective(Comment{Author: Actor{Login: test.author, Type: "User"}, Body: test.body})
			if ok != test.ok || policy != test.policy || reason != test.reason {
				t.Fatalf("ParseDirective() = %q, %q, %v; want %q, %q, %v", policy, reason, ok, test.policy, test.reason, test.ok)
			}
		})
	}
}

func TestIssueLinkFixtures(t *testing.T) {
	t.Parallel()

	for _, fixture := range loadFixtureCases(t, "links") {
		fixture := fixture
		t.Run(fixture.Name, func(t *testing.T) {
			snapshot := baseSnapshot()
			switch fixture.Scenario {
			case "local":
				snapshot.ClosingIssues = []ClosingIssue{{Owner: "shruggietech", Repository: "cueson", Number: 9}}
			case "cross-repository":
				snapshot.ClosingIssues = []ClosingIssue{{Owner: "shruggietech", Repository: "planning", Number: 42}}
			case "operator-exception":
				snapshot.Comments = []Comment{{ID: 1, Author: operator, Body: "skip: issue-link - issue is tracked outside this repository"}}
			case "untrusted-exception":
				snapshot.Comments = []Comment{{ID: 1, Author: Actor{Login: "outside", ID: 99, Type: "User"}, Body: "skip: issue-link - let me through"}}
			case "dependabot":
				snapshot.Author = Actor{Login: DependabotLogin, ID: DependabotUserID, Type: "Bot"}
				snapshot.AuthorLogin = DependabotLogin
			case "missing":
			default:
				t.Fatalf("unknown fixture scenario %q", fixture.Scenario)
			}

			result := EvaluateIssueLink(snapshot)
			assertResult(t, result, fixture.ExpectedState, fixture.RequestSecond)
			if !result.Publish || result.HeadSHA != testHead || result.Context != IssueLinkContext {
				t.Fatalf("result lacks current-head publication identity: %+v", result)
			}
		})
	}
}

func TestOperatorWaiversRequireStableCommentIdentifiers(t *testing.T) {
	t.Parallel()

	for _, policy := range []string{"issue-link", "codex-review"} {
		snapshot := baseSnapshot()
		snapshot.Comments = []Comment{{Author: operator, Body: "skip: " + policy + " - emergency recovery"}}
		var result PolicyResult
		if policy == "issue-link" {
			result = EvaluateIssueLink(snapshot)
		} else {
			result = EvaluateCodexReview(snapshot)
		}
		if result.State != StatusError || result.RequestSecond {
			t.Fatalf("%s accepted unstable waiver evidence: %+v", policy, result)
		}
	}
}

func TestCodexSummaryParsing(t *testing.T) {
	t.Parallel()

	body := "<!-- codex-pull-request-review-summary -->\n\n| Review | Status | Commit | Review trigger |\n| --- | --- | --- | --- |\n| Code Review | ✅ **Completed** | `" + testHead[:7] + "` | Automatic |\n"
	summary, ok, err := ParseCodexSummary(Comment{ID: 3, Author: codex, Body: body, UpdatedAt: testStart.Add(time.Minute)})
	if err != nil || !ok {
		t.Fatalf("ParseCodexSummary() failed: ok=%v err=%v", ok, err)
	}
	if summary.Status != SummaryCompleted || summary.CommitPrefix != testHead[:7] {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	runningBody := strings.Replace(body, "✅ **Completed**", "🔄 **Running**", 1)
	running, ok, err := ParseCodexSummary(Comment{ID: 6, Author: codex, Body: runningBody, UpdatedAt: testStart.Add(time.Minute)})
	if err != nil || !ok || running.Status != SummaryPending {
		t.Fatalf("empirical running summary failed: summary=%+v ok=%v err=%v", running, ok, err)
	}

	_, ok, err = ParseCodexSummary(Comment{ID: 4, Author: codex, Body: "<!-- codex-pull-request-review-summary -->\nunknown"})
	if !ok || err == nil {
		t.Fatalf("malformed marked summary = ok %v, err %v; want recognized error", ok, err)
	}

	_, ok, err = ParseCodexSummary(Comment{ID: 5, Author: operator, Body: body})
	if ok || err != nil {
		t.Fatalf("untrusted summary = ok %v, err %v; want ignored", ok, err)
	}
}

func TestStaleSummaryReactionDoesNotApproveEditedCurrentHead(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	addSummary(&snapshot, SummaryCompleted, testHead, testStart.Add(2*time.Minute))
	commentID := snapshot.Comments[0].ID
	snapshot.Reactions = []Reaction{{Author: codex, Content: "+1", Target: fmt.Sprintf("comment:%d", commentID), CreatedAt: testStart.Add(time.Minute)}}
	result := EvaluateCodexReview(snapshot)
	assertResult(t, result, StatusPending, false)
}

func TestEmpiricalFailureCommentParsesFullUnquotedRef(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	snapshot.Comments = []Comment{{
		ID:        9,
		Author:    codex,
		CreatedAt: testStart,
		Body:      "Codex Review: Something went wrong. Please try again.\n\n```\nProvided git ref " + testHead + " does not exist\n```",
	}}
	result := EvaluateCodexReview(snapshot)
	assertResult(t, result, StatusFailure, false)
}

func TestCleanTerminalCommentStillRequiresCurrentRoundThumbsUp(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	snapshot.Comments = []Comment{{
		ID:        24,
		Author:    codex,
		Body:      "Codex Review: Didn't find any major issues.\n\n**Reviewed commit:** `" + testHead[:10] + "`",
		CreatedAt: testStart.Add(time.Minute),
	}}
	withoutReaction := EvaluateCodexReview(snapshot)
	assertResult(t, withoutReaction, StatusPending, false)

	snapshot.Reactions = []Reaction{{Author: codex, Content: "+1", Target: "pull_request", CreatedAt: testStart.Add(2 * time.Minute)}}
	withReaction := EvaluateCodexReview(snapshot)
	assertResult(t, withReaction, StatusSuccess, false)
}

func TestReviewStateFixtures(t *testing.T) {
	t.Parallel()

	for _, fixture := range loadFixtureCases(t, "reviews") {
		fixture := fixture
		t.Run(fixture.Name, func(t *testing.T) {
			snapshot := reviewScenario(t, fixture.Scenario)
			result := EvaluateCodexReview(snapshot)
			assertResult(t, result, fixture.ExpectedState, fixture.RequestSecond)
			if !result.Publish || result.Context != CodexReviewContext || result.HeadSHA != testHead {
				t.Fatalf("result lacks current-head publication identity: %+v", result)
			}
		})
	}
}

func TestOverrideFixtures(t *testing.T) {
	t.Parallel()

	for _, fixture := range loadFixtureCases(t, "overrides") {
		fixture := fixture
		t.Run(fixture.Name, func(t *testing.T) {
			snapshot := baseSnapshot()
			switch fixture.Scenario {
			case "operator-review-waiver":
				snapshot.Comments = []Comment{{ID: 10, Author: operator, Body: "skip: codex-review - two connector attempts failed"}}
			case "other-owner-waiver":
				snapshot.Comments = []Comment{{ID: 11, Author: Actor{Login: "other-owner", ID: 2, Type: "User"}, Body: "skip: codex-review - owner association is insufficient"}}
			case "author-self-waiver":
				snapshot.Author = Actor{Login: "contributor", ID: 3, Type: "User"}
				snapshot.AuthorLogin = "contributor"
				snapshot.Comments = []Comment{{ID: 12, Author: snapshot.Author, Body: "skip: codex-review - self approved"}}
			case "missing-reason":
				snapshot.Comments = []Comment{{ID: 13, Author: operator, Body: "skip: codex-review -   "}}
			case "dependabot":
				snapshot.Author = Actor{Login: DependabotLogin, ID: DependabotUserID, Type: "Bot"}
				snapshot.AuthorLogin = DependabotLogin
			default:
				t.Fatalf("unknown fixture scenario %q", fixture.Scenario)
			}

			result := EvaluateCodexReview(snapshot)
			assertResult(t, result, fixture.ExpectedState, fixture.RequestSecond)
		})
	}
}

func TestRoundTwoPublicationFormat(t *testing.T) {
	t.Parallel()

	body := RoundTwoComment(19, testHead)
	if strings.Count(body, "@codex review") != 1 {
		t.Fatalf("request contains %d invocations, want one: %q", strings.Count(body, "@codex review"), body)
	}
	comment := Comment{ID: 20, Author: Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"}, Body: body}
	pr, head, ok := ParseRoundTwoMarker(comment)
	if !ok || pr != 19 || head != testHead {
		t.Fatalf("ParseRoundTwoMarker() = %d, %q, %v", pr, head, ok)
	}
	if got := RoundTwoReservationDescription(19, testHead); len(got) > 140 || !strings.Contains(got, testHead) {
		t.Fatalf("invalid reservation description %q", got)
	}
}

func TestUntrustedRoundTwoMarkerFailsClosed(t *testing.T) {
	t.Parallel()

	for _, scenario := range []string{"round-one-pending", "round-one-clean", "round-two-ready"} {
		t.Run(scenario, func(t *testing.T) {
			snapshot := reviewScenario(t, scenario)
			snapshot.Comments = append(snapshot.Comments, Comment{
				ID:        20,
				Author:    snapshot.Author,
				Body:      RoundTwoComment(snapshot.Number, snapshot.HeadSHA),
				CreatedAt: testStart,
			})
			result := EvaluateCodexReview(snapshot)
			if result.State != StatusError || result.RequestSecond {
				t.Fatalf("untrusted marker established a round boundary: %+v", result)
			}
		})
	}
}

func TestWhitespaceSurroundedOperatorInvocationConsumesRoundTwo(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	snapshot.Comments = []Comment{{ID: 21, Author: operator, Body: "context\n  @codex review  \n\nReviewed commit: `" + testHead[:10] + "`", CreatedAt: testStart}}
	result := EvaluateCodexReview(snapshot)
	assertResult(t, result, StatusPending, false)
}

func TestOperatorReviewRequestRequiresOneHistoricalHead(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		body string
	}{
		{name: "missing commit", body: "@codex review"},
		{name: "unknown commit", body: "@codex review\n\nReviewed commit: `abcdef0`"},
		{name: "multiple commits", body: "@codex review\n\nReviewed commits: `" + testHead[:10] + "` and `" + testOld[:10] + "`"},
	} {
		t.Run(test.name, func(t *testing.T) {
			snapshot := baseSnapshot()
			snapshot.Comments = []Comment{{ID: 21, Author: operator, Body: test.body, CreatedAt: testStart}}
			result := EvaluateCodexReview(snapshot)
			if result.State != StatusError || result.RequestSecond {
				t.Fatalf("unbound operator request did not fail closed: %+v", result)
			}
		})
	}

	snapshot := baseSnapshot()
	snapshot.HistoricalHeads = []string{testOld, testHead}
	snapshot.Comments = []Comment{{ID: 21, Author: operator, Body: "@codex review\n\nReviewed commit: `" + testOld[:10] + "`", CreatedAt: testStart}}
	result := EvaluateCodexReview(snapshot)
	if result.State != StatusFailure || !strings.Contains(result.Description, "stale") || result.RequestSecond {
		t.Fatalf("historical operator request was rebound to current head: %+v", result)
	}
}

func TestEditedOperatorRequestUsesEditTimeAsRoundBoundary(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	addCompletedSummary(&snapshot, testHead, testStart.Add(2*time.Minute), "+1")
	snapshot.Comments = append(snapshot.Comments, Comment{
		ID:        21,
		Author:    operator,
		Body:      "@codex review\n\nReviewed commit: `" + testHead[:10] + "`",
		CreatedAt: testStart,
		UpdatedAt: testStart.Add(3 * time.Minute),
	})
	reservation := EvaluateCodexReview(snapshot)
	if reservation.State != StatusPending || reservation.Description != RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA) {
		t.Fatalf("edited request did not reserve consumption: %+v", reservation)
	}
	snapshot.Checks = append(snapshot.Checks, CheckResult{
		Context:     CodexReviewContext,
		State:       string(StatusPending),
		Description: reservation.Description,
		SHA:         snapshot.HeadSHA,
		Creator:     Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"},
	})
	result := EvaluateCodexReview(snapshot)
	if result.State != StatusPending || result.RequestSecond {
		t.Fatalf("pre-edit terminal evidence completed round two: %+v", result)
	}
}

func TestDeletedOperatorRequestCannotRestoreReviewAllowance(t *testing.T) {
	t.Parallel()

	snapshot := baseSnapshot()
	snapshot.Comments = []Comment{{
		ID:        21,
		Author:    operator,
		Body:      "@codex review\n\nReviewed commit: `" + testHead[:10] + "`",
		CreatedAt: testStart,
	}}
	reservation := EvaluateCodexReview(snapshot)
	if reservation.State != StatusPending || reservation.Description != RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA) {
		t.Fatalf("operator request did not produce durable reservation: %+v", reservation)
	}
	snapshot.Checks = append(snapshot.Checks, CheckResult{
		Context:     CodexReviewContext,
		State:       string(StatusPending),
		Description: reservation.Description,
		SHA:         snapshot.HeadSHA,
		Creator:     Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"},
	})
	snapshot.Comments = nil
	result := EvaluateCodexReview(snapshot)
	if result.State != StatusFailure || result.RequestSecond {
		t.Fatalf("deleted request restored review allowance: %+v", result)
	}
}

func TestTenfoldReplayRequestsAtMostOneSecondRound(t *testing.T) {
	t.Parallel()

	snapshot := reviewScenario(t, "round-two-ready")
	requests := 0
	for replay := 0; replay < 10; replay++ {
		result := EvaluateCodexReview(snapshot)
		if result.RequestSecond {
			requests++
			snapshot.Comments = append(snapshot.Comments, Comment{
				ID:        100,
				Author:    Actor{Login: "github-actions[bot]", ID: 41898282, Type: "Bot"},
				Body:      RoundTwoComment(snapshot.Number, snapshot.HeadSHA),
				CreatedAt: testStart.Add(3 * time.Minute),
			})
		}
	}
	if requests != 1 {
		t.Fatalf("ten replays requested round two %d times, want exactly one", requests)
	}
}

func TestReservationInterleavingCannotIssueDuplicateRequest(t *testing.T) {
	t.Parallel()

	snapshot := reviewScenario(t, "round-two-ready")
	first := EvaluateCodexReview(snapshot)
	if !first.RequestSecond {
		t.Fatalf("initial delivery did not request round two: %+v", first)
	}
	snapshot.Checks = append(snapshot.Checks, CheckResult{
		Context:     CodexReviewContext,
		State:       string(StatusPending),
		Description: first.Description,
		SHA:         snapshot.HeadSHA,
		Creator:     Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"},
	})
	interleaved := EvaluateCodexReview(snapshot)
	if interleaved.RequestSecond || interleaved.State != StatusFailure {
		t.Fatalf("delivery interleaved after reservation could request again: %+v", interleaved)
	}
	addRoundTwoRequest(&snapshot)
	afterRequest := EvaluateCodexReview(snapshot)
	if afterRequest.RequestSecond || afterRequest.State != StatusPending {
		t.Fatalf("delivery after marked request could request again: %+v", afterRequest)
	}
}

func TestNoStateAfterRoundTwoRequestsAnotherReview(t *testing.T) {
	t.Parallel()

	for _, scenario := range []string{"round-two-pending", "round-two-clean", "round-two-unresolved", "round-two-resolved", "round-two-failed", "round-two-stale", "protocol-violation"} {
		t.Run(scenario, func(t *testing.T) {
			snapshot := reviewScenario(t, scenario)
			for replay := 0; replay < 10; replay++ {
				if result := EvaluateCodexReview(snapshot); result.RequestSecond {
					t.Fatalf("replay %d requested a forbidden later review: %+v", replay, result)
				}
			}
		})
	}
}

func TestRemediationBridgeFailsClosedWithoutAncestryAndGreenCI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Snapshot)
		state  StatusState
	}{
		{name: "missing ancestry", mutate: func(snapshot *Snapshot) { snapshot.Comparisons = nil }, state: StatusFailure},
		{name: "diverged ancestry", mutate: func(snapshot *Snapshot) { snapshot.Comparisons[0].Status = "diverged" }, state: StatusFailure},
		{name: "same reviewed head", mutate: func(snapshot *Snapshot) {
			snapshot.Reviews[0].CommitSHA = snapshot.HeadSHA
			snapshot.ReviewThreads[0].CommitSHA = snapshot.HeadSHA
		}, state: StatusFailure},
		{name: "missing CI", mutate: func(snapshot *Snapshot) { snapshot.Checks = snapshot.Checks[1:] }, state: StatusPending},
		{name: "failed CI", mutate: func(snapshot *Snapshot) { snapshot.Checks[0].State = "failure" }, state: StatusPending},
		{name: "stale CI", mutate: func(snapshot *Snapshot) { snapshot.Checks[0].SHA = testOld }, state: StatusPending},
		{name: "ambiguous equal-time CI", mutate: func(snapshot *Snapshot) {
			check := snapshot.Checks[0]
			check.State = "failure"
			snapshot.Checks = append(snapshot.Checks, check)
		}, state: StatusPending},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			snapshot := reviewScenario(t, "round-two-ready")
			test.mutate(&snapshot)
			result := EvaluateCodexReview(snapshot)
			assertResult(t, result, test.state, false)
		})
	}
}

func TestReservationRequiresAuthenticPendingStatus(t *testing.T) {
	t.Parallel()

	for _, mutate := range []func(*CheckResult){
		func(check *CheckResult) { check.State = "success" },
		func(check *CheckResult) { check.Creator.ID = 1 },
		func(check *CheckResult) { check.Creator.Login = "lookalike[bot]" },
		func(check *CheckResult) { check.Creator.Type = "User" },
	} {
		snapshot := reviewScenario(t, "round-two-reserved")
		mutate(&snapshot.Checks[len(snapshot.Checks)-1])
		result := EvaluateCodexReview(snapshot)
		if result.State != StatusError || result.RequestSecond {
			t.Fatalf("inauthentic reservation did not fail closed: %+v", result)
		}
	}
}

func TestReservationAndMarkerMustTargetSameHead(t *testing.T) {
	t.Parallel()

	snapshot := reviewScenario(t, "round-two-pending")
	snapshot.Checks = append(snapshot.Checks, CheckResult{
		Context:     CodexReviewContext,
		State:       "pending",
		Description: RoundTwoReservationDescription(snapshot.Number, testOld),
		SHA:         testOld,
		Creator:     Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"},
	})
	result := EvaluateCodexReview(snapshot)
	if result.State != StatusError || result.RequestSecond {
		t.Fatalf("mismatched reservation and marker did not fail closed: %+v", result)
	}
}

func TestEvaluationFailsClosedOutsideManagedBoundary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Snapshot)
	}{
		{name: "wrong base", mutate: func(snapshot *Snapshot) { snapshot.BaseRef = "release" }},
		{name: "closed", mutate: func(snapshot *Snapshot) { snapshot.State = "closed" }},
		{name: "incomplete", mutate: func(snapshot *Snapshot) { snapshot.Complete = false }},
		{name: "collector error", mutate: func(snapshot *Snapshot) { snapshot.Errors = []string{"comments page incomplete"} }},
		{name: "invalid head", mutate: func(snapshot *Snapshot) { snapshot.HeadSHA = "abc" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			snapshot := baseSnapshot()
			test.mutate(&snapshot)
			issue, review := Evaluate(snapshot)
			for _, result := range []PolicyResult{issue, review} {
				if result.Publish || result.RequestSecond || result.State != StatusError {
					t.Fatalf("evaluation did not fail closed: %+v", result)
				}
			}
		})
	}
}

func reviewScenario(t *testing.T, scenario string) Snapshot {
	t.Helper()

	snapshot := baseSnapshot()
	switch scenario {
	case "draft":
		snapshot.Draft = true
	case "round-one-pending":
	case "round-one-eyes":
		snapshot.Reactions = []Reaction{{Author: codex, Content: "eyes", Target: "pull_request", CreatedAt: testStart.Add(time.Minute)}}
	case "round-one-clean":
		addCompletedSummary(&snapshot, testHead, testStart.Add(time.Minute), "+1")
	case "round-one-stale-clean":
		addCompletedSummary(&snapshot, testOld, testStart.Add(time.Minute), "+1")
	case "round-one-failed":
		addSummary(&snapshot, SummaryFailed, testHead, testStart.Add(time.Minute))
	case "round-one-unresolved":
		addFirstFinding(&snapshot, testOld, false)
	case "round-two-ready":
		addFirstFinding(&snapshot, testOld, true)
		addGreenBridge(&snapshot)
	case "round-two-reserved":
		addFirstFinding(&snapshot, testOld, true)
		addGreenBridge(&snapshot)
		snapshot.Checks = append(snapshot.Checks, CheckResult{Context: CodexReviewContext, State: "pending", Description: RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA), SHA: snapshot.HeadSHA, Creator: Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"}})
	case "round-two-pending":
		addRoundTwoRequest(&snapshot)
	case "round-two-clean":
		addRoundTwoRequest(&snapshot)
		addCompletedSummary(&snapshot, testHead, testStart.Add(4*time.Minute), "+1")
	case "round-two-unresolved":
		addRoundTwoRequest(&snapshot)
		addSecondFinding(&snapshot, false)
	case "round-two-resolved":
		addRoundTwoRequest(&snapshot)
		addSecondFinding(&snapshot, true)
	case "round-two-failed":
		addRoundTwoRequest(&snapshot)
		addSummary(&snapshot, SummaryFailed, testHead, testStart.Add(4*time.Minute))
	case "round-two-stale":
		addRoundTwoRequest(&snapshot)
		snapshot.Comments[len(snapshot.Comments)-1].Body = RoundTwoComment(snapshot.Number, testOld)
		snapshot.Checks[len(snapshot.Checks)-1].Description = RoundTwoReservationDescription(snapshot.Number, testOld)
		snapshot.Checks[len(snapshot.Checks)-1].SHA = testOld
	case "protocol-violation":
		addRoundTwoRequest(&snapshot)
		snapshot.Comments = append(snapshot.Comments, Comment{ID: 31, Author: operator, Body: "@codex review\n\nReviewed commit: `" + testHead[:10] + "`", CreatedAt: testStart.Add(4 * time.Minute)})
	default:
		t.Fatalf("unknown review scenario %q", scenario)
	}
	return snapshot
}

func baseSnapshot() Snapshot {
	return Snapshot{
		Complete:        true,
		Number:          19,
		HeadSHA:         testHead,
		BaseRef:         "main",
		AuthorLogin:     "contributor",
		Author:          Actor{Login: "contributor", ID: 7, Type: "User"},
		State:           "open",
		HistoricalHeads: []string{testHead},
	}
}

func addSummary(snapshot *Snapshot, status SummaryStatus, commit string, at time.Time) {
	body := fmt.Sprintf("<!-- codex-pull-request-review-summary -->\n\n| Review | Status | Commit | Review trigger |\n| --- | --- | --- | --- |\n| Code Review | %s | `%s` | Automatic |\n", summaryLabel(status), commit[:7])
	snapshot.Comments = append(snapshot.Comments, Comment{ID: int64(100 + len(snapshot.Comments)), Author: codex, Body: body, CreatedAt: at, UpdatedAt: at})
}

func addCompletedSummary(snapshot *Snapshot, commit string, at time.Time, reaction string) {
	addSummary(snapshot, SummaryCompleted, commit, at)
	commentID := snapshot.Comments[len(snapshot.Comments)-1].ID
	snapshot.Reactions = append(snapshot.Reactions, Reaction{Author: codex, Content: reaction, Target: fmt.Sprintf("comment:%d", commentID), CreatedAt: at.Add(time.Second)})
}

func summaryLabel(status SummaryStatus) string {
	switch status {
	case SummaryCompleted:
		return "✅ **Completed**"
	case SummaryFailed:
		return "⚠️ **Failed**"
	case SummaryPending:
		return "👀 **In progress**"
	default:
		return "❓ **Unknown**"
	}
}

func addFirstFinding(snapshot *Snapshot, commit string, resolved bool) {
	reviewID := "round-one-review"
	snapshot.Reviews = append(snapshot.Reviews, Review{ID: reviewID, Author: codex, CommitSHA: commit, SubmittedAt: testStart.Add(time.Minute)})
	snapshot.ReviewThreads = append(snapshot.ReviewThreads, ReviewThread{ID: "round-one-thread", ReviewID: reviewID, CommitSHA: commit, Author: codex, CreatedAt: testStart.Add(time.Minute), Resolved: resolved, Outdated: true})
}

func addGreenBridge(snapshot *Snapshot) {
	snapshot.Comparisons = append(snapshot.Comparisons, Comparison{BaseSHA: testOld, HeadSHA: testHead, Status: "ahead", AheadBy: 1})
	for _, context := range RequiredS006Checks {
		snapshot.Checks = append(snapshot.Checks, CheckResult{Context: context, State: "success", SHA: testHead})
	}
}

func addRoundTwoRequest(snapshot *Snapshot) {
	snapshot.Comments = append(snapshot.Comments, Comment{ID: 30, Author: Actor{Login: "github-actions[bot]", ID: 41898282, Type: "Bot"}, Body: RoundTwoComment(snapshot.Number, snapshot.HeadSHA), CreatedAt: testStart.Add(3 * time.Minute), UpdatedAt: testStart.Add(3 * time.Minute)})
	for _, check := range snapshot.Checks {
		if check.Context == CodexReviewContext && check.Description == RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA) {
			return
		}
	}
	snapshot.Checks = append(snapshot.Checks, CheckResult{Context: CodexReviewContext, State: string(StatusPending), Description: RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA), SHA: snapshot.HeadSHA, Creator: Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"}})
}

func addSecondFinding(snapshot *Snapshot, resolved bool) {
	reviewID := "round-two-review"
	snapshot.Reviews = append(snapshot.Reviews, Review{ID: reviewID, Author: codex, CommitSHA: testHead, SubmittedAt: testStart.Add(4 * time.Minute)})
	snapshot.ReviewThreads = append(snapshot.ReviewThreads, ReviewThread{ID: "round-two-thread", ReviewID: reviewID, CommitSHA: testHead, Author: codex, CreatedAt: testStart.Add(4 * time.Minute), Resolved: resolved})
}

func assertResult(t *testing.T, result PolicyResult, state StatusState, requestSecond bool) {
	t.Helper()
	if result.State != state || result.RequestSecond != requestSecond {
		t.Fatalf("result = state %q request_second %v (%s); want %q, %v; evidence=%v", result.State, result.RequestSecond, result.Description, state, requestSecond, result.Evidence)
	}
}

func loadFixtureCases(t *testing.T, group string) []fixtureCase {
	t.Helper()
	path := filepath.Join("testdata", group, "cases.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var cases []fixtureCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return cases
}
