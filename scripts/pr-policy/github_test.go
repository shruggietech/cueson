package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowPreservesTrustedMutationBoundary(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "pr-policy.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(contents)
	for _, required := range []string{
		"pull_request_target:",
		"issue_comment:",
		"schedule:",
		"contents: read",
		"checks: read",
		"issues: write",
		"pull-requests: read",
		"statuses: write",
		"group: cueson-pr-policy",
		"cancel-in-progress: false",
		"ref: main",
		"persist-credentials: false",
		"cache: false",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("workflow is missing trusted-boundary fragment %q", required)
		}
	}
	for _, forbidden := range []string{
		"pull_request_review:",
		"pull_request_review_comment:",
		"workflow_dispatch:",
		"github.event.pull_request.head",
		"refs/pull/",
		"secrets.",
	} {
		if strings.Contains(workflow, forbidden) {
			t.Errorf("workflow contains forbidden fragment %q", forbidden)
		}
	}
}

func TestRESTPaginationStaysOnConfiguredOrigin(t *testing.T) {
	t.Parallel()

	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization = %q", got)
		}
		if got := r.Header.Get("X-GitHub-Api-Version"); got != githubAPIVersion {
			t.Fatalf("X-GitHub-Api-Version = %q", got)
		}
		switch r.URL.Query().Get("page") {
		case "", "1":
			w.Header().Set("Link", fmt.Sprintf(`<%s/items?page=2&per_page=100>; rel="next"`, serverURL(r)))
			_, _ = io.WriteString(w, `[{"id":1}]`)
		case "2":
			_, _ = io.WriteString(w, `[{"id":2}]`)
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	var items []struct {
		ID int64 `json:"id"`
	}
	if err := client.getRESTPages(context.Background(), "/items?per_page=100", &items); err != nil {
		t.Fatal(err)
	}
	if requests != 2 || len(items) != 2 || items[0].ID != 1 || items[1].ID != 2 {
		t.Fatalf("requests=%d items=%+v", requests, items)
	}
}

func TestRESTPaginationRejectsCrossOriginNextLink(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", `<https://attacker.invalid/steal>; rel="next"`)
		_, _ = io.WriteString(w, `[]`)
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	var items []json.RawMessage
	err := client.getRESTPages(context.Background(), "/items", &items)
	if err == nil || !strings.Contains(err.Error(), "different origin") {
		t.Fatalf("error = %v", err)
	}
}

func TestGraphQLRejectsPartialData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"data":{"repository":{}},"errors":[{"message":"forbidden"}]}`)
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	var result struct {
		Repository struct{} `json:"repository"`
	}
	err := client.graphQL(context.Background(), "query { viewer { login } }", nil, &result)
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("error = %v", err)
	}
}

func TestClosingIssuesGraphQLPagination(t *testing.T) {
	t.Parallel()

	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		var request struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(request.Query, "excludeUserLinked:true") {
			t.Fatalf("query did not exclude user-linked issues: %s", request.Query)
		}
		if request.Variables["cursor"] == nil {
			_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"closingIssuesReferences":{"pageInfo":{"hasNextPage":true,"endCursor":"next"},"nodes":[{"number":7,"repository":{"nameWithOwner":"shruggietech/cueson"}}]}}}}}`)
			return
		}
		if request.Variables["cursor"] != "next" {
			t.Fatalf("cursor = %#v", request.Variables["cursor"])
		}
		_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"closingIssuesReferences":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"number":8,"repository":{"nameWithOwner":"shruggietech/cueson"}}]}}}}}`)
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	issues, err := client.listClosingIssues(context.Background(), 9)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 || len(issues) != 2 || issues[0].Number != 7 || issues[1].Number != 8 {
		t.Fatalf("requests=%d issues=%+v", requests, issues)
	}
}

func TestCollectSnapshotJoinsCurrentReviewEvidence(t *testing.T) {
	t.Parallel()

	const (
		head = "0123456789abcdef0123456789abcdef01234567"
		old  = "fedcba9876543210fedcba9876543210fedcba98"
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/pulls/9"):
			_, _ = io.WriteString(w, `{"number":9,"state":"open","draft":false,"updated_at":"2026-09-09T12:00:00Z","user":{"login":"contributor","id":1,"type":"User"},"head":{"ref":"topic","sha":"`+head+`"},"base":{"ref":"main","sha":"`+strings.Repeat("a", 40)+`"}}`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/issues/9/comments"):
			body := "<!-- codex-pull-request-review-summary -->\n\n| Review | Status | Commit |\n| - | - | - |\n| Code Review | **Completed** | `0123456` |"
			_, _ = io.WriteString(w, `[{"id":77,"body":`+mustJSON(t, body)+`,"created_at":"2026-09-09T12:01:00Z","updated_at":"2026-09-09T12:02:00Z","user":{"login":"chatgpt-codex-connector[bot]","id":199175422,"type":"Bot"}}]`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/pulls/9/reviews"):
			_, _ = io.WriteString(w, `[{"id":8,"node_id":"PRR_example","body":"finding","commit_id":"`+old+`","submitted_at":"2026-09-09T12:01:00Z","user":{"login":"chatgpt-codex-connector[bot]","id":199175422,"type":"Bot"}}]`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/issues/9/reactions"):
			_, _ = io.WriteString(w, `[{"id":1,"content":"eyes","created_at":"2026-09-09T12:01:00Z","user":{"login":"chatgpt-codex-connector[bot]","id":199175422,"type":"User"}}]`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/issues/comments/77/reactions"):
			_, _ = io.WriteString(w, `[{"id":2,"content":"+1","created_at":"2026-09-09T12:03:00Z","user":{"login":"chatgpt-codex-connector[bot]","id":199175422,"type":"User"}}]`)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/check-runs"):
			_, _ = io.WriteString(w, `{"check_runs":[{"id":3,"name":"Formatting","head_sha":"`+head+`","status":"completed","conclusion":"success","started_at":"2026-09-09T12:00:00Z","completed_at":"2026-09-09T12:01:00Z"}]}`)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/statuses"):
			_, _ = io.WriteString(w, `[]`)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/compare/"):
			_, _ = io.WriteString(w, `{"status":"ahead","ahead_by":1}`)
		case r.Method == http.MethodPost && r.URL.Path == "/graphql":
			var request struct {
				Query string `json:"query"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			switch {
			case strings.Contains(request.Query, "closingIssuesReferences"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"closingIssuesReferences":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"number":9,"repository":{"nameWithOwner":"shruggietech/cueson"}}]}}}}}`)
			case strings.Contains(request.Query, "reviewThreads"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"reviewThreads":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"id":"PRRT_example","isResolved":true,"isOutdated":true,"comments":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"author":{"__typename":"Bot","login":"chatgpt-codex-connector","databaseId":199175422},"createdAt":"2026-09-09T12:01:00Z","pullRequestReview":{"id":"PRR_example","commit":{"oid":"`+old+`"}}}]}}]}}}}}`)
			case strings.Contains(request.Query, "timelineItems"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"timelineItems":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"__typename":"PullRequestCommit","commit":{"oid":"`+old+`"}}]}}}}}`)
			default:
				t.Fatalf("unexpected GraphQL query: %s", request.Query)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	snapshot, err := client.collectSnapshot(context.Background(), 9)
	if err != nil {
		t.Fatal(err)
	}
	if !snapshot.Complete || len(snapshot.ClosingIssues) != 1 || len(snapshot.Reviews) != 1 || len(snapshot.ReviewThreads) != 1 || len(snapshot.Comparisons) != 1 {
		t.Fatalf("incomplete joined snapshot: %+v", snapshot)
	}
	if snapshot.Reviews[0].ID != "PRR_example" || snapshot.ReviewThreads[0].ReviewID != "PRR_example" || snapshot.Reviews[0].UnthreadedFinding {
		t.Fatalf("review/thread join failed: reviews=%+v threads=%+v", snapshot.Reviews, snapshot.ReviewThreads)
	}
	var pullReaction, summaryReaction bool
	for _, reaction := range snapshot.Reactions {
		pullReaction = pullReaction || reaction.Target == "pull_request" && reaction.Content == "eyes"
		summaryReaction = summaryReaction || reaction.Target == "comment:77" && reaction.Content == "+1"
	}
	if !pullReaction || !summaryReaction {
		t.Fatalf("reaction targets were not preserved: %+v", snapshot.Reactions)
	}
	for _, reaction := range snapshot.Reactions {
		if !IsCodexActor(reaction.Author) {
			t.Fatalf("empirical REST reaction actor was not narrowly normalized: %+v", reaction.Author)
		}
	}
}

func TestExecuteRequiresDryRunOutsideGitHubActions(t *testing.T) {
	t.Parallel()

	getenv := func(string) string { return "" }
	_, err := execute(context.Background(), []string{"-repo", "shruggietech/cueson", "-pr", "9"}, getenv)
	if err == nil || !strings.Contains(err.Error(), "use -dry-run locally") {
		t.Fatalf("error = %v", err)
	}
}

func TestPostStatusIsNoOpWhenLatestMatches(t *testing.T) {
	t.Parallel()

	const sha = "0123456789abcdef0123456789abcdef01234567"
	var posts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/statuses"):
			_, _ = io.WriteString(w, `[{"id":7,"state":"success","context":"Cueson PR policy / Issue link","description":"linked to #9","target_url":"https://github.example/run/2","creator":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}]`)
		case r.Method == http.MethodPost:
			posts++
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	changed, err := client.publishStatus(context.Background(), sha, commitStatus{
		Context:     IssueLinkContext,
		State:       "success",
		Description: "linked to #9",
		TargetURL:   "https://github.example/run/2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if changed || posts != 0 {
		t.Fatalf("changed=%v posts=%d", changed, posts)
	}
}

func TestPostStatusReadsMutationBack(t *testing.T) {
	t.Parallel()

	const sha = "0123456789abcdef0123456789abcdef01234567"
	var posted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/statuses"):
			if !posted {
				_, _ = io.WriteString(w, `[]`)
				return
			}
			_, _ = io.WriteString(w, `[{"id":8,"state":"failure","context":"Cueson PR policy / Issue link","description":"missing closing reference","target_url":"https://github.example/run/2","creator":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}]`)
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/statuses/"):
			posted = true
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":8,"state":"failure","context":"Cueson PR policy / Issue link","description":"missing closing reference","target_url":"https://github.example/run/2","creator":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	changed, err := client.publishStatus(context.Background(), sha, commitStatus{
		Context:     IssueLinkContext,
		State:       "failure",
		Description: "missing closing reference",
		TargetURL:   "https://github.example/run/2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !changed || !posted {
		t.Fatalf("changed=%v posted=%v", changed, posted)
	}
}

func TestPostStatusReplacesMatchingUntrustedStatus(t *testing.T) {
	t.Parallel()

	const sha = "0123456789abcdef0123456789abcdef01234567"
	var posted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/statuses"):
			if !posted {
				_, _ = io.WriteString(w, `[{"id":7,"state":"success","context":"Cueson PR policy / Issue link","description":"linked to #9","target_url":"https://github.example/run/2","creator":{"login":"mallory","id":9001,"type":"User"}}]`)
				return
			}
			_, _ = io.WriteString(w, `[{"id":8,"state":"success","context":"Cueson PR policy / Issue link","description":"linked to #9","target_url":"https://github.example/run/2","creator":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}]`)
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/statuses/"):
			posted = true
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":8,"state":"success","context":"Cueson PR policy / Issue link","description":"linked to #9","target_url":"https://github.example/run/2","creator":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	changed, err := client.publishStatus(context.Background(), sha, commitStatus{
		Context:     IssueLinkContext,
		State:       "success",
		Description: "linked to #9",
		TargetURL:   "https://github.example/run/2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !changed || !posted {
		t.Fatalf("changed=%v posted=%v", changed, posted)
	}
}

func TestCreateReviewRequestReadsExactCommentBack(t *testing.T) {
	t.Parallel()

	const body = "@codex review\n\n<!-- cueson-codex-review-request:v1 round=2 pr=9 head=0123456789abcdef0123456789abcdef01234567 -->"
	var created bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/issues/9/comments"):
			var request struct {
				Body string `json:"body"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			if request.Body != body {
				t.Fatalf("body = %q", request.Body)
			}
			created = true
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":41,"body":`+mustJSON(t, body)+`,"user":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/issues/comments/41"):
			if !created {
				t.Fatal("read-back occurred before create")
			}
			_, _ = io.WriteString(w, `{"id":41,"body":`+mustJSON(t, body)+`,"user":{"login":"github-actions[bot]","id":41898282,"type":"Bot"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	comment, err := client.createReviewRequest(context.Background(), 9, body)
	if err != nil {
		t.Fatal(err)
	}
	if comment.ID != 41 || comment.Body != body {
		t.Fatalf("comment = %+v", comment)
	}
}

func TestMutationRequiresMainBaseAndStableHead(t *testing.T) {
	t.Parallel()

	client := newTestGitHubClient(t, "https://api.github.invalid", "test-token")
	const sha = "0123456789abcdef0123456789abcdef01234567"
	for _, test := range []struct {
		name      string
		pr        githubPullRequest
		wantError string
	}{
		{name: "wrong base", pr: githubPullRequest{Number: 9, State: "open", Base: githubRef{Ref: "release"}, Head: githubRef{SHA: sha}}, wantError: "base branch"},
		{name: "wrong head", pr: githubPullRequest{Number: 9, State: "open", Base: githubRef{Ref: "main"}, Head: githubRef{SHA: strings.Repeat("f", 40)}}, wantError: "head changed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := client.requireMutationTarget(test.pr, 9, sha)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("error = %v, want %q", err, test.wantError)
			}
		})
	}
}

func TestReservationRefetchRevalidatesSecondReviewEligibility(t *testing.T) {
	t.Parallel()

	snapshot := reviewScenario(t, "round-two-ready")
	description := RoundTwoReservationDescription(snapshot.Number, snapshot.HeadSHA)
	snapshot.Checks = append(snapshot.Checks, CheckResult{
		Context:     CodexReviewContext,
		State:       string(StatusPending),
		Description: description,
		SHA:         snapshot.HeadSHA,
		Creator:     Actor{Login: GitHubActionsLogin, ID: GitHubActionsUserID, Type: "Bot"},
	})

	withoutReservation, err := withoutExpectedReservation(snapshot, description, snapshot.HeadSHA)
	if err != nil {
		t.Fatal(err)
	}
	result := EvaluateCodexReview(withoutReservation)
	if !result.RequestSecond || result.Description != description {
		t.Fatalf("unchanged evidence lost eligibility: %+v", result)
	}

	snapshot.ReviewThreads[0].Resolved = false
	withoutReservation, err = withoutExpectedReservation(snapshot, description, snapshot.HeadSHA)
	if err != nil {
		t.Fatal(err)
	}
	result = EvaluateCodexReview(withoutReservation)
	if result.RequestSecond || result.State != StatusFailure {
		t.Fatalf("changed evidence remained eligible: %+v", result)
	}
}

func TestReviewUsesGraphQLNodeIDForThreadJoin(t *testing.T) {
	t.Parallel()

	review, err := reviewFromGitHub(githubReview{
		ID:       123,
		NodeID:   "PRR_kwDOExample",
		CommitID: "0123456789abcdef0123456789abcdef01234567",
		User:     githubActor{Login: CodexRESTLogin, ID: CodexUserID, Type: "Bot"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if review.ID != "PRR_kwDOExample" {
		t.Fatalf("review ID = %q", review.ID)
	}
	if _, err := reviewFromGitHub(githubReview{ID: 123}); err == nil {
		t.Fatal("missing node ID was accepted")
	}
}

func TestSnapshotHeadRaceRetriesAndFailsClosed(t *testing.T) {
	t.Parallel()

	const first = "0123456789abcdef0123456789abcdef01234567"
	const second = "fedcba9876543210fedcba9876543210fedcba98"
	var pullReads int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/pulls/9"):
			pullReads++
			head := first
			if pullReads%2 == 0 {
				head = second
			}
			_, _ = fmt.Fprintf(w, `{"number":9,"state":"open","draft":false,"user":{"login":"contributor","id":1,"type":"User"},"head":{"ref":"topic","sha":"%s"},"base":{"ref":"main","sha":"%s"}}`, head, strings.Repeat("a", 40))
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/check-runs"):
			_, _ = io.WriteString(w, `{"check_runs":[]}`)
		case r.Method == http.MethodGet:
			_, _ = io.WriteString(w, `[]`)
		case r.Method == http.MethodPost && r.URL.Path == "/graphql":
			var request struct {
				Query string `json:"query"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			switch {
			case strings.Contains(request.Query, "closingIssuesReferences"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"closingIssuesReferences":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[]}}}}}`)
			case strings.Contains(request.Query, "reviewThreads"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"reviewThreads":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[]}}}}}`)
			case strings.Contains(request.Query, "timelineItems"):
				_, _ = io.WriteString(w, `{"data":{"repository":{"pullRequest":{"timelineItems":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[]}}}}}`)
			default:
				t.Fatalf("unexpected GraphQL query: %s", request.Query)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestGitHubClient(t, server.URL, "test-token")
	_, err := client.collectSnapshot(context.Background(), 9)
	if err == nil || !strings.Contains(err.Error(), "did not stabilize") {
		t.Fatalf("error = %v", err)
	}
	if pullReads != 2*maximumSnapshotTry {
		t.Fatalf("pull reads = %d, want %d", pullReads, 2*maximumSnapshotTry)
	}
}

func newTestGitHubClient(t *testing.T, endpoint, token string) *githubClient {
	t.Helper()
	base, err := url.Parse(endpoint + "/")
	if err != nil {
		t.Fatal(err)
	}
	return &githubClient{
		httpClient: http.DefaultClient,
		restBase:   base,
		graphqlURL: endpoint + "/graphql",
		token:      token,
		owner:      "shruggietech",
		repository: "cueson",
	}
}

func serverURL(r *http.Request) string {
	return "http://" + r.Host
}

func mustJSON(t *testing.T, value string) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
