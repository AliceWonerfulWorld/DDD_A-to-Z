package http

import (
	"log/slog"
	stdhttp "net/http"

	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
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

func guildResponses(guilds []guilddomain.Guild) []map[string]any {
	responses := make([]map[string]any, 0, len(guilds))
	for _, guild := range guilds {
		responses = append(responses, map[string]any{
			"id":           guild.ID,
			"slug":         guild.Slug,
			"name":         guild.Name,
			"description":  guild.Description,
			"icon":         guild.Icon,
			"color":        guild.Color,
			"member_count": guild.MemberCount,
		})
	}

	return responses
}
