package contributionpoint

import (
	"context"
	"errors"
	"time"

	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

var ErrInsufficientBalance = errors.New("contribution point balance is insufficient")

type UseCase struct {
	ledger LedgerRepository
	ids    IDGenerator
	now    func() time.Time
}

type Service = UseCase

func NewUseCase(ledger LedgerRepository, ids IDGenerator) *UseCase {
	return &UseCase{
		ledger: ledger,
		ids:    ids,
		now:    time.Now,
	}
}

func NewService(ledger LedgerRepository, ids IDGenerator) *Service {
	return NewUseCase(ledger, ids)
}

func (u *UseCase) Apply(ctx context.Context, command ApplyCommand) (contributionpointdomain.LedgerEntry, error) {
	if err := ctx.Err(); err != nil {
		return contributionpointdomain.LedgerEntry{}, err
	}

	id, err := u.ids.NewID()
	if err != nil {
		return contributionpointdomain.LedgerEntry{}, err
	}

	entry, err := contributionpointdomain.NewLedgerEntry(
		id,
		command.UserID,
		command.PointType,
		command.Amount,
		command.Type,
		command.Reason,
		command.SourceType,
		command.SourceID,
		u.now(),
	)
	if err != nil {
		return contributionpointdomain.LedgerEntry{}, err
	}

	return u.ledger.Record(ctx, entry)
}

func (u *UseCase) Earn(ctx context.Context, command EarnCommand) (contributionpointdomain.LedgerEntry, error) {
	return u.Apply(ctx, ApplyCommand{
		UserID:     command.UserID,
		PointType:  command.PointType,
		Amount:     command.Amount,
		Type:       contributionpointdomain.EntryTypeEarn,
		Reason:     command.Reason,
		SourceType: command.SourceType,
		SourceID:   command.SourceID,
	})
}

func (u *UseCase) Spend(ctx context.Context, command SpendCommand) (contributionpointdomain.LedgerEntry, error) {
	return u.Apply(ctx, ApplyCommand{
		UserID:     command.UserID,
		PointType:  command.PointType,
		Amount:     -command.Amount,
		Type:       contributionpointdomain.EntryTypeSpend,
		Reason:     command.Reason,
		SourceType: command.SourceType,
		SourceID:   command.SourceID,
	})
}

func (u *UseCase) GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return u.ledger.GetBalance(ctx, userID, pointType)
}

type ApplyCommand struct {
	UserID     user.ID
	PointType  contributionpointdomain.PointType
	Amount     int64
	Type       contributionpointdomain.EntryType
	Reason     string
	SourceType string
	SourceID   string
}

type EarnCommand struct {
	UserID     user.ID
	PointType  contributionpointdomain.PointType
	Amount     int64
	Reason     string
	SourceType string
	SourceID   string
}

type SpendCommand struct {
	UserID     user.ID
	PointType  contributionpointdomain.PointType
	Amount     int64
	Reason     string
	SourceType string
	SourceID   string
}
