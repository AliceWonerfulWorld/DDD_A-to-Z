package home

import (
	"context"
	"errors"
	"time"

	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type UseCase struct {
	auth CurrentUserRepository
	cp   ContributionPointReader
	now  func() time.Time
}

func NewUseCase(auth CurrentUserRepository, cp ContributionPointReader) *UseCase {
	return &UseCase{auth: auth, cp: cp, now: time.Now}
}

type HomeData struct {
	TotalCP                  int64
	TodayCP                  int64
	PlayerLevel              int32
	PlayerLevelTotalCP       int64
	NextPlayerLevel          int32
	NextPlayerLevelTotalCP   int64
	NextPlayerLevelRemaining int64
	LifetimeTotalEarnedCP    int64
}

func (u *UseCase) GetHome(ctx context.Context, sessionToken string) (HomeData, error) {
	if sessionToken == "" {
		return HomeData{}, ErrUnauthenticated
	}

	appUser, ok, err := u.auth.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return HomeData{}, err
	}
	if !ok {
		return HomeData{}, ErrUnauthenticated
	}

	balance, err := u.cp.GetBalance(ctx, appUser.ID)
	if err != nil {
		return HomeData{}, err
	}

	todayEarned, err := u.cp.GetTodayEarned(ctx, appUser.ID)
	if err != nil {
		return HomeData{}, err
	}

	totalEarned, err := u.cp.GetTotalEarned(ctx, appUser.ID)
	if err != nil {
		return HomeData{}, err
	}

	playerLevel := contributionpointdomain.PlayerLevelFromTotalEarned(totalEarned)
	currentLevelTotalEarned := contributionpointdomain.TotalEarnedForPlayerLevel(playerLevel)
	nextLevel, nextLevelTotalEarned, remaining := contributionpointdomain.NextPlayerLevelProgress(totalEarned)

	return HomeData{
		TotalCP:                  balance,
		TodayCP:                  todayEarned,
		PlayerLevel:              int32(playerLevel),
		PlayerLevelTotalCP:       currentLevelTotalEarned,
		NextPlayerLevel:          int32(nextLevel),
		NextPlayerLevelTotalCP:   nextLevelTotalEarned,
		NextPlayerLevelRemaining: remaining,
		LifetimeTotalEarnedCP:    totalEarned,
	}, nil
}
