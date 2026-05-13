package cp

import (
	"context"

	cpdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/cp"
)

type LedgerRepository interface {
	Record(ctx context.Context, entry cpdomain.LedgerEntry) (cpdomain.LedgerEntry, error)
}

type IDGenerator interface {
	NewID() (string, error)
}
