package http

import (
	"log/slog"
	"strconv"

	seasonapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/season"
	seasondomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/season"
	stdhttp "net/http"
)

type SeasonController struct {
	usecase *seasonapp.UseCase
	logger  *slog.Logger
}

func NewSeasonController(usecase *seasonapp.UseCase, logger *slog.Logger) *SeasonController {
	return &SeasonController{
		usecase: usecase,
		logger:  logger,
	}
}

func (c *SeasonController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /seasons", c.listSeasons)
	mux.HandleFunc("GET /seasons/current", c.getCurrentSeason)
	mux.HandleFunc("GET /seasons/{seasonNumber}", c.getSeasonByNumber)
	mux.HandleFunc("GET /seasons/{seasonNumber}/guild-rankings", c.listGuildRankings)
	mux.HandleFunc("GET /seasons/{seasonNumber}/guilds/{guildID}/member-rankings", c.listGuildMemberRankings)
}

func (c *SeasonController) getCurrentSeason(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s, err := c.usecase.GetCurrentSeason()
	if err != nil {
		c.logger.Error("failed to get current season", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, seasonResponse(s)); err != nil {
		c.logger.Error("failed to write current season response", "error", err)
	}
}

func (c *SeasonController) listSeasons(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	seasons, err := c.usecase.ListSeasons()
	if err != nil {
		c.logger.Error("failed to list seasons", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	items := make([]map[string]any, 0, len(seasons))
	for _, s := range seasons {
		items = append(items, seasonResponse(s))
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"seasons": items,
	}); err != nil {
		c.logger.Error("failed to write seasons list response", "error", err)
	}
}

func (c *SeasonController) getSeasonByNumber(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	numberStr := r.PathValue("seasonNumber")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_season_number", "Season number must be an integer", 0, nil)
		return
	}

	s, err := c.usecase.GetSeasonByNumber(number)
	if err != nil {
		c.logger.Error("failed to get season", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, seasonResponse(s)); err != nil {
		c.logger.Error("failed to write season response", "error", err)
	}
}

func (c *SeasonController) listGuildRankings(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	seasonID, err := c.resolveSeasonID(w, r)
	if err != nil {
		return
	}

	rankings, err := c.usecase.ListGuildRankings(seasonID)
	if err != nil {
		c.logger.Error("failed to list guild rankings", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"rankings": guildSeasonRankingResponses(rankings),
	}); err != nil {
		c.logger.Error("failed to write guild rankings response", "error", err)
	}
}

func (c *SeasonController) listGuildMemberRankings(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	seasonID, err := c.resolveSeasonID(w, r)
	if err != nil {
		return
	}

	guildID := r.PathValue("guildID")

	rankings, err := c.usecase.ListGuildMemberRankings(seasonID, guildID)
	if err != nil {
		c.logger.Error("failed to list guild member rankings", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"rankings": guildSeasonMemberRankingResponses(rankings),
	}); err != nil {
		c.logger.Error("failed to write guild member rankings response", "error", err)
	}
}

func (c *SeasonController) resolveSeasonID(w stdhttp.ResponseWriter, r *stdhttp.Request) (seasondomain.ID, error) {
	numberStr := r.PathValue("seasonNumber")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		writeAPIError(w, stdhttp.StatusBadRequest, "invalid_season_number", "Season number must be an integer", 0, nil)
		return "", err
	}

	s, err := c.usecase.GetSeasonByNumber(number)
	if err != nil {
		c.logger.Error("failed to resolve season", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
		return "", err
	}

	return s.ID, nil
}

func seasonResponse(s seasondomain.Season) map[string]any {
	return map[string]any{
		"id":         string(s.ID),
		"number":     s.Number,
		"starts_at":  s.StartsAt,
		"ends_at":    s.EndsAt,
		"is_current": s.IsCurrent(),
	}
}

func guildSeasonRankingResponses(rankings []seasondomain.GuildSeasonRanking) []map[string]any {
	result := make([]map[string]any, 0, len(rankings))
	for _, r := range rankings {
		result = append(result, map[string]any{
			"id":           r.ID,
			"season_id":    string(r.SeasonID),
			"guild_id":     r.GuildID,
			"total_cp":     r.TotalCP,
			"rank":         r.Rank,
			"member_count": r.MemberCount,
		})
	}
	return result
}

func guildSeasonMemberRankingResponses(rankings []seasondomain.GuildSeasonMemberRanking) []map[string]any {
	result := make([]map[string]any, 0, len(rankings))
	for _, r := range rankings {
		result = append(result, map[string]any{
			"id":             r.ID,
			"season_id":      string(r.SeasonID),
			"guild_id":       r.GuildID,
			"user_id":        r.UserID,
			"user_name":      r.UserName,
			"contributed_cp": r.ContributedCP,
			"rank":           r.Rank,
		})
	}
	return result
}
