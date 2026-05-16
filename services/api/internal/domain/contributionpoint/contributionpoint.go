// Package contributionpoint は ContributionPoint アカウントと台帳のルールを管理する。
package contributionpoint

import (
	"errors"
	"fmt"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type EntryType string

const (
	EntryTypeEarn   EntryType = "earn"
	EntryTypeSpend  EntryType = "spend"
	EntryTypeAdjust EntryType = "adjust"
)

// PointType はポイントの種別を表す。有効値は point_types マスターテーブルで管理し、
// 新しいSP種別の追加は INSERT のみで完結する。
type PointType string

const (
	PointTypeCP       PointType = "CP"
	PointTypeGolangSP PointType = "Golang_SP"
	PointTypeTSSP     PointType = "TypeScript_SP"
)

type LedgerEntry struct {
	ID           string
	UserID       user.ID
	PointType    PointType
	Amount       int64
	Type         EntryType
	Reason       string
	SourceType   string
	SourceID     string
	BalanceAfter int64
	CreatedAt    time.Time
}

func NewLedgerEntry(
	id string,
	userID user.ID,
	pointType PointType,
	amount int64,
	entryType EntryType,
	reason string,
	sourceType string,
	sourceID string,
	createdAt time.Time,
) (LedgerEntry, error) {
	if err := validateRequiredStrings(
		requiredString{name: "contribution point ledger id", value: id},
		requiredString{name: "user id", value: string(userID)},
		requiredString{name: "contribution point reason", value: reason},
		requiredString{name: "contribution point source type", value: sourceType},
		requiredString{name: "contribution point source id", value: sourceID},
	); err != nil {
		return LedgerEntry{}, err
	}
	if string(pointType) == "" {
		return LedgerEntry{}, fmt.Errorf("point type is required")
	}
	if amount == 0 {
		return LedgerEntry{}, errors.New("contribution point amount must not be zero")
	}

	switch entryType {
	case EntryTypeEarn:
		if amount < 0 {
			return LedgerEntry{}, errors.New("earn contribution point amount must be positive")
		}
	case EntryTypeSpend:
		if amount > 0 {
			return LedgerEntry{}, errors.New("spend contribution point amount must be negative")
		}
	case EntryTypeAdjust:
	default:
		return LedgerEntry{}, errors.New("contribution point entry type is invalid")
	}

	return LedgerEntry{
		ID:         id,
		UserID:     userID,
		PointType:  pointType,
		Amount:     amount,
		Type:       entryType,
		Reason:     reason,
		SourceType: sourceType,
		SourceID:   sourceID,
		CreatedAt:  createdAt,
	}, nil
}

func (e LedgerEntry) WithBalanceAfter(balanceAfter int64) LedgerEntry {
	e.BalanceAfter = balanceAfter
	return e
}

type requiredString struct {
	name  string
	value string
}

func validateRequiredStrings(fields ...requiredString) error {
	for _, field := range fields {
		if field.value == "" {
			return fmt.Errorf("%s is required", field.name)
		}
	}
	return nil
}
