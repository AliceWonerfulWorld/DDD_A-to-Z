package postgres

import (
	"context"
	"errors"
	"time"

	badgeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/badge"
	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	badgedomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/badge"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"

	"github.com/jackc/pgx/v5/pgconn"
)

type BadgeStore struct {
	db *gorm.DB
}

func NewBadgeStore(db *gorm.DB) *BadgeStore {
	return &BadgeStore{db: db}
}

func (s *BadgeStore) FindByConditionType(ctx context.Context, conditionType badgedomain.ConditionType) ([]badgedomain.Badge, error) {
	var records []badgeRecord
	err := s.db.WithContext(ctx).
		Where("condition_type = ?", string(conditionType)).
		Order("threshold ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	badges := make([]badgedomain.Badge, 0, len(records))
	for _, r := range records {
		badges = append(badges, r.toDomain())
	}
	return badges, nil
}

func (s *BadgeStore) Save(ctx context.Context, userBadge badgedomain.UserBadge) (badgedomain.UserBadge, error) {
	record := userBadgeRecord{
		ID:        userBadge.ID,
		UserID:    string(userBadge.UserID),
		BadgeSlug: userBadge.BadgeSlug,
		EarnedAt:  userBadge.EarnedAt,
		CreatedAt: userBadge.CreatedAt,
		UpdatedAt: userBadge.UpdatedAt,
	}
	result := s.db.WithContext(ctx).Create(&record)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return badgedomain.UserBadge{}, badgeapp.ErrBadgeAlreadyGranted
		}
		return badgedomain.UserBadge{}, result.Error
	}
	return userBadge, nil
}

func (s *BadgeStore) FindByUser(ctx context.Context, userID user.ID) ([]badgedomain.UserBadgeWithBadge, error) {
	type joinedRecord struct {
		ID        string    `gorm:"column:id"`
		UserID    string    `gorm:"column:user_id"`
		BadgeSlug string    `gorm:"column:badge_slug"`
		EarnedAt  time.Time `gorm:"column:earned_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		// badge fields
		BName          string    `gorm:"column:b_name"`
		BDescription   string    `gorm:"column:b_description"`
		BIcon          string    `gorm:"column:b_icon"`
		BConditionType string    `gorm:"column:b_condition_type"`
		BThreshold     int64     `gorm:"column:b_threshold"`
		BCreatedAt     time.Time `gorm:"column:b_created_at"`
		BUpdatedAt     time.Time `gorm:"column:b_updated_at"`
	}

	var rows []joinedRecord
	err := s.db.WithContext(ctx).Raw(`
		SELECT ub.id, ub.user_id, ub.badge_slug, ub.earned_at, ub.created_at, ub.updated_at,
		       b.name AS b_name, b.description AS b_description, b.icon AS b_icon,
		       b.condition_type AS b_condition_type, b.threshold AS b_threshold,
		       b.created_at AS b_created_at, b.updated_at AS b_updated_at
		FROM user_badges ub
		JOIN badges b ON b.slug = ub.badge_slug
		WHERE ub.user_id = ?
		ORDER BY ub.earned_at ASC
	`, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]badgedomain.UserBadgeWithBadge, 0, len(rows))
	for _, r := range rows {
		results = append(results, badgedomain.UserBadgeWithBadge{
			UserBadge: badgedomain.UserBadge{
				ID:        r.ID,
				UserID:    user.ID(r.UserID),
				BadgeSlug: r.BadgeSlug,
				EarnedAt:  r.EarnedAt,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			},
			Badge: badgedomain.Badge{
				Slug:          r.BadgeSlug,
				Name:          r.BName,
				Description:   r.BDescription,
				Icon:          r.BIcon,
				ConditionType: badgedomain.ConditionType(r.BConditionType),
				Threshold:     r.BThreshold,
				CreatedAt:     r.BCreatedAt,
				UpdatedAt:     r.BUpdatedAt,
			},
		})
	}
	return results, nil
}

func (s *BadgeStore) FindGrantedSlugsByUser(ctx context.Context, userID user.ID) ([]string, error) {
	var slugs []string
	err := s.db.WithContext(ctx).
		Model(&userBadgeRecord{}).
		Where("user_id = ?", userID).
		Pluck("badge_slug", &slugs).Error
	if err != nil {
		return nil, err
	}
	return slugs, nil
}

func (s *BadgeStore) FindBadgeSummariesByUser(ctx context.Context, userID user.ID) ([]mypageapp.BadgeSummary, error) {
	type row struct {
		BadgeSlug    string    `gorm:"column:badge_slug"`
		BName        string    `gorm:"column:b_name"`
		BDescription string    `gorm:"column:b_description"`
		BIcon        string    `gorm:"column:b_icon"`
		EarnedAt     time.Time `gorm:"column:earned_at"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT ub.badge_slug, b.name AS b_name, b.description AS b_description, b.icon AS b_icon, ub.earned_at
		FROM user_badges ub
		JOIN badges b ON b.slug = ub.badge_slug
		WHERE ub.user_id = ?
		ORDER BY ub.earned_at ASC
	`, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	summaries := make([]mypageapp.BadgeSummary, 0, len(rows))
	for _, r := range rows {
		summaries = append(summaries, mypageapp.BadgeSummary{
			Slug:        r.BadgeSlug,
			Name:        r.BName,
			Description: r.BDescription,
			Icon:        r.BIcon,
			EarnedAt:    r.EarnedAt.Format(time.RFC3339),
		})
	}
	return summaries, nil
}

type badgeRecord struct {
	Slug          string    `gorm:"column:slug"`
	Name          string    `gorm:"column:name"`
	Description   string    `gorm:"column:description"`
	Icon          string    `gorm:"column:icon"`
	ConditionType string    `gorm:"column:condition_type"`
	Threshold     int64     `gorm:"column:threshold"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (badgeRecord) TableName() string { return "badges" }

func (r badgeRecord) toDomain() badgedomain.Badge {
	return badgedomain.Badge{
		Slug:          r.Slug,
		Name:          r.Name,
		Description:   r.Description,
		Icon:          r.Icon,
		ConditionType: badgedomain.ConditionType(r.ConditionType),
		Threshold:     r.Threshold,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

type userBadgeRecord struct {
	ID        string    `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	BadgeSlug string    `gorm:"column:badge_slug"`
	EarnedAt  time.Time `gorm:"column:earned_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (userBadgeRecord) TableName() string { return "user_badges" }
