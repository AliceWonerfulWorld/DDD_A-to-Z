package postgres

import (
	"context"
	"time"

	domainprofile "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/profile"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

// PostgreSQLのテーブルを操作する。GORMのDBインスタンス。
type ProfileStore struct {
	db *gorm.DB
}

func NewProfileStore(db *gorm.DB) *ProfileStore {
	return &ProfileStore{db: db}
}

type profileRecord struct {
	UserID            user.ID   `gorm:"column:user_id"`
	DisplayName       string    `gorm:"column:display_name"`
	SelectedBadgeSlug *string   `gorm:"column:selected_badge_slug"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (r profileRecord) toDomain() domainprofile.Profile {
	return domainprofile.Profile{
		UserID:            r.UserID,
		DisplayName:       r.DisplayName,
		SelectedBadgeSlug: r.SelectedBadgeSlug,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

// レコード保存(Upsert)
func (s *ProfileStore) Save(ctx context.Context, p domainprofile.Profile) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO user_profiles (user_id, display_name, selected_badge_slug, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (user_id) DO UPDATE
		SET display_name = EXCLUDED.display_name,
		    selected_badge_slug = EXCLUDED.selected_badge_slug,
		    updated_at   = EXCLUDED.updated_at
	`, p.UserID, p.DisplayName, p.SelectedBadgeSlug, p.CreatedAt, p.UpdatedAt).Error
}

func (s *ProfileStore) GetSelectedBadgeSlug(ctx context.Context, userID user.ID) (*string, error) {
	var slug *string
	err := s.db.WithContext(ctx).Raw(`
		SELECT selected_badge_slug FROM user_profiles WHERE user_id = ?
	`, userID).Scan(&slug).Error
	if err != nil {
		return nil, err
	}
	return slug, nil
}

func (s *ProfileStore) UpdateSelectedBadgeSlug(ctx context.Context, userID user.ID, badgeSlug *string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Exec(`
		UPDATE user_profiles SET selected_badge_slug = ?, updated_at = ? WHERE user_id = ?
	`, badgeSlug, now, userID).Error
}

// レコード取得
func (s *ProfileStore) FindByUserID(ctx context.Context, userID user.ID) (domainprofile.Profile, bool, error) {
	var rec profileRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT user_id, display_name, selected_badge_slug, created_at, updated_at
		FROM user_profiles
		WHERE user_id = ?
	`, userID).Scan(&rec)
	if result.Error != nil {
		return domainprofile.Profile{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return domainprofile.Profile{}, false, nil
	}
	return rec.toDomain(), true, nil
}
