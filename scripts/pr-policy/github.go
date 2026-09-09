package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	githubAPIVersion   = "2022-11-28"
	maximumAPIPages    = 1000
	maximumOpenPRs     = 100
	maximumSnapshotTry = 3
)

var fullSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

type githubClient struct {
	httpClient *http.Client
	restBase   *url.URL
	graphqlURL string
	token      string
	owner      string
	repository string
}

type githubActor struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Type  string `json:"type"`
}

type githubRef struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

type githubPullRequest struct {
	Number    int         `json:"number"`
	State     string      `json:"state"`
	Draft     bool        `json:"draft"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      githubActor `json:"user"`
	Head      githubRef   `json:"head"`
	Base      githubRef   `json:"base"`
}

type githubComment struct {
	ID        int64       `json:"id"`
	Body      string      `json:"body"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      githubActor `json:"user"`
}

type githubReaction struct {
	ID        int64       `json:"id"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	User      githubActor `json:"user"`
}

type githubReview struct {
	ID          int64       `json:"id"`
	NodeID      string      `json:"node_id"`
	Body        string      `json:"body"`
	CommitID    string      `json:"commit_id"`
	SubmittedAt time.Time   `json:"submitted_at"`
	User        githubActor `json:"user"`
}

type githubStatus struct {
	ID          int64       `json:"id"`
	SHA         string      `json:"sha"`
	State       string      `json:"state"`
	Context     string      `json:"context"`
	Description string      `json:"description"`
	TargetURL   string      `json:"target_url"`
	CreatedAt   time.Time   `json:"created_at"`
	Creator     githubActor `json:"creator"`
}

type githubCheckRun struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	HeadSHA     string    `json:"head_sha"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type commitStatus struct {
	Context     string `json:"context"`
	State       string `json:"state"`
	Description string `json:"description"`
	TargetURL   string `json:"target_url,omitempty"`
}

type comparisonResponse struct {
	Status  string `json:"status"`
	AheadBy int    `json:"ahead_by"`
}

func newGitHubClient(repository, token string) (*githubClient, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.ContainsAny(repository, "?#\\") {
		return nil, fmt.Errorf("repository must be owner/name")
	}
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("GitHub token is empty")
	}
	base, err := url.Parse("https://api.github.com/")
	if err != nil {
		return nil, fmt.Errorf("parse GitHub API URL: %w", err)
	}
	return &githubClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		restBase:   base,
		graphqlURL: "https://api.github.com/graphql",
		token:      token,
		owner:      parts[0],
		repository: parts[1],
	}, nil
}

func (c *githubClient) repositoryPath(suffix string) string {
	return "/repos/" + url.PathEscape(c.owner) + "/" + url.PathEscape(c.repository) + suffix
}

func (c *githubClient) request(ctx context.Context, method, target string, body any) (*http.Response, error) {
	requestURL, err := c.resolveRESTTarget(target)
	if err != nil {
		return nil, err
	}
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode GitHub request: %w", err)
		}
		reader = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("create GitHub request: %w", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	request.Header.Set("User-Agent", "cueson-pr-policy")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GitHub %s %s: %w", method, requestURL.Redacted(), err)
	}
	return response, nil
}

func (c *githubClient) resolveRESTTarget(target string) (*url.URL, error) {
	parsed, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("parse GitHub URL: %w", err)
	}
	if !parsed.IsAbs() {
		parsed = c.restBase.ResolveReference(parsed)
	}
	if !sameOrigin(c.restBase, parsed) {
		return nil, fmt.Errorf("refusing GitHub pagination URL on a different origin: %s", parsed.Redacted())
	}
	return parsed, nil
}

func sameOrigin(left, right *url.URL) bool {
	return strings.EqualFold(left.Scheme, right.Scheme) && strings.EqualFold(left.Host, right.Host)
}

func (c *githubClient) getJSON(ctx context.Context, target string, output any) (*http.Response, error) {
	response, err := c.request(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	if err := decodeResponse(response, http.StatusOK, output); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *githubClient) postJSON(ctx context.Context, target string, input, output any) error {
	response, err := c.request(ctx, http.MethodPost, target, input)
	if err != nil {
		return err
	}
	return decodeResponse(response, http.StatusCreated, output)
}

func decodeResponse(response *http.Response, expected int, output any) error {
	defer response.Body.Close()
	if response.StatusCode != expected {
		limited, _ := io.ReadAll(io.LimitReader(response.Body, 16<<10))
		return fmt.Errorf("GitHub API returned %s: %s", response.Status, strings.TrimSpace(string(limited)))
	}
	if output == nil {
		_, _ = io.Copy(io.Discard, response.Body)
		return nil
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(output); err != nil {
		return fmt.Errorf("decode GitHub response: %w", err)
	}
	return nil
}

func (c *githubClient) getRESTPages(ctx context.Context, target string, output any) error {
	destination := reflect.ValueOf(output)
	if destination.Kind() != reflect.Pointer || destination.IsNil() || destination.Elem().Kind() != reflect.Slice {
		return errors.New("REST pagination output must be a non-nil pointer to a slice")
	}
	destination.Elem().Set(reflect.MakeSlice(destination.Elem().Type(), 0, 0))
	next := target
	for page := 0; next != ""; page++ {
		if page >= maximumAPIPages {
			return fmt.Errorf("REST pagination exceeded %d pages", maximumAPIPages)
		}
		pageValue := reflect.New(destination.Elem().Type())
		response, err := c.getJSON(ctx, next, pageValue.Interface())
		if err != nil {
			return err
		}
		destination.Elem().Set(reflect.AppendSlice(destination.Elem(), pageValue.Elem()))
		next, err = nextPage(response.Header.Values("Link"))
		if err != nil {
			return err
		}
		if next != "" {
			if _, err := c.resolveRESTTarget(next); err != nil {
				return err
			}
		}
	}
	return nil
}

func nextPage(headers []string) (string, error) {
	var next string
	for _, header := range headers {
		for _, part := range strings.Split(header, ",") {
			part = strings.TrimSpace(part)
			pieces := strings.Split(part, ";")
			if len(pieces) < 2 {
				continue
			}
			isNext := false
			for _, parameter := range pieces[1:] {
				if strings.TrimSpace(parameter) == `rel="next"` {
					isNext = true
				}
			}
			if !isNext {
				continue
			}
			candidate := strings.TrimSpace(pieces[0])
			if len(candidate) < 3 || candidate[0] != '<' || candidate[len(candidate)-1] != '>' {
				return "", fmt.Errorf("malformed GitHub next Link %q", candidate)
			}
			candidate = candidate[1 : len(candidate)-1]
			if next != "" && next != candidate {
				return "", errors.New("multiple conflicting GitHub next links")
			}
			next = candidate
		}
	}
	return next, nil
}

func (c *githubClient) graphQL(ctx context.Context, query string, variables map[string]any, output any) error {
	requestBody, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return fmt.Errorf("encode GraphQL request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.graphqlURL, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("create GraphQL request: %w", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	request.Header.Set("User-Agent", "cueson-pr-policy")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("GitHub GraphQL request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		limited, _ := io.ReadAll(io.LimitReader(response.Body, 16<<10))
		return fmt.Errorf("GitHub GraphQL returned %s: %s", response.Status, strings.TrimSpace(string(limited)))
	}
	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode GraphQL response: %w", err)
	}
	if len(envelope.Errors) != 0 {
		messages := make([]string, 0, len(envelope.Errors))
		for _, item := range envelope.Errors {
			messages = append(messages, item.Message)
		}
		return fmt.Errorf("GitHub GraphQL errors: %s", strings.Join(messages, "; "))
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return errors.New("GitHub GraphQL response omitted data")
	}
	if err := json.Unmarshal(envelope.Data, output); err != nil {
		return fmt.Errorf("decode GraphQL data: %w", err)
	}
	return nil
}

func (c *githubClient) getPullRequest(ctx context.Context, number int) (githubPullRequest, error) {
	if number <= 0 {
		return githubPullRequest{}, errors.New("pull request number must be positive")
	}
	var result githubPullRequest
	_, err := c.getJSON(ctx, c.repositoryPath("/pulls/"+strconv.Itoa(number)), &result)
	if err != nil {
		return githubPullRequest{}, err
	}
	if result.Number != number || !fullSHA.MatchString(result.Head.SHA) {
		return githubPullRequest{}, errors.New("GitHub returned an invalid pull request identity")
	}
	return result, nil
}

func (c *githubClient) listOpenPullRequests(ctx context.Context) ([]int, error) {
	var pulls []githubPullRequest
	path := c.repositoryPath("/pulls?state=open&base=main&per_page=100")
	if err := c.getRESTPages(ctx, path, &pulls); err != nil {
		return nil, err
	}
	if len(pulls) > maximumOpenPRs {
		return nil, fmt.Errorf("open pull request count %d exceeds limit %d", len(pulls), maximumOpenPRs)
	}
	numbers := make([]int, 0, len(pulls))
	seen := make(map[int]struct{}, len(pulls))
	for _, pull := range pulls {
		if pull.Number <= 0 || pull.Base.Ref != "main" {
			return nil, errors.New("open pull request listing contained invalid evidence")
		}
		if _, exists := seen[pull.Number]; exists {
			return nil, fmt.Errorf("duplicate pull request %d in listing", pull.Number)
		}
		seen[pull.Number] = struct{}{}
		numbers = append(numbers, pull.Number)
	}
	sort.Ints(numbers)
	return numbers, nil
}

func (c *githubClient) listComments(ctx context.Context, number int) ([]githubComment, error) {
	var result []githubComment
	err := c.getRESTPages(ctx, c.repositoryPath("/issues/"+strconv.Itoa(number)+"/comments?per_page=100"), &result)
	return result, err
}

func (c *githubClient) listReactions(ctx context.Context, target string) ([]githubReaction, error) {
	var result []githubReaction
	err := c.getRESTPages(ctx, target+querySeparator(target)+"per_page=100", &result)
	return result, err
}

func querySeparator(value string) string {
	if strings.Contains(value, "?") {
		return "&"
	}
	return "?"
}

func (c *githubClient) listReviews(ctx context.Context, number int) ([]githubReview, error) {
	var result []githubReview
	err := c.getRESTPages(ctx, c.repositoryPath("/pulls/"+strconv.Itoa(number)+"/reviews?per_page=100"), &result)
	return result, err
}

func (c *githubClient) listStatuses(ctx context.Context, sha string) ([]githubStatus, error) {
	if !fullSHA.MatchString(sha) {
		return nil, errors.New("status SHA must be a full lowercase hexadecimal identifier")
	}
	var result []githubStatus
	err := c.getRESTPages(ctx, c.repositoryPath("/commits/"+sha+"/statuses?per_page=100"), &result)
	return result, err
}

func (c *githubClient) listCheckRuns(ctx context.Context, sha string) ([]githubCheckRun, error) {
	if !fullSHA.MatchString(sha) {
		return nil, errors.New("check-run SHA must be a full lowercase hexadecimal identifier")
	}
	var result []githubCheckRun
	next := c.repositoryPath("/commits/" + sha + "/check-runs?per_page=100")
	for page := 0; next != ""; page++ {
		if page >= maximumAPIPages {
			return nil, fmt.Errorf("check-run pagination exceeded %d pages", maximumAPIPages)
		}
		var responseBody struct {
			CheckRuns []githubCheckRun `json:"check_runs"`
		}
		response, err := c.getJSON(ctx, next, &responseBody)
		if err != nil {
			return nil, err
		}
		result = append(result, responseBody.CheckRuns...)
		next, err = nextPage(response.Header.Values("Link"))
		if err != nil {
			return nil, err
		}
		if next != "" {
			if _, err := c.resolveRESTTarget(next); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func (c *githubClient) compare(ctx context.Context, base, head string) (Comparison, error) {
	if !fullSHA.MatchString(base) || !fullSHA.MatchString(head) {
		return Comparison{}, errors.New("comparison requires full lowercase SHAs")
	}
	var response comparisonResponse
	path := c.repositoryPath("/compare/" + base + "..." + head)
	if _, err := c.getJSON(ctx, path, &response); err != nil {
		return Comparison{}, err
	}
	return Comparison{BaseSHA: base, HeadSHA: head, Status: response.Status, AheadBy: response.AheadBy}, nil
}

func (c *githubClient) publishStatus(ctx context.Context, sha string, desired commitStatus) (bool, error) {
	if !fullSHA.MatchString(sha) {
		return false, errors.New("status SHA must be a full lowercase hexadecimal identifier")
	}
	if desired.Context != IssueLinkContext && desired.Context != CodexReviewContext {
		return false, fmt.Errorf("unsupported status context %q", desired.Context)
	}
	if desired.State != "pending" && desired.State != "success" && desired.State != "failure" && desired.State != "error" {
		return false, fmt.Errorf("unsupported status state %q", desired.State)
	}
	if len(desired.Context) > 100 || len(desired.Description) == 0 || len(desired.Description) > 140 {
		return false, errors.New("status context or description violates GitHub limits")
	}
	statuses, err := c.listStatuses(ctx, sha)
	if err != nil {
		return false, err
	}
	for _, status := range statuses {
		if status.Context != desired.Context {
			continue
		}
		if status.State == desired.State && status.Description == desired.Description {
			return false, nil
		}
		break
	}
	var created githubStatus
	if err := c.postJSON(ctx, c.repositoryPath("/statuses/"+sha), desired, &created); err != nil {
		return false, err
	}
	statuses, err = c.listStatuses(ctx, sha)
	if err != nil {
		return true, fmt.Errorf("status mutation accepted but read-back failed: %w", err)
	}
	for _, status := range statuses {
		if status.Context != desired.Context {
			continue
		}
		if status.ID != created.ID || status.SHA != sha || status.State != desired.State || status.Description != desired.Description {
			return true, errors.New("status mutation read-back did not match the accepted result")
		}
		return true, nil
	}
	return true, errors.New("status mutation was absent from read-back")
}

func (c *githubClient) createReviewRequest(ctx context.Context, number int, body string) (githubComment, error) {
	if number <= 0 {
		return githubComment{}, errors.New("pull request number must be positive")
	}
	if strings.Count(body, "@codex review") != 1 || strings.Count(body, RoundTwoMarkerPrefix) != 1 {
		return githubComment{}, errors.New("round-two comment must contain one invocation and one marker")
	}
	var created githubComment
	if err := c.postJSON(ctx, c.repositoryPath("/issues/"+strconv.Itoa(number)+"/comments"), map[string]string{"body": body}, &created); err != nil {
		return githubComment{}, fmt.Errorf("round-two comment result is ambiguous and must not be retried automatically: %w", err)
	}
	var readBack githubComment
	if _, err := c.getJSON(ctx, c.repositoryPath("/issues/comments/"+strconv.FormatInt(created.ID, 10)), &readBack); err != nil {
		return githubComment{}, fmt.Errorf("round-two comment was accepted but read-back failed: %w", err)
	}
	if readBack.ID != created.ID || readBack.Body != body || !IsGitHubActionsActor(actorFromGitHub(readBack.User)) {
		return githubComment{}, errors.New("round-two comment read-back did not match the authenticated mutation")
	}
	return readBack, nil
}

func (c *githubClient) requireMutationTarget(pull githubPullRequest, number int, head string) error {
	if pull.Number != number || pull.Number <= 0 {
		return errors.New("mutation pull request number changed")
	}
	if pull.Base.Ref != "main" {
		return fmt.Errorf("refusing mutation for base branch %q", pull.Base.Ref)
	}
	if pull.State != "open" {
		return fmt.Errorf("refusing mutation for pull request state %q", pull.State)
	}
	if !fullSHA.MatchString(head) || pull.Head.SHA != head {
		return errors.New("refusing mutation because pull request head changed")
	}
	return nil
}

func (c *githubClient) collectSnapshot(ctx context.Context, number int) (Snapshot, error) {
	var lastError error
	for attempt := 0; attempt < maximumSnapshotTry; attempt++ {
		start, err := c.getPullRequest(ctx, number)
		if err != nil {
			return Snapshot{}, err
		}
		snapshot, err := c.collectSnapshotAt(ctx, start)
		if err != nil {
			return Snapshot{}, err
		}
		end, err := c.getPullRequest(ctx, number)
		if err != nil {
			return Snapshot{}, err
		}
		if samePullRequestVersion(start, end) {
			snapshot.Complete = true
			return snapshot, nil
		}
		lastError = errors.New("pull request evidence changed during snapshot collection")
	}
	return Snapshot{}, fmt.Errorf("snapshot did not stabilize after %d attempts: %w", maximumSnapshotTry, lastError)
}

func samePullRequestVersion(left, right githubPullRequest) bool {
	return left.Number == right.Number && left.Head.SHA == right.Head.SHA && left.Base.Ref == right.Base.Ref && left.State == right.State && left.Draft == right.Draft && left.UpdatedAt.Equal(right.UpdatedAt)
}

func (c *githubClient) collectSnapshotAt(ctx context.Context, pull githubPullRequest) (Snapshot, error) {
	closingIssues, err := c.listClosingIssues(ctx, pull.Number)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect closing issues: %w", err)
	}
	apiComments, err := c.listComments(ctx, pull.Number)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect comments: %w", err)
	}
	apiReviews, err := c.listReviews(ctx, pull.Number)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect reviews: %w", err)
	}
	reviewThreads, err := c.listReviewThreads(ctx, pull.Number)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect review threads: %w", err)
	}

	snapshot := Snapshot{
		Number:        pull.Number,
		HeadSHA:       pull.Head.SHA,
		BaseRef:       pull.Base.Ref,
		AuthorLogin:   pull.User.Login,
		Author:        actorFromGitHub(pull.User),
		State:         pull.State,
		Draft:         pull.Draft,
		ClosingIssues: closingIssues,
		ReviewThreads: reviewThreads,
	}
	for _, item := range apiComments {
		snapshot.Comments = append(snapshot.Comments, commentFromGitHub(item))
	}
	for _, item := range apiReviews {
		review, err := reviewFromGitHub(item)
		if err != nil {
			return Snapshot{}, fmt.Errorf("collect reviews: %w", err)
		}
		snapshot.Reviews = append(snapshot.Reviews, review)
	}

	pullReactions, err := c.listReactions(ctx, c.repositoryPath("/issues/"+strconv.Itoa(pull.Number)+"/reactions"))
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect pull request reactions: %w", err)
	}
	for _, item := range pullReactions {
		snapshot.Reactions = append(snapshot.Reactions, reactionFromGitHub(item, "pull_request"))
	}
	for _, comment := range apiComments {
		if !strings.Contains(comment.Body, "<!-- codex-pull-request-review-summary -->") || !isRawCodexActor(comment.User) {
			continue
		}
		target := c.repositoryPath("/issues/comments/" + strconv.FormatInt(comment.ID, 10) + "/reactions")
		reactions, err := c.listReactions(ctx, target)
		if err != nil {
			return Snapshot{}, fmt.Errorf("collect summary comment reactions: %w", err)
		}
		for _, item := range reactions {
			snapshot.Reactions = append(snapshot.Reactions, reactionFromGitHub(item, "comment:"+strconv.FormatInt(comment.ID, 10)))
		}
	}

	checkRuns, err := c.listCheckRuns(ctx, pull.Head.SHA)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect check runs: %w", err)
	}
	for _, check := range checkRuns {
		state := check.Conclusion
		if state == "" {
			state = check.Status
		}
		createdAt := check.CompletedAt
		if createdAt.IsZero() {
			createdAt = check.StartedAt
		}
		snapshot.Checks = append(snapshot.Checks, CheckResult{Context: check.Name, State: state, SHA: check.HeadSHA, CreatedAt: createdAt})
	}
	historicalHeads, err := c.listHistoricalHeads(ctx, pull.Number, pull.Head.SHA)
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect pull request history: %w", err)
	}
	for _, sha := range historicalHeads {
		statuses, err := c.listStatuses(ctx, sha)
		if err != nil {
			return Snapshot{}, fmt.Errorf("collect status history for %s: %w", sha, err)
		}
		for _, status := range statuses {
			if status.Context != IssueLinkContext && status.Context != CodexReviewContext {
				continue
			}
			snapshot.Checks = append(snapshot.Checks, CheckResult{
				Context:     status.Context,
				State:       status.State,
				Description: status.Description,
				SHA:         sha,
				Creator:     actorFromGitHub(status.Creator),
				CreatedAt:   status.CreatedAt,
			})
		}
	}

	reviewIDsWithThreads := make(map[string]struct{})
	for _, thread := range snapshot.ReviewThreads {
		if thread.ReviewID != "" {
			reviewIDsWithThreads[thread.ReviewID] = struct{}{}
		}
	}
	comparisonSHAs := make(map[string]struct{})
	for index := range snapshot.Reviews {
		review := &snapshot.Reviews[index]
		if IsCodexActor(review.Author) && strings.TrimSpace(review.Body) != "" {
			if _, exists := reviewIDsWithThreads[review.ID]; !exists {
				review.UnthreadedFinding = true
			}
		}
		if IsCodexActor(review.Author) && fullSHA.MatchString(review.CommitSHA) && review.CommitSHA != pull.Head.SHA {
			comparisonSHAs[review.CommitSHA] = struct{}{}
		}
	}
	for _, thread := range snapshot.ReviewThreads {
		if IsCodexActor(thread.Author) && fullSHA.MatchString(thread.CommitSHA) && thread.CommitSHA != pull.Head.SHA {
			comparisonSHAs[thread.CommitSHA] = struct{}{}
		}
	}
	orderedComparisonSHAs := make([]string, 0, len(comparisonSHAs))
	for sha := range comparisonSHAs {
		orderedComparisonSHAs = append(orderedComparisonSHAs, sha)
	}
	sort.Strings(orderedComparisonSHAs)
	for _, sha := range orderedComparisonSHAs {
		comparison, err := c.compare(ctx, sha, pull.Head.SHA)
		if err != nil {
			return Snapshot{}, fmt.Errorf("compare review head %s: %w", sha, err)
		}
		snapshot.Comparisons = append(snapshot.Comparisons, comparison)
	}
	return snapshot, nil
}

func actorFromGitHub(actor githubActor) Actor {
	return Actor(actor)
}

func commentFromGitHub(comment githubComment) Comment {
	return Comment{ID: comment.ID, Author: actorFromGitHub(comment.User), Body: comment.Body, CreatedAt: comment.CreatedAt, UpdatedAt: comment.UpdatedAt}
}

func reviewFromGitHub(review githubReview) (Review, error) {
	if strings.TrimSpace(review.NodeID) == "" {
		return Review{}, fmt.Errorf("review %d has no GraphQL node ID", review.ID)
	}
	return Review{ID: review.NodeID, Author: actorFromGitHub(review.User), CommitSHA: review.CommitID, Body: review.Body, SubmittedAt: review.SubmittedAt}, nil
}

func reactionFromGitHub(reaction githubReaction, target string) Reaction {
	author := actorFromGitHub(reaction.User)
	if author.Login == CodexRESTLogin && author.ID == CodexUserID && author.Type == "User" {
		author.Type = "Bot"
	}
	return Reaction{Author: author, Content: reaction.Content, Target: target, CreatedAt: reaction.CreatedAt}
}

func isRawCodexActor(actor githubActor) bool {
	return actor.Login == CodexRESTLogin && actor.ID == CodexUserID && actor.Type == "Bot"
}

func (c *githubClient) listClosingIssues(ctx context.Context, number int) ([]ClosingIssue, error) {
	const query = `query($owner:String!,$repository:String!,$number:Int!,$cursor:String){repository(owner:$owner,name:$repository){pullRequest(number:$number){closingIssuesReferences(first:100,after:$cursor,excludeUserLinked:true){pageInfo{hasNextPage endCursor}nodes{number repository{nameWithOwner}}}}}}`
	var result []ClosingIssue
	var cursor any
	for page := 0; ; page++ {
		if page >= maximumAPIPages {
			return nil, fmt.Errorf("closing-reference pagination exceeded %d pages", maximumAPIPages)
		}
		var response struct {
			Repository *struct {
				PullRequest *struct {
					Closing struct {
						PageInfo graphQLPageInfo `json:"pageInfo"`
						Nodes    []struct {
							Number     int `json:"number"`
							Repository struct {
								NameWithOwner string `json:"nameWithOwner"`
							} `json:"repository"`
						} `json:"nodes"`
					} `json:"closingIssuesReferences"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		variables := map[string]any{"owner": c.owner, "repository": c.repository, "number": number, "cursor": cursor}
		if err := c.graphQL(ctx, query, variables, &response); err != nil {
			return nil, err
		}
		if response.Repository == nil || response.Repository.PullRequest == nil {
			return nil, errors.New("closing-reference query did not return the pull request")
		}
		connection := response.Repository.PullRequest.Closing
		for _, node := range connection.Nodes {
			parts := strings.Split(node.Repository.NameWithOwner, "/")
			if len(parts) != 2 || node.Number <= 0 {
				return nil, errors.New("closing-reference query returned an invalid issue")
			}
			result = append(result, ClosingIssue{Owner: parts[0], Repository: parts[1], Number: node.Number})
		}
		if !connection.PageInfo.HasNextPage {
			return result, nil
		}
		if connection.PageInfo.EndCursor == "" {
			return nil, errors.New("closing-reference query omitted its next cursor")
		}
		cursor = connection.PageInfo.EndCursor
	}
}

type graphQLPageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type graphQLThreadComment struct {
	Author *struct {
		Login      string `json:"login"`
		DatabaseID int64  `json:"databaseId"`
		Type       string `json:"__typename"`
	} `json:"author"`
	CreatedAt         time.Time `json:"createdAt"`
	PullRequestReview *struct {
		ID     string `json:"id"`
		Commit *struct {
			OID string `json:"oid"`
		} `json:"commit"`
	} `json:"pullRequestReview"`
}

type graphQLThread struct {
	ID         string `json:"id"`
	IsResolved bool   `json:"isResolved"`
	IsOutdated bool   `json:"isOutdated"`
	Comments   struct {
		PageInfo graphQLPageInfo        `json:"pageInfo"`
		Nodes    []graphQLThreadComment `json:"nodes"`
	} `json:"comments"`
}

func (c *githubClient) listReviewThreads(ctx context.Context, number int) ([]ReviewThread, error) {
	const query = `query($owner:String!,$repository:String!,$number:Int!,$cursor:String){repository(owner:$owner,name:$repository){pullRequest(number:$number){reviewThreads(first:100,after:$cursor){pageInfo{hasNextPage endCursor}nodes{id isResolved isOutdated comments(first:100){pageInfo{hasNextPage endCursor}nodes{author{__typename login ... on Bot{databaseId} ... on User{databaseId}}createdAt pullRequestReview{id commit{oid}}}}}}}}}`
	var rawThreads []graphQLThread
	var cursor any
	for page := 0; ; page++ {
		if page >= maximumAPIPages {
			return nil, fmt.Errorf("review-thread pagination exceeded %d pages", maximumAPIPages)
		}
		var response struct {
			Repository *struct {
				PullRequest *struct {
					ReviewThreads struct {
						PageInfo graphQLPageInfo `json:"pageInfo"`
						Nodes    []graphQLThread `json:"nodes"`
					} `json:"reviewThreads"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		variables := map[string]any{"owner": c.owner, "repository": c.repository, "number": number, "cursor": cursor}
		if err := c.graphQL(ctx, query, variables, &response); err != nil {
			return nil, err
		}
		if response.Repository == nil || response.Repository.PullRequest == nil {
			return nil, errors.New("review-thread query did not return the pull request")
		}
		connection := response.Repository.PullRequest.ReviewThreads
		rawThreads = append(rawThreads, connection.Nodes...)
		if !connection.PageInfo.HasNextPage {
			break
		}
		if connection.PageInfo.EndCursor == "" {
			return nil, errors.New("review-thread query omitted its next cursor")
		}
		cursor = connection.PageInfo.EndCursor
	}
	result := make([]ReviewThread, 0, len(rawThreads))
	for _, thread := range rawThreads {
		comments := append([]graphQLThreadComment(nil), thread.Comments.Nodes...)
		cursor := thread.Comments.PageInfo.EndCursor
		for page := 1; thread.Comments.PageInfo.HasNextPage; page++ {
			if page >= maximumAPIPages {
				return nil, fmt.Errorf("review-thread comment pagination exceeded %d pages", maximumAPIPages)
			}
			if cursor == "" {
				return nil, errors.New("review-thread comment query omitted its next cursor")
			}
			pageComments, pageInfo, err := c.listReviewThreadCommentPage(ctx, thread.ID, cursor)
			if err != nil {
				return nil, err
			}
			comments = append(comments, pageComments...)
			thread.Comments.PageInfo = pageInfo
			cursor = pageInfo.EndCursor
		}
		var finding *graphQLThreadComment
		for index := range comments {
			comment := &comments[index]
			if comment.Author == nil {
				continue
			}
			actor := Actor{Login: comment.Author.Login, ID: comment.Author.DatabaseID, Type: comment.Author.Type}
			if IsCodexActor(actor) {
				finding = comment
				break
			}
		}
		if finding == nil {
			continue
		}
		reviewThread := ReviewThread{ID: thread.ID, Author: Actor{Login: finding.Author.Login, ID: finding.Author.DatabaseID, Type: finding.Author.Type}, CreatedAt: finding.CreatedAt, Resolved: thread.IsResolved, Outdated: thread.IsOutdated}
		if finding.PullRequestReview != nil {
			reviewThread.ReviewID = finding.PullRequestReview.ID
			if finding.PullRequestReview.Commit != nil {
				reviewThread.CommitSHA = finding.PullRequestReview.Commit.OID
			}
		}
		result = append(result, reviewThread)
	}
	return result, nil
}

func (c *githubClient) listReviewThreadCommentPage(ctx context.Context, threadID, cursor string) ([]graphQLThreadComment, graphQLPageInfo, error) {
	const query = `query($id:ID!,$cursor:String!){node(id:$id){... on PullRequestReviewThread{comments(first:100,after:$cursor){pageInfo{hasNextPage endCursor}nodes{author{__typename login ... on Bot{databaseId} ... on User{databaseId}}createdAt pullRequestReview{id commit{oid}}}}}}}`
	var response struct {
		Node *struct {
			Comments struct {
				PageInfo graphQLPageInfo        `json:"pageInfo"`
				Nodes    []graphQLThreadComment `json:"nodes"`
			} `json:"comments"`
		} `json:"node"`
	}
	if err := c.graphQL(ctx, query, map[string]any{"id": threadID, "cursor": cursor}, &response); err != nil {
		return nil, graphQLPageInfo{}, err
	}
	if response.Node == nil {
		return nil, graphQLPageInfo{}, errors.New("review-thread comment query did not return its thread")
	}
	return response.Node.Comments.Nodes, response.Node.Comments.PageInfo, nil
}

func (c *githubClient) listHistoricalHeads(ctx context.Context, number int, current string) ([]string, error) {
	const query = `query($owner:String!,$repository:String!,$number:Int!,$cursor:String){repository(owner:$owner,name:$repository){pullRequest(number:$number){timelineItems(first:100,after:$cursor,itemTypes:[PULL_REQUEST_COMMIT,HEAD_REF_FORCE_PUSHED_EVENT]){pageInfo{hasNextPage endCursor}nodes{__typename ... on PullRequestCommit{commit{oid}} ... on HeadRefForcePushedEvent{beforeCommit{oid}afterCommit{oid}}}}}}}`
	shas := map[string]struct{}{current: {}}
	var cursor any
	for page := 0; ; page++ {
		if page >= maximumAPIPages {
			return nil, fmt.Errorf("pull request history pagination exceeded %d pages", maximumAPIPages)
		}
		var response struct {
			Repository *struct {
				PullRequest *struct {
					Timeline struct {
						PageInfo graphQLPageInfo `json:"pageInfo"`
						Nodes    []struct {
							Type   string `json:"__typename"`
							Commit *struct {
								OID string `json:"oid"`
							} `json:"commit"`
							Before *struct {
								OID string `json:"oid"`
							} `json:"beforeCommit"`
							After *struct {
								OID string `json:"oid"`
							} `json:"afterCommit"`
						} `json:"nodes"`
					} `json:"timelineItems"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		variables := map[string]any{"owner": c.owner, "repository": c.repository, "number": number, "cursor": cursor}
		if err := c.graphQL(ctx, query, variables, &response); err != nil {
			return nil, err
		}
		if response.Repository == nil || response.Repository.PullRequest == nil {
			return nil, errors.New("pull request history query did not return the pull request")
		}
		connection := response.Repository.PullRequest.Timeline
		for _, node := range connection.Nodes {
			for _, commit := range []*struct {
				OID string `json:"oid"`
			}{node.Commit, node.Before, node.After} {
				if commit != nil && fullSHA.MatchString(commit.OID) {
					shas[commit.OID] = struct{}{}
				}
			}
		}
		if !connection.PageInfo.HasNextPage {
			break
		}
		if connection.PageInfo.EndCursor == "" {
			return nil, errors.New("pull request history query omitted its next cursor")
		}
		cursor = connection.PageInfo.EndCursor
	}
	result := make([]string, 0, len(shas))
	for sha := range shas {
		result = append(result, sha)
	}
	sort.Strings(result)
	return result, nil
}
