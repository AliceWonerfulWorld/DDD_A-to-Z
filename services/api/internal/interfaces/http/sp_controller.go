package http

import (
	"errors"
	"log/slog"
	stdhttp "net/http"

	spapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/sp"
)

// SPController handles GET /me/sp.
type SPController struct {
	usecase *spapp.UseCase
	logger  *slog.Logger
}

// NewSPController creates a new SPController.
func NewSPController(usecase *spapp.UseCase, logger *slog.Logger) *SPController {
	return &SPController{usecase: usecase, logger: logger}
}

// RegisterRoutes registers the /me/sp route.
func (c *SPController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /me/sp", c.getLoginUserSP)
}

func (c *SPController) getLoginUserSP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, spapp.ErrUnauthenticated)
		return
	}

	balances, err := c.usecase.GetSP(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(balances))
	for _, b := range balances {
		items = append(items, map[string]any{
			"language": b.Language,
			"balance":  b.Balance,
		})
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{"skill_points": items}); err != nil {
		c.logger.Error("failed to write sp response", "error", err)
	}
}

func (c *SPController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, spapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	default:
		c.logger.Error("sp request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}
