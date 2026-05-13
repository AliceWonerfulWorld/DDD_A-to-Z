package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/database"
	"gorm.io/gorm"
)

func TestAuthStoreFindOrCreateByGitHub(t *testing.T) {
	t.Run("GitHubプロフィールからユーザーを作成または更新できる", func(t *testing.T) {
		ctx := context.Background()
		tx := beginPostgresTestTransaction(t, ctx)
		store := NewAuthStore(tx)

		githubID := uniqueGitHubID()
		now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
		profile := user.GitHubProfile{
			GitHubID:  githubID,
			Username:  "first-user",
			AvatarURL: "https://example.com/first.png",
		}

		created, err := store.FindOrCreateByGitHub(ctx, profile, now)
		if err != nil {
			t.Fatalf("FindOrCreateByGitHub() の作成でエラーが発生しました: %v", err)
		}

		wantUserID := user.ID(fmt.Sprintf("github_%d", githubID))
		if created.ID != wantUserID {
			t.Fatalf("created.ID = %q, 期待値 %q", created.ID, wantUserID)
		}
		if created.GitHubAccount != (user.GitHubAccount(profile)) {
			t.Fatalf("created.GitHubAccount = %#v, 期待値 %#v", created.GitHubAccount, user.GitHubAccount(profile))
		}

		updatedAt := now.Add(time.Hour)
		updatedProfile := user.GitHubProfile{
			GitHubID:  githubID,
			Username:  "updated-user",
			AvatarURL: "https://example.com/updated.png",
		}

		updated, err := store.FindOrCreateByGitHub(ctx, updatedProfile, updatedAt)
		if err != nil {
			t.Fatalf("FindOrCreateByGitHub() の更新でエラーが発生しました: %v", err)
		}

		if updated.ID != wantUserID {
			t.Fatalf("updated.ID = %q, 期待値 %q", updated.ID, wantUserID)
		}
		if updated.GitHubAccount != (user.GitHubAccount(updatedProfile)) {
			t.Fatalf("updated.GitHubAccount = %#v, 期待値 %#v", updated.GitHubAccount, user.GitHubAccount(updatedProfile))
		}
		if !updated.UpdatedAt.Equal(updatedAt) {
			t.Fatalf("updated.UpdatedAt = %s, 期待値 %s", updated.UpdatedAt, updatedAt)
		}

		var balance int64
		if err := tx.WithContext(ctx).Raw("SELECT balance FROM cp_accounts WHERE user_id = ?", wantUserID).Scan(&balance).Error; err != nil {
			t.Fatalf("cp_accounts の初期残高取得でエラーが発生しました: %v", err)
		}
		if balance != 0 {
			t.Fatalf("cp_accounts.balance = %d, 期待値 0", balance)
		}
	})
}

func TestAuthStoreSessionLifecycle(t *testing.T) {
	t.Run("セッションの保存と取得ができる", func(t *testing.T) {
		ctx := context.Background()
		tx := beginPostgresTestTransaction(t, ctx)
		store := NewAuthStore(tx)

		now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
		profile := user.GitHubProfile{
			GitHubID:  uniqueGitHubID(),
			Username:  "session-user",
			AvatarURL: "https://example.com/session.png",
		}
		appUser, err := store.FindOrCreateByGitHub(ctx, profile, now)
		if err != nil {
			t.Fatalf("FindOrCreateByGitHub() がエラーを返しました: %v", err)
		}

		activeToken := fmt.Sprintf("active-token-%d", profile.GitHubID)
		if err := store.Save(ctx, authapp.Session{
			Token:     activeToken,
			UserID:    appUser.ID,
			ExpiresAt: now.Add(time.Hour),
		}); err != nil {
			t.Fatalf("Save() の有効なセッション保存でエラーが発生しました: %v", err)
		}

		foundUser, ok, err := store.FindUserBySessionToken(ctx, activeToken, now)
		if err != nil {
			t.Fatalf("FindUserBySessionToken() の有効なセッション取得でエラーが発生しました: %v", err)
		}
		if !ok {
			t.Fatal("FindUserBySessionToken() の有効なセッション ok = false, 期待値 true")
		}
		if foundUser.ID != appUser.ID {
			t.Fatalf("foundUser.ID = %q, 期待値 %q", foundUser.ID, appUser.ID)
		}
		if foundUser.GitHubAccount != (user.GitHubAccount(profile)) {
			t.Fatalf("foundUser.GitHubAccount = %#v, 期待値 %#v", foundUser.GitHubAccount, user.GitHubAccount(profile))
		}

		expiredToken := fmt.Sprintf("expired-token-%d", profile.GitHubID)
		if err := store.Save(ctx, authapp.Session{
			Token:     expiredToken,
			UserID:    appUser.ID,
			ExpiresAt: now.Add(-time.Minute),
		}); err != nil {
			t.Fatalf("Save() の期限切れセッション保存でエラーが発生しました: %v", err)
		}

		_, ok, err = store.FindUserBySessionToken(ctx, expiredToken, now)
		if err != nil {
			t.Fatalf("FindUserBySessionToken() の期限切れセッション取得でエラーが発生しました: %v", err)
		}
		if ok {
			t.Fatal("FindUserBySessionToken() の期限切れセッション ok = true, 期待値 false")
		}
	})
}

func beginPostgresTestTransaction(t *testing.T, ctx context.Context) *gorm.DB {
	t.Helper()

	db, err := database.Open(ctx, config.DatabaseURLFromEnv())
	if err != nil {
		if os.Getenv("DATABASE_URL") == "" {
			t.Skipf("PostgreSQL 結合テストをスキップします: DATABASE_URL が未設定で、デフォルトDBへ接続できません: %v", err)
		}
		t.Fatalf("DATABASE_URL で指定された PostgreSQL へ接続できません: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() がエラーを返しました: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("sqlDB.Close() がエラーを返しました: %v", err)
		}
	})

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		if isMissingAuthSchemaError(tx.Error) {
			t.Skipf("PostgreSQL 結合テストをスキップします: auth schema が migrate されていません: %v", tx.Error)
		}
		t.Fatalf("Begin() がエラーを返しました: %v", tx.Error)
	}
	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil {
			t.Errorf("Rollback() がエラーを返しました: %v", err)
		}
	})

	if err := verifyAuthSchema(tx); err != nil {
		if isMissingAuthSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: auth schema が migrate されていません: %v", err)
		}
		t.Fatalf("auth schema の検証でエラーが発生しました: %v", err)
	}

	return tx
}

func verifyAuthSchema(db *gorm.DB) error {
	if err := db.Exec("SELECT 1 FROM users LIMIT 1").Error; err != nil {
		return err
	}
	if err := db.Exec("SELECT 1 FROM cp_accounts LIMIT 1").Error; err != nil {
		return err
	}
	return db.Exec("SELECT 1 FROM cp_ledger LIMIT 1").Error
}

func uniqueGitHubID() int64 {
	return 900_000_000_000 + time.Now().UnixNano()%100_000_000_000
}

func isMissingAuthSchemaError(err error) bool {
	message := err.Error()
	return strings.Contains(message, `relation "users" does not exist`) ||
		strings.Contains(message, `relation "cp_accounts" does not exist`) ||
		strings.Contains(message, `relation "cp_ledger" does not exist`) ||
		strings.Contains(message, "SQLSTATE 42P01")
}
