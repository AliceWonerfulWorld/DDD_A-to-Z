package badge

import (
	"context"
	"time"

	badgedomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/badge"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

var ErrBadgeAlreadyGranted = errBadgeAlreadyGranted{}

type errBadgeAlreadyGranted struct{}

func (errBadgeAlreadyGranted) Error() string { return "badge already granted" }

type UseCase struct {
	badges     BadgeRepository
	userBadges UserBadgeRepository
	ids        IDGenerator
	now        func() time.Time
}

func NewUseCase(
	badges BadgeRepository,
	userBadges UserBadgeRepository,
	ids IDGenerator,
) *UseCase {
	return &UseCase{
		badges:     badges,
		userBadges: userBadges,
		ids:        ids,
		now:        time.Now,
	}
}

type GrantResult struct {
	Badge      badgedomain.Badge
	UserBadge  badgedomain.UserBadge
	JustEarned bool
}

func (u *UseCase) CheckAndGrantBadges(ctx context.Context, userID user.ID, conditionType badgedomain.ConditionType, value int64) ([]GrantResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	allBadges, err := u.badges.FindByConditionType(ctx, conditionType)
	if err != nil {
		return nil, err
	}

	var eligible []badgedomain.Badge
	for _, b := range allBadges {
		if value >= b.Threshold {
			eligible = append(eligible, b)
		}
	}
	if len(eligible) == 0 {
		return nil, nil
	}

	grantedSlugs, err := u.userBadges.FindGrantedSlugsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	grantedSet := make(map[string]struct{}, len(grantedSlugs))
	for _, s := range grantedSlugs {
		grantedSet[s] = struct{}{}
	}

	var results []GrantResult
	for _, b := range eligible {
		if _, ok := grantedSet[b.Slug]; ok {
			results = append(results, GrantResult{
				Badge:      b,
				JustEarned: false,
			})
			continue
		}

		id, err := u.ids.NewID()
		if err != nil {
			return nil, err
		}

		now := u.now()
		ub, err := badgedomain.NewUserBadge(id, userID, b.Slug, now, now, now)
		if err != nil {
			return nil, err
		}

		saved, err := u.userBadges.Save(ctx, ub)
		if err != nil {
			return nil, err
		}

		results = append(results, GrantResult{
			Badge:      b,
			UserBadge:  saved,
			JustEarned: true,
		})
	}

	return results, nil
}

func (u *UseCase) ListUserBadges(ctx context.Context, userID user.ID) ([]badgedomain.UserBadgeWithBadge, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return u.userBadges.FindByUser(ctx, userID)
}

type BadgeReader interface {
	FindGrantedSlugsByUser(ctx context.Context, userID user.ID) ([]string, error)
}
