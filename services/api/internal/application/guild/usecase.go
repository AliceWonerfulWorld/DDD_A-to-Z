package guild

import (
	"context"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type UseCase struct {
	repository Repository
}

func NewUseCase(repository Repository) *UseCase {
	if repository == nil {
		panic("guild repository is required")
	}

	return &UseCase{repository: repository}
}

func (u *UseCase) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return u.repository.ListGuilds(ctx)
}
