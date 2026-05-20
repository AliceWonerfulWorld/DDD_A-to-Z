package contributionpoint

import (
	"context"
	"errors"
	"testing"
	"time"

	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeLedgerRepository struct {
	recorded []contributionpointdomain.LedgerEntry
	balance  int64
}

func (r *fakeLedgerRepository) Record(ctx context.Context, entry contributionpointdomain.LedgerEntry) (contributionpointdomain.LedgerEntry, error) {
	if err := ctx.Err(); err != nil {
		return contributionpointdomain.LedgerEntry{}, err
	}
	if r.balance+entry.Amount < 0 {
		return contributionpointdomain.LedgerEntry{}, ErrInsufficientBalance
	}

	r.recorded = append(r.recorded, entry)
	r.balance += entry.Amount
	return entry.WithBalanceAfter(r.balance), nil
}

func (r *fakeLedgerRepository) GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return r.balance, nil
}

func (r *fakeLedgerRepository) GetLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (*time.Time, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *fakeLedgerRepository) UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType, at time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

type fakeIDGenerator struct {
	id string
}

func (g fakeIDGenerator) NewID() (string, error) {
	return g.id, nil
}

type fakeIDGeneratorError struct{}

func (g fakeIDGeneratorError) NewID() (string, error) {
	return "", errors.New("ID生成に失敗しました")
}

func TestUseCaseApply(t *testing.T) {
	t.Run("コンテキストがキャンセル済みならエラーを返し履歴を残さない", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_1"})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := usecase.Apply(ctx, ApplyCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     100,
			Type:       contributionpointdomain.EntryTypeEarn,
			Reason:     "reward",
			SourceType: "test",
			SourceID:   "source_1",
		})
		if err == nil {
			t.Fatal("Apply() error = nil, エラーを期待")
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})

	t.Run("ID生成に失敗したらエラーを返し履歴を残さない", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGeneratorError{})

		_, err := usecase.Apply(context.Background(), ApplyCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     100,
			Type:       contributionpointdomain.EntryTypeEarn,
			Reason:     "reward",
			SourceType: "test",
			SourceID:   "source_1",
		})
		if err == nil {
			t.Fatal("Apply() error = nil, エラーを期待")
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})

	t.Run("ドメインバリデーションエラーは呼び出し元に伝播する", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_1"})

		_, err := usecase.Apply(context.Background(), ApplyCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     0,
			Type:       contributionpointdomain.EntryTypeEarn,
			Reason:     "reward",
			SourceType: "test",
			SourceID:   "source_1",
		})
		if err == nil {
			t.Fatal("Apply() error = nil, エラーを期待")
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})

	t.Run("ContributionPoint増減を履歴として記録しDBが確定した残高を返す", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_1"})
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
		}

		entry, err := usecase.Apply(context.Background(), ApplyCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     120,
			Type:       contributionpointdomain.EntryTypeEarn,
			Reason:     "repository analysis reward",
			SourceType: "repository_analysis",
			SourceID:   "analysis_1",
		})
		if err != nil {
			t.Fatalf("Apply() がエラーを返しました: %v", err)
		}

		if entry.ID != "contribution_point_ledger_1" {
			t.Fatalf("entry.ID = %q, 期待値 contribution_point_ledger_1", entry.ID)
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

func TestUseCaseEarn(t *testing.T) {
	t.Run("ContributionPointServiceとして生成して獲得処理を呼び出せる", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		service := NewService(ledger, fakeIDGenerator{id: "contribution_point_ledger_service_1"})

		entry, err := service.Earn(context.Background(), EarnCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     10,
			Reason:     "service api test",
			SourceType: "test",
			SourceID:   "service_1",
		})
		if err != nil {
			t.Fatalf("Earn() がエラーを返しました: %v", err)
		}
		if entry.BalanceAfter != 10 {
			t.Fatalf("entry.BalanceAfter = %d, 期待値 10", entry.BalanceAfter)
		}
	})

	t.Run("ContributionPoint獲得を共通サービス経由で履歴に記録する", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 10}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_earn_1"})
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
		}

		entry, err := usecase.Earn(context.Background(), EarnCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     120,
			Reason:     "repository analysis reward",
			SourceType: "repository_analysis",
			SourceID:   "analysis_1",
		})
		if err != nil {
			t.Fatalf("Earn() がエラーを返しました: %v", err)
		}

		if entry.Type != contributionpointdomain.EntryTypeEarn {
			t.Fatalf("entry.Type = %q, 期待値 %q", entry.Type, contributionpointdomain.EntryTypeEarn)
		}
		if entry.Amount != 120 {
			t.Fatalf("entry.Amount = %d, 期待値 120", entry.Amount)
		}
		if entry.BalanceAfter != 130 {
			t.Fatalf("entry.BalanceAfter = %d, 期待値 130", entry.BalanceAfter)
		}
	})

	t.Run("コンテキストがキャンセル済みなら Earn はエラーを返す", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_earn_cancel"})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := usecase.Earn(ctx, EarnCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     10,
			Reason:     "reward",
			SourceType: "test",
			SourceID:   "source_1",
		})
		if err == nil {
			t.Fatal("Earn() error = nil, エラーを期待")
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})

	t.Run("Go SP の獲得は CP とは独立して記録される", func(t *testing.T) {
		ledger := &fakeLedgerRepository{}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "point_ledger_go_sp_1"})
		goSP := contributionpointdomain.SPType("Go")

		entry, err := usecase.Earn(context.Background(), EarnCommand{
			UserID:     user.ID("user_1"),
			PointType:  goSP,
			Amount:     50,
			Reason:     "go commit reward",
			SourceType: "commit",
			SourceID:   "commit_1",
		})
		if err != nil {
			t.Fatalf("Earn() がエラーを返しました: %v", err)
		}
		if entry.PointType != goSP {
			t.Fatalf("entry.PointType = %v, 期待値 %v", entry.PointType, goSP)
		}
		if entry.BalanceAfter != 50 {
			t.Fatalf("entry.BalanceAfter = %d, 期待値 50", entry.BalanceAfter)
		}
	})
}

func TestUseCaseSpend(t *testing.T) {
	t.Run("ContributionPoint消費を共通サービス経由で負の履歴として記録する", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 100}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_spend_1"})
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
		}

		entry, err := usecase.Spend(context.Background(), SpendCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     40,
			Reason:     "guild build",
			SourceType: "guild",
			SourceID:   "guild_1",
		})
		if err != nil {
			t.Fatalf("Spend() がエラーを返しました: %v", err)
		}

		if entry.Type != contributionpointdomain.EntryTypeSpend {
			t.Fatalf("entry.Type = %q, 期待値 %q", entry.Type, contributionpointdomain.EntryTypeSpend)
		}
		if entry.Amount != -40 {
			t.Fatalf("entry.Amount = %d, 期待値 -40", entry.Amount)
		}
		if entry.BalanceAfter != 60 {
			t.Fatalf("entry.BalanceAfter = %d, 期待値 60", entry.BalanceAfter)
		}
	})

	t.Run("ContributionPoint不足なら不足エラーを返し履歴を残さない", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 30}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_spend_2"})

		_, err := usecase.Spend(context.Background(), SpendCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     40,
			Reason:     "guild build",
			SourceType: "guild",
			SourceID:   "guild_1",
		})
		if !errors.Is(err, ErrInsufficientBalance) {
			t.Fatalf("Spend() error = %v, 期待値 ErrInsufficientBalance", err)
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})

	t.Run("コンテキストがキャンセル済みなら Spend はエラーを返す", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 100}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_spend_cancel"})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := usecase.Spend(ctx, SpendCommand{
			UserID:     user.ID("user_1"),
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     10,
			Reason:     "cost",
			SourceType: "test",
			SourceID:   "source_1",
		})
		if err == nil {
			t.Fatal("Spend() error = nil, エラーを期待")
		}
		if len(ledger.recorded) != 0 {
			t.Fatalf("記録された履歴数 = %d, 期待値 0", len(ledger.recorded))
		}
	})
}

func TestUseCaseGetBalance(t *testing.T) {
	t.Run("コンテキストがキャンセル済みなら GetBalance はエラーを返す", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 75}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_1"})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := usecase.GetBalance(ctx, user.ID("user_1"), contributionpointdomain.PointTypeCP)
		if err == nil {
			t.Fatal("GetBalance() error = nil, エラーを期待")
		}
	})

	t.Run("現在残高を共通サービス経由で取得する", func(t *testing.T) {
		ledger := &fakeLedgerRepository{balance: 75}
		usecase := NewUseCase(ledger, fakeIDGenerator{id: "contribution_point_ledger_1"})

		balance, err := usecase.GetBalance(context.Background(), user.ID("user_1"), contributionpointdomain.PointTypeCP)
		if err != nil {
			t.Fatalf("GetBalance() がエラーを返しました: %v", err)
		}
		if balance != 75 {
			t.Fatalf("balance = %d, 期待値 75", balance)
		}
	})
}
