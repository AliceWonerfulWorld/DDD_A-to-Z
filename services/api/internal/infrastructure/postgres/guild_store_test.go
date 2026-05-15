package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

func TestGuildStoreListGuilds(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}

	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	rustID := fmt.Sprintf("guild_rust_test_%d", uniqueGitHubID())
	goID := fmt.Sprintf("guild_go_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          rustID,
		Slug:        rustID,
		Name:        "Rust",
		Description: "安全性と速度を掲げるギルド。",
		Icon:        "RS",
		Color:       "#ff6b35",
		SortOrder:   2,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          goID,
		Slug:        goID,
		Name:        "Go",
		Description: "シンプルさと並列処理で前に進むギルド。",
		Icon:        "GO",
		Color:       "#00acd7",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	activeUser := createPostgresTestUser(t, ctx, tx)
	leftUser := createPostgresTestUser(t, ctx, tx)
	insertPostgresTestMembership(t, ctx, tx, "membership_active_"+string(activeUser.ID), activeUser.ID, goID, now, nil)
	leftAt := now.Add(time.Hour)
	insertPostgresTestMembership(t, ctx, tx, "membership_left_"+string(leftUser.ID), leftUser.ID, goID, now, &leftAt)

	guilds, err := store.ListGuilds(ctx)
	if err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild schema が migrate されていません: %v", err)
		}
		t.Fatalf("ListGuilds() がエラーを返しました: %v", err)
	}

	var goGuildIndex = -1
	var rustGuildIndex = -1
	for i, guild := range guilds {
		switch string(guild.ID) {
		case goID:
			goGuildIndex = i
			if guild.MemberCount != 1 {
				t.Fatalf("Go guild MemberCount = %d, 期待値 1", guild.MemberCount)
			}
		case rustID:
			rustGuildIndex = i
			if guild.MemberCount != 0 {
				t.Fatalf("Rust guild MemberCount = %d, 期待値 0", guild.MemberCount)
			}
		}
	}
	if goGuildIndex < 0 || rustGuildIndex < 0 {
		t.Fatalf("テスト用 guild が一覧に含まれていません: go=%d rust=%d", goGuildIndex, rustGuildIndex)
	}
	if goGuildIndex > rustGuildIndex {
		t.Fatalf("guild の並び順が sort_order ASC ではありません: go=%d rust=%d", goGuildIndex, rustGuildIndex)
	}
}

func TestGuildMembershipAllowsOnlyOneActiveGuildPerUser(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)

	now := time.Date(2026, 5, 15, 1, 0, 0, 0, time.UTC)
	firstGuildID := fmt.Sprintf("guild_first_test_%d", uniqueGitHubID())
	secondGuildID := fmt.Sprintf("guild_second_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          firstGuildID,
		Slug:        firstGuildID,
		Name:        "First",
		Description: "最初のテストギルド。",
		Icon:        "F",
		Color:       "#111111",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          secondGuildID,
		Slug:        secondGuildID,
		Name:        "Second",
		Description: "2つ目のテストギルド。",
		Icon:        "S",
		Color:       "#222222",
		SortOrder:   2,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	appUser := createPostgresTestUser(t, ctx, tx)
	insertPostgresTestMembership(t, ctx, tx, "membership_first_"+string(appUser.ID), appUser.ID, firstGuildID, now, nil)

	expectPostgresStatementError(t, tx, func() error {
		return tx.WithContext(ctx).Exec(`
			INSERT INTO guild_memberships (id, user_id, guild_id, joined_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "membership_second_"+string(appUser.ID), appUser.ID, secondGuildID, now, now, now).Error
	})
}

func TestGuildsRejectInvalidColor(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)

	now := time.Date(2026, 5, 15, 2, 0, 0, 0, time.UTC)
	validGuildID := fmt.Sprintf("guild_valid_color_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          validGuildID,
		Slug:        validGuildID,
		Name:        "Valid Color",
		Description: "色制約の前提を確認するテストギルド。",
		Icon:        "V",
		Color:       "#abcdef",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	expectPostgresStatementError(t, tx, func() error {
		return tx.WithContext(ctx).Exec(`
			INSERT INTO guilds (id, slug, name, description, icon, color, sort_order, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, "guild_invalid_color_"+validGuildID, "invalid_color_"+validGuildID, "Invalid Color", "不正な色を持つテストギルド。", "I", "3178c6", 2, now, now).Error
	})
}

func TestNewGuildStoreRejectsNilDB(t *testing.T) {
	store, err := NewGuildStore(nil)
	if err == nil {
		t.Fatal("NewGuildStore(nil) error = nil, 期待値 エラー")
	}
	if err.Error() != "db is nil" {
		t.Fatalf("NewGuildStore(nil) error = %q, 期待値 db is nil", err.Error())
	}
	if store != nil {
		t.Fatalf("NewGuildStore(nil) store = %#v, 期待値 nil", store)
	}
}

type testGuild struct {
	ID          string
	Slug        string
	Name        string
	Description string
	Icon        string
	Color       string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func insertPostgresTestGuild(t *testing.T, ctx context.Context, tx *gorm.DB, guild testGuild) {
	t.Helper()

	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO guilds (id, slug, name, description, icon, color, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, guild.ID, guild.Slug, guild.Name, guild.Description, guild.Icon, guild.Color, guild.SortOrder, guild.CreatedAt, guild.UpdatedAt).Error; err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild schema が migrate されていません: %v", err)
		}
		t.Fatalf("guilds INSERT でエラーが発生しました: %v", err)
	}
}

func insertPostgresTestMembership(t *testing.T, ctx context.Context, tx *gorm.DB, id string, userID user.ID, guildID string, joinedAt time.Time, leftAt *time.Time) {
	t.Helper()

	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO guild_memberships (id, user_id, guild_id, joined_at, left_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, userID, guildID, joinedAt, leftAt, joinedAt, joinedAt).Error; err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild membership schema が migrate されていません: %v", err)
		}
		t.Fatalf("guild_memberships INSERT でエラーが発生しました: %v", err)
	}
}
