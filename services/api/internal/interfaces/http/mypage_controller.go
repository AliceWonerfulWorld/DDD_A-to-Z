package http

import (
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
)

// MypageController handles the GET /mypage endpoint.
type MypageController struct {
	usecase *mypageapp.UseCase
	logger  *slog.Logger
}

// NewMypageController creates a new MypageController.
func NewMypageController(usecase *mypageapp.UseCase, logger *slog.Logger) *MypageController {
	return &MypageController{
		usecase: usecase,
		logger:  logger,
	}
}

// RegisterRoutes registers the /mypage route.
func (c *MypageController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /mypage", c.getMyPage)
}

func (c *MypageController) getMyPage(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, mypageapp.ErrUnauthenticated)
		return
	}

	data, err := c.usecase.GetMyPage(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, c.toResponse(data)); err != nil {
		c.logger.Error("failed to write mypage response", "error", err)
	}
}

func (c *MypageController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, mypageapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	default:
		c.logger.Error("mypage request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func (c *MypageController) toResponse(data mypageapp.MyPageData) map[string]any {
	recent := make([]map[string]any, 0, len(data.Repositories.Recent))
	for _, r := range data.Repositories.Recent {
		repo := map[string]any{
			"github_id": r.GitHubID,
			"full_name": r.FullName,
			"language":  r.Language,
			"html_url":  r.HTMLURL,
		}
		if r.PushedAt != nil {
			repo["pushed_at"] = *r.PushedAt
		} else {
			repo["pushed_at"] = nil
		}
		recent = append(recent, repo)
	}

	langSummary := data.Repositories.LanguageSummary
	if langSummary == nil {
		langSummary = map[string]int{}
	}

	badges := make([]map[string]any, 0, len(data.Badges))
	for _, b := range data.Badges {
		badges = append(badges, map[string]any{
			"slug":        b.Slug,
			"name":        b.Name,
			"description": b.Description,
			"icon":        b.Icon,
			"earned_at":   b.EarnedAt,
		})
	}

	resp := map[string]any{
		"user": map[string]any{
			"id":         data.User.ID,
			"github_id":  data.User.GitHubAccount.GitHubID,
			"username":   data.User.GitHubAccount.Username,
			"avatar_url": data.User.GitHubAccount.AvatarURL,
			"created_at": data.User.CreatedAt.Format(time.RFC3339),
		},
		"contribution_points": map[string]any{
			"balance":      data.CP.Balance,
			"total_earned": data.CP.TotalEarned,
			"total_spent":  data.CP.TotalSpent,
		},
		"repositories": map[string]any{
			"total_count":      data.Repositories.TotalCount,
			"language_summary": langSummary,
			"recent":           recent,
		},
		"badges":              badges,
		"selected_badge_slug": data.SelectedBadgeSlug,
	}

	if data.GitHubStats != nil {
		resp["github_stats"] = map[string]any{
			"total_stars":          data.GitHubStats.TotalStars,
			"total_prs":            data.GitHubStats.TotalPRs,
			"total_issues":         data.GitHubStats.TotalIssues,
			"contributed_to":       data.GitHubStats.ContributedTo,
			"public_repos":         data.GitHubStats.PublicRepos,
			"github_created_at":    data.GitHubStats.GitHubCreatedAt,
			"yearly_commits":       data.GitHubStats.YearlyCommits,
			"yearly_contributions": data.GitHubStats.YearlyContributions,
		}
	} else {
		resp["github_stats"] = nil
	}

	if data.Guild != nil {
		resp["guild"] = map[string]any{
			"id":           data.Guild.ID,
			"name":         data.Guild.Name,
			"slug":         data.Guild.Slug,
			"icon":         data.Guild.Icon,
			"color":        data.Guild.Color,
			"description":  data.Guild.Description,
			"member_count": data.Guild.MemberCount,
			"rank":         data.Guild.Rank,
			"total_guilds": data.Guild.TotalGuilds,
			"cp":           data.Guild.CP,
		}
	} else {
		resp["guild"] = nil
	}

	return resp
}
