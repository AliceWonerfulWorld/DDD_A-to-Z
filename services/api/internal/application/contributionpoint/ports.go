package contributionpoint

import (
	"context"
	"time"

	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type LedgerRepository interface {
	Record(ctx context.Context, entry contributionpointdomain.LedgerEntry) (contributionpointdomain.LedgerEntry, error)
	GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
	GetLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (*time.Time, error)
	UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType, at time.Time) error
}

type IDGenerator interface {
	NewID() (string, error)
}
