package main

import (
	"context"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type cpManager struct {
	inner *contributionpointapp.UseCase
}

func newCPManager(inner *contributionpointapp.UseCase) *cpManager {
	return &cpManager{inner: inner}
}

func (m *cpManager) Earn(ctx context.Context, userID user.ID, amount int64, reason, sourceType, sourceID string) error {
	_, err := m.inner.Earn(ctx, contributionpointapp.EarnCommand{
		UserID:     userID,
		PointType:  contributionpointdomain.PointTypeCP,
		Amount:     amount,
		Reason:     reason,
		SourceType: sourceType,
		SourceID:   sourceID,
	})
	return err
}

func (m *cpManager) GetBalance(ctx context.Context, userID user.ID) (int64, error) {
	return m.inner.GetBalance(ctx, userID, contributionpointdomain.PointTypeCP)
}

func (m *cpManager) GetLastAnalyzedAt(ctx context.Context, userID user.ID) (*time.Time, error) {
	return m.inner.GetLastAnalyzedAt(ctx, userID, contributionpointdomain.PointTypeCP)
}

func (m *cpManager) UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, at time.Time) error {
	return m.inner.UpdateLastAnalyzedAt(ctx, userID, contributionpointdomain.PointTypeCP, at)
}

type mypageCPReader struct {
	balance interface {
		GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
	}
	totals interface {
		GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
		GetTotalSpent(ctx context.Context, userID user.ID) (int64, error)
	}
}

func newMypageCPReader(
	balance interface {
		GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
	},
	totals interface {
		GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
		GetTotalSpent(ctx context.Context, userID user.ID) (int64, error)
	},
) *mypageCPReader {
	return &mypageCPReader{
		balance: balance,
		totals:  totals,
	}
}

func (r *mypageCPReader) GetBalance(ctx context.Context, userID user.ID) (int64, error) {
	return r.balance.GetBalance(ctx, userID, contributionpointdomain.PointTypeCP)
}

func (r *mypageCPReader) GetTotalEarned(ctx context.Context, userID user.ID) (int64, error) {
	return r.totals.GetTotalEarned(ctx, userID)
}

func (r *mypageCPReader) GetTotalSpent(ctx context.Context, userID user.ID) (int64, error) {
	return r.totals.GetTotalSpent(ctx, userID)
}
