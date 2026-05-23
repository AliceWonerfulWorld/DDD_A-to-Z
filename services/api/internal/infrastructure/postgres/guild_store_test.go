package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
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
	cpStore := NewContributionPointStore(tx)
	activeCP := recordEarnedCP(t, ctx, cpStore, activeUser.ID, 349, now)
	recordEarnedCP(t, ctx, cpStore, leftUser.ID, 999, now)
	insertPostgresTestGuildCPContribution(t, ctx, tx, "guild_cp_list_"+string(activeUser.ID), activeUser.ID, goID, activeCP.ID, 349, now)

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
			if guild.TotalContributedCP != 349 {
				t.Fatalf("Go guild TotalContributedCP = %d, 期待値 349", guild.TotalContributedCP)
			}
		case rustID:
			rustGuildIndex = i
			if guild.MemberCount != 0 {
				t.Fatalf("Rust guild MemberCount = %d, 期待値 0", guild.MemberCount)
			}
			if guild.TotalContributedCP != 0 {
				t.Fatalf("Rust guild TotalContributedCP = %d, 期待値 0", guild.TotalContributedCP)
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

func TestGuildStoreUpdateMembershipLeavesMembership(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}

	now := time.Date(2026, 5, 15, 3, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_leave_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Leave Test",
		Description: "脱退更新のテストギルド。",
		Icon:        "L",
		Color:       "#123456",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	appUser := createPostgresTestUser(t, ctx, tx)
	membershipID := "membership_leave_" + string(appUser.ID)
	insertPostgresTestMembership(t, ctx, tx, membershipID, appUser.ID, guildID, now, nil)

	leftAt := now.Add(2 * time.Hour)
	err = store.UpdateMembership(ctx, guilddomain.Membership{
		ID:        guilddomain.MembershipID(membershipID),
		UserID:    appUser.ID,
		GuildID:   guilddomain.ID(guildID),
		JoinedAt:  now,
		LeftAt:    &leftAt,
		CreatedAt: now,
		UpdatedAt: leftAt,
	})
	if err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild membership schema が migrate されていません: %v", err)
		}
		t.Fatalf("UpdateMembership() がエラーを返しました: %v", err)
	}

	if _, ok, err := store.FindActiveMembershipByUserID(ctx, appUser.ID); err != nil {
		t.Fatalf("FindActiveMembershipByUserID() がエラーを返しました: %v", err)
	} else if ok {
		t.Fatal("脱退後は active membership が存在しない必要があります")
	}

	var row struct {
		LeftAt    *time.Time `gorm:"column:left_at"`
		UpdatedAt time.Time  `gorm:"column:updated_at"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT left_at, updated_at
		FROM guild_memberships
		WHERE id = ?
	`, membershipID).Scan(&row).Error; err != nil {
		t.Fatalf("guild_memberships の確認に失敗しました: %v", err)
	}
	if row.LeftAt == nil {
		t.Fatal("left_at が保存されている必要があります")
	}
	if !row.LeftAt.Equal(leftAt) {
		t.Fatalf("left_at = %v, 期待値 %v", row.LeftAt, leftAt)
	}
	if !row.UpdatedAt.Equal(leftAt) {
		t.Fatalf("updated_at = %v, 期待値 %v", row.UpdatedAt, leftAt)
	}
}

func TestGuildStoreUpdateMembershipRejectsAlreadyLeftMembership(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}

	now := time.Date(2026, 5, 15, 4, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_already_left_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Already Left Test",
		Description: "脱退済み membership の更新拒否テスト。",
		Icon:        "A",
		Color:       "#654321",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	appUser := createPostgresTestUser(t, ctx, tx)
	membershipID := "membership_already_left_" + string(appUser.ID)
	firstLeftAt := now.Add(time.Hour)
	insertPostgresTestMembership(t, ctx, tx, membershipID, appUser.ID, guildID, now, &firstLeftAt)

	secondLeftAt := now.Add(2 * time.Hour)
	err = store.UpdateMembership(ctx, guilddomain.Membership{
		ID:        guilddomain.MembershipID(membershipID),
		UserID:    appUser.ID,
		GuildID:   guilddomain.ID(guildID),
		JoinedAt:  now,
		LeftAt:    &secondLeftAt,
		CreatedAt: now,
		UpdatedAt: secondLeftAt,
	})
	if !errors.Is(err, guildapp.ErrActiveMembershipNotFound) {
		t.Fatalf("UpdateMembership() error = %v, 期待値 ErrActiveMembershipNotFound", err)
	}
}

func TestGuildStoreCreatePetRejectsDuplicateGuildPet(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}

	now := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_pet_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Pet Guild",
		Description: "ペット配布用のテストギルド。",
		Icon:        "PG",
		Color:       "#00acd7",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	appUser := createPostgresTestUser(t, ctx, tx)
	firstPet := petdomain.Pet{
		ID:        petdomain.ID("pet_first_" + string(appUser.ID)),
		UserID:    appUser.ID,
		GuildID:   guilddomain.ID(guildID),
		Attribute: petdomain.AttributeGo,
		Stats:     petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := store.CreatePet(ctx, firstPet); err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: player_pets schema が migrate されていません: %v", err)
		}
		t.Fatalf("CreatePet() がエラーを返しました: %v", err)
	}

	foundPet, ok, err := store.FindPetByUserAndGuild(ctx, appUser.ID, guilddomain.ID(guildID))
	if err != nil {
		t.Fatalf("FindPetByUserAndGuild() がエラーを返しました: %v", err)
	}
	if !ok {
		t.Fatal("FindPetByUserAndGuild() ok = false, 期待値 true")
	}
	if foundPet.ID != firstPet.ID {
		t.Fatalf("found pet id = %q, 期待値 %q", foundPet.ID, firstPet.ID)
	}

	duplicatePet := firstPet
	duplicatePet.ID = petdomain.ID("pet_duplicate_" + string(appUser.ID))
	if err := store.CreatePet(ctx, duplicatePet); !errors.Is(err, guildapp.ErrPetAlreadyOwned) {
		t.Fatalf("CreatePet() error = %v, 期待値 ErrPetAlreadyOwned", err)
	}
}

func TestGuildStoreCPContributions(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}
	cpStore := NewContributionPointStore(tx)

	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_cp_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "CP Test",
		Description: "CP投入履歴用のテストギルド。",
		Icon:        "CP",
		Color:       "#123abc",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	appUser := createPostgresTestUserWithPointAccount(t, ctx, tx)
	insertPostgresTestMembership(t, ctx, tx, "membership_cp_"+string(appUser.ID), appUser.ID, guildID, now, nil)

	earned, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_guild_cp_earn_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		100,
		contributionpointdomain.EntryTypeEarn,
		"guild cp contribution setup",
		"test",
		"guild_cp_setup_"+string(appUser.ID),
		now,
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	if _, err := cpStore.Record(ctx, earned); err != nil {
		t.Fatalf("CP獲得履歴保存でエラーが発生しました: %v", err)
	}
	spent, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_guild_cp_spend_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		-40,
		contributionpointdomain.EntryTypeSpend,
		"guild cp contribution",
		"guild_cp_contribution",
		"guild_cp_contribution_"+string(appUser.ID),
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	recordedSpend, err := cpStore.Record(ctx, spent)
	if err != nil {
		t.Fatalf("CP消費履歴保存でエラーが発生しました: %v", err)
	}

	contribution := guilddomain.CPContribution{
		ID:            guilddomain.CPContributionID("guild_cp_contribution_" + string(appUser.ID)),
		GuildID:       guilddomain.ID(guildID),
		UserID:        appUser.ID,
		PointLedgerID: recordedSpend.ID,
		Amount:        40,
		CreatedAt:     recordedSpend.CreatedAt,
	}
	if err := store.CreateCPContribution(ctx, contribution); err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild cp contribution schema が migrate されていません: %v", err)
		}
		t.Fatalf("CreateCPContribution() がエラーを返しました: %v", err)
	}

	byGuild, err := store.ListCPContributionsByGuild(ctx, guilddomain.ID(guildID), 10)
	if err != nil {
		t.Fatalf("ListCPContributionsByGuild() がエラーを返しました: %v", err)
	}
	if len(byGuild) != 1 {
		t.Fatalf("guild contributions length = %d, 期待値 1", len(byGuild))
	}
	if byGuild[0].Amount != 40 {
		t.Fatalf("guild contribution amount = %d, 期待値 40", byGuild[0].Amount)
	}

	byUser, err := store.ListCPContributionsByUser(ctx, appUser.ID, 10)
	if err != nil {
		t.Fatalf("ListCPContributionsByUser() がエラーを返しました: %v", err)
	}
	if len(byUser) != 1 {
		t.Fatalf("user contributions length = %d, 期待値 1", len(byUser))
	}

	guilds, err := store.ListGuilds(ctx)
	if err != nil {
		t.Fatalf("ListGuilds() がエラーを返しました: %v", err)
	}
	for _, guild := range guilds {
		if string(guild.ID) == guildID {
			if guild.TotalContributedCP != 60 {
				t.Fatalf("TotalContributedCP = %d, 期待値 60", guild.TotalContributedCP)
			}
			return
		}
	}
	t.Fatalf("テスト用 guild %q が一覧に含まれていません", guildID)
}

func TestGuildStoreCreateCPContributionRejectsMismatchedLedger(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}
	cpStore := NewContributionPointStore(tx)

	now := time.Date(2026, 5, 18, 13, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_cp_invalid_ledger_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Invalid Ledger Test",
		Description: "CP投入履歴とledgerの整合性テストギルド。",
		Icon:        "IL",
		Color:       "#123abc",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	appUser := createPostgresTestUserWithPointAccount(t, ctx, tx)

	earned, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_invalid_setup_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		100,
		contributionpointdomain.EntryTypeEarn,
		"setup",
		"test",
		"invalid_setup_"+string(appUser.ID),
		now,
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	if _, err := cpStore.Record(ctx, earned); err != nil {
		t.Fatalf("CP獲得履歴保存でエラーが発生しました: %v", err)
	}
	spent, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_invalid_spend_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		-40,
		contributionpointdomain.EntryTypeSpend,
		"guild cp contribution",
		"guild_cp_contribution",
		"different_contribution_id",
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	recordedSpend, err := cpStore.Record(ctx, spent)
	if err != nil {
		t.Fatalf("CP消費履歴保存でエラーが発生しました: %v", err)
	}

	err = store.CreateCPContribution(ctx, guilddomain.CPContribution{
		ID:            guilddomain.CPContributionID("guild_cp_contribution_invalid_" + string(appUser.ID)),
		GuildID:       guilddomain.ID(guildID),
		UserID:        appUser.ID,
		PointLedgerID: recordedSpend.ID,
		Amount:        40,
		CreatedAt:     recordedSpend.CreatedAt,
	})
	if !errors.Is(err, guildapp.ErrInvalidCPContributionLedger) {
		t.Fatalf("CreateCPContribution() error = %v, 期待値 ErrInvalidCPContributionLedger", err)
	}
}

func TestGuildStoreListActiveMembersByGuild(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}
	cpStore := NewContributionPointStore(tx)

	now := time.Date(2026, 5, 19, 8, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_members_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Members Test",
		Description: "ギルドメンバー一覧のテストギルド。",
		Icon:        "MT",
		Color:       "#0fbcf9",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	firstUser := createPostgresTestUserWithPointAccount(t, ctx, tx)
	secondUser := createPostgresTestUserWithPointAccount(t, ctx, tx)
	leftUser := createPostgresTestUserWithPointAccount(t, ctx, tx)
	insertPostgresTestMembership(t, ctx, tx, "membership_first_"+string(firstUser.ID), firstUser.ID, guildID, now.Add(-2*time.Hour), nil)
	insertPostgresTestMembership(t, ctx, tx, "membership_second_"+string(secondUser.ID), secondUser.ID, guildID, now.Add(-time.Hour), nil)
	leftAt := now
	insertPostgresTestMembership(t, ctx, tx, "membership_left_"+string(leftUser.ID), leftUser.ID, guildID, now.Add(-3*time.Hour), &leftAt)

	recordEarnedCP(t, ctx, cpStore, firstUser.ID, 70, now.Add(-90*time.Minute))
	recordEarnedCP(t, ctx, cpStore, secondUser.ID, 120, now.Add(-30*time.Minute))
	recordEarnedCP(t, ctx, cpStore, leftUser.ID, 999, now.Add(-15*time.Minute))

	members, err := store.ListActiveMembersByGuild(ctx, guilddomain.ID(guildID))
	if err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild member schema が migrate されていません: %v", err)
		}
		t.Fatalf("ListActiveMembersByGuild() がエラーを返しました: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("members length = %d, 期待値 2", len(members))
	}
	if members[0].UserID != secondUser.ID {
		t.Fatalf("members[0].UserID = %q, 期待値 %q", members[0].UserID, secondUser.ID)
	}
	if members[0].TotalEarnedCP != 120 {
		t.Fatalf("members[0].TotalEarnedCP = %d, 期待値 120", members[0].TotalEarnedCP)
	}
	if members[1].UserID != firstUser.ID {
		t.Fatalf("members[1].UserID = %q, 期待値 %q", members[1].UserID, firstUser.ID)
	}
}

func TestGuildCPContributionTransactionRollsBackSpendWhenContributionInsertFails(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	guildStore, err := NewGuildStore(tx)
	if err != nil {
		t.Fatalf("NewGuildStore() がエラーを返しました: %v", err)
	}
	cpStore := NewContributionPointStore(tx)
	authStore := NewAuthStore(tx, newTestTokenCipher(t))

	now := time.Date(2026, 5, 18, 14, 0, 0, 0, time.UTC)
	guildID := fmt.Sprintf("guild_cp_rollback_test_%d", uniqueGitHubID())
	insertPostgresTestGuild(t, ctx, tx, testGuild{
		ID:          guildID,
		Slug:        guildID,
		Name:        "Rollback Test",
		Description: "CP投入transactionのrollbackテストギルド。",
		Icon:        "RB",
		Color:       "#123abc",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	appUser := createPostgresTestUserWithPointAccount(t, ctx, tx)
	insertPostgresTestMembership(t, ctx, tx, "membership_cp_rollback_"+string(appUser.ID), appUser.ID, guildID, now, nil)
	if err := authStore.Save(ctx, authapp.Session{
		Token:     "session-token",
		UserID:    appUser.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	}); err != nil {
		t.Fatalf("session 保存でエラーが発生しました: %v", err)
	}

	earned, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_rollback_earn_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		100,
		contributionpointdomain.EntryTypeEarn,
		"setup",
		"test",
		"rollback_setup_"+string(appUser.ID),
		now,
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	if _, err := cpStore.Record(ctx, earned); err != nil {
		t.Fatalf("CP獲得履歴保存でエラーが発生しました: %v", err)
	}

	duplicateContributionID := "guild_cp_contribution_rollback_" + string(appUser.ID)
	existingSpend, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_rollback_existing_spend_"+string(appUser.ID),
		appUser.ID,
		contributionpointdomain.PointTypeCP,
		-10,
		contributionpointdomain.EntryTypeSpend,
		"guild cp contribution",
		"guild_cp_contribution",
		duplicateContributionID,
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	recordedExistingSpend, err := cpStore.Record(ctx, existingSpend)
	if err != nil {
		t.Fatalf("既存CP消費履歴保存でエラーが発生しました: %v", err)
	}
	if err := guildStore.CreateCPContribution(ctx, guilddomain.CPContribution{
		ID:            guilddomain.CPContributionID(duplicateContributionID),
		GuildID:       guilddomain.ID(guildID),
		UserID:        appUser.ID,
		PointLedgerID: recordedExistingSpend.ID,
		Amount:        10,
		CreatedAt:     recordedExistingSpend.CreatedAt,
	}); err != nil {
		t.Fatalf("既存CP投入履歴保存でエラーが発生しました: %v", err)
	}

	transactioner := NewGuildCPContributionTransactioner(tx, fixedPostgresTestIDGenerator{id: "point_ledger_rollback_new_spend_" + string(appUser.ID)})
	usecase := guildapp.NewUseCaseWithCPTransaction(
		guildStore,
		authStore,
		security.NewIDGenerator("membership_unused"),
		fixedPostgresTestIDGenerator{id: duplicateContributionID},
		contributionpointapp.NewUseCase(cpStore, fixedPostgresTestIDGenerator{id: "point_ledger_unused"}),
		transactioner,
	)

	_, err = usecase.ContributeCP(ctx, "session-token", 40)
	if err == nil {
		t.Fatal("ContributeCP() error = nil, 期待値 duplicate contribution error")
	}

	balance, err := cpStore.GetBalance(ctx, appUser.ID, contributionpointdomain.PointTypeCP)
	if err != nil {
		t.Fatalf("GetBalance() がエラーを返しました: %v", err)
	}
	if balance != 90 {
		t.Fatalf("balance = %d, 期待値 90", balance)
	}

	var newLedgerCount int64
	if err := tx.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM point_ledger
		WHERE id = ?
	`, "point_ledger_rollback_new_spend_"+string(appUser.ID)).Scan(&newLedgerCount).Error; err != nil {
		t.Fatalf("point_ledger 件数確認でエラーが発生しました: %v", err)
	}
	if newLedgerCount != 0 {
		t.Fatalf("rollback後の新規ledger件数 = %d, 期待値 0", newLedgerCount)
	}
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

type fixedPostgresTestIDGenerator struct {
	id string
}

func (g fixedPostgresTestIDGenerator) NewID() (string, error) {
	return g.id, nil
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

func insertPostgresTestGuildCPContribution(t *testing.T, ctx context.Context, tx *gorm.DB, id string, userID user.ID, guildID string, pointLedgerID string, amount int64, createdAt time.Time) {
	t.Helper()

	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO guild_cp_contributions (id, guild_id, user_id, point_ledger_id, amount, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, guildID, userID, pointLedgerID, amount, createdAt).Error; err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: guild cp contribution schema が migrate されていません: %v", err)
		}
		t.Fatalf("guild_cp_contributions INSERT でエラーが発生しました: %v", err)
	}
}

func recordEarnedCP(t *testing.T, ctx context.Context, store *ContributionPointStore, userID user.ID, amount int64, createdAt time.Time) contributionpointdomain.LedgerEntry {
	t.Helper()

	entry, err := contributionpointdomain.NewLedgerEntry(
		"point_ledger_earn_"+string(userID),
		userID,
		contributionpointdomain.PointTypeCP,
		amount,
		contributionpointdomain.EntryTypeEarn,
		"guild member cp setup",
		"test",
		"guild_member_cp_"+string(userID),
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewLedgerEntry() がエラーを返しました: %v", err)
	}
	recorded, err := store.Record(ctx, entry)
	if err != nil {
		t.Fatalf("CP獲得履歴保存でエラーが発生しました: %v", err)
	}
	return recorded
}
