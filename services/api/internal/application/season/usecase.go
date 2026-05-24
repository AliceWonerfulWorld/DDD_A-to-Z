package season

import (
	"errors"
	"fmt"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/season"
)

type UseCase struct {
	repo  Repository
	idGen IDGenerator
}

func NewUseCase(repo Repository, idGen IDGenerator) *UseCase {
	return &UseCase{repo: repo, idGen: idGen}
}

func (uc *UseCase) getOrCreateCurrentSeason() (season.Season, error) {
	s, err := uc.repo.FindCurrent()
	if err == nil {
		return s, nil
	}
	if !errors.Is(err, season.ErrSeasonNotFound) {
		return season.Season{}, fmt.Errorf("find current season: %w", err)
	}

	allSeasons, err := uc.repo.ListAll()
	if err != nil {
		return season.Season{}, fmt.Errorf("list seasons for auto-advance: %w", err)
	}

	now := time.Now()
	newIDStr, err := uc.idGen.NewID()
	if err != nil {
		return season.Season{}, fmt.Errorf("generate season id: %w", err)
	}
	newID := season.ID(newIDStr)

	var newNumber int
	var newStartsAt time.Time

	if len(allSeasons) > 0 {
		lastSeason := allSeasons[0]
		newNumber = lastSeason.Number + 1
		newStartsAt = lastSeason.EndsAt
		if now.After(lastSeason.EndsAt) {
			newStartsAt = now
		}
	} else {
		newNumber = 1
		newStartsAt = now
	}

	newEndsAt := newStartsAt.AddDate(0, 3, 0)

	newSeason, err := season.NewSeason(season.Season{
		ID:        newID,
		Number:    newNumber,
		StartsAt:  newStartsAt,
		EndsAt:    newEndsAt,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return season.Season{}, fmt.Errorf("validate new season: %w", err)
	}

	if err := uc.repo.Create(newSeason); err != nil {
		return season.Season{}, fmt.Errorf("save new season: %w", err)
	}

	return newSeason, nil
}

func (uc *UseCase) GetCurrentSeason() (season.Season, error) {
	return uc.getOrCreateCurrentSeason()
}

func (uc *UseCase) GetSeasonByNumber(number int) (season.Season, error) {
	s, err := uc.repo.FindByNumber(number)
	if err != nil {
		return season.Season{}, fmt.Errorf("get season by number: %w", err)
	}

	return s, nil
}

func (uc *UseCase) ListSeasons() ([]season.Season, error) {
	_, _ = uc.getOrCreateCurrentSeason()

	seasons, err := uc.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("list seasons: %w", err)
	}

	return seasons, nil
}

func (uc *UseCase) ListGuildRankings(seasonID season.ID) ([]season.GuildSeasonRanking, error) {
	rankings, err := uc.repo.ListGuildRankings(seasonID)
	if err != nil {
		return nil, fmt.Errorf("list guild rankings: %w", err)
	}

	return rankings, nil
}

func (uc *UseCase) ListGuildMemberRankings(seasonID season.ID, guildID string) ([]season.GuildSeasonMemberRanking, error) {
	rankings, err := uc.repo.ListGuildMemberRankings(seasonID, guildID)
	if err != nil {
		return nil, fmt.Errorf("list guild member rankings: %w", err)
	}

	return rankings, nil
}
