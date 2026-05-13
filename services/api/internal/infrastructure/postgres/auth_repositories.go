package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

type AuthStore struct {
	db *gorm.DB
}

func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{db: db}
}

func (s *AuthStore) FindOrCreateByGitHub(ctx context.Context, profile user.GitHubProfile, now time.Time) (user.User, error) {
	if _, err := user.NewGitHubAccount(profile); err != nil {
		return user.User{}, err
	}

	var appUser user.User
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userID := user.ID(fmt.Sprintf("github_%d", profile.GitHubID))
		if err := tx.Exec(`
			INSERT INTO users (id, created_at, updated_at)
			VALUES (?, ?, ?)
			ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at
		`, userID, now, now).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO github_accounts (github_id, user_id, username, avatar_url, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT (github_id) DO UPDATE
			SET username = EXCLUDED.username,
				avatar_url = EXCLUDED.avatar_url,
				updated_at = EXCLUDED.updated_at
		`, profile.GitHubID, userID, profile.Username, profile.AvatarURL, now, now).Error; err != nil {
			return err
		}

		foundUser, ok, err := findUserByID(tx, userID)
		if err != nil {
			return err
		}
		if !ok {
			return gorm.ErrRecordNotFound
		}

		appUser = foundUser
		return nil
	})
	if err != nil {
		return user.User{}, err
	}

	return appUser, nil
}

func (s *AuthStore) Save(ctx context.Context, session authapp.Session) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO sessions (token_hash, user_id, expires_at, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (token_hash) DO UPDATE
		SET user_id = EXCLUDED.user_id,
			expires_at = EXCLUDED.expires_at
	`, tokenHash(session.Token), session.UserID, session.ExpiresAt, time.Now()).Error
}

func (s *AuthStore) Delete(ctx context.Context, sessionToken string) error {
	return s.db.WithContext(ctx).Exec(`
		DELETE FROM sessions
		WHERE token_hash = ?
	`, tokenHash(sessionToken)).Error
}

func (s *AuthStore) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	var record userRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT u.id, ga.github_id, ga.username, ga.avatar_url, u.created_at, u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		JOIN github_accounts ga ON ga.user_id = u.id
		WHERE s.token_hash = ?
			AND s.expires_at > ?
	`, tokenHash(sessionToken), now).Scan(&record)
	if result.Error != nil {
		return user.User{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return user.User{}, false, nil
	}

	return record.toDomain(), true, nil
}

type userRecord struct {
	ID        user.ID   `gorm:"column:id"`
	GitHubID  int64     `gorm:"column:github_id"`
	Username  string    `gorm:"column:username"`
	AvatarURL string    `gorm:"column:avatar_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (r userRecord) toDomain() user.User {
	return user.User{
		ID: r.ID,
		GitHubAccount: user.GitHubAccount{
			GitHubID:  r.GitHubID,
			Username:  r.Username,
			AvatarURL: r.AvatarURL,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func findUserByID(db *gorm.DB, id user.ID) (user.User, bool, error) {
	var record userRecord
	result := db.Raw(`
		SELECT u.id, ga.github_id, ga.username, ga.avatar_url, u.created_at, u.updated_at
		FROM users u
		JOIN github_accounts ga ON ga.user_id = u.id
		WHERE u.id = ?
	`, id).Scan(&record)
	if result.Error != nil {
		return user.User{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return user.User{}, false, nil
	}

	return record.toDomain(), true, nil
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
