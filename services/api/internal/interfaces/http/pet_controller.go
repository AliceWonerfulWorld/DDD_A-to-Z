package http

import (
	"context"
	"encoding/json"
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
	mux.HandleFunc("GET /pets/battle/opponents", c.listBattleOpponents)
	mux.HandleFunc("POST /pets/{petId}/train", c.trainPet)
	mux.HandleFunc("POST /pets/{petId}/battle", c.battlePet)
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

func (c *PetController) trainPet(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, petapp.ErrUnauthenticated)
		return
	}

	var request struct {
		Stat string `json:"stat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_request", "Invalid request body", 0, nil)
		return
	}

	result, err := c.usecase.TrainPet(r.Context(), petapp.TrainPetCommand{
		SessionToken: cookie.Value,
		PetID:        r.PathValue("petId"),
		Stat:         request.Stat,
	})
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"pet":       petToResponse(result.Pet),
		"spentCp":   result.SpentCP,
		"cpBalance": result.CPBalance,
	}); err != nil {
		c.logger.Error("failed to write pet training response", "error", err)
	}
}

func (c *PetController) listBattleOpponents(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, petapp.ErrUnauthenticated)
		return
	}

	data, err := c.usecase.ListBattleOpponents(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	opponents := make([]map[string]any, 0, len(data.Opponents))
	for _, opponent := range data.Opponents {
		opponents = append(opponents, battleOpponentToResponse(opponent))
	}
	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{"opponents": opponents}); err != nil {
		c.logger.Error("failed to write pet battle opponents response", "error", err)
	}
}

func (c *PetController) battlePet(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, petapp.ErrUnauthenticated)
		return
	}

	var request struct {
		OpponentPetID string `json:"opponentPetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_request", "Invalid request body", 0, nil)
		return
	}

	result, err := c.usecase.BattlePet(r.Context(), petapp.BattlePetCommand{
		SessionToken:  cookie.Value,
		PetID:         r.PathValue("petId"),
		OpponentPetID: request.OpponentPetID,
	})
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, battleResultToResponse(result)); err != nil {
		c.logger.Error("failed to write pet battle response", "error", err)
	}
}

func (c *PetController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, petapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, petapp.ErrInvalidTrainStat):
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_training_stat", "invalid_training_stat", 0, nil)
	case errors.Is(err, petapp.ErrInsufficientCP):
		writeAPIError(w, stdhttp.StatusConflict, "insufficient_cp", "CP balance is not enough to train this pet.", 0, nil)
	case errors.Is(err, petapp.ErrPetNotFound):
		writeAPIError(w, stdhttp.StatusNotFound, "pet_not_found", "pet_not_found", 0, nil)
	case errors.Is(err, petapp.ErrOpponentPetNotFound):
		writeAPIError(w, stdhttp.StatusNotFound, "opponent_pet_not_found", "opponent_pet_not_found", 0, nil)
	case errors.Is(err, petapp.ErrInvalidBattleTarget):
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_battle_target", "invalid_battle_target", 0, nil)
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

func battleOpponentToResponse(p petapp.OpponentSummary) map[string]any {
	return map[string]any{
		"petId":       p.ID,
		"guildId":     p.GuildID,
		"displayName": p.DisplayName,
		"guildName":   p.GuildName,
		"name":        p.Name,
		"species":     p.Species,
		"attribute":   p.Attribute,
		"level":       p.Level,
		"maxHp":       p.MaxHP,
		"power":       p.Power,
		"guard":       p.Guard,
		"speed":       p.Speed,
	}
}

func battleResultToResponse(result petapp.BattleResult) map[string]any {
	turns := make([]map[string]any, 0, len(result.Turns))
	for _, turn := range result.Turns {
		turns = append(turns, map[string]any{
			"turn":              turn.Turn,
			"actorPetId":        turn.ActorPetID,
			"targetPetId":       turn.TargetPetID,
			"damage":            turn.Damage,
			"targetRemainingHp": turn.TargetRemainingHP,
			"message":           turn.Message,
		})
	}

	var winnerPetID any
	if result.WinnerPetID != "" {
		winnerPetID = result.WinnerPetID
	}

	return map[string]any{
		"result":      result.Result,
		"winnerPetId": winnerPetID,
		"turns":       turns,
		"attacker": map[string]any{
			"petId":       result.Attacker.PetID,
			"name":        result.Attacker.Name,
			"remainingHp": result.Attacker.RemainingHP,
		},
		"defender": map[string]any{
			"petId":       result.Defender.PetID,
			"name":        result.Defender.Name,
			"remainingHp": result.Defender.RemainingHP,
		},
	}
}
