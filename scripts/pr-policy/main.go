package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type commandOptions struct {
	Repository string
	PRNumber   int
	AllOpen    bool
	DryRun     bool
}

type commandReport struct {
	Repository string       `json:"repository"`
	DryRun     bool         `json:"dry_run"`
	Results    []pullReport `json:"results,omitempty"`
	Error      string       `json:"error,omitempty"`
}

type pullReport struct {
	Number             int          `json:"number"`
	HeadSHA            string       `json:"head_sha,omitempty"`
	BaseRef            string       `json:"base_ref,omitempty"`
	IssueLink          policyReport `json:"issue_link"`
	CodexReview        policyReport `json:"codex_review"`
	IssueStatusChanged bool         `json:"issue_status_changed,omitempty"`
	CodexStatusChanged bool         `json:"codex_status_changed,omitempty"`
	ReviewCommentID    int64        `json:"review_comment_id,omitempty"`
	Error              string       `json:"error,omitempty"`
}

type policyReport struct {
	Context       string   `json:"context,omitempty"`
	State         string   `json:"state,omitempty"`
	Description   string   `json:"description,omitempty"`
	HeadSHA       string   `json:"head_sha,omitempty"`
	Publish       bool     `json:"publish"`
	RequestSecond bool     `json:"request_second"`
	Evidence      []string `json:"evidence,omitempty"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	report, err := execute(ctx, os.Args[1:], os.Getenv)
	if err != nil {
		report.Error = err.Error()
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(true)
	if encodeErr := encoder.Encode(report); encodeErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "encode result: %v\n", encodeErr)
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
}

func execute(ctx context.Context, arguments []string, getenv func(string) string) (commandReport, error) {
	options, err := parseOptions(arguments, io.Discard)
	if err != nil {
		return commandReport{}, err
	}
	report := commandReport{Repository: options.Repository, DryRun: options.DryRun}
	if !options.DryRun && getenv("GITHUB_ACTIONS") != "true" {
		return report, errors.New("mutation is allowed only in GitHub Actions; use -dry-run locally")
	}
	token := strings.TrimSpace(getenv("CUESON_GITHUB_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(getenv("GITHUB_TOKEN"))
	}
	client, err := newGitHubClient(options.Repository, token)
	if err != nil {
		return report, err
	}
	numbers := []int{options.PRNumber}
	if options.AllOpen {
		numbers, err = client.listOpenPullRequests(ctx)
		if err != nil {
			return report, fmt.Errorf("list open pull requests: %w", err)
		}
	}
	targetURL := workflowRunURL(options.Repository, getenv)
	var failures []error
	for _, number := range numbers {
		result, err := reconcilePullRequest(ctx, client, number, options.DryRun, targetURL)
		if err != nil {
			result.Error = err.Error()
			failures = append(failures, fmt.Errorf("pull request %d: %w", number, err))
		}
		report.Results = append(report.Results, result)
	}
	return report, errors.Join(failures...)
}

func parseOptions(arguments []string, output io.Writer) (commandOptions, error) {
	var options commandOptions
	set := flag.NewFlagSet("pr-policy", flag.ContinueOnError)
	set.SetOutput(output)
	set.StringVar(&options.Repository, "repo", "", "GitHub repository in owner/name form")
	set.IntVar(&options.PRNumber, "pr", 0, "one pull request number to reconcile")
	set.BoolVar(&options.AllOpen, "all-open", false, "reconcile every open pull request targeting main")
	set.BoolVar(&options.DryRun, "dry-run", false, "collect and evaluate without mutation")
	if err := set.Parse(arguments); err != nil {
		return commandOptions{}, err
	}
	if set.NArg() != 0 {
		return commandOptions{}, errors.New("positional arguments are not accepted")
	}
	if options.Repository == "" {
		return commandOptions{}, errors.New("-repo owner/name is required")
	}
	if (options.PRNumber > 0) == options.AllOpen {
		return commandOptions{}, errors.New("specify exactly one of -pr or -all-open")
	}
	return options, nil
}

func workflowRunURL(repository string, getenv func(string) string) string {
	server := strings.TrimSuffix(getenv("GITHUB_SERVER_URL"), "/")
	if server == "" {
		server = "https://github.com"
	}
	runID := getenv("GITHUB_RUN_ID")
	if _, err := strconv.ParseUint(runID, 10, 64); err != nil {
		return ""
	}
	return server + "/" + repository + "/actions/runs/" + runID
}

func reconcilePullRequest(ctx context.Context, client *githubClient, number int, dryRun bool, targetURL string) (pullReport, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()
	snapshot, err := client.collectSnapshot(ctx, number)
	if err != nil {
		return pullReport{Number: number}, err
	}
	issueResult, codexResult := Evaluate(snapshot)
	report := pullReport{
		Number:      number,
		HeadSHA:     snapshot.HeadSHA,
		BaseRef:     snapshot.BaseRef,
		IssueLink:   reportPolicy(issueResult),
		CodexReview: reportPolicy(codexResult),
	}
	if dryRun {
		return report, nil
	}
	if snapshot.BaseRef != "main" || snapshot.State != "open" {
		return report, nil
	}
	if issueResult.Publish {
		changed, err := client.publishGuardedStatus(ctx, number, snapshot.HeadSHA, statusFromPolicy(issueResult, targetURL))
		if err != nil {
			return report, err
		}
		report.IssueStatusChanged = changed
	}
	if !codexResult.Publish {
		return report, nil
	}
	if !codexResult.RequestSecond {
		changed, err := client.publishGuardedStatus(ctx, number, snapshot.HeadSHA, statusFromPolicy(codexResult, targetURL))
		if err != nil {
			return report, err
		}
		report.CodexStatusChanged = changed
		return report, nil
	}

	reservation := statusFromPolicy(codexResult, targetURL)
	changed, err := client.publishGuardedStatus(ctx, number, snapshot.HeadSHA, reservation)
	if err != nil {
		return report, fmt.Errorf("reserve second review: %w", err)
	}
	report.CodexStatusChanged = changed
	if !changed {
		return report, errors.New("second-review reservation already existed; refusing another request")
	}
	refreshed, err := client.collectSnapshot(ctx, number)
	if err != nil {
		return report, fmt.Errorf("read back second-review reservation: %w", err)
	}
	if refreshed.HeadSHA != snapshot.HeadSHA || refreshed.BaseRef != "main" || refreshed.State != "open" {
		return report, errors.New("pull request changed after second-review reservation")
	}
	if !hasReservation(refreshed, reservation.Description, snapshot.HeadSHA) {
		return report, errors.New("second-review reservation was absent from the complete snapshot")
	}
	preRequestSnapshot, err := withoutExpectedReservation(refreshed, reservation.Description, snapshot.HeadSHA)
	if err != nil {
		return report, fmt.Errorf("validate second-review reservation: %w", err)
	}
	preRequestResult := EvaluateCodexReview(preRequestSnapshot)
	if !preRequestResult.Publish || !preRequestResult.RequestSecond || preRequestResult.Description != reservation.Description {
		return report, fmt.Errorf("second-review eligibility changed after reservation: %s", preRequestResult.Description)
	}
	if countRoundTwoRequests(refreshed) != 0 {
		return report, errors.New("a second-review request appeared after reservation; refusing a duplicate")
	}
	latestPull, err := client.getPullRequest(ctx, number)
	if err != nil {
		return report, err
	}
	if err := client.requireMutationTarget(latestPull, number, snapshot.HeadSHA); err != nil {
		return report, err
	}
	body := RoundTwoComment(number, snapshot.HeadSHA)
	created, err := client.createReviewRequest(ctx, number, body)
	if err != nil {
		return report, err
	}
	report.ReviewCommentID = created.ID
	finalSnapshot, err := client.collectSnapshot(ctx, number)
	if err != nil {
		return report, fmt.Errorf("read back second-review evidence: %w", err)
	}
	if finalSnapshot.HeadSHA != snapshot.HeadSHA || countAuthenticRoundTwoMarkers(finalSnapshot, number, snapshot.HeadSHA, created.ID, body) != 1 {
		return report, errors.New("second-review comment did not survive complete read-back exactly once")
	}
	_, finalCodex := Evaluate(finalSnapshot)
	report.CodexReview = reportPolicy(finalCodex)
	if finalCodex.Publish {
		changed, err := client.publishGuardedStatus(ctx, number, snapshot.HeadSHA, statusFromPolicy(finalCodex, targetURL))
		if err != nil {
			return report, fmt.Errorf("publish post-request review status: %w", err)
		}
		report.CodexStatusChanged = report.CodexStatusChanged || changed
	}
	return report, nil
}

func reportPolicy(result PolicyResult) policyReport {
	return policyReport{
		Context:       result.Context,
		State:         string(result.State),
		Description:   result.Description,
		HeadSHA:       result.HeadSHA,
		Publish:       result.Publish,
		RequestSecond: result.RequestSecond,
		Evidence:      append([]string(nil), result.Evidence...),
	}
}

func statusFromPolicy(result PolicyResult, targetURL string) commitStatus {
	return commitStatus{Context: result.Context, State: string(result.State), Description: result.Description, TargetURL: targetURL}
}

func (c *githubClient) publishGuardedStatus(ctx context.Context, number int, head string, status commitStatus) (bool, error) {
	pull, err := c.getPullRequest(ctx, number)
	if err != nil {
		return false, err
	}
	if err := c.requireMutationTarget(pull, number, head); err != nil {
		return false, err
	}
	return c.publishStatus(ctx, head, status)
}

func hasReservation(snapshot Snapshot, description, head string) bool {
	for _, check := range snapshot.Checks {
		if check.Context == CodexReviewContext && check.SHA == head && check.Description == description && check.State == string(StatusPending) && IsGitHubActionsActor(check.Creator) {
			return true
		}
	}
	return false
}

func withoutExpectedReservation(snapshot Snapshot, description, head string) (Snapshot, error) {
	filtered := append([]CheckResult(nil), snapshot.Checks...)
	kept := filtered[:0]
	removed := 0
	for _, check := range filtered {
		if check.Context == CodexReviewContext && check.SHA == head && check.Description == description && check.State == string(StatusPending) && IsGitHubActionsActor(check.Creator) {
			removed++
			continue
		}
		kept = append(kept, check)
	}
	if removed != 1 {
		return Snapshot{}, fmt.Errorf("expected one authenticated reservation, found %d", removed)
	}
	snapshot.Checks = kept
	return snapshot, nil
}

func countRoundTwoRequests(snapshot Snapshot) int {
	count := 0
	for _, comment := range snapshot.Comments {
		if _, _, ok := ParseRoundTwoMarker(comment); ok {
			count++
			continue
		}
		if comment.Author.Login == OperatorLogin && containsStandaloneCodexReview(comment.Body) {
			count++
		}
	}
	return count
}

func countAuthenticRoundTwoMarkers(snapshot Snapshot, number int, head string, commentID int64, body string) int {
	count := 0
	for _, comment := range snapshot.Comments {
		markerPR, markerHead, ok := ParseRoundTwoMarker(comment)
		if !ok || markerPR != number || markerHead != head {
			continue
		}
		if comment.ID != commentID || comment.Body != body || !IsGitHubActionsActor(comment.Author) {
			continue
		}
		if strings.Count(comment.Body, "@codex review") == 1 {
			count++
		}
	}
	return count
}

func containsStandaloneCodexReview(body string) bool {
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "@codex review" {
			return true
		}
	}
	return false
}
