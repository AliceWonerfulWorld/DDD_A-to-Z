package mypage

import (
	"context"
	"errors"
	"time"
)

const defaultRecentLimit = 5

var ErrUnauthenticated = errors.New("unauthenticated")

// UseCase orchestrates the my-page data aggregation.
type UseCase struct {
	current      CurrentUserRepository
	cp           ContributionPointReader
	repositories RepositorySummaryReader
	now          func() time.Time
}

// NewUseCase creates a new my-page use case.
func NewUseCase(
	current CurrentUserRepository,
	cp ContributionPointReader,
	repositories RepositorySummaryReader,
) *UseCase {
	return &UseCase{
		current:      current,
		cp:           cp,
		repositories: repositories,
		now:          time.Now,
	}
}

// GetMyPage aggregates user information, CP summary, and repository summary
// into a single MyPageData for the authenticated user.
func (u *UseCase) GetMyPage(ctx context.Context, sessionToken string) (MyPageData, error) {
	if sessionToken == "" {
		return MyPageData{}, ErrUnauthenticated
	}

	// 1. Resolve session → user
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return MyPageData{}, err
	}
	if !ok {
		return MyPageData{}, ErrUnauthenticated
	}

	// 2. Fetch CP summary
	balance, err := u.cp.GetBalance(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}
	totalEarned, err := u.cp.GetTotalEarned(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}
	totalSpent, err := u.cp.GetTotalSpent(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}

	// 3. Fetch repository summary
	repoSummary, err := u.repositories.GetRepositorySummary(ctx, appUser.ID, defaultRecentLimit)
	if err != nil {
		return MyPageData{}, err
	}

	// 4. Guild is nil until the guild feature is implemented.
	return MyPageData{
		User: appUser,
		CP: CPSummary{
			Balance:     balance,
			TotalEarned: totalEarned,
			TotalSpent:  totalSpent,
		},
		Repositories: repoSummary,
		Guild:        nil,
	}, nil
}
