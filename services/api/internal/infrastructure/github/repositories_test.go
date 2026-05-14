package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

func TestRepositoryClientListRepositories(t *testing.T) {
	t.Run("pagination を辿って repository を正規化する", func(t *testing.T) {
		var requests int
		httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			requests++
			if r.Header.Get("Authorization") != "Bearer github-token" {
				t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
			}
			if r.Header.Get("User-Agent") != userAgent {
				t.Fatalf("User-Agent = %q", r.Header.Get("User-Agent"))
			}

			if requests == 1 {
				return responseWithBody(http.StatusOK, map[string]string{
					"Link": `<https://api.github.test/user/repos?page=2>; rel="next"`,
				}, `[{
					"id": 1,
					"name": "private-repo",
					"full_name": "octocat/private-repo",
					"private": true,
					"fork": false,
					"archived": false,
					"default_branch": "main",
					"language": "Go",
					"html_url": "https://github.com/octocat/private-repo",
					"pushed_at": "2026-05-14T10:00:00Z",
					"updated_at": "2026-05-14T11:00:00Z",
					"owner": {"login": "octocat"}
				}]`), nil
			}

			return responseWithBody(http.StatusOK, nil, `[]`), nil
		})}

		client := NewRepositoryClient(httpClient)
		client.baseURL = "https://api.github.test"
		syncedAt := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
		repositories, err := client.ListRepositories(context.Background(), "github-token", user.ID("user_1"), syncedAt)
		if err != nil {
			t.Fatalf("ListRepositories() がエラーを返しました: %v", err)
		}

		if requests != 2 {
			t.Fatalf("requests = %d, 期待値 2", requests)
		}
		if len(repositories) != 1 {
			t.Fatalf("repositories length = %d, 期待値 1", len(repositories))
		}
		if !repositories[0].Private {
			t.Fatal("private repository が private=true で正規化される必要があります")
		}
		if repositories[0].Language != "Go" {
			t.Fatalf("Language = %q, 期待値 Go", repositories[0].Language)
		}
	})
}

func TestRepositoryClientErrors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		header map[string]string
		body   string
		kind   githubapp.ErrorKind
	}{
		{
			name:   "primary rate limit",
			status: http.StatusForbidden,
			header: map[string]string{"X-RateLimit-Remaining": "0", "X-RateLimit-Reset": "1770000000"},
			body:   `{"message":"API rate limit exceeded"}`,
			kind:   githubapp.ErrorKindRateLimited,
		},
		{
			name:   "secondary rate limit",
			status: http.StatusTooManyRequests,
			header: map[string]string{"Retry-After": "30"},
			body:   `{"message":"You have exceeded a secondary rate limit"}`,
			kind:   githubapp.ErrorKindRateLimited,
		},
		{
			name:   "too many requests without rate limit headers",
			status: http.StatusTooManyRequests,
			body:   `{"message":"Too Many Requests"}`,
			kind:   githubapp.ErrorKindRateLimited,
		},
		{name: "token invalid", status: http.StatusUnauthorized, body: `{"message":"Bad credentials"}`, kind: githubapp.ErrorKindTokenInvalid},
		{name: "permission denied", status: http.StatusForbidden, body: `{"message":"Resource not accessible"}`, kind: githubapp.ErrorKindPermissionDenied},
		{name: "not found", status: http.StatusNotFound, body: `{"message":"Not Found"}`, kind: githubapp.ErrorKindPermissionDeniedOrNotFound},
		{name: "unavailable", status: http.StatusBadGateway, body: `{"message":"Server Error"}`, kind: githubapp.ErrorKindUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRepositoryClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return responseWithBody(tt.status, tt.header, tt.body), nil
			})})
			client.baseURL = "https://api.github.test"
			_, err := client.ListRepositories(context.Background(), "github-token", "user_1", time.Now())
			var apiErr *githubapp.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v, 期待値 APIError", err)
			}
			if apiErr.Kind != tt.kind {
				t.Fatalf("Kind = %q, 期待値 %q", apiErr.Kind, tt.kind)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func responseWithBody(statusCode int, headers map[string]string, body string) *http.Response {
	header := http.Header{}
	for key, value := range headers {
		header.Set(key, value)
	}
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
