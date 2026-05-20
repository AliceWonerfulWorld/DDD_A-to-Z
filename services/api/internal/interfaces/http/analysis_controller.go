package http

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	analysisapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/repositoryanalysis"
)

type Analyzer interface {
	Analyze(ctx context.Context, sessionToken string) (analysisapp.AnalysisResult, error)
}

type AnalysisController struct {
	usecase Analyzer
	logger  *slog.Logger
}

func NewAnalysisController(usecase Analyzer, logger *slog.Logger) *AnalysisController {
	return &AnalysisController{
		usecase: usecase,
		logger:  logger,
	}
}

func (c *AnalysisController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("POST /analysis/contribution", c.analyzeContribution)
}

func (c *AnalysisController) analyzeContribution(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	c.logger.Info("analysis request received")

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	result, err := c.usecase.Analyze(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	contributions := make([]map[string]any, 0, len(result.Contributions))
	for _, c := range result.Contributions {
		contributions = append(contributions, map[string]any{
			"repo":        c.Repo,
			"type":        c.Type,
			"external_id": c.ExternalID,
			"message":     c.Message,
			"language":    c.Language,
			"cp":          c.CP,
			"timestamp":   c.Timestamp.Format(time.RFC3339),
		})
	}

	breakdown := make([]map[string]any, 0, len(result.LanguageBreakdown))
	for _, lb := range result.LanguageBreakdown {
		breakdown = append(breakdown, map[string]any{
			"name": lb.Name,
			"cp":   lb.CP,
		})
	}

	c.logger.Info("analysis response",
		"totalCommits", result.TotalCommits,
		"totalPRs", result.TotalPRs,
		"totalCP", result.TotalCP,
		"totalBalance", result.TotalBalance,
		"contributions", len(contributions),
	)

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"totalCommits":      result.TotalCommits,
		"totalPRs":          result.TotalPRs,
		"totalCP":           result.TotalCP,
		"totalBalance":      result.TotalBalance,
		"languageBreakdown": breakdown,
		"contributions":     contributions,
	}); err != nil {
		c.logger.Error("failed to write analysis response", "error", err)
	}
}

func (c *AnalysisController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, analysisapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, analysisapp.ErrMissingGitHubToken):
		writeAPIError(w, stdhttp.StatusUnauthorized, "github_token_invalid", "GitHub access token is missing", 0, nil)
	default:
		var apiErr *githubapp.APIError
		if errors.As(err, &apiErr) {
			c.writeGitHubError(w, apiErr)
			return
		}

		c.logger.Error("analysis request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func (c *AnalysisController) writeGitHubError(w stdhttp.ResponseWriter, err *githubapp.APIError) {
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
