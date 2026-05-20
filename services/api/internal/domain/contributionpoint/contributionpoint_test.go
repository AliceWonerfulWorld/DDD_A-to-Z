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

func TestPlayerLevelFromTotalEarned(t *testing.T) {
	tests := []struct {
		name        string
		totalEarned int64
		want        int
	}{
		{name: "獲得 CP が 0 なら LV1", totalEarned: 0, want: 1},
		{name: "獲得 CP が負でも LV1", totalEarned: -100, want: 1},
		{name: "100 CP で LV2", totalEarned: 100, want: 2},
		{name: "400 CP で LV3", totalEarned: 400, want: 3},
		{name: "2500 CP で LV6", totalEarned: 2500, want: 6},
		{name: "10000 CP で LV11", totalEarned: 10000, want: 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlayerLevelFromTotalEarned(tt.totalEarned); got != tt.want {
				t.Fatalf("PlayerLevelFromTotalEarned(%d) = %d, 期待値 %d", tt.totalEarned, got, tt.want)
			}
		})
	}
}

func TestTotalEarnedForPlayerLevel(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  int64
	}{
		{name: "LV1 は 0 CP", level: 1, want: 0},
		{name: "LV0 も 0 CP", level: 0, want: 0},
		{name: "LV2 は 100 CP", level: 2, want: 100},
		{name: "LV6 は 2500 CP", level: 6, want: 2500},
		{name: "LV11 は 10000 CP", level: 11, want: 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TotalEarnedForPlayerLevel(tt.level); got != tt.want {
				t.Fatalf("TotalEarnedForPlayerLevel(%d) = %d, 期待値 %d", tt.level, got, tt.want)
			}
		})
	}
}

func TestNextPlayerLevelProgress(t *testing.T) {
	nextLevel, nextLevelTotalEarned, remaining := NextPlayerLevelProgress(2400)

	if nextLevel != 6 {
		t.Fatalf("nextLevel = %d, 期待値 6", nextLevel)
	}
	if nextLevelTotalEarned != 2500 {
		t.Fatalf("nextLevelTotalEarned = %d, 期待値 2500", nextLevelTotalEarned)
	}
	if remaining != 100 {
		t.Fatalf("remaining = %d, 期待値 100", remaining)
	}
}
