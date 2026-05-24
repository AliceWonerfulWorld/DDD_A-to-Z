package postgres

import (
	"errors"
	"fmt"
	"time"

	seasondomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/season"
	"gorm.io/gorm"
)

type SeasonStore struct {
	db *gorm.DB
}

func NewSeasonStore(db *gorm.DB) *SeasonStore {
	return &SeasonStore{db: db}
}

func (s *SeasonStore) FindCurrent() (seasondomain.Season, error) {
	var row seasonRow
	if err := s.db.Where("starts_at <= NOW() AND ends_at >= NOW()").First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return seasondomain.Season{}, seasondomain.ErrSeasonNotFound
		}
		return seasondomain.Season{}, fmt.Errorf("find current season: %w", err)
	}

	return toSeason(row), nil
}

func (s *SeasonStore) FindByID(id seasondomain.ID) (seasondomain.Season, error) {
	var row seasonRow
	if err := s.db.Where("id = ?", string(id)).First(&row).Error; err != nil {
		return seasondomain.Season{}, fmt.Errorf("find season by id: %w", err)
	}

	return toSeason(row), nil
}

func (s *SeasonStore) FindByNumber(number int) (seasondomain.Season, error) {
	var row seasonRow
	if err := s.db.Where("number = ?", number).First(&row).Error; err != nil {
		return seasondomain.Season{}, fmt.Errorf("find season by number: %w", err)
	}

	return toSeason(row), nil
}

func (s *SeasonStore) ListAll() ([]seasondomain.Season, error) {
	var rows []seasonRow
	if err := s.db.Order("number DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list seasons: %w", err)
	}

	seasons := make([]seasondomain.Season, 0, len(rows))
	for _, r := range rows {
		seasons = append(seasons, toSeason(r))
	}

	return seasons, nil
}

func (s *SeasonStore) Create(szn seasondomain.Season) error {
	row := seasonRow{
		ID:        string(szn.ID),
		Number:    szn.Number,
		StartsAt:  szn.StartsAt,
		EndsAt:    szn.EndsAt,
		CreatedAt: szn.CreatedAt,
		UpdatedAt: szn.UpdatedAt,
	}

	if err := s.db.Create(&row).Error; err != nil {
		return fmt.Errorf("create season: %w", err)
	}

	return nil
}

func (s *SeasonStore) ListGuildRankings(seasonID seasondomain.ID) ([]seasondomain.GuildSeasonRanking, error) {
	var rows []guildSeasonRankingRow
	if err := s.db.Where("season_id = ?", string(seasonID)).Order("rank ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list guild season rankings: %w", err)
	}

	rankings := make([]seasondomain.GuildSeasonRanking, 0, len(rows))
	for _, r := range rows {
		rankings = append(rankings, toGuildSeasonRanking(r))
	}

	return rankings, nil
}

func (s *SeasonStore) ListGuildMemberRankings(seasonID seasondomain.ID, guildID string) ([]seasondomain.GuildSeasonMemberRanking, error) {
	var rows []guildSeasonMemberRankingRow
	if err := s.db.Where("season_id = ? AND guild_id = ?", string(seasonID), guildID).Order("rank ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list guild season member rankings: %w", err)
	}

	rankings := make([]seasondomain.GuildSeasonMemberRanking, 0, len(rows))
	for _, r := range rows {
		rankings = append(rankings, toGuildSeasonMemberRanking(r))
	}

	return rankings, nil
}

type seasonRow struct {
	ID        string `gorm:"primaryKey"`
	Number    int
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type guildSeasonRankingRow struct {
	ID          string `gorm:"primaryKey"`
	SeasonID    string
	GuildID     string
	TotalCP     int64 `gorm:"column:total_cp"`
	Rank        int
	MemberCount int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type guildSeasonMemberRankingRow struct {
	ID            string `gorm:"primaryKey"`
	SeasonID      string
	GuildID       string
	UserID        string
	UserName      string
	ContributedCP int64 `gorm:"column:contributed_cp"`
	Rank          int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (seasonRow) TableName() string                   { return "seasons" }
func (guildSeasonRankingRow) TableName() string       { return "guild_season_rankings" }
func (guildSeasonMemberRankingRow) TableName() string { return "guild_season_member_rankings" }

func toSeason(r seasonRow) seasondomain.Season {
	return seasondomain.Season{
		ID:        seasondomain.ID(r.ID),
		Number:    r.Number,
		StartsAt:  r.StartsAt,
		EndsAt:    r.EndsAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func toGuildSeasonRanking(r guildSeasonRankingRow) seasondomain.GuildSeasonRanking {
	return seasondomain.GuildSeasonRanking{
		ID:          r.ID,
		SeasonID:    seasondomain.ID(r.SeasonID),
		GuildID:     r.GuildID,
		TotalCP:     r.TotalCP,
		Rank:        r.Rank,
		MemberCount: r.MemberCount,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toGuildSeasonMemberRanking(r guildSeasonMemberRankingRow) seasondomain.GuildSeasonMemberRanking {
	return seasondomain.GuildSeasonMemberRanking{
		ID:            r.ID,
		SeasonID:      seasondomain.ID(r.SeasonID),
		GuildID:       r.GuildID,
		UserID:        r.UserID,
		UserName:      r.UserName,
		ContributedCP: r.ContributedCP,
		Rank:          r.Rank,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}
