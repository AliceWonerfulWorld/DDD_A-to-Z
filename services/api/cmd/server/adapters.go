package main

import (
	"context"
	"time"

	badgeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/badge"
	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	badgedomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/badge"
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

func (m *cpManager) EarnSP(ctx context.Context, userID user.ID, language string, amount int64, reason, sourceType, sourceID string) error {
	_, err := m.inner.Earn(ctx, contributionpointapp.EarnCommand{
		UserID:     userID,
		PointType:  contributionpointdomain.SPType(language),
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

type homeCPDataProvider struct {
	balance interface {
		GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
	}
	totals interface {
		GetTodayEarned(ctx context.Context, userID user.ID) (int64, error)
		GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
	}
}

func newHomeCPDataProvider(
	balance interface {
		GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
	},
	totals interface {
		GetTodayEarned(ctx context.Context, userID user.ID) (int64, error)
		GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
	},
) *homeCPDataProvider {
	return &homeCPDataProvider{
		balance: balance,
		totals:  totals,
	}
}

func (p *homeCPDataProvider) GetBalance(ctx context.Context, userID user.ID) (int64, error) {
	return p.balance.GetBalance(ctx, userID, contributionpointdomain.PointTypeCP)
}

func (p *homeCPDataProvider) GetTodayEarned(ctx context.Context, userID user.ID) (int64, error) {
	return p.totals.GetTodayEarned(ctx, userID)
}

func (p *homeCPDataProvider) GetTotalEarned(ctx context.Context, userID user.ID) (int64, error) {
	return p.totals.GetTotalEarned(ctx, userID)
}

type mypageBadgeReader struct {
	inner interface {
		FindBadgeSummariesByUser(ctx context.Context, userID user.ID) ([]mypageapp.BadgeSummary, error)
	}
}

func newMypageBadgeReader(inner interface {
	FindBadgeSummariesByUser(ctx context.Context, userID user.ID) ([]mypageapp.BadgeSummary, error)
}) *mypageBadgeReader {
	return &mypageBadgeReader{inner: inner}
}

func (r *mypageBadgeReader) ListUserBadges(ctx context.Context, userID user.ID) ([]mypageapp.BadgeSummary, error) {
	return r.inner.FindBadgeSummariesByUser(ctx, userID)
}

type mypageBadgeGrantingChecker struct {
	inner *badgeapp.UseCase
}

func newMypageBadgeGrantingChecker(inner *badgeapp.UseCase) *mypageBadgeGrantingChecker {
	return &mypageBadgeGrantingChecker{inner: inner}
}

func (c *mypageBadgeGrantingChecker) CheckAndGrantBadges(ctx context.Context, userID user.ID, conditionType badgedomain.ConditionType, value int64) (int, error) {
	results, err := c.inner.CheckAndGrantBadges(ctx, userID, conditionType, value)
	if err != nil {
		return 0, err
	}
	return len(results), nil
}
