package cp

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
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "cp ledger id is required",
		},
		{
			name:      "source_id が空ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			amount:    1,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "",
			wantError: "cp source id is required",
		},
		{
			name:      "amount が 0 ならエラー",
			id:        "ledger_1",
			userID:    "user_1",
			amount:    0,
			entryType: EntryTypeEarn,
			reason:    "reward",
			source:    "test",
			sourceID:  "source_1",
			wantError: "cp amount must not be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLedgerEntry(
				tt.id,
				tt.userID,
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
