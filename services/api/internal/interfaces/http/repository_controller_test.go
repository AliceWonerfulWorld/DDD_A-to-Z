package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type repositoryTestCurrentUser struct {
	appUser user.User
	ok      bool
}

func (r repositoryTestCurrentUser) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return r.appUser, r.ok, nil
}

type repositoryTestTokenStore struct {
	token string
	ok    bool
}

func (s repositoryTestTokenStore) GitHubAccessToken(ctx context.Context, userID user.ID) (string, bool, error) {
	return s.token, s.ok, nil
}

type repositoryTestClient struct {
	repositories []repositoryanalysis.Repository
	err          error
}

func (c repositoryTestClient) ListRepositories(ctx context.Context, accessToken string, userID user.ID, syncedAt time.Time) ([]repositoryanalysis.Repository, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.repositories, nil
}

type repositoryTestStore struct {
	repositories []repositoryanalysis.Repository
}

func (s *repositoryTestStore) UpsertRepositories(ctx context.Context, repositories []repositoryanalysis.Repository) error {
	s.repositories = repositories
	return nil
}

func (s *repositoryTestStore) ListRepositories(ctx context.Context, userID user.ID) ([]repositoryanalysis.Repository, error) {
	return s.repositories, nil
}

func TestRepositoryControllerSyncRepositories(t *testing.T) {
	syncedAt := time.Date(2026, 5, 14, 16, 0, 0, 0, time.UTC)
	repository := repositoryanalysis.Repository{
		GitHubID:        1,
		UserID:          "user_1",
		Owner:           "octocat",
		Name:            "hello-world",
		FullName:        "octocat/hello-world",
		DefaultBranch:   "main",
		Language:        "Go",
		HTMLURL:         "https://github.com/octocat/hello-world",
		GitHubUpdatedAt: syncedAt,
		SyncedAt:        syncedAt,
	}
	store := &repositoryTestStore{}
	controller := newTestRepositoryController(repositoryTestClient{repositories: []repositoryanalysis.Repository{repository}}, store)
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/github/repositories/sync", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		SyncedCount  int `json:"synced_count"`
		Repositories []struct {
			FullName string `json:"full_name"`
			Language string `json:"language"`
		} `json:"repositories"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.SyncedCount != 1 {
		t.Fatalf("synced_count = %d, 期待値 1", body.SyncedCount)
	}
	if body.Repositories[0].FullName != "octocat/hello-world" {
		t.Fatalf("full_name = %q", body.Repositories[0].FullName)
	}
}

func TestRepositoryControllerErrorResponse(t *testing.T) {
	resetAt := time.Date(2026, 5, 14, 17, 0, 0, 0, time.UTC)
	controller := newTestRepositoryController(repositoryTestClient{
		err: &githubapp.APIError{
			Kind:       githubapp.ErrorKindRateLimited,
			Message:    "GitHub API rate limit exceeded",
			RetryAfter: 30 * time.Second,
			ResetAt:    &resetAt,
		},
	}, &repositoryTestStore{})
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/github/repositories/sync", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusTooManyRequests {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusTooManyRequests)
	}

	var body struct {
		Error struct {
			Code              string `json:"code"`
			RetryAfterSeconds int64  `json:"retry_after_seconds"`
			ResetAt           string `json:"reset_at"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.Error.Code != "github_rate_limited" {
		t.Fatalf("error.code = %q", body.Error.Code)
	}
	if body.Error.RetryAfterSeconds != 30 {
		t.Fatalf("retry_after_seconds = %d, 期待値 30", body.Error.RetryAfterSeconds)
	}
	if body.Error.ResetAt != resetAt.Format(time.RFC3339) {
		t.Fatalf("reset_at = %q", body.Error.ResetAt)
	}
}

func newTestRepositoryController(client repositoryTestClient, store *repositoryTestStore) *RepositoryController {
	usecase := githubapp.NewUseCase(
		repositoryTestCurrentUser{appUser: user.User{ID: "user_1"}, ok: true},
		repositoryTestTokenStore{token: "github-token", ok: true},
		client,
		store,
	)

	return NewRepositoryController(usecase, slog.New(slog.NewTextHandler(io.Discard, nil)))
}
