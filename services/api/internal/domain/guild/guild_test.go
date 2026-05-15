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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewGuild(tt.guild); err == nil {
				t.Fatal("NewGuild() error = nil, 期待値 エラー")
			}
		})
	}
}
