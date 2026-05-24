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
	var season seasonRow
	if err := s.db.Where("id = ?", string(seasonID)).First(&season).Error; err != nil {
		return nil, fmt.Errorf("find season for rankings: %w", err)
	}

	type aggRow struct {
		GuildID     string
		TotalCP     int64
		MemberCount int64
	}

	var results []aggRow
	if err := s.db.Raw(`
		SELECT gcc.guild_id, SUM(gcc.amount) AS total_cp, COUNT(DISTINCT gcc.user_id) AS member_count
		FROM guild_cp_contributions gcc
		WHERE gcc.created_at >= ? AND gcc.created_at <= ?
		GROUP BY gcc.guild_id
		ORDER BY total_cp DESC
	`, season.StartsAt, season.EndsAt).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("aggregate guild season rankings: %w", err)
	}

	now := time.Now()
	rankings := make([]seasondomain.GuildSeasonRanking, 0, len(results))
	for i, r := range results {
		rankings = append(rankings, seasondomain.GuildSeasonRanking{
			ID:          fmt.Sprintf("%s_%s", string(seasonID), r.GuildID),
			SeasonID:    seasonID,
			GuildID:     r.GuildID,
			TotalCP:     r.TotalCP,
			Rank:        i + 1,
			MemberCount: int(r.MemberCount),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return rankings, nil
}

func (s *SeasonStore) ListGuildMemberRankings(seasonID seasondomain.ID, guildID string) ([]seasondomain.GuildSeasonMemberRanking, error) {
	var season seasonRow
	if err := s.db.Where("id = ?", string(seasonID)).First(&season).Error; err != nil {
		return nil, fmt.Errorf("find season for member rankings: %w", err)
	}

	type aggRow struct {
		UserID        string
		UserName      string
		ContributedCP int64
	}

	var results []aggRow
	if err := s.db.Raw(`
		SELECT gcc.user_id, COALESCE(ga.username, '') AS user_name, SUM(gcc.amount) AS contributed_cp
		FROM guild_cp_contributions gcc
		LEFT JOIN github_accounts ga ON ga.user_id = gcc.user_id
		WHERE gcc.created_at >= ? AND gcc.created_at <= ? AND gcc.guild_id = ?
		GROUP BY gcc.user_id, ga.username
		ORDER BY contributed_cp DESC
	`, season.StartsAt, season.EndsAt, guildID).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("aggregate guild season member rankings: %w", err)
	}

	now := time.Now()
	rankings := make([]seasondomain.GuildSeasonMemberRanking, 0, len(results))
	for i, r := range results {
		rankings = append(rankings, seasondomain.GuildSeasonMemberRanking{
			ID:            fmt.Sprintf("%s_%s_%s", string(seasonID), guildID, r.UserID),
			SeasonID:      seasonID,
			GuildID:       guildID,
			UserID:        r.UserID,
			UserName:      r.UserName,
			ContributedCP: r.ContributedCP,
			Rank:          i + 1,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	return rankings, nil
}

func (s *SeasonStore) GetGuildSeasonCP(seasonID seasondomain.ID, guildID string) (int64, error) {
	var season seasonRow
	if err := s.db.Where("id = ?", string(seasonID)).First(&season).Error; err != nil {
		return 0, fmt.Errorf("find season for guild cp: %w", err)
	}

	var totalCP struct {
		Total int64
	}
	if err := s.db.Raw(`
		SELECT COALESCE(SUM(amount), 0) AS total
		FROM guild_cp_contributions
		WHERE guild_id = ? AND created_at >= ? AND created_at <= ?
	`, guildID, season.StartsAt, season.EndsAt).Scan(&totalCP).Error; err != nil {
		return 0, fmt.Errorf("aggregate guild season cp: %w", err)
	}

	return totalCP.Total, nil
}

type seasonRow struct {
	ID        string `gorm:"primaryKey"`
	Number    int
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (seasonRow) TableName() string { return "seasons" }

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
