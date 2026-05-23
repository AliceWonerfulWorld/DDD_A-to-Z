package http

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
)

// PetController handles GET /pets/me.
type PetController struct {
	usecase *petapp.UseCase
	logger  *slog.Logger
}

// NewPetController creates a new PetController.
func NewPetController(usecase *petapp.UseCase, logger *slog.Logger) *PetController {
	return &PetController{usecase: usecase, logger: logger}
}

// RegisterRoutes registers the /pets/me route.
func (c *PetController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /pets/me", c.getMyPets)
}

func (c *PetController) getMyPets(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, petapp.ErrUnauthenticated)
		return
	}

	data, err := c.usecase.GetMyPets(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, c.toResponse(data)); err != nil {
		c.logger.Error("failed to write pets response", "error", err)
	}
}

func (c *PetController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, petapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, context.Canceled):
		return
	default:
		c.logger.Error("pets request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func (c *PetController) toResponse(data petapp.MyPetsData) map[string]any {
	pets := make([]map[string]any, 0, len(data.Pets))
	for _, p := range data.Pets {
		pets = append(pets, petToResponse(p))
	}

	var currentGuildPet any
	if data.CurrentGuildPet != nil {
		currentGuildPet = petToResponse(*data.CurrentGuildPet)
	}

	return map[string]any{
		"cpBalance":       data.CPBalance,
		"currentGuildPet": currentGuildPet,
		"pets":            pets,
	}
}

func petToResponse(p petapp.PetSummary) map[string]any {
	return map[string]any{
		"id":         p.ID,
		"guildId":    p.GuildID,
		"guildName":  p.GuildName,
		"name":       p.Name,
		"species":    p.Species,
		"attribute":  p.Attribute,
		"level":      p.Level,
		"exp":        p.Exp,
		"maxHp":      p.MaxHP,
		"power":      p.Power,
		"guard":      p.Guard,
		"speed":      p.Speed,
		"acquiredAt": p.AcquiredAt.Format(time.RFC3339),
	}
}
