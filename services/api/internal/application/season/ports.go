package season

import "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/season"

type Repository interface {
	FindCurrent() (season.Season, error)
	FindByID(id season.ID) (season.Season, error)
	FindByNumber(number int) (season.Season, error)
	ListAll() ([]season.Season, error)
	Create(s season.Season) error
	ListGuildRankings(seasonID season.ID) ([]season.GuildSeasonRanking, error)
	ListGuildMemberRankings(seasonID season.ID, guildID string) ([]season.GuildSeasonMemberRanking, error)
	GetGuildSeasonCP(seasonID season.ID, guildID string) (int64, error)
}

type IDGenerator interface {
	NewID() (string, error)
}
