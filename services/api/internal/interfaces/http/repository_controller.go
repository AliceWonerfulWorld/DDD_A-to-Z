package http

import (
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
)

type RepositoryController struct {
	usecase *githubapp.UseCase
	logger  *slog.Logger
}

func NewRepositoryController(usecase *githubapp.UseCase, logger *slog.Logger) *RepositoryController {
	return &RepositoryController{
		usecase: usecase,
		logger:  logger,
	}
}

func (c *RepositoryController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("POST /github/repositories/sync", c.syncRepositories)
	mux.HandleFunc("GET /repositories", c.listRepositories)
}

func (c *RepositoryController) syncRepositories(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, githubapp.ErrUnauthenticated)
		return
	}

	result, err := c.usecase.SyncRepositories(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"synced_count": result.SyncedCount,
		"repositories": repositoryResponses(result.Repositories),
	}); err != nil {
		c.logger.Error("failed to write repository sync response", "error", err)
	}
}

func (c *RepositoryController) listRepositories(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, githubapp.ErrUnauthenticated)
		return
	}

	repositories, err := c.usecase.ListRepositories(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"repositories": repositoryResponses(repositories),
	}); err != nil {
		c.logger.Error("failed to write repository list response", "error", err)
	}
}

func (c *RepositoryController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, githubapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, githubapp.ErrMissingGitHubToken):
		writeAPIError(w, stdhttp.StatusUnauthorized, "github_token_invalid", "GitHub access token is missing", 0, nil)
	default:
		var apiErr *githubapp.APIError
		if errors.As(err, &apiErr) {
			c.writeGitHubAPIError(w, apiErr)
			return
		}

		c.logger.Error("repository request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func (c *RepositoryController) writeGitHubAPIError(w stdhttp.ResponseWriter, err *githubapp.APIError) {
	switch err.Kind {
	case githubapp.ErrorKindRateLimited:
		writeAPIError(w, stdhttp.StatusTooManyRequests, string(err.Kind), err.Error(), int64(err.RetryAfter.Seconds()), err.ResetAt)
	case githubapp.ErrorKindTokenInvalid:
		writeAPIError(w, stdhttp.StatusUnauthorized, string(err.Kind), err.Error(), 0, nil)
	case githubapp.ErrorKindPermissionDenied, githubapp.ErrorKindPermissionDeniedOrNotFound:
		writeAPIError(w, stdhttp.StatusForbidden, string(err.Kind), err.Error(), 0, nil)
	case githubapp.ErrorKindUnavailable:
		writeAPIError(w, stdhttp.StatusBadGateway, string(err.Kind), err.Error(), 0, nil)
	default:
		c.logger.Error("unknown github api error", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func writeAPIError(w stdhttp.ResponseWriter, statusCode int, code, message string, retryAfterSeconds int64, resetAt *time.Time) {
	body := map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
	errorBody := body["error"].(map[string]any)
	if retryAfterSeconds > 0 {
		errorBody["retry_after_seconds"] = retryAfterSeconds
	}
	if resetAt != nil {
		errorBody["reset_at"] = resetAt.Format(time.RFC3339)
	}

	_ = writeJSON(w, statusCode, body)
}

func repositoryResponses(repositories []repositoryanalysis.Repository) []map[string]any {
	responses := make([]map[string]any, 0, len(repositories))
	for _, repository := range repositories {
		responses = append(responses, map[string]any{
			"github_id":         repository.GitHubID,
			"owner":             repository.Owner,
			"name":              repository.Name,
			"full_name":         repository.FullName,
			"private":           repository.Private,
			"fork":              repository.Fork,
			"archived":          repository.Archived,
			"default_branch":    repository.DefaultBranch,
			"language":          repository.Language,
			"html_url":          repository.HTMLURL,
			"pushed_at":         optionalTime(repository.PushedAt),
			"github_updated_at": repository.GitHubUpdatedAt.Format(time.RFC3339),
			"synced_at":         repository.SyncedAt.Format(time.RFC3339),
		})
	}

	return responses
}

func optionalTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.Format(time.RFC3339)
}
