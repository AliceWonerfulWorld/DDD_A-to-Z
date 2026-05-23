package http

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	badgeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/badge"
	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	analysisapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/repositoryanalysis"
	badgedomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/badge"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type Analyzer interface {
	Analyze(ctx context.Context, sessionToken string) (analysisapp.AnalysisResult, error)
}

type BadgeGrantingChecker interface {
	CheckAndGrantBadges(ctx context.Context, userID user.ID, conditionType badgedomain.ConditionType, value int64) ([]badgeapp.GrantResult, error)
}

type TotalEarnedReader interface {
	GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
}

type AnalysisController struct {
	usecase      Analyzer
	badgeChecker BadgeGrantingChecker
	totalEarned  TotalEarnedReader
	session      interface {
		FindUserBySessionToken(ctx context.Context, token string, now time.Time) (user.User, bool, error)
	}
	logger *slog.Logger
}

func NewAnalysisController(
	usecase Analyzer,
	badgeChecker BadgeGrantingChecker,
	totalEarned TotalEarnedReader,
	session interface {
		FindUserBySessionToken(ctx context.Context, token string, now time.Time) (user.User, bool, error)
	},
	logger *slog.Logger,
) *AnalysisController {
	return &AnalysisController{
		usecase:      usecase,
		badgeChecker: badgeChecker,
		totalEarned:  totalEarned,
		session:      session,
		logger:       logger,
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

	if result.TotalCP > 0 {
		appUser, ok, userErr := c.session.FindUserBySessionToken(r.Context(), cookie.Value, time.Now())
		if userErr == nil && ok {
			totalEarned, earnedErr := c.totalEarned.GetTotalEarned(r.Context(), appUser.ID)
			if earnedErr == nil {
				grantResults, badgeErr := c.badgeChecker.CheckAndGrantBadges(r.Context(), appUser.ID, badgedomain.ConditionTypeCPEarned, totalEarned)
				if badgeErr != nil {
					c.logger.WarnContext(r.Context(), "failed to check/grant badges", "error", badgeErr, "user_id", appUser.ID)
				} else {
					for _, gr := range grantResults {
						if gr.JustEarned {
							c.logger.Info("badge earned", "badge", gr.Badge.Slug, "user_id", appUser.ID)
						}
					}
				}
			} else {
				c.logger.WarnContext(r.Context(), "failed to get total earned for badge check", "error", earnedErr, "user_id", appUser.ID)
			}
		} else if userErr != nil {
			c.logger.WarnContext(r.Context(), "failed to resolve user for badge check", "error", userErr)
		}
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
			"sp":   lb.SP,
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
