package postgres

import (
	"context"
	"fmt"
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
		t.Skipf("PostgreSQL 結合テストをスキップします: %v", err)
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
	return db.Exec("SELECT 1 FROM users LIMIT 1").Error
}

func uniqueGitHubID() int64 {
	return 900_000_000_000 + time.Now().UnixNano()%100_000_000_000
}

func isMissingAuthSchemaError(err error) bool {
	message := err.Error()
	return strings.Contains(message, `relation "users" does not exist`) ||
		strings.Contains(message, "SQLSTATE 42P01")
}
