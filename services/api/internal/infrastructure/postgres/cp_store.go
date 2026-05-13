package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	cpapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/cp"
	cpdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/cp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CPStore struct {
	db *gorm.DB
}

func NewCPStore(db *gorm.DB) *CPStore {
	return &CPStore{db: db}
}

func (s *CPStore) Record(ctx context.Context, entry cpdomain.LedgerEntry) (cpdomain.LedgerEntry, error) {
	record := cpLedgerRecord{
		ID:         entry.ID,
		UserID:     string(entry.UserID),
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
		return cpdomain.LedgerEntry{}, mapCPStoreError(result.Error)
	}
	if result.RowsAffected == 0 {
		return cpdomain.LedgerEntry{}, gorm.ErrRecordNotFound
	}

	return record.toDomain(), nil
}

func (s *CPStore) GetBalance(ctx context.Context, userID user.ID) (int64, error) {
	var record cpAccountRecord
	result := s.db.WithContext(ctx).
		Select("balance").
		Where("user_id = ?", userID).
		Take(&record)
	if result.Error != nil {
		return 0, result.Error
	}

	return record.Balance, nil
}

func mapCPStoreError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) &&
		pgErr.Code == "23514" &&
		strings.Contains(pgErr.Message, "cp balance cannot be negative") {
		return fmt.Errorf("%w: %v", cpapp.ErrInsufficientBalance, err)
	}

	return err
}

type cpLedgerRecord struct {
	ID           string             `gorm:"column:id"`
	UserID       string             `gorm:"column:user_id"`
	Amount       int64              `gorm:"column:amount"`
	Type         cpdomain.EntryType `gorm:"column:type"`
	Reason       string             `gorm:"column:reason"`
	SourceType   string             `gorm:"column:source_type"`
	SourceID     string             `gorm:"column:source_id"`
	BalanceAfter int64              `gorm:"column:balance_after"`
	CreatedAt    time.Time          `gorm:"column:created_at"`
}

func (cpLedgerRecord) TableName() string {
	return "cp_ledger"
}

func (r cpLedgerRecord) toDomain() cpdomain.LedgerEntry {
	return cpdomain.LedgerEntry{
		ID:           r.ID,
		UserID:       user.ID(r.UserID),
		Amount:       r.Amount,
		Type:         r.Type,
		Reason:       r.Reason,
		SourceType:   r.SourceType,
		SourceID:     r.SourceID,
		BalanceAfter: r.BalanceAfter,
		CreatedAt:    r.CreatedAt,
	}
}

type cpAccountRecord struct {
	UserID  string `gorm:"column:user_id"`
	Balance int64  `gorm:"column:balance"`
}

func (cpAccountRecord) TableName() string {
	return "cp_accounts"
}
