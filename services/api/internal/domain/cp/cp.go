// Package cp owns Contribution Point account and ledger rules.
package cp

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

type LedgerEntry struct {
	ID           string
	UserID       user.ID
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
	amount int64,
	entryType EntryType,
	reason string,
	sourceType string,
	sourceID string,
	createdAt time.Time,
) (LedgerEntry, error) {
	if err := validateRequiredStrings(
		requiredString{name: "cp ledger id", value: id},
		requiredString{name: "user id", value: string(userID)},
		requiredString{name: "cp reason", value: reason},
		requiredString{name: "cp source type", value: sourceType},
		requiredString{name: "cp source id", value: sourceID},
	); err != nil {
		return LedgerEntry{}, err
	}
	if amount == 0 {
		return LedgerEntry{}, errors.New("cp amount must not be zero")
	}

	switch entryType {
	case EntryTypeEarn:
		if amount < 0 {
			return LedgerEntry{}, errors.New("earn cp amount must be positive")
		}
	case EntryTypeSpend:
		if amount > 0 {
			return LedgerEntry{}, errors.New("spend cp amount must be negative")
		}
	case EntryTypeAdjust:
	default:
		return LedgerEntry{}, errors.New("cp entry type is invalid")
	}

	return LedgerEntry{
		ID:         id,
		UserID:     userID,
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
