package badge

import (
	"context"

	badgedomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/badge"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type BadgeRepository interface {
	FindByConditionType(ctx context.Context, conditionType badgedomain.ConditionType) ([]badgedomain.Badge, error)
}

type UserBadgeRepository interface {
	Save(ctx context.Context, userBadge badgedomain.UserBadge) (badgedomain.UserBadge, error)
	FindByUser(ctx context.Context, userID user.ID) ([]badgedomain.UserBadgeWithBadge, error)
	FindGrantedSlugsByUser(ctx context.Context, userID user.ID) ([]string, error)
}

type IDGenerator interface {
	NewID() (string, error)
}

type BadgeGrantingChecker interface {
	CheckAndGrantBadges(ctx context.Context, userID user.ID, conditionType badgedomain.ConditionType, value int64) ([]GrantResult, error)
}
