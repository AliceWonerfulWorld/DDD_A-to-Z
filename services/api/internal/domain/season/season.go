package season

import (
	"errors"
	"strings"
	"time"
)

var ErrSeasonNotFound = errors.New("season not found")

type ID string

type Season struct {
	ID        ID
	Number    int
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSeason(s Season) (Season, error) {
	if s.ID == "" {
		return Season{}, errors.New("season id is required")
	}
	if s.Number < 1 {
		return Season{}, errors.New("season number must be positive")
	}
	if s.StartsAt.IsZero() {
		return Season{}, errors.New("season starts at is required")
	}
	if s.EndsAt.IsZero() {
		return Season{}, errors.New("season ends at is required")
	}
	if !s.EndsAt.After(s.StartsAt) {
		return Season{}, errors.New("season ends at must be after starts at")
	}
	if s.CreatedAt.IsZero() {
		return Season{}, errors.New("season created at is required")
	}
	if s.UpdatedAt.IsZero() {
		return Season{}, errors.New("season updated at is required")
	}

	return s, nil
}

func (s Season) IsCurrent() bool {
	now := time.Now()
	return !now.Before(s.StartsAt) && !now.After(s.EndsAt)
}

type GuildSeasonRanking struct {
	ID          string
	SeasonID    ID
	GuildID     string
	TotalCP     int64
	Rank        int
	MemberCount int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewGuildSeasonRanking(r GuildSeasonRanking) (GuildSeasonRanking, error) {
	if strings.TrimSpace(r.ID) == "" {
		return GuildSeasonRanking{}, errors.New("guild season ranking id is required")
	}
	if r.SeasonID == "" {
		return GuildSeasonRanking{}, errors.New("guild season ranking season id is required")
	}
	if strings.TrimSpace(r.GuildID) == "" {
		return GuildSeasonRanking{}, errors.New("guild season ranking guild id is required")
	}
	if r.Rank < 1 {
		return GuildSeasonRanking{}, errors.New("guild season ranking rank must be positive")
	}
	if r.TotalCP < 0 {
		return GuildSeasonRanking{}, errors.New("guild season ranking total cp cannot be negative")
	}
	if r.CreatedAt.IsZero() {
		return GuildSeasonRanking{}, errors.New("guild season ranking created at is required")
	}
	if r.UpdatedAt.IsZero() {
		return GuildSeasonRanking{}, errors.New("guild season ranking updated at is required")
	}

	return r, nil
}

type GuildSeasonMemberRanking struct {
	ID            string
	SeasonID      ID
	GuildID       string
	UserID        string
	UserName      string
	ContributedCP int64
	Rank          int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewGuildSeasonMemberRanking(r GuildSeasonMemberRanking) (GuildSeasonMemberRanking, error) {
	if strings.TrimSpace(r.ID) == "" {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking id is required")
	}
	if r.SeasonID == "" {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking season id is required")
	}
	if strings.TrimSpace(r.GuildID) == "" {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking guild id is required")
	}
	if strings.TrimSpace(r.UserID) == "" {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking user id is required")
	}
	if strings.TrimSpace(r.UserName) == "" {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking user name is required")
	}
	if r.Rank < 1 {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking rank must be positive")
	}
	if r.ContributedCP < 0 {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking contributed cp cannot be negative")
	}
	if r.CreatedAt.IsZero() {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking created at is required")
	}
	if r.UpdatedAt.IsZero() {
		return GuildSeasonMemberRanking{}, errors.New("guild season member ranking updated at is required")
	}

	return r, nil
}
