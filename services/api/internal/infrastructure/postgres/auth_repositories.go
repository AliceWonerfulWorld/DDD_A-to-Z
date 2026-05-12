package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{db: db}
}

func (s *AuthStore) FindOrCreateByGitHub(ctx context.Context, profile user.GitHubProfile, now time.Time) (user.User, error) {
	if _, err := user.NewGitHubAccount(profile); err != nil {
		return user.User{}, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return user.User{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	userID := user.ID(fmt.Sprintf("github_%d", profile.GitHubID))
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO users (id, created_at, updated_at)
		VALUES ($1, $2, $2)
		ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at
	`, userID, now); err != nil {
		return user.User{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO github_accounts (github_id, user_id, username, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		ON CONFLICT (github_id) DO UPDATE
		SET username = EXCLUDED.username,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = EXCLUDED.updated_at
	`, profile.GitHubID, userID, profile.Username, profile.AvatarURL, now); err != nil {
		return user.User{}, err
	}

	appUser, ok, err := findUserByID(ctx, tx, userID)
	if err != nil {
		return user.User{}, err
	}
	if !ok {
		return user.User{}, sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return user.User{}, err
	}

	return appUser, nil
}

func (s *AuthStore) Save(ctx context.Context, session authapp.Session) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (token_hash, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (token_hash) DO UPDATE
		SET user_id = EXCLUDED.user_id,
			expires_at = EXCLUDED.expires_at
	`, tokenHash(session.Token), session.UserID, session.ExpiresAt, time.Now())
	return err
}

func (s *AuthStore) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	var appUser user.User
	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, ga.github_id, ga.username, ga.avatar_url, u.created_at, u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		JOIN github_accounts ga ON ga.user_id = u.id
		WHERE s.token_hash = $1
			AND s.expires_at > $2
	`, tokenHash(sessionToken), now).Scan(
		&appUser.ID,
		&appUser.GitHubAccount.GitHubID,
		&appUser.GitHubAccount.Username,
		&appUser.GitHubAccount.AvatarURL,
		&appUser.CreatedAt,
		&appUser.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, false, nil
	}
	if err != nil {
		return user.User{}, false, err
	}

	return appUser, true, nil
}

type userQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func findUserByID(ctx context.Context, db userQuerier, id user.ID) (user.User, bool, error) {
	var appUser user.User
	err := db.QueryRowContext(ctx, `
		SELECT u.id, ga.github_id, ga.username, ga.avatar_url, u.created_at, u.updated_at
		FROM users u
		JOIN github_accounts ga ON ga.user_id = u.id
		WHERE u.id = $1
	`, id).Scan(
		&appUser.ID,
		&appUser.GitHubAccount.GitHubID,
		&appUser.GitHubAccount.Username,
		&appUser.GitHubAccount.AvatarURL,
		&appUser.CreatedAt,
		&appUser.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, false, nil
	}
	if err != nil {
		return user.User{}, false, err
	}

	return appUser, true, nil
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
