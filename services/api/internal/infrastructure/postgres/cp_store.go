package postgres

import (
	"context"
	"time"

	cpdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/cp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

type CPStore struct {
	db *gorm.DB
}

func NewCPStore(db *gorm.DB) *CPStore {
	return &CPStore{db: db}
}

func (s *CPStore) Record(ctx context.Context, entry cpdomain.LedgerEntry) (cpdomain.LedgerEntry, error) {
	var record cpLedgerRecord
	result := s.db.WithContext(ctx).Raw(`
		INSERT INTO cp_ledger (id, user_id, amount, type, reason, source_type, source_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, user_id, amount, type, reason, source_type, source_id, balance_after, created_at
	`, entry.ID, entry.UserID, entry.Amount, entry.Type, entry.Reason, entry.SourceType, entry.SourceID, entry.CreatedAt).Scan(&record)
	if result.Error != nil {
		return cpdomain.LedgerEntry{}, result.Error
	}
	if result.RowsAffected == 0 {
		return cpdomain.LedgerEntry{}, gorm.ErrRecordNotFound
	}

	return record.toDomain(), nil
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
