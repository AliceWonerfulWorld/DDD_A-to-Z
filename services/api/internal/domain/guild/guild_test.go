package guild

import (
	"testing"
	"time"
)

func TestNewGuild(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	valid := Guild{
		ID:          "guild_go",
		Slug:        "go",
		Name:        "Go",
		Description: "シンプルさと並列処理で前に進むギルド。",
		Icon:        "GO",
		Color:       "#00acd7",
		SortOrder:   1,
		MemberCount: 3,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if _, err := NewGuild(valid); err != nil {
		t.Fatalf("NewGuild() がエラーを返しました: %v", err)
	}

	tests := []struct {
		name  string
		guild Guild
	}{
		{name: "id が必須", guild: func() Guild {
			guild := valid
			guild.ID = ""
			return guild
		}()},
		{name: "slug が必須", guild: func() Guild {
			guild := valid
			guild.Slug = " "
			return guild
		}()},
		{name: "name が必須", guild: func() Guild {
			guild := valid
			guild.Name = ""
			return guild
		}()},
		{name: "description が必須", guild: func() Guild {
			guild := valid
			guild.Description = ""
			return guild
		}()},
		{name: "icon が必須", guild: func() Guild {
			guild := valid
			guild.Icon = ""
			return guild
		}()},
		{name: "color は hex 形式", guild: func() Guild {
			guild := valid
			guild.Color = "00acd7"
			return guild
		}()},
		{name: "sort_order は非負", guild: func() Guild {
			guild := valid
			guild.SortOrder = -1
			return guild
		}()},
		{name: "member_count は非負", guild: func() Guild {
			guild := valid
			guild.MemberCount = -1
			return guild
		}()},
		{name: "guild experience は非負", guild: func() Guild {
			guild := valid
			guild.GuildExperience = -1
			return guild
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewGuild(tt.guild); err == nil {
				t.Fatal("NewGuild() error = nil, 期待値 エラー")
			}
		})
	}
}

func TestGuildLevelProgressFromExperience(t *testing.T) {
	tests := []struct {
		name        string
		experience  int64
		wantLevel   int
		wantCurrent int64
		wantNext    int64
	}{
		{name: "0 exp is level 1", experience: 0, wantLevel: 1, wantCurrent: 0, wantNext: 1000},
		{name: "threshold reaches level 2", experience: 1000, wantLevel: 2, wantCurrent: 1000, wantNext: 3000},
		{name: "between thresholds stays level 3", experience: 4500, wantLevel: 3, wantCurrent: 3000, wantNext: 7000},
		{name: "max level caps at level 5", experience: 25000, wantLevel: 5, wantCurrent: 15000, wantNext: 15000},
		{name: "negative exp is treated as zero", experience: -10, wantLevel: 1, wantCurrent: 0, wantNext: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := GuildLevelProgressFromExperience(tt.experience)
			if progress.Level != tt.wantLevel {
				t.Fatalf("Level = %d, 期待値 %d", progress.Level, tt.wantLevel)
			}
			if progress.CurrentLevelExperience != tt.wantCurrent {
				t.Fatalf("CurrentLevelExperience = %d, 期待値 %d", progress.CurrentLevelExperience, tt.wantCurrent)
			}
			if progress.NextLevelExperience != tt.wantNext {
				t.Fatalf("NextLevelExperience = %d, 期待値 %d", progress.NextLevelExperience, tt.wantNext)
			}
		})
	}
}

func TestMembershipLeave(t *testing.T) {
	joinedAt := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	leftAt := joinedAt.Add(2 * time.Hour)

	membership, err := NewMembership(Membership{
		ID:        "membership_1",
		UserID:    "user_1",
		GuildID:   "guild_go",
		JoinedAt:  joinedAt,
		CreatedAt: joinedAt,
		UpdatedAt: joinedAt,
	})
	if err != nil {
		t.Fatalf("NewMembership() がエラーを返しました: %v", err)
	}

	leftMembership, err := membership.Leave(leftAt)
	if err != nil {
		t.Fatalf("Leave() がエラーを返しました: %v", err)
	}
	if leftMembership.LeftAt == nil {
		t.Fatal("left_at が設定されている必要があります")
	}
	if !leftMembership.LeftAt.Equal(leftAt) {
		t.Fatalf("left_at = %v, 期待値 %v", leftMembership.LeftAt, leftAt)
	}
	if !leftMembership.UpdatedAt.Equal(leftAt) {
		t.Fatalf("updated_at = %v, 期待値 %v", leftMembership.UpdatedAt, leftAt)
	}
}

func TestNewActivityLogRejectsInvalidType(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	_, err := NewActivityLog(ActivityLog{
		ID:         "user_1:issue:repo:1",
		UserID:     "user_1",
		Player:     "Alice",
		Type:       "issue",
		Repo:       "jyogi-web/DDD_A-to-Z",
		Message:    "Close issue",
		Language:   "Go",
		CP:         1,
		OccurredAt: now,
	})
	if err == nil {
		t.Fatal("NewActivityLog() error = nil, 期待値 invalid type error")
	}
	if err.Error() != "invalid guild activity log type: must be 'commit' or 'pull_request'" {
		t.Fatalf("error = %v, 期待値 invalid guild activity log type", err)
	}
}
