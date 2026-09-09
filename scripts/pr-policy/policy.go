package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	IssueLinkContext   = "Cueson PR policy / Issue link"
	CodexReviewContext = "Cueson PR policy / Codex review"

	OperatorLogin = "h8rt3rmin8r"

	CodexRESTLogin = "chatgpt-codex-connector[bot]"
	CodexUserID    = int64(199175422)

	DependabotLogin     = "dependabot[bot]"
	DependabotUserID    = int64(49699333)
	GitHubActionsLogin  = "github-actions[bot]"
	GitHubActionsUserID = int64(41898282)

	RoundTwoMarkerPrefix      = "<!-- cueson-codex-review-request:v1"
	roundTwoReservationPrefix = "cueson:v1 round=2 reserved"
)

var (
	commitSHAExpression = regexp.MustCompile(`^[0-9a-f]{40}$`)
	commitPrefixPattern = regexp.MustCompile("`([0-9a-f]{7,40})`")
	bareSHAExpression   = regexp.MustCompile(`(?:^|[^0-9a-f])([0-9a-f]{40})(?:$|[^0-9a-f])`)
	markerPattern       = regexp.MustCompile(`^<!-- cueson-codex-review-request:v1 round=2 pr=([1-9][0-9]*) head=([0-9a-f]{40}) -->$`)
	reservationPattern  = regexp.MustCompile(`^cueson:v1 round=2 reserved pr=([1-9][0-9]*) head=([0-9a-f]{40})$`)

	// RequiredS006Checks is the closed CI workflow set that must be green on a
	// remediation head before automation spends the second review round.
	RequiredS006Checks = []string{
		"Formatting",
		"Repository text",
		"Vet",
		"Schema and conformance",
		"Static analysis",
		"Vulnerability scan",
		"Native tests (Linux)",
		"Native tests (Windows)",
		"Native tests (macOS)",
		"Race detection",
		"Pure-Go build (linux-amd64)",
		"Pure-Go build (linux-arm64)",
		"Pure-Go build (windows-amd64)",
		"Pure-Go build (windows-arm64)",
		"Pure-Go build (darwin-amd64)",
		"Pure-Go build (darwin-arm64)",
	}
)

type StatusState string

const (
	StatusPending StatusState = "pending"
	StatusSuccess StatusState = "success"
	StatusFailure StatusState = "failure"
	StatusError   StatusState = "error"
)

type SummaryStatus string

const (
	SummaryPending   SummaryStatus = "pending"
	SummaryCompleted SummaryStatus = "completed"
	SummaryFailed    SummaryStatus = "failed"
)

type Actor struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Type  string `json:"type"`
}

type ClosingIssue struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Number     int    `json:"number"`
}

type Comment struct {
	ID        int64     `json:"id"`
	Author    Actor     `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Reaction struct {
	Author    Actor     `json:"author"`
	Content   string    `json:"content"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID                string    `json:"id"`
	Author            Actor     `json:"author"`
	CommitSHA         string    `json:"commit_sha"`
	Body              string    `json:"body"`
	SubmittedAt       time.Time `json:"submitted_at"`
	UnthreadedFinding bool      `json:"unthreaded_finding"`
}

type ReviewThread struct {
	ID        string    `json:"id"`
	ReviewID  string    `json:"review_id"`
	CommitSHA string    `json:"commit_sha"`
	Author    Actor     `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Resolved  bool      `json:"resolved"`
	Outdated  bool      `json:"outdated"`
}

type CodexSummary struct {
	CommentID    int64         `json:"comment_id"`
	CommitPrefix string        `json:"commit_prefix"`
	Status       SummaryStatus `json:"status"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Comparison struct {
	BaseSHA string `json:"base_sha"`
	HeadSHA string `json:"head_sha"`
	Status  string `json:"status"`
	AheadBy int    `json:"ahead_by"`
}

type CheckResult struct {
	Context     string    `json:"context"`
	State       string    `json:"state"`
	Description string    `json:"description"`
	SHA         string    `json:"sha"`
	Creator     Actor     `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
}

type Snapshot struct {
	Complete      bool           `json:"complete"`
	Errors        []string       `json:"errors,omitempty"`
	Number        int            `json:"number"`
	HeadSHA       string         `json:"head_sha"`
	BaseRef       string         `json:"base_ref"`
	AuthorLogin   string         `json:"author_login"`
	Author        Actor          `json:"author"`
	Draft         bool           `json:"draft"`
	State         string         `json:"state"`
	ClosingIssues []ClosingIssue `json:"closing_issues"`
	Comments      []Comment      `json:"comments"`
	Reactions     []Reaction     `json:"reactions"`
	Reviews       []Review       `json:"reviews"`
	ReviewThreads []ReviewThread `json:"review_threads"`
	Comparisons   []Comparison   `json:"comparisons"`
	Checks        []CheckResult  `json:"checks"`
}

type PolicyResult struct {
	Context       string      `json:"context"`
	State         StatusState `json:"state"`
	Description   string      `json:"description"`
	HeadSHA       string      `json:"head_sha"`
	Publish       bool        `json:"publish"`
	RequestSecond bool        `json:"request_second"`
	Evidence      []string    `json:"evidence,omitempty"`
}

func Evaluate(snapshot Snapshot) (PolicyResult, PolicyResult) {
	if err := validateManagedSnapshot(snapshot); err != nil {
		return closedResult(IssueLinkContext, snapshot.HeadSHA, err), closedResult(CodexReviewContext, snapshot.HeadSHA, err)
	}
	return EvaluateIssueLink(snapshot), EvaluateCodexReview(snapshot)
}

func EvaluateIssueLink(snapshot Snapshot) PolicyResult {
	if err := validateManagedSnapshot(snapshot); err != nil {
		return closedResult(IssueLinkContext, snapshot.HeadSHA, err)
	}
	if IsDependabotActor(snapshot.Author) {
		return result(IssueLinkContext, snapshot.HeadSHA, StatusSuccess, "Dependabot is exempt from issue linkage", "dependabot identity verified")
	}

	if len(snapshot.ClosingIssues) != 0 {
		targets := make([]string, 0, len(snapshot.ClosingIssues))
		for _, issue := range snapshot.ClosingIssues {
			if issue.Owner == "" || issue.Repository == "" || issue.Number <= 0 {
				return errorResult(IssueLinkContext, snapshot.HeadSHA, "GitHub returned an invalid closing issue")
			}
			targets = append(targets, fmt.Sprintf("%s/%s#%d", issue.Owner, issue.Repository, issue.Number))
		}
		sort.Strings(targets)
		description := "Closes " + strings.Join(targets, ", ")
		if len(description) > 140 {
			description = fmt.Sprintf("Closes %d GitHub-resolved issues", len(targets))
		}
		return result(IssueLinkContext, snapshot.HeadSHA, StatusSuccess, description, targets...)
	}

	if comment, reason, ok := findOperatorDirective(snapshot.Comments, "issue-link"); ok {
		if comment.ID <= 0 {
			return errorResult(IssueLinkContext, snapshot.HeadSHA, "Issue-link waiver has an invalid comment identifier")
		}
		return result(IssueLinkContext, snapshot.HeadSHA, StatusSuccess, fmt.Sprintf("Issue-link waiver in comment %d", comment.ID), reason)
	}
	return result(IssueLinkContext, snapshot.HeadSHA, StatusFailure, "Missing GitHub-resolved closing issue", "add a complete closing reference or operator waiver")
}

func EvaluateCodexReview(snapshot Snapshot) PolicyResult {
	if err := validateManagedSnapshot(snapshot); err != nil {
		return closedResult(CodexReviewContext, snapshot.HeadSHA, err)
	}
	if IsDependabotActor(snapshot.Author) {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusSuccess, "Dependabot is excluded from Codex review", "dependabot identity verified")
	}
	if comment, reason, ok := findOperatorDirective(snapshot.Comments, "codex-review"); ok {
		if comment.ID <= 0 {
			return errorResult(CodexReviewContext, snapshot.HeadSHA, "Codex-review waiver has an invalid comment identifier")
		}
		return result(CodexReviewContext, snapshot.HeadSHA, StatusSuccess, fmt.Sprintf("Codex review waived in comment %d", comment.ID), reason)
	}
	if snapshot.Draft {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Draft pull request is not reviewable", "native round one has not started")
	}

	evidence, err := classifyReviewEvidence(snapshot)
	if err != nil {
		return errorResult(CodexReviewContext, snapshot.HeadSHA, err.Error())
	}
	if evidence.requestCount > 1 {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Multiple second-round requests violate policy", evidence.requestEvidence...)
	}
	if evidence.requestCount == 1 || evidence.reservation != nil {
		return evaluateRoundTwo(snapshot, evidence)
	}
	return evaluateRoundOne(snapshot, evidence)
}

func IsCodexActor(actor Actor) bool {
	login := actor.Login
	if login == CodexRESTLogin {
		login = strings.TrimSuffix(login, "[bot]")
	}
	return login == strings.TrimSuffix(CodexRESTLogin, "[bot]") && actor.ID == CodexUserID && actor.Type == "Bot"
}

func IsDependabotActor(actor Actor) bool {
	return actor.Login == DependabotLogin && actor.ID == DependabotUserID && actor.Type == "Bot"
}

func IsGitHubActionsActor(actor Actor) bool {
	return actor.Login == GitHubActionsLogin && actor.ID == GitHubActionsUserID && actor.Type == "Bot"
}

func ParseDirective(comment Comment) (string, string, bool) {
	if comment.Author.Login != OperatorLogin || comment.Author.Type != "User" {
		return "", "", false
	}
	for _, line := range logicalLines(comment.Body) {
		for _, policy := range []string{"issue-link", "codex-review"} {
			prefix := "skip: " + policy + " - "
			if !strings.HasPrefix(line, prefix) {
				continue
			}
			reason := strings.TrimSpace(strings.TrimPrefix(line, prefix))
			if reason != "" && line == prefix+reason {
				return policy, reason, true
			}
		}
	}
	return "", "", false
}

func ParseRoundTwoMarker(comment Comment) (int, string, bool) {
	if !IsGitHubActionsActor(comment.Author) {
		return 0, "", false
	}
	var number int
	var head string
	count := 0
	for _, line := range logicalLines(comment.Body) {
		matches := markerPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		parsed, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, "", false
		}
		number, head = parsed, matches[2]
		count++
	}
	return number, head, count == 1
}

func RoundTwoComment(number int, headSHA string) string {
	return fmt.Sprintf("@codex review\n\n<!-- cueson-codex-review-request:v1 round=2 pr=%d head=%s -->", number, headSHA)
}

func RoundTwoReservationDescription(number int, headSHA string) string {
	return fmt.Sprintf("%s pr=%d head=%s", roundTwoReservationPrefix, number, headSHA)
}

func ParseCodexSummary(comment Comment) (CodexSummary, bool, error) {
	if !strings.Contains(comment.Body, "<!-- codex-pull-request-review-summary -->") {
		return CodexSummary{}, false, nil
	}
	if !IsCodexActor(comment.Author) {
		if canonicalCodexLogin(comment.Author.Login) {
			return CodexSummary{}, true, errors.New("codex summary actor identity did not match configured bot")
		}
		return CodexSummary{}, false, nil
	}

	status := SummaryStatus("")
	for candidate, marker := range map[SummaryStatus]string{
		SummaryCompleted: "**Completed**",
		SummaryFailed:    "**Failed**",
		SummaryPending:   "**In progress**",
	} {
		if strings.Contains(comment.Body, marker) {
			if status != "" {
				return CodexSummary{}, true, errors.New("codex summary contains contradictory states")
			}
			status = candidate
		}
	}
	matches := commitPrefixPattern.FindAllStringSubmatch(comment.Body, -1)
	if status == "" || len(matches) != 1 {
		return CodexSummary{}, true, errors.New("codex summary has an unsupported state or commit shape")
	}
	return CodexSummary{CommentID: comment.ID, CommitPrefix: matches[0][1], Status: status, UpdatedAt: latestTime(comment.UpdatedAt, comment.CreatedAt)}, true, nil
}

type reviewEvidence struct {
	summary         *CodexSummary
	clean           []terminalEvidence
	failed          []terminalEvidence
	threads         []ReviewThread
	unthreaded      []Review
	requestCount    int
	requestHead     string
	requestAt       time.Time
	requestEvidence []string
	reservation     *CheckResult
}

type terminalEvidence struct {
	commit string
	at     time.Time
	label  string
}

func classifyReviewEvidence(snapshot Snapshot) (reviewEvidence, error) {
	var evidence reviewEvidence
	seenComments := make(map[int64]struct{}, len(snapshot.Comments))
	for _, comment := range snapshot.Comments {
		if comment.ID <= 0 {
			return reviewEvidence{}, errors.New("comment evidence has an invalid identifier")
		}
		if _, duplicate := seenComments[comment.ID]; duplicate {
			return reviewEvidence{}, fmt.Errorf("duplicate comment %d", comment.ID)
		}
		seenComments[comment.ID] = struct{}{}

		summary, recognized, err := ParseCodexSummary(comment)
		if err != nil {
			return reviewEvidence{}, err
		}
		if recognized {
			if evidence.summary != nil {
				return reviewEvidence{}, errors.New("multiple Codex summary comments are ambiguous")
			}
			evidence.summary = &summary
		}

		markerNumber, markerHead, marker := ParseRoundTwoMarker(comment)
		markerOccurrences := strings.Count(comment.Body, RoundTwoMarkerPrefix)
		if markerOccurrences != 0 {
			if !marker || markerOccurrences != 1 || markerNumber != snapshot.Number {
				return reviewEvidence{}, fmt.Errorf("malformed or mismatched round-two marker in comment %d", comment.ID)
			}
			if standaloneInvocationCount(comment.Body) != 1 {
				return reviewEvidence{}, fmt.Errorf("round-two marker comment %d must contain exactly one standalone invocation", comment.ID)
			}
			evidence.requestCount++
			evidence.requestHead = markerHead
			evidence.requestAt = comment.CreatedAt
			evidence.requestEvidence = append(evidence.requestEvidence, fmt.Sprintf("marked request comment %d", comment.ID))
		} else if comment.Author.Login == OperatorLogin && comment.Author.Type == "User" && standaloneInvocationCount(comment.Body) != 0 {
			evidence.requestCount += standaloneInvocationCount(comment.Body)
			evidence.requestAt = comment.CreatedAt
			evidence.requestEvidence = append(evidence.requestEvidence, fmt.Sprintf("operator request comment %d", comment.ID))
		}

		terminal, kind, ok, err := parseTerminalComment(comment)
		if err != nil {
			return reviewEvidence{}, err
		}
		if ok {
			switch kind {
			case SummaryCompleted:
				evidence.clean = append(evidence.clean, terminal)
			case SummaryFailed:
				evidence.failed = append(evidence.failed, terminal)
			}
		}
	}

	seenReviews := make(map[string]Review, len(snapshot.Reviews))
	for _, review := range snapshot.Reviews {
		if !IsCodexActor(review.Author) {
			continue
		}
		if review.ID == "" || !commitSHAExpression.MatchString(review.CommitSHA) {
			return reviewEvidence{}, errors.New("codex review has an invalid identity or commit")
		}
		if _, duplicate := seenReviews[review.ID]; duplicate {
			return reviewEvidence{}, fmt.Errorf("duplicate Codex review %q", review.ID)
		}
		seenReviews[review.ID] = review
		if review.UnthreadedFinding {
			evidence.unthreaded = append(evidence.unthreaded, review)
		}
	}

	seenThreads := make(map[string]struct{}, len(snapshot.ReviewThreads))
	for _, thread := range snapshot.ReviewThreads {
		if !IsCodexActor(thread.Author) {
			continue
		}
		if thread.ID == "" || thread.ReviewID == "" || !commitSHAExpression.MatchString(thread.CommitSHA) {
			return reviewEvidence{}, errors.New("codex review thread has invalid attribution")
		}
		if _, duplicate := seenThreads[thread.ID]; duplicate {
			return reviewEvidence{}, fmt.Errorf("duplicate Codex thread %q", thread.ID)
		}
		seenThreads[thread.ID] = struct{}{}
		review, exists := seenReviews[thread.ReviewID]
		if !exists || review.CommitSHA != thread.CommitSHA {
			return reviewEvidence{}, fmt.Errorf("codex thread %q does not match its originating review", thread.ID)
		}
		evidence.threads = append(evidence.threads, thread)
	}

	for index := range snapshot.Checks {
		check := &snapshot.Checks[index]
		if check.Context != CodexReviewContext || !strings.HasPrefix(check.Description, roundTwoReservationPrefix) {
			continue
		}
		number, head, ok := parseReservation(check.Description)
		if !ok || number != snapshot.Number || check.SHA != head || check.State != string(StatusPending) || !IsGitHubActionsActor(check.Creator) {
			return reviewEvidence{}, errors.New("round-two reservation status is malformed")
		}
		if evidence.reservation != nil {
			return reviewEvidence{}, errors.New("multiple round-two reservations are ambiguous")
		}
		evidence.reservation = check
		if evidence.requestHead != "" && evidence.requestHead != head {
			return reviewEvidence{}, errors.New("round-two reservation and request target different heads")
		}
		if evidence.requestHead == "" {
			evidence.requestHead = head
		}
	}

	if evidence.requestCount == 1 && evidence.requestHead == "" {
		evidence.requestHead = snapshot.HeadSHA
	}
	return evidence, nil
}

func evaluateRoundOne(snapshot Snapshot, evidence reviewEvidence) PolicyResult {
	for _, thread := range evidence.threads {
		if !thread.Resolved {
			return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-one Codex finding remains unresolved", thread.ID)
		}
	}
	if len(evidence.unthreaded) != 0 {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Codex finding has no resolvable thread", evidence.unthreaded[0].ID)
	}
	if len(evidence.threads) != 0 {
		commits := uniqueThreadCommits(evidence.threads)
		for _, commit := range commits {
			if commit == snapshot.HeadSHA || !isProvenDescendant(snapshot, commit) {
				return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Remediation head is not a proven newer descendant", commit)
			}
		}
		if missing := missingSuccessfulChecks(snapshot.Checks, snapshot.HeadSHA); len(missing) != 0 {
			return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Remediation waits for required CI", missing...)
		}
		description := RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA)
		answer := result(CodexReviewContext, snapshot.HeadSHA, StatusPending, description, "all round-one findings resolved", "remediation ancestry proven", "required CI green")
		answer.RequestSecond = true
		return answer
	}

	terminal, err := currentRoundTerminal(snapshot, evidence, time.Time{})
	if err != nil {
		return errorResult(CodexReviewContext, snapshot.HeadSHA, err.Error())
	}
	if terminal.failed {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Current-head Codex review failed", terminal.label)
	}
	if terminal.clean {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusSuccess, "Current-head Codex review is clean", terminal.label)
	}
	if terminal.stale {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Codex result is stale for the current head", terminal.label)
	}
	if hasCodexReaction(snapshot.Reactions, "eyes", "", time.Time{}) {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Codex acknowledged round one", "eyes is acknowledgement only")
	}
	return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Waiting for native Codex round one", "automation emits no first-round request")
}

func evaluateRoundTwo(snapshot Snapshot, evidence reviewEvidence) PolicyResult {
	requestHead := evidence.requestHead
	if requestHead == "" && evidence.reservation != nil {
		_, requestHead, _ = parseReservation(evidence.reservation.Description)
	}
	if requestHead != snapshot.HeadSHA {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-two evidence is stale after a head change", requestHead)
	}
	if evidence.reservation != nil && evidence.requestCount == 0 {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-two reservation has ambiguous publication", evidence.reservation.Description)
	}

	var roundTwoThreads []ReviewThread
	for _, thread := range evidence.threads {
		if evidence.requestAt.IsZero() || !thread.CreatedAt.Before(evidence.requestAt) {
			roundTwoThreads = append(roundTwoThreads, thread)
		}
		if !thread.Resolved {
			return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Codex finding remains unresolved after round two", thread.ID)
		}
	}
	for _, review := range evidence.unthreaded {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Codex finding has no resolvable thread", review.ID)
	}
	if len(roundTwoThreads) != 0 {
		for _, thread := range roundTwoThreads {
			if thread.CommitSHA != snapshot.HeadSHA {
				return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-two finding belongs to a stale head", thread.CommitSHA)
			}
		}
		return result(CodexReviewContext, snapshot.HeadSHA, StatusSuccess, "Round-two findings are resolved", "no third review is permitted")
	}

	terminal, err := currentRoundTerminal(snapshot, evidence, evidence.requestAt)
	if err != nil {
		return errorResult(CodexReviewContext, snapshot.HeadSHA, err.Error())
	}
	if terminal.failed {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-two Codex review failed", terminal.label, "no third review is permitted")
	}
	if terminal.clean {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusSuccess, "Round-two Codex review is clean", terminal.label, "no third review is permitted")
	}
	if terminal.stale {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusFailure, "Round-two result is stale", terminal.label, "no third review is permitted")
	}
	if hasCodexReaction(snapshot.Reactions, "eyes", "", evidence.requestAt) {
		return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Codex acknowledged round two", "eyes is acknowledgement only")
	}
	return result(CodexReviewContext, snapshot.HeadSHA, StatusPending, "Waiting for the single second Codex review", evidence.requestEvidence...)
}

type terminalState struct {
	clean  bool
	failed bool
	stale  bool
	label  string
}

func currentRoundTerminal(snapshot Snapshot, evidence reviewEvidence, after time.Time) (terminalState, error) {
	var clean, failed, stale []string
	if evidence.summary != nil && !evidence.summary.UpdatedAt.Before(after) {
		matchesHead := strings.HasPrefix(snapshot.HeadSHA, evidence.summary.CommitPrefix)
		switch evidence.summary.Status {
		case SummaryCompleted:
			if matchesHead && hasCodexReaction(snapshot.Reactions, "+1", fmt.Sprintf("comment:%d", evidence.summary.CommentID), latestTime(after, evidence.summary.UpdatedAt)) {
				clean = append(clean, fmt.Sprintf("summary comment %d", evidence.summary.CommentID))
			} else if !matchesHead {
				stale = append(stale, fmt.Sprintf("summary commit %s", evidence.summary.CommitPrefix))
			}
		case SummaryFailed:
			if matchesHead {
				failed = append(failed, fmt.Sprintf("summary comment %d", evidence.summary.CommentID))
			} else {
				stale = append(stale, fmt.Sprintf("summary commit %s", evidence.summary.CommitPrefix))
			}
		}
	}
	for _, item := range evidence.clean {
		if item.at.Before(after) {
			continue
		}
		if strings.HasPrefix(snapshot.HeadSHA, item.commit) && hasCodexReaction(snapshot.Reactions, "+1", "pull_request", latestTime(after, item.at)) {
			clean = append(clean, item.label)
		} else if !strings.HasPrefix(snapshot.HeadSHA, item.commit) {
			stale = append(stale, item.label)
		}
	}
	for _, item := range evidence.failed {
		if item.at.Before(after) {
			continue
		}
		if strings.HasPrefix(snapshot.HeadSHA, item.commit) {
			failed = append(failed, item.label)
		} else {
			stale = append(stale, item.label)
		}
	}
	if len(clean) != 0 && len(failed) != 0 {
		return terminalState{}, errors.New("codex supplied contradictory current-head terminal outcomes")
	}
	if len(clean) != 0 {
		return terminalState{clean: true, label: clean[len(clean)-1]}, nil
	}
	if len(failed) != 0 {
		return terminalState{failed: true, label: failed[len(failed)-1]}, nil
	}
	if len(stale) != 0 {
		return terminalState{stale: true, label: stale[len(stale)-1]}, nil
	}
	return terminalState{}, nil
}

func parseTerminalComment(comment Comment) (terminalEvidence, SummaryStatus, bool, error) {
	if !IsCodexActor(comment.Author) {
		return terminalEvidence{}, "", false, nil
	}
	kind := SummaryStatus("")
	switch {
	case strings.HasPrefix(comment.Body, "Codex Review: Didn't find any major issues."):
		kind = SummaryCompleted
	case strings.HasPrefix(comment.Body, "Codex Review: Something went wrong."):
		kind = SummaryFailed
	default:
		return terminalEvidence{}, "", false, nil
	}
	matches := commitPrefixPattern.FindAllStringSubmatch(comment.Body, -1)
	if kind == SummaryFailed && len(matches) == 0 {
		matches = bareSHAExpression.FindAllStringSubmatch(comment.Body, -1)
	}
	if len(matches) != 1 {
		return terminalEvidence{}, "", true, fmt.Errorf("codex terminal comment %d lacks one reviewed commit", comment.ID)
	}
	return terminalEvidence{commit: matches[0][1], at: latestTime(comment.UpdatedAt, comment.CreatedAt), label: fmt.Sprintf("Codex comment %d", comment.ID)}, kind, true, nil
}

func findOperatorDirective(comments []Comment, wanted string) (Comment, string, bool) {
	for _, comment := range comments {
		policy, reason, ok := ParseDirective(comment)
		if ok && policy == wanted {
			return comment, reason, true
		}
	}
	return Comment{}, "", false
}

func validateManagedSnapshot(snapshot Snapshot) error {
	if !snapshot.Complete || len(snapshot.Errors) != 0 {
		return errors.New("pull-request evidence is incomplete")
	}
	if snapshot.Number <= 0 || !commitSHAExpression.MatchString(snapshot.HeadSHA) {
		return errors.New("pull-request identity is invalid")
	}
	if snapshot.BaseRef != "main" {
		return fmt.Errorf("base branch %q is outside policy", snapshot.BaseRef)
	}
	if snapshot.State != "open" {
		return fmt.Errorf("pull request state %q is not mutable", snapshot.State)
	}
	if snapshot.Author.Login == "" || snapshot.Author.Type == "" || snapshot.AuthorLogin != snapshot.Author.Login {
		return errors.New("pull-request author identity is incomplete")
	}
	return nil
}

func isProvenDescendant(snapshot Snapshot, base string) bool {
	for _, comparison := range snapshot.Comparisons {
		if comparison.BaseSHA == base && comparison.HeadSHA == snapshot.HeadSHA && comparison.Status == "ahead" && comparison.AheadBy > 0 {
			return true
		}
	}
	return false
}

func missingSuccessfulChecks(checks []CheckResult, head string) []string {
	latest := make(map[string]CheckResult, len(checks))
	ambiguous := make(map[string]bool, len(checks))
	for _, check := range checks {
		if check.SHA != head {
			continue
		}
		current, exists := latest[check.Context]
		if !exists || check.CreatedAt.After(current.CreatedAt) {
			latest[check.Context] = check
			ambiguous[check.Context] = false
		} else if check.CreatedAt.Equal(current.CreatedAt) && check.State != current.State {
			ambiguous[check.Context] = true
		}
	}
	var missing []string
	for _, context := range RequiredS006Checks {
		check, ok := latest[context]
		if !ok || ambiguous[context] || check.State != "success" {
			missing = append(missing, context)
		}
	}
	return missing
}

func uniqueThreadCommits(threads []ReviewThread) []string {
	seen := make(map[string]struct{}, len(threads))
	var commits []string
	for _, thread := range threads {
		if _, ok := seen[thread.CommitSHA]; ok {
			continue
		}
		seen[thread.CommitSHA] = struct{}{}
		commits = append(commits, thread.CommitSHA)
	}
	sort.Strings(commits)
	return commits
}

func parseReservation(description string) (int, string, bool) {
	matches := reservationPattern.FindStringSubmatch(description)
	if matches == nil {
		return 0, "", false
	}
	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", false
	}
	return number, matches[2], true
}

func standaloneInvocationCount(body string) int {
	count := 0
	for _, line := range logicalLines(body) {
		if strings.TrimSpace(line) == "@codex review" {
			count++
		}
	}
	return count
}

func hasCodexReaction(reactions []Reaction, content, target string, after time.Time) bool {
	for _, reaction := range reactions {
		if !IsCodexActor(reaction.Author) || reaction.Content != content || reaction.CreatedAt.Before(after) {
			continue
		}
		if target == "" || reaction.Target == target || reaction.Target == "pull_request" {
			return true
		}
	}
	return false
}

func canonicalCodexLogin(login string) bool {
	return login == CodexRESTLogin || login == strings.TrimSuffix(CodexRESTLogin, "[bot]")
}

func logicalLines(body string) []string {
	return strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
}

func latestTime(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}
	return right
}

func result(context, head string, state StatusState, description string, evidence ...string) PolicyResult {
	if len(description) > 140 {
		description = description[:140]
	}
	return PolicyResult{Context: context, State: state, Description: description, HeadSHA: head, Publish: true, Evidence: evidence}
}

func errorResult(context, head, description string) PolicyResult {
	return result(context, head, StatusError, description, "automation failed closed")
}

func closedResult(context, head string, err error) PolicyResult {
	return PolicyResult{Context: context, State: StatusError, Description: err.Error(), HeadSHA: head, Publish: false, RequestSecond: false, Evidence: []string{"no mutation permitted"}}
}
