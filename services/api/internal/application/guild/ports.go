package guild

import (
	"context"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type Repository interface {
	ListGuilds(ctx context.Context) ([]guilddomain.Guild, error)
}
