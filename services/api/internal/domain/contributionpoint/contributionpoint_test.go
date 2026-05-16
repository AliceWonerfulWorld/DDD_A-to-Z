package contributionpoint

import (
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

func TestNewLedgerEntryValidation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		userID    user.ID
		pointType PointType
		amount    int64
		entryType EntryType
		reason    string
		source    string
		sourceID  string
		wantError string
	}{
		{
			name:      "id が空ならエラー",
			id:        "",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "contribution point ledger id is required",
		},
		{
			name:      "source_id が空ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "",
			wantError: "contribution point source id is required",
		},
		{
			name:      "amount が 0 ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    0,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "contribution point amount must not be zero",
		},
		{
			name:      "point_type が空ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: "",
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "point type is required",
		},
		{
			name:      "user_id が空ならエラー",
			id:        "ledger_1",
			userID:    "",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "user id is required",
		},
		{
			name:      "reason が空ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "",
			source:    "test",
			sourceID:  "source_1",
			wantError: "contribution point reason is required",
		},
		{
			name:      "source_type が空ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "",
			sourceID:  "source_1",
			wantError: "contribution point source type is required",
		},
		{
			name:      "entry_type が不正ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: "invalid",
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "contribution point entry type is invalid",
		},
		{
			name:      "earn なのに amount が負ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    -1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "earn contribution point amount must be positive",
		},
		{
			name:      "spend なのに amount が正ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			pointType: PointTypeCP,
			amount:    1,
			entryType: EntryTypeSpend,
			reason:    "cost",
			source:    "test",
			sourceID:  "source_1",
			wantError: "spend contribution point amount must be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLedgerEntry(
				tt.id,
				tt.userID,
				tt.pointType,
				tt.amount,
				tt.entryType,
				tt.reason,
				tt.source,
				tt.sourceID,
				time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC),
			)
			if err == nil {
				t.Fatal("NewLedgerEntry() error = nil, エラーを期待")
			}
			if err.Error() != tt.wantError {
				t.Fatalf("NewLedgerEntry() error = %q, 期待値 %q", err.Error(), tt.wantError)
			}
		})
	}
}
