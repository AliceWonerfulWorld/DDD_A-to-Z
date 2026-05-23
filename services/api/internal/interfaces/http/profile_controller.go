package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	stdhttp "net/http"

	profileapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/profile"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// プロフィールコントローラー本体
type SelectedBadgeSlugUpdater interface {
	UpdateSelectedBadgeSlug(ctx context.Context, userID user.ID, badgeSlug *string) error
}

type ProfileController struct {
	usecase *profileapp.UseCase
	badges  SelectedBadgeSlugUpdater
	logger  *slog.Logger
}

func NewProfileController(usecase *profileapp.UseCase, badges SelectedBadgeSlugUpdater, logger *slog.Logger) *ProfileController {
	return &ProfileController{
		usecase: usecase,
		badges:  badges,
		logger:  logger,
	}
}

// プロフィール関連のエンドポイントを登録
func (c *ProfileController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("POST /profile/complete", c.completeInitialProfile)
	mux.HandleFunc("GET /profile", c.getProfile)
	mux.HandleFunc("PATCH /profile/badge", c.updateBadge)
}

// bodyの構造体
type completeInitialProfileRequest struct {
	DisplayName string `json:"display_name"`
}

// プロフィール初期設定
func (c *ProfileController) completeInitialProfile(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	var req completeInitialProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_request", "invalid json body", 0, nil)
		return
	}

	input := profileapp.CompleteInitialProfileInput{
		SessionToken: cookie.Value,
		DisplayName:  req.DisplayName,
	}

	if err := c.usecase.CompleteInitialProfile(r.Context(), input); err != nil {
		switch {
		case errors.Is(err, profileapp.ErrUnauthenticated):
			writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		case errors.Is(err, profileapp.ErrProfileExists):
			// 既に完了している場合は成功(冪等)とみなすことも、競合とみなすことも可能。
			// 200 OKまたは409 Conflictを返す。
			// フロントエンドに通知するために409を返すことにする。
			writeAPIError(w, stdhttp.StatusConflict, "profile_exists", "profile already completed", 0, nil)
		case errors.Is(err, profileapp.ErrInvalidDisplayName):
			writeAPIError(w, stdhttp.StatusBadRequest, "invalid_display_name", err.Error(), 0, nil)
		default:
			c.logger.Error("failed to complete initial profile", "error", err)
			writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		}
		return
	}

	w.WriteHeader(stdhttp.StatusOK)
}

func (c *ProfileController) getProfile(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	appUser, ok, err := c.usecase.FindUser(r.Context(), cookie.Value)
	if err != nil || !ok {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	profile, exists, err := c.usecase.GetProfile(r.Context(), appUser.ID)
	if err != nil {
		c.logger.Error("failed to get profile", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if !exists {
		writeAPIError(w, stdhttp.StatusNotFound, "not_found", "profile not found", 0, nil)
		return
	}

	resp := map[string]any{
		"display_name":        profile.DisplayName,
		"selected_badge_slug": profile.SelectedBadgeSlug,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c.logger.Error("failed to write profile response", "error", err)
	}
}

type updateBadgeRequest struct {
	BadgeSlug *string `json:"badge_slug"`
}

func (c *ProfileController) updateBadge(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	appUser, ok, err := c.usecase.FindUser(r.Context(), cookie.Value)
	if err != nil || !ok {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	var req updateBadgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_request", "invalid json body", 0, nil)
		return
	}

	if err := c.badges.UpdateSelectedBadgeSlug(r.Context(), appUser.ID, req.BadgeSlug); err != nil {
		c.logger.Error("failed to update badge", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	w.WriteHeader(stdhttp.StatusNoContent)
}
