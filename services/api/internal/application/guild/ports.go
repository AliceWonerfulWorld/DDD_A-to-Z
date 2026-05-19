package guild

import (
	"context"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type Repository interface {
	ListGuilds(ctx context.Context) ([]guilddomain.Guild, error)
	FindGuildByID(ctx context.Context, guildID guilddomain.ID) (guilddomain.Guild, bool, error)
	FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error)
	ListActiveMembersByGuild(ctx context.Context, guildID guilddomain.ID) ([]guilddomain.MemberContribution, error)
	CreateMembership(ctx context.Context, membership guilddomain.Membership) error
	UpdateMembership(ctx context.Context, membership guilddomain.Membership) error
	CreateCPContribution(ctx context.Context, contribution guilddomain.CPContribution) error
	ListCPContributionsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.CPContribution, error)
	ListCPContributionsByUser(ctx context.Context, userID user.ID, limit int) ([]guilddomain.CPContribution, error)
}

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type IDGenerator interface {
	NewID() (string, error)
}

type CPSpender interface {
	Spend(ctx context.Context, command contributionpointapp.SpendCommand) (contributionpointdomain.LedgerEntry, error)
}

type CPContributionTransactioner interface {
	WithinCPContribution(ctx context.Context, run func(ctx context.Context, repository Repository, cp CPSpender) error) error
}
