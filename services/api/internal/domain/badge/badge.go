package badge

import (
	"errors"
	"fmt"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type ConditionType string

const (
	ConditionTypeCPEarned ConditionType = "cp_earned"
)

type Badge struct {
	Slug          string
	Name          string
	Description   string
	Icon          string
	ConditionType ConditionType
	Threshold     int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBadge(
	slug string,
	name string,
	description string,
	icon string,
	conditionType ConditionType,
	threshold int64,
	createdAt time.Time,
	updatedAt time.Time,
) (Badge, error) {
	if err := validateRequiredStrings(
		requiredString{name: "badge slug", value: slug},
		requiredString{name: "badge name", value: name},
		requiredString{name: "badge description", value: description},
		requiredString{name: "badge icon", value: icon},
	); err != nil {
		return Badge{}, err
	}
	if conditionType == "" {
		return Badge{}, fmt.Errorf("condition type is required")
	}
	if threshold <= 0 {
		return Badge{}, fmt.Errorf("threshold must be positive")
	}

	return Badge{
		Slug:          slug,
		Name:          name,
		Description:   description,
		Icon:          icon,
		ConditionType: conditionType,
		Threshold:     threshold,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

type UserBadge struct {
	ID        string
	UserID    user.ID
	BadgeSlug string
	EarnedAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserBadge(
	id string,
	userID user.ID,
	badgeSlug string,
	earnedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (UserBadge, error) {
	if err := validateRequiredStrings(
		requiredString{name: "user badge id", value: id},
		requiredString{name: "user id", value: string(userID)},
		requiredString{name: "badge slug", value: badgeSlug},
	); err != nil {
		return UserBadge{}, err
	}
	if earnedAt.IsZero() {
		return UserBadge{}, errors.New("earned_at is required")
	}

	return UserBadge{
		ID:        id,
		UserID:    userID,
		BadgeSlug: badgeSlug,
		EarnedAt:  earnedAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

type UserBadgeWithBadge struct {
	UserBadge UserBadge
	Badge     Badge
}

type requiredString struct {
	name  string
	value string
}

func validateRequiredStrings(fields ...requiredString) error {
	for _, field := range fields {
		if field.value == "" {
			return fmt.Errorf("%s is required", field.name)
		}
	}
	return nil
}
