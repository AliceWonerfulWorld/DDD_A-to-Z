package guild

import (
	"context"
	"errors"
	"strings"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
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
	ErrPetAlreadyOwned             = errors.New("guild pet already owned")
)

const (
	defaultContributionHistoryLimit = 50
	defaultActivityLogLimit         = 20
	maxActivityLogLimit             = 50
)

type MyGuildDetails struct {
	Membership guilddomain.Membership
	Guild      guilddomain.Guild
	Members    []guilddomain.MemberContribution
}

type JoinGuildResult struct {
	Membership      guilddomain.MembershipWithGuild
	GrantedPet      *petdomain.Pet
	PetAlreadyOwned bool
}

type UseCase struct {
	repository      Repository
	current         CurrentUserRepository
	ids             IDGenerator
	petIDs          IDGenerator
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
	return NewUseCaseWithPetAndCPTransaction(repository, current, ids, ids, contributionIDs, cp, cpTransactioner)
}

func NewUseCaseWithPetAndCPTransaction(
	repository Repository,
	current CurrentUserRepository,
	ids IDGenerator,
	petIDs IDGenerator,
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
	if petIDs == nil {
		panic("guild pet id generator is required")
	}
	if contributionIDs == nil {
		panic("guild cp contribution id generator is required")
	}

	return &UseCase{
		repository:      repository,
		current:         current,
		ids:             ids,
		petIDs:          petIDs,
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

func (u *UseCase) JoinGuild(ctx context.Context, sessionToken string, guildID guilddomain.ID) (JoinGuildResult, error) {
	if strings.TrimSpace(sessionToken) == "" {
		return JoinGuildResult{}, ErrUnauthenticated
	}
	if strings.TrimSpace(string(guildID)) == "" {
		return JoinGuildResult{}, ErrGuildNotFound
	}

	now := u.now()
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, now)
	if err != nil {
		return JoinGuildResult{}, err
	}
	if !ok {
		return JoinGuildResult{}, ErrUnauthenticated
	}

	if membership, ok, err := u.repository.FindActiveMembershipByUserID(ctx, appUser.ID); err != nil {
		return JoinGuildResult{}, err
	} else if ok {
		result, grantErr := u.joinGuildResultWithPet(ctx, membership, now)
		if grantErr != nil {
			return JoinGuildResult{}, grantErr
		}
		return result, ErrAlreadyJoined
	}

	foundGuild, ok, err := u.repository.FindGuildByID(ctx, guildID)
	if err != nil {
		return JoinGuildResult{}, err
	}
	if !ok {
		return JoinGuildResult{}, ErrGuildNotFound
	}

	membershipID, err := u.ids.NewID()
	if err != nil {
		return JoinGuildResult{}, err
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
		return JoinGuildResult{}, err
	}

	if err := u.repository.CreateMembership(ctx, membership); err != nil {
		if errors.Is(err, ErrAlreadyJoined) {
			if existing, ok, findErr := u.repository.FindActiveMembershipByUserID(ctx, appUser.ID); findErr != nil {
				return JoinGuildResult{}, findErr
			} else if ok {
				result, grantErr := u.joinGuildResultWithPet(ctx, existing, now)
				if grantErr != nil {
					return JoinGuildResult{}, grantErr
				}
				return result, ErrAlreadyJoined
			}
		}
		return JoinGuildResult{}, err
	}

	result := JoinGuildResult{
		Membership: guilddomain.MembershipWithGuild{
			Membership: membership,
			Guild:      foundGuild,
		},
	}
	grantedPet, alreadyOwned, err := u.grantGuildPet(ctx, appUser.ID, foundGuild.ID, now)
	if err != nil {
		return JoinGuildResult{}, err
	}
	result.GrantedPet = grantedPet
	result.PetAlreadyOwned = alreadyOwned

	return result, nil
}

func (u *UseCase) joinGuildResultWithPet(ctx context.Context, membership guilddomain.MembershipWithGuild, now time.Time) (JoinGuildResult, error) {
	grantedPet, alreadyOwned, err := u.grantGuildPet(ctx, membership.Membership.UserID, membership.Membership.GuildID, now)
	if err != nil {
		return JoinGuildResult{}, err
	}

	return JoinGuildResult{
		Membership:      membership,
		GrantedPet:      grantedPet,
		PetAlreadyOwned: alreadyOwned,
	}, nil
}

func (u *UseCase) grantGuildPet(ctx context.Context, userID user.ID, guildID guilddomain.ID, now time.Time) (*petdomain.Pet, bool, error) {
	if _, ok, err := u.repository.FindPetByUserAndGuild(ctx, userID, guildID); err != nil {
		return nil, false, err
	} else if ok {
		return nil, true, nil
	}

	petID, err := u.petIDs.NewID()
	if err != nil {
		return nil, false, err
	}
	grantedPet, err := petdomain.NewPetFromGuild(petdomain.ID(petID), userID, guildID, now)
	if err != nil {
		return nil, false, err
	}

	if err := u.repository.CreatePet(ctx, grantedPet); err != nil {
		if errors.Is(err, ErrPetAlreadyOwned) {
			return nil, true, nil
		}
		return nil, false, err
	}

	return &grantedPet, false, nil
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

func (u *UseCase) ListGuildActivityLogs(ctx context.Context, sessionToken string, guildID guilddomain.ID, limit int) ([]guilddomain.ActivityLog, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(guildID)) == "" {
		return nil, ErrGuildNotFound
	}

	membership, ok, err := u.GetMyGuild(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrActiveMembershipNotFound
	}
	if membership.Membership.GuildID != guildID {
		return nil, ErrGuildAccessDenied
	}

	if limit <= 0 {
		limit = defaultActivityLogLimit
	}
	if limit > maxActivityLogLimit {
		limit = maxActivityLogLimit
	}

	return u.repository.ListActivityLogsByGuild(ctx, guildID, limit)
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
