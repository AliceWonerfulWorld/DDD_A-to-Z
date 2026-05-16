package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContributionPointStore struct {
	db *gorm.DB
}

func NewContributionPointStore(db *gorm.DB) *ContributionPointStore {
	return &ContributionPointStore{db: db}
}

func (s *ContributionPointStore) Record(ctx context.Context, entry contributionpointdomain.LedgerEntry) (contributionpointdomain.LedgerEntry, error) {
	record := pointLedgerRecord{
		ID:         entry.ID,
		UserID:     string(entry.UserID),
		PointType:  entry.PointType,
		Amount:     entry.Amount,
		Type:       entry.Type,
		Reason:     entry.Reason,
		SourceType: entry.SourceType,
		SourceID:   entry.SourceID,
		CreatedAt:  entry.CreatedAt,
	}
	result := s.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&record)
	if result.Error != nil {
		return contributionpointdomain.LedgerEntry{}, mapContributionPointStoreError(result.Error)
	}
	if result.RowsAffected == 0 {
		return contributionpointdomain.LedgerEntry{}, gorm.ErrRecordNotFound
	}

	return record.toDomain(), nil
}

func (s *ContributionPointStore) GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error) {
	var record pointAccountRecord
	result := s.db.WithContext(ctx).
		Select("balance").
		Where("user_id = ? AND point_type = ?", userID, pointType).
		Take(&record)
	if result.Error != nil {
		return 0, result.Error
	}

	return record.Balance, nil
}

func mapContributionPointStoreError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) &&
		pgErr.Code == "23514" &&
		strings.Contains(pgErr.Message, "point balance cannot be negative") {
		return fmt.Errorf("%w: %v", contributionpointapp.ErrInsufficientBalance, err)
	}

	return err
}

type pointLedgerRecord struct {
	ID           string                            `gorm:"column:id"`
	UserID       string                            `gorm:"column:user_id"`
	PointType    contributionpointdomain.PointType `gorm:"column:point_type"`
	Amount       int64                             `gorm:"column:amount"`
	Type         contributionpointdomain.EntryType `gorm:"column:type"`
	Reason       string                            `gorm:"column:reason"`
	SourceType   string                            `gorm:"column:source_type"`
	SourceID     string                            `gorm:"column:source_id"`
	BalanceAfter int64                             `gorm:"column:balance_after"`
	CreatedAt    time.Time                         `gorm:"column:created_at"`
}

func (pointLedgerRecord) TableName() string {
	return "point_ledger"
}

func (r pointLedgerRecord) toDomain() contributionpointdomain.LedgerEntry {
	return contributionpointdomain.LedgerEntry{
		ID:           r.ID,
		UserID:       user.ID(r.UserID),
		PointType:    r.PointType,
		Amount:       r.Amount,
		Type:         r.Type,
		Reason:       r.Reason,
		SourceType:   r.SourceType,
		SourceID:     r.SourceID,
		BalanceAfter: r.BalanceAfter,
		CreatedAt:    r.CreatedAt,
	}
}

type pointAccountRecord struct {
	UserID    string                            `gorm:"column:user_id"`
	PointType contributionpointdomain.PointType `gorm:"column:point_type"`
	Balance   int64                             `gorm:"column:balance"`
}

func (pointAccountRecord) TableName() string {
	return "point_accounts"
}
