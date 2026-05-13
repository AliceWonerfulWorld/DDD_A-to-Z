package cp

import (
	"context"
	"time"

	cpdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/cp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type UseCase struct {
	ledger LedgerRepository
	ids    IDGenerator
	now    func() time.Time
}

func NewUseCase(ledger LedgerRepository, ids IDGenerator) *UseCase {
	return &UseCase{
		ledger: ledger,
		ids:    ids,
		now:    time.Now,
	}
}

func (u *UseCase) Apply(ctx context.Context, command ApplyCommand) (cpdomain.LedgerEntry, error) {
	if err := ctx.Err(); err != nil {
		return cpdomain.LedgerEntry{}, err
	}

	id, err := u.ids.NewID()
	if err != nil {
		return cpdomain.LedgerEntry{}, err
	}

	entry, err := cpdomain.NewLedgerEntry(
		id,
		command.UserID,
		command.Amount,
		command.Type,
		command.Reason,
		command.SourceType,
		command.SourceID,
		u.now(),
	)
	if err != nil {
		return cpdomain.LedgerEntry{}, err
	}

	return u.ledger.Record(ctx, entry)
}

type ApplyCommand struct {
	UserID     user.ID
	Amount     int64
	Type       cpdomain.EntryType
	Reason     string
	SourceType string
	SourceID   string
}
