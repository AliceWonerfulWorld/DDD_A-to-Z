package cp

import (
	"context"
	"testing"
	"time"

	cpdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/cp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeLedgerRepository struct {
	recorded []cpdomain.LedgerEntry
}

func (r *fakeLedgerRepository) Record(ctx context.Context, entry cpdomain.LedgerEntry) (cpdomain.LedgerEntry, error) {
	if err := ctx.Err(); err != nil {
		return cpdomain.LedgerEntry{}, err
	}

	r.recorded = append(r.recorded, entry)
	return entry.WithBalanceAfter(120), nil
}

type fakeIDGenerator struct {
	id string
}

func (g fakeIDGenerator) NewID() (string, error) {
	return g.id, nil
}

func TestUseCaseApply(t *testing.T) {
	t.Run("CP増減を履歴として記録しDBが確定した残高を返す", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "cp_ledger_1"})
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
		}

		entry, err := usecase.Apply(context.Background(), ApplyCommand{
			UserID:     user.ID("user_1"),
			Amount:     120,
			Type:       cpdomain.EntryTypeEarn,
			Reason:     "repository analysis reward",
			SourceType: "repository_analysis",
			SourceID:   "analysis_1",
		})
		if err != nil {
			t.Fatalf("Apply() がエラーを返しました: %v", err)
		}

		if entry.ID != "cp_ledger_1" {
			t.Fatalf("entry.ID = %q, 期待値 cp_ledger_1", entry.ID)
		}
		if entry.BalanceAfter != 120 {
			t.Fatalf("entry.BalanceAfter = %d, 期待値 120", entry.BalanceAfter)
		}
		if len(ledger.recorded) != 1 {
			t.Fatalf("記録された履歴数 = %d, 期待値 1", len(ledger.recorded))
		}
		if ledger.recorded[0].BalanceAfter != 0 {
			t.Fatalf("repository へ渡す前の BalanceAfter = %d, 期待値 0", ledger.recorded[0].BalanceAfter)
		}
	})
}
