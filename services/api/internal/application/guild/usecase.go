package guild

import (
	"context"
	"errors"
	"strings"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

var (
	ErrUnauthenticated             = errors.New("unauthenticated")
	ErrGuildNotFound               = errors.New("guild not found")
	ErrAlreadyJoined               = errors.New("user already joined a guild")
	ErrActiveMembershipNotFound    = errors.New("active guild membership not found")
	ErrGuildAccessDenied           = errors.New("guild access denied")
	ErrInvalidCPContribution       = errors.New("guild cp contribution amount must be positive")
	ErrInvalidCPContributionLedger = errors.New("guild cp contribution ledger is invalid")
	ErrCPServiceUnavailable        = errors.New("contribution point service is unavailable")
)

const defaultContributionHistoryLimit = 50

type MyGuildDetails struct {
	Membership guilddomain.Membership
	Guild      guilddomain.Guild
	Members    []guilddomain.MemberContribution
}

type UseCase struct {
	repository      Repository
	current         CurrentUserRepository
	ids             IDGenerator
	contributionIDs IDGenerator
	cp              CPSpender
	cpTransactioner CPContributionTransactioner
	now             func() time.Time
}

func NewUseCase(repository Repository, current CurrentUserRepository, ids IDGenerator) *UseCase {
	return NewUseCaseWithCP(repository, current, ids, ids, nil)
}

func NewUseCaseWithCP(
	repository Repository,
	current CurrentUserRepository,
	ids IDGenerator,
	contributionIDs IDGenerator,
	cp CPSpender,
) *UseCase {
	return NewUseCaseWithCPTransaction(repository, current, ids, contributionIDs, cp, nil)
}

func NewUseCaseWithCPTransaction(
	repository Repository,
	current CurrentUserRepository,
	ids IDGenerator,
	contributionIDs IDGenerator,
	cp CPSpender,
	cpTransactioner CPContributionTransactioner,
) *UseCase {
	if repository == nil {
		panic("guild repository is required")
	}
	if current == nil {
		panic("current user repository is required")
	}
	if ids == nil {
		panic("guild membership id generator is required")
	}
	if contributionIDs == nil {
		panic("guild cp contribution id generator is required")
	}

	return &UseCase{
		repository:      repository,
		current:         current,
		ids:             ids,
		contributionIDs: contributionIDs,
		cp:              cp,
		cpTransactioner: cpTransactioner,
		now:             time.Now,
	}
}

func (u *UseCase) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return u.repository.ListGuilds(ctx)
}

func (u *UseCase) JoinGuild(ctx context.Context, sessionToken string, guildID guilddomain.ID) (guilddomain.MembershipWithGuild, error) {
	if strings.TrimSpace(sessionToken) == "" {
		return guilddomain.MembershipWithGuild{}, ErrUnauthenticated
	}
	if strings.TrimSpace(string(guildID)) == "" {
		return guilddomain.MembershipWithGuild{}, ErrGuildNotFound
	}

	now := u.now()
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, now)
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}
	if !ok {
		return guilddomain.MembershipWithGuild{}, ErrUnauthenticated
	}

	if membership, ok, err := u.repository.FindActiveMembershipByUserID(ctx, appUser.ID); err != nil {
		return guilddomain.MembershipWithGuild{}, err
	} else if ok {
		return membership, ErrAlreadyJoined
	}

	foundGuild, ok, err := u.repository.FindGuildByID(ctx, guildID)
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}
	if !ok {
		return guilddomain.MembershipWithGuild{}, ErrGuildNotFound
	}

	membershipID, err := u.ids.NewID()
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}
	membership, err := guilddomain.NewMembership(guilddomain.Membership{
		ID:        guilddomain.MembershipID(membershipID),
		UserID:    appUser.ID,
		GuildID:   foundGuild.ID,
		JoinedAt:  now,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}

	if err := u.repository.CreateMembership(ctx, membership); err != nil {
		if errors.Is(err, ErrAlreadyJoined) {
			if existing, ok, findErr := u.repository.FindActiveMembershipByUserID(ctx, appUser.ID); findErr != nil {
				return guilddomain.MembershipWithGuild{}, findErr
			} else if ok {
				return existing, ErrAlreadyJoined
			}
		}
		return guilddomain.MembershipWithGuild{}, err
	}

	return guilddomain.MembershipWithGuild{
		Membership: membership,
		Guild:      foundGuild,
	}, nil
}

func (u *UseCase) GetMyGuild(ctx context.Context, sessionToken string) (guilddomain.MembershipWithGuild, bool, error) {
	if strings.TrimSpace(sessionToken) == "" {
		return guilddomain.MembershipWithGuild{}, false, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return guilddomain.MembershipWithGuild{}, false, err
	}
	if !ok {
		return guilddomain.MembershipWithGuild{}, false, ErrUnauthenticated
	}

	return u.repository.FindActiveMembershipByUserID(ctx, appUser.ID)
}

func (u *UseCase) GetMyGuildDetails(ctx context.Context, sessionToken string) (MyGuildDetails, bool, error) {
	membership, ok, err := u.GetMyGuild(ctx, sessionToken)
	if err != nil || !ok {
		return MyGuildDetails{}, ok, err
	}

	members, err := u.repository.ListActiveMembersByGuild(ctx, membership.Membership.GuildID)
	if err != nil {
		return MyGuildDetails{}, false, err
	}

	return MyGuildDetails{
		Membership: membership.Membership,
		Guild:      membership.Guild,
		Members:    members,
	}, true, nil
}

func (u *UseCase) GetGuildDashboard(ctx context.Context, sessionToken string, guildID guilddomain.ID) (MyGuildDetails, error) {
	if strings.TrimSpace(string(guildID)) == "" {
		return MyGuildDetails{}, ErrGuildNotFound
	}

	membership, ok, err := u.GetMyGuild(ctx, sessionToken)
	if err != nil {
		return MyGuildDetails{}, err
	}
	if !ok {
		return MyGuildDetails{}, ErrActiveMembershipNotFound
	}
	if membership.Membership.GuildID != guildID {
		return MyGuildDetails{}, ErrGuildAccessDenied
	}

	members, err := u.repository.ListActiveMembersByGuild(ctx, membership.Membership.GuildID)
	if err != nil {
		return MyGuildDetails{}, err
	}

	return MyGuildDetails{
		Membership: membership.Membership,
		Guild:      membership.Guild,
		Members:    members,
	}, nil
}

func (u *UseCase) LeaveMyGuild(ctx context.Context, sessionToken string) error {
	if strings.TrimSpace(sessionToken) == "" {
		return ErrUnauthenticated
	}

	now := u.now()
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, now)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnauthenticated
	}

	membershipWithGuild, ok, err := u.repository.FindActiveMembershipByUserID(ctx, appUser.ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrActiveMembershipNotFound
	}

	leftMembership, err := membershipWithGuild.Membership.Leave(now)
	if err != nil {
		return err
	}

	return u.repository.UpdateMembership(ctx, leftMembership)
}

func (u *UseCase) ContributeCP(ctx context.Context, sessionToken string, amount int64) (guilddomain.CPContribution, error) {
	if err := ctx.Err(); err != nil {
		return guilddomain.CPContribution{}, err
	}
	if strings.TrimSpace(sessionToken) == "" {
		return guilddomain.CPContribution{}, ErrUnauthenticated
	}
	if amount <= 0 {
		return guilddomain.CPContribution{}, ErrInvalidCPContribution
	}
	if u.cp == nil {
		return guilddomain.CPContribution{}, ErrCPServiceUnavailable
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return guilddomain.CPContribution{}, err
	}
	if !ok {
		return guilddomain.CPContribution{}, ErrUnauthenticated
	}

	contributionID, err := u.contributionIDs.NewID()
	if err != nil {
		return guilddomain.CPContribution{}, err
	}

	var contribution guilddomain.CPContribution
	err = u.withCPContributionTransaction(ctx, func(ctx context.Context, repository Repository, cp CPSpender) error {
		membership, ok, err := repository.FindActiveMembershipByUserID(ctx, appUser.ID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrActiveMembershipNotFound
		}

		ledgerEntry, err := cp.Spend(ctx, contributionpointapp.SpendCommand{
			UserID:     appUser.ID,
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     amount,
			Reason:     "guild cp contribution",
			SourceType: "guild_cp_contribution",
			SourceID:   contributionID,
		})
		if err != nil {
			return err
		}

		contribution, err = guilddomain.NewCPContribution(guilddomain.CPContribution{
			ID:            guilddomain.CPContributionID(contributionID),
			GuildID:       membership.Membership.GuildID,
			UserID:        appUser.ID,
			PointLedgerID: ledgerEntry.ID,
			Amount:        amount,
			CreatedAt:     ledgerEntry.CreatedAt,
		})
		if err != nil {
			return err
		}

		return repository.CreateCPContribution(ctx, contribution)
	})
	if err != nil {
		return guilddomain.CPContribution{}, err
	}

	return contribution, nil
}

func (u *UseCase) withCPContributionTransaction(
	ctx context.Context,
	run func(ctx context.Context, repository Repository, cp CPSpender) error,
) error {
	if u.cpTransactioner != nil {
		return u.cpTransactioner.WithinCPContribution(ctx, run)
	}

	return run(ctx, u.repository, u.cp)
}

func (u *UseCase) ListGuildCPContributions(ctx context.Context, guildID guilddomain.ID) ([]guilddomain.CPContribution, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(guildID)) == "" {
		return nil, ErrGuildNotFound
	}
	if _, ok, err := u.repository.FindGuildByID(ctx, guildID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrGuildNotFound
	}

	return u.repository.ListCPContributionsByGuild(ctx, guildID, defaultContributionHistoryLimit)
}

func (u *UseCase) ListMyGuildCPContributions(ctx context.Context, sessionToken string) ([]guilddomain.CPContribution, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(sessionToken) == "" {
		return nil, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUnauthenticated
	}

	return u.repository.ListCPContributionsByUser(ctx, appUser.ID, defaultContributionHistoryLimit)
}
