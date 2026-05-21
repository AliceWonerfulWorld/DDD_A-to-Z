package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	stdhttp "net/http"

	guildtownapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guildtown"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
)

const guildTownPlacementsRequestMaxBytes = 1 << 20

type GuildTownController struct {
	usecase *guildtownapp.UseCase
	logger  *slog.Logger
}

func NewGuildTownController(usecase *guildtownapp.UseCase, logger *slog.Logger) *GuildTownController {
	return &GuildTownController{usecase: usecase, logger: logger}
}

func (c *GuildTownController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /me/guild/town", c.getTown)
	mux.HandleFunc("GET /me/guild/town/inventory", c.getInventory)
	mux.HandleFunc("GET /me/guild/town/placements", c.getPlacements)
	mux.HandleFunc("PUT /me/guild/town/placements", c.savePlacements)
}

func (c *GuildTownController) getTown(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	state, err := c.townState(r)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, townStateResponse(state)); err != nil {
		c.logger.Error("failed to write guild town response", "error", err)
	}
}

func (c *GuildTownController) getInventory(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	state, err := c.townState(r)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"buildings": buildingResponses(state.Buildings),
		"inventory": inventoryResponses(state.Inventory),
	}); err != nil {
		c.logger.Error("failed to write guild town inventory response", "error", err)
	}
}

func (c *GuildTownController) getPlacements(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	state, err := c.townState(r)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"placements": placementResponses(state.Placements),
	}); err != nil {
		c.logger.Error("failed to write guild town placements response", "error", err)
	}
}

func (c *GuildTownController) savePlacements(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildtownapp.ErrUnauthenticated)
		return
	}

	var request struct {
		Placements []struct {
			ID           string  `json:"id"`
			BuildingType string  `json:"building_type"`
			Type         string  `json:"type"`
			X            float64 `json:"x"`
			Y            float64 `json:"y"`
			Width        float64 `json:"width"`
		} `json:"placements"`
	}
	r.Body = stdhttp.MaxBytesReader(w, r.Body, guildTownPlacementsRequestMaxBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		var maxBytesErr *stdhttp.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeAPIError(w, stdhttp.StatusRequestEntityTooLarge, "invalid_request_body", "request body too large", 0, nil)
			return
		}
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_json", "invalid json", 0, nil)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_json", "invalid json", 0, nil)
		return
	}

	commands := make([]guildtownapp.SavePlacementCommand, 0, len(request.Placements))
	for _, placement := range request.Placements {
		buildingType := placement.BuildingType
		if buildingType == "" {
			buildingType = placement.Type
		}
		commands = append(commands, guildtownapp.SavePlacementCommand{
			ID:           guildtowndomain.PlacementID(placement.ID),
			BuildingType: guildtowndomain.BuildingType(buildingType),
			X:            placement.X,
			Y:            placement.Y,
			Width:        placement.Width,
		})
	}

	state, err := c.usecase.SavePlacements(r.Context(), cookie.Value, commands)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, townStateResponse(state)); err != nil {
		c.logger.Error("failed to write saved guild town response", "error", err)
	}
}

func (c *GuildTownController) townState(r *stdhttp.Request) (guildtownapp.TownState, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return guildtownapp.TownState{}, guildtownapp.ErrUnauthenticated
	}

	return c.usecase.GetTown(r.Context(), cookie.Value)
}

func (c *GuildTownController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, guildtownapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, guildtownapp.ErrActiveMembershipNotFound):
		writeAPIError(w, stdhttp.StatusNotFound, "guild_membership_not_found", "active guild membership not found", 0, nil)
	case errors.Is(err, guildtownapp.ErrUnknownBuildingType):
		writeAPIError(w, stdhttp.StatusBadRequest, "unknown_building_type", "unknown building type", 0, nil)
	case errors.Is(err, guildtownapp.ErrInsufficientInventory):
		writeAPIError(w, stdhttp.StatusBadRequest, "insufficient_inventory", "insufficient inventory", 0, nil)
	default:
		c.logger.Error("guild town request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func townStateResponse(state guildtownapp.TownState) map[string]any {
	progress := guilddomain.GuildLevelProgressFromExperience(state.Guild.GuildExperience)
	guildLevel := state.Guild.GuildLevel
	if guildLevel <= 0 {
		guildLevel = progress.Level
	}
	currentLevelExperience := state.Guild.CurrentGuildLevelExperience
	if currentLevelExperience == 0 {
		currentLevelExperience = progress.CurrentLevelExperience
	}
	nextLevelExperience := state.Guild.NextGuildLevelExperience
	if nextLevelExperience == 0 {
		nextLevelExperience = progress.NextLevelExperience
	}

	return map[string]any{
		"guild_level":                    guildLevel,
		"guild_experience":               state.Guild.GuildExperience,
		"current_guild_level_experience": currentLevelExperience,
		"next_guild_level_experience":    nextLevelExperience,
		"buildings":                      buildingResponses(state.Buildings),
		"inventory":                      inventoryResponses(state.Inventory),
		"placements":                     placementResponses(state.Placements),
	}
}

func buildingResponses(buildings []guildtowndomain.BuildingMaster) []map[string]any {
	responses := make([]map[string]any, 0, len(buildings))
	for _, building := range buildings {
		responses = append(responses, map[string]any{
			"type":          building.Type,
			"name":          building.Name,
			"title":         building.Title,
			"description":   building.Description,
			"src":           building.Src,
			"min_map_width": building.MinMapWidth,
			"map_width_vw":  building.MapWidthVW,
			"max_map_width": building.MaxMapWidth,
		})
	}

	return responses
}

func inventoryResponses(inventory []guildtowndomain.InventoryItem) []map[string]any {
	responses := make([]map[string]any, 0, len(inventory))
	for _, item := range inventory {
		responses = append(responses, map[string]any{
			"type":  item.BuildingType,
			"count": item.Quantity,
		})
	}

	return responses
}

func placementResponses(placements []guildtowndomain.Placement) []map[string]any {
	responses := make([]map[string]any, 0, len(placements))
	for _, placement := range placements {
		responses = append(responses, map[string]any{
			"id":            placement.ID,
			"type":          placement.BuildingType,
			"building_type": placement.BuildingType,
			"x":             placement.X,
			"y":             placement.Y,
			"width":         placement.Width,
			"z_index":       placement.ZIndex,
		})
	}

	return responses
}
