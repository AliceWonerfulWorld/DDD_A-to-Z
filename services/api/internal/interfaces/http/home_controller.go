package http

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"

	homeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/home"
)

type homeUseCase interface {
	GetHome(ctx context.Context, sessionToken string) (homeapp.HomeData, error)
}

type HomeController struct {
	uc     homeUseCase
	logger *slog.Logger
}

func NewHomeController(uc homeUseCase, logger *slog.Logger) *HomeController {
	return &HomeController{uc: uc, logger: logger}
}

func (c *HomeController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /home", c.getHome)
}

func (c *HomeController) getHome(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	data, err := c.uc.GetHome(r.Context(), cookie.Value)
	if errors.Is(err, homeapp.ErrUnauthenticated) {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}
	if err != nil {
		c.logger.Error("failed to get home data", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	resp := map[string]any{
		"total_cp":                    data.TotalCP,
		"today_cp":                    data.TodayCP,
		"player_level":                data.PlayerLevel,
		"player_level_total_cp":       data.PlayerLevelTotalCP,
		"next_player_level":           data.NextPlayerLevel,
		"next_player_level_total_cp":  data.NextPlayerLevelTotalCP,
		"next_player_level_remaining": data.NextPlayerLevelRemaining,
		"lifetime_total_earned_cp":    data.LifetimeTotalEarnedCP,
	}

	if err := writeJSON(w, stdhttp.StatusOK, resp); err != nil {
		c.logger.Error("failed to write home response", "error", err)
	}
}
