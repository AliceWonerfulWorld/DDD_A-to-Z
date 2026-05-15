package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

func TestRepositoryStoreUpsertAndListRepositories(t *testing.T) {
	ctx := context.Background()
	tx := beginPostgresTestTransaction(t, ctx)
	store := NewRepositoryStore(tx)
	appUser := insertRepositoryTestUser(t, ctx, tx)

	now := time.Date(2026, 5, 14, 16, 0, 0, 0, time.UTC)
	pushedAt := now.Add(-time.Hour)
	repository := repositoryanalysis.Repository{
		GitHubID:        uniqueGitHubID(),
		UserID:          appUser.ID,
		Owner:           "octocat",
		Name:            "hello-world",
		FullName:        "octocat/hello-world",
		Private:         true,
		DefaultBranch:   "main",
		Language:        "Go",
		HTMLURL:         "https://github.com/octocat/hello-world",
		PushedAt:        &pushedAt,
		GitHubUpdatedAt: now,
		SyncedAt:        now,
	}

	if err := store.UpsertRepositories(ctx, []repositoryanalysis.Repository{repository}); err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: repository schema が migrate されていません: %v", err)
		}
		t.Fatalf("UpsertRepositories() の作成でエラーが発生しました: %v", err)
	}

	repository.Language = "TypeScript"
	repository.Private = false
	repository.SyncedAt = now.Add(time.Minute)
	if err := store.UpsertRepositories(ctx, []repositoryanalysis.Repository{repository}); err != nil {
		t.Fatalf("UpsertRepositories() の更新でエラーが発生しました: %v", err)
	}

	repositories, err := store.ListRepositories(ctx, appUser.ID)
	if err != nil {
		t.Fatalf("ListRepositories() がエラーを返しました: %v", err)
	}
	if len(repositories) != 1 {
		t.Fatalf("repositories length = %d, 期待値 1", len(repositories))
	}
	if repositories[0].Language != "TypeScript" {
		t.Fatalf("Language = %q, 期待値 TypeScript", repositories[0].Language)
	}
	if repositories[0].Private {
		t.Fatal("Private = true, 期待値 false")
	}
}

func insertRepositoryTestUser(t *testing.T, ctx context.Context, tx *gorm.DB) user.User {
	t.Helper()
	now := time.Date(2026, 5, 14, 15, 0, 0, 0, time.UTC)
	githubID := uniqueGitHubID()
	appUser := user.User{
		ID: user.ID(fmt.Sprintf("github_%d", githubID)),
		GitHubAccount: user.GitHubAccount{
			GitHubID:  githubID,
			Username:  fmt.Sprintf("repository-user-%d", githubID),
			AvatarURL: "https://example.com/repository.png",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO users (id, created_at, updated_at)
		VALUES (?, ?, ?)
	`, appUser.ID, now, now).Error; err != nil {
		t.Fatalf("users INSERT でエラーが発生しました: %v", err)
	}
	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO github_accounts (github_id, user_id, username, avatar_url, access_token_ciphertext, created_at, updated_at)
		VALUES (?, ?, ?, ?, '', ?, ?)
	`, appUser.GitHubAccount.GitHubID, appUser.ID, appUser.GitHubAccount.Username, appUser.GitHubAccount.AvatarURL, now, now).Error; err != nil {
		if isMissingSchemaError(err) {
			t.Skipf("PostgreSQL 結合テストをスキップします: repository schema が migrate されていません: %v", err)
		}
		t.Fatalf("github_accounts INSERT でエラーが発生しました: %v", err)
	}
	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO contribution_point_accounts (user_id, balance, created_at, updated_at)
		VALUES (?, 0, ?, ?)
	`, appUser.ID, now, now).Error; err != nil {
		t.Fatalf("contribution_point_accounts INSERT でエラーが発生しました: %v", err)
	}

	return appUser
}
