package sp

import (
	"context"
	"errors"
	"time"
)

var ErrUnauthenticated = errors.New("unauthenticated")

// UseCase handles SP data retrieval for the authenticated user.
type UseCase struct {
	current CurrentUserRepository
	sp      SPReader
	now     func() time.Time
}

// NewUseCase creates a new SP use case.
func NewUseCase(current CurrentUserRepository, sp SPReader) *UseCase {
	return &UseCase{
		current: current,
		sp:      sp,
		now:     time.Now,
	}
}

// GetSP returns the SP balances for all languages of the authenticated user.
// Languages with zero balance are excluded.
func (u *UseCase) GetSP(ctx context.Context, sessionToken string) ([]SPBalance, error) {
	if sessionToken == "" {
		return nil, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUnauthenticated
	}

	return u.sp.GetSPBalances(ctx, appUser.ID)
}
