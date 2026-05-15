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
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
	"gorm.io/gorm"
)

func TestAuthStoreFindOrCreateByGitHub(t *testing.T) {
	t.Run("GitHubプロフィールからユーザーを作成または更新できる", func(t *testing.T) {
		ctx := context.Background()
		tx := beginPostgresTestTransaction(t, ctx)
		store := NewAuthStore(tx, newTestTokenCipher(t))

		githubID := uniqueGitHubID()
		now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
		profile := user.GitHubProfile{
			GitHubID:  githubID,
			Username:  "first-user",
			AvatarURL: "https://example.com/first.png",
		}

		created, err := store.FindOrCreateByGitHub(ctx, authapp.GitHubLogin{Profile: profile}, now)
		if err != nil {
			if isMissingSchemaError(err) {
				t.Skipf("PostgreSQL 結合テストをスキップします: auth schema が migrate されていません: %v", err)
			}
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

		updated, err := store.FindOrCreateByGitHub(ctx, authapp.GitHubLogin{Profile: updatedProfile}, updatedAt)
		if err != nil {
			if isMissingSchemaError(err) {
				t.Skipf("PostgreSQL 結合テストをスキップします: auth schema が migrate されていません: %v", err)
			}
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
		if err := tx.WithContext(ctx).Raw("SELECT balance FROM contribution_point_accounts WHERE user_id = ?", wantUserID).Scan(&balance).Error; err != nil {
			t.Fatalf("contribution_point_accounts の初期残高取得でエラーが発生しました: %v", err)
		}
		if balance != 0 {
			t.Fatalf("contribution_point_accounts.balance = %d, 期待値 0", balance)
		}
	})
}

func TestAuthStoreSessionLifecycle(t *testing.T) {
	t.Run("セッションの保存と取得ができる", func(t *testing.T) {
		ctx := context.Background()
		tx := beginPostgresTestTransaction(t, ctx)
		store := NewAuthStore(tx, newTestTokenCipher(t))

		now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
		profile := user.GitHubProfile{
			GitHubID:  uniqueGitHubID(),
			Username:  "session-user",
			AvatarURL: "https://example.com/session.png",
		}
		appUser, err := store.FindOrCreateByGitHub(ctx, authapp.GitHubLogin{Profile: profile}, now)
		if err != nil {
			if isMissingSchemaError(err) {
				t.Skipf("PostgreSQL 結合テストをスキップします: auth schema が migrate されていません: %v", err)
			}
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

func TestAuthStoreGitHubAccessToken(t *testing.T) {
	t.Run("GitHub access token を暗号化保存して復号取得できる", func(t *testing.T) {
		ctx := context.Background()
		tx := beginPostgresTestTransaction(t, ctx)
		store := NewAuthStore(tx, newTestTokenCipher(t))

		now := time.Date(2026, 5, 14, 16, 0, 0, 0, time.UTC)
		profile := user.GitHubProfile{
			GitHubID:  uniqueGitHubID(),
			Username:  "token-user",
			AvatarURL: "https://example.com/token.png",
		}
		appUser, err := store.FindOrCreateByGitHub(ctx, authapp.GitHubLogin{
			Profile:     profile,
			AccessToken: "secret-github-token",
		}, now)
		if err != nil {
			if isMissingSchemaError(err) {
				t.Skipf("PostgreSQL 結合テストをスキップします: token schema が migrate されていません: %v", err)
			}
			t.Fatalf("FindOrCreateByGitHub() がエラーを返しました: %v", err)
		}

		var stored string
		if err := tx.WithContext(ctx).Raw(`
			SELECT access_token_ciphertext
			FROM github_accounts
			WHERE user_id = ?
		`, appUser.ID).Scan(&stored).Error; err != nil {
			t.Fatalf("access_token_ciphertext の取得でエラーが発生しました: %v", err)
		}
		if stored == "" || stored == "secret-github-token" {
			t.Fatalf("access_token_ciphertext = %q, 暗号化された値が必要です", stored)
		}

		accessToken, ok, err := store.GitHubAccessToken(ctx, appUser.ID)
		if err != nil {
			t.Fatalf("GitHubAccessToken() がエラーを返しました: %v", err)
		}
		if !ok {
			t.Fatal("GitHubAccessToken() ok = false, 期待値 true")
		}
		if accessToken != "secret-github-token" {
			t.Fatalf("accessToken = %q, 期待値 secret-github-token", accessToken)
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
		if isMissingSchemaError(tx.Error) {
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
		if isMissingSchemaError(err) {
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
	if err := db.Exec("SELECT 1 FROM contribution_point_accounts LIMIT 1").Error; err != nil {
		return err
	}
	return db.Exec("SELECT 1 FROM contribution_point_ledger LIMIT 1").Error
}

func uniqueGitHubID() int64 {
	return 900_000_000_000 + time.Now().UnixNano()%100_000_000_000
}

func newTestTokenCipher(t *testing.T) *security.TokenCipher {
	t.Helper()

	cipher, err := security.NewTokenCipher("test-token-secret")
	if err != nil {
		t.Fatalf("NewTokenCipher() がエラーを返しました: %v", err)
	}

	return cipher
}

func isMissingSchemaError(err error) bool {
	message := err.Error()
	return strings.Contains(message, `relation "users" does not exist`) ||
		strings.Contains(message, `relation "contribution_point_accounts" does not exist`) ||
		strings.Contains(message, `relation "contribution_point_ledger" does not exist`) ||
		strings.Contains(message, `relation "github_repositories" does not exist`) ||
		strings.Contains(message, `relation "guilds" does not exist`) ||
		strings.Contains(message, `relation "guild_memberships" does not exist`) ||
		strings.Contains(message, `column "access_token_ciphertext" does not exist`) ||
		strings.Contains(message, "SQLSTATE 42P01") ||
		strings.Contains(message, "SQLSTATE 42703")
}
