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
	UserID      user.ID   `gorm:"column:user_id"`
	DisplayName string    `gorm:"column:display_name"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (r profileRecord) toDomain() domainprofile.Profile {
	return domainprofile.Profile{
		UserID:      r.UserID,
		DisplayName: r.DisplayName,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// レコード保存(Upsert)
func (s *ProfileStore) Save(ctx context.Context, p domainprofile.Profile) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO user_profiles (user_id, display_name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (user_id) DO UPDATE
		SET display_name = EXCLUDED.display_name,
		    updated_at   = EXCLUDED.updated_at
	`, p.UserID, p.DisplayName, p.CreatedAt, p.UpdatedAt).Error
}

// レコード取得
func (s *ProfileStore) FindByUserID(ctx context.Context, userID user.ID) (domainprofile.Profile, bool, error) {
	var rec profileRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT user_id, display_name, created_at, updated_at
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
