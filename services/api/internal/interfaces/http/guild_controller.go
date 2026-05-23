package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"strconv"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
)

type GuildController struct {
	usecase *guildapp.UseCase
	logger  *slog.Logger
}

func NewGuildController(usecase *guildapp.UseCase, logger *slog.Logger) *GuildController {
	return &GuildController{
		usecase: usecase,
		logger:  logger,
	}
}

func (c *GuildController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /guilds", c.listGuilds)
	mux.HandleFunc("POST /guilds/{guildID}/join", c.joinGuild)
	mux.HandleFunc("GET /guilds/{guildID}/dashboard", c.getGuildDashboard)
	mux.HandleFunc("GET /guilds/{guildID}/activity-logs", c.listGuildActivityLogs)
	mux.HandleFunc("GET /guilds/{guildID}/cp-contributions", c.listGuildCPContributions)
	mux.HandleFunc("GET /me/guild", c.getMyGuild)
	mux.HandleFunc("DELETE /me/guild", c.leaveMyGuild)
	mux.HandleFunc("POST /me/guild/cp-contributions", c.contributeCP)
	mux.HandleFunc("GET /me/guild/cp-contributions", c.listMyGuildCPContributions)
}

func (c *GuildController) listGuilds(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	guilds, err := c.usecase.ListGuilds(r.Context())
	if err != nil {
		c.logger.Error("failed to list guilds", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"guilds": guildResponses(guilds),
	}); err != nil {
		c.logger.Error("failed to write guild list response", "error", err)
	}
}

func (c *GuildController) joinGuild(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	result, err := c.usecase.JoinGuild(r.Context(), cookie.Value, guilddomain.ID(r.PathValue("guildID")))
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusCreated, joinGuildResponse(result)); err != nil {
		c.logger.Error("failed to write guild join response", "error", err)
	}
}

func (c *GuildController) getMyGuild(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	details, ok, err := c.usecase.GetMyGuildDetails(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}
	if !ok {
		if err := writeJSON(w, stdhttp.StatusOK, map[string]any{"guild": nil}); err != nil {
			c.logger.Error("failed to write empty my guild response", "error", err)
		}
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, myGuildDetailsResponse(details)); err != nil {
		c.logger.Error("failed to write my guild response", "error", err)
	}
}

func (c *GuildController) getGuildDashboard(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	details, err := c.usecase.GetGuildDashboard(r.Context(), cookie.Value, guilddomain.ID(r.PathValue("guildID")))
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, guildDashboardResponse(details)); err != nil {
		c.logger.Error("failed to write guild dashboard response", "error", err)
	}
}

func (c *GuildController) leaveMyGuild(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	if err := c.usecase.LeaveMyGuild(r.Context(), cookie.Value); err != nil {
		c.writeError(w, err)
		return
	}

	w.WriteHeader(stdhttp.StatusNoContent)
}

func (c *GuildController) listGuildActivityLogs(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	limit := 0
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil {
			writeAPIError(w, stdhttp.StatusBadRequest, "invalid_limit", "limit must be a number", 0, nil)
			return
		}
		limit = parsedLimit
	}

	logs, err := c.usecase.ListGuildActivityLogs(r.Context(), cookie.Value, guilddomain.ID(r.PathValue("guildID")), limit)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"logs": activityLogResponses(logs),
	}); err != nil {
		c.logger.Error("failed to write guild activity log response", "error", err)
	}
}

func (c *GuildController) contributeCP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	var request struct {
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_request", "invalid request body", 0, nil)
		return
	}

	contribution, err := c.usecase.ContributeCP(r.Context(), cookie.Value, request.Amount)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusCreated, map[string]any{
		"contribution": cpContributionResponse(contribution),
	}); err != nil {
		c.logger.Error("failed to write guild cp contribution response", "error", err)
	}
}

func (c *GuildController) listGuildCPContributions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	contributions, err := c.usecase.ListGuildCPContributions(r.Context(), guilddomain.ID(r.PathValue("guildID")))
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"contributions": cpContributionResponses(contributions),
	}); err != nil {
		c.logger.Error("failed to write guild cp contribution list response", "error", err)
	}
}

func (c *GuildController) listMyGuildCPContributions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeError(w, guildapp.ErrUnauthenticated)
		return
	}

	contributions, err := c.usecase.ListMyGuildCPContributions(r.Context(), cookie.Value)
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"contributions": cpContributionResponses(contributions),
	}); err != nil {
		c.logger.Error("failed to write my guild cp contribution list response", "error", err)
	}
}

func (c *GuildController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, guildapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, guildapp.ErrGuildNotFound):
		writeAPIError(w, stdhttp.StatusNotFound, "guild_not_found", "guild not found", 0, nil)
	case errors.Is(err, guildapp.ErrAlreadyJoined):
		writeAPIError(w, stdhttp.StatusConflict, "already_joined_guild", "user already joined a guild", 0, nil)
	case errors.Is(err, guildapp.ErrActiveMembershipNotFound):
		writeAPIError(w, stdhttp.StatusNotFound, "guild_membership_not_found", "active guild membership not found", 0, nil)
	case errors.Is(err, guildapp.ErrGuildAccessDenied):
		writeAPIError(w, stdhttp.StatusForbidden, "guild_access_denied", "guild access denied", 0, nil)
	case errors.Is(err, guildapp.ErrInvalidCPContribution):
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_cp_contribution", "guild cp contribution amount must be positive", 0, nil)
	case errors.Is(err, contributionpointapp.ErrInsufficientBalance):
		writeAPIError(w, stdhttp.StatusConflict, "insufficient_cp_balance", "contribution point balance is insufficient", 0, nil)
	default:
		c.logger.Error("guild request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}

func guildResponses(guilds []guilddomain.Guild) []map[string]any {
	responses := make([]map[string]any, 0, len(guilds))
	for _, guild := range guilds {
		responses = append(responses, guildResponse(guild))
	}

	return responses
}

func membershipResponse(membership guilddomain.MembershipWithGuild) map[string]any {
	return map[string]any{
		"guild": guildResponse(membership.Guild),
		"membership": map[string]any{
			"id":        membership.Membership.ID,
			"user_id":   membership.Membership.UserID,
			"joined_at": membership.Membership.JoinedAt.Format(time.RFC3339),
		},
	}
}

func joinGuildResponse(result guildapp.JoinGuildResult) map[string]any {
	response := membershipResponse(result.Membership)
	response["granted_pet"] = petResponse(result.GrantedPet)
	response["pet_already_owned"] = result.PetAlreadyOwned

	return response
}

func petResponse(pet *petdomain.Pet) any {
	if pet == nil {
		return nil
	}

	return map[string]any{
		"id":         pet.ID,
		"user_id":    pet.UserID,
		"guild_id":   pet.GuildID,
		"attribute":  pet.Attribute,
		"vitality":   pet.Stats.Vitality,
		"strength":   pet.Stats.Strength,
		"agility":    pet.Stats.Agility,
		"created_at": pet.CreatedAt.Format(time.RFC3339),
		"updated_at": pet.UpdatedAt.Format(time.RFC3339),
	}
}

func myGuildDetailsResponse(details guildapp.MyGuildDetails) map[string]any {
	return map[string]any{
		"guild":   guildResponse(details.Guild),
		"members": memberContributionResponses(details.Members),
		"membership": map[string]any{
			"id":        details.Membership.ID,
			"user_id":   details.Membership.UserID,
			"joined_at": details.Membership.JoinedAt.Format(time.RFC3339),
		},
	}
}

func guildDashboardResponse(details guildapp.MyGuildDetails) map[string]any {
	response := myGuildDetailsResponse(details)
	response["state"] = "joined"

	return response
}

func guildResponse(guild guilddomain.Guild) map[string]any {
	return map[string]any{
		"id":                             guild.ID,
		"slug":                           guild.Slug,
		"name":                           guild.Name,
		"description":                    guild.Description,
		"icon":                           guild.Icon,
		"color":                          guild.Color,
		"member_count":                   guild.MemberCount,
		"total_contributed_cp":           guild.TotalContributedCP,
		"guild_experience":               guild.GuildExperience,
		"current_exp":                    guild.GuildExperience,
		"currentExp":                     guild.GuildExperience,
		"guild_level":                    guild.GuildLevel,
		"guildLevel":                     guild.GuildLevel,
		"current_guild_level_experience": guild.CurrentGuildLevelExperience,
		"next_guild_level_experience":    guild.NextGuildLevelExperience,
	}
}

func memberContributionResponses(members []guilddomain.MemberContribution) []map[string]any {
	responses := make([]map[string]any, 0, len(members))
	for _, member := range members {
		responses = append(responses, map[string]any{
			"user_id":              member.UserID,
			"name":                 member.Name,
			"total_earned_cp":      member.TotalEarnedCP,
			"total_contributed_cp": member.TotalContributedCP,
			"joined_at":            member.JoinedAt.Format(time.RFC3339),
		})
	}

	return responses
}

func activityLogResponses(logs []guilddomain.ActivityLog) []map[string]any {
	responses := make([]map[string]any, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, map[string]any{
			"id":          log.ID,
			"user_id":     log.UserID,
			"player":      log.Player,
			"type":        log.Type,
			"repo":        log.Repo,
			"message":     log.Message,
			"language":    log.Language,
			"cp":          log.CP,
			"occurred_at": log.OccurredAt.Format(time.RFC3339),
		})
	}

	return responses
}

func cpContributionResponses(contributions []guilddomain.CPContribution) []map[string]any {
	responses := make([]map[string]any, 0, len(contributions))
	for _, contribution := range contributions {
		responses = append(responses, cpContributionResponse(contribution))
	}

	return responses
}

func cpContributionResponse(contribution guilddomain.CPContribution) map[string]any {
	return map[string]any{
		"id":              contribution.ID,
		"guild_id":        contribution.GuildID,
		"user_id":         contribution.UserID,
		"point_ledger_id": contribution.PointLedgerID,
		"amount":          contribution.Amount,
		"created_at":      contribution.CreatedAt.Format(time.RFC3339),
	}
}
