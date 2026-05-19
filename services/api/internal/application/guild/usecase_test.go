package guild

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type testRepository struct {
	guilds           []guilddomain.Guild
	activeMembership *guilddomain.MembershipWithGuild
	created          *guilddomain.Membership
	updated          *guilddomain.Membership
	contribution     *guilddomain.CPContribution
	contributions    []guilddomain.CPContribution
}

func (r testRepository) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	return r.guilds, nil
}

func (r testRepository) FindGuildByID(ctx context.Context, guildID guilddomain.ID) (guilddomain.Guild, bool, error) {
	for _, guild := range r.guilds {
		if guild.ID == guildID {
			return guild, true, nil
		}
	}

	return guilddomain.Guild{}, false, nil
}

func (r testRepository) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	if r.activeMembership == nil {
		return guilddomain.MembershipWithGuild{}, false, nil
	}
	if r.activeMembership.Membership.UserID != userID {
		return guilddomain.MembershipWithGuild{}, false, nil
	}

	return *r.activeMembership, true, nil
}

func (r *testRepository) CreateMembership(ctx context.Context, membership guilddomain.Membership) error {
	r.created = &membership
	return nil
}

func (r *testRepository) UpdateMembership(ctx context.Context, membership guilddomain.Membership) error {
	r.updated = &membership
	return nil
}

func (r *testRepository) CreateCPContribution(ctx context.Context, contribution guilddomain.CPContribution) error {
	r.contribution = &contribution
	r.contributions = append(r.contributions, contribution)
	return nil
}

func (r testRepository) ListCPContributionsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.CPContribution, error) {
	var contributions []guilddomain.CPContribution
	for _, contribution := range r.contributions {
		if contribution.GuildID == guildID {
			contributions = append(contributions, contribution)
		}
	}

	return contributions, nil
}

func (r testRepository) ListCPContributionsByUser(ctx context.Context, userID user.ID, limit int) ([]guilddomain.CPContribution, error) {
	var contributions []guilddomain.CPContribution
	for _, contribution := range r.contributions {
		if contribution.UserID == userID {
			contributions = append(contributions, contribution)
		}
	}

	return contributions, nil
}

type testCurrentUserRepository struct {
	appUser user.User
	ok      bool
}

func (r testCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return r.appUser, r.ok, nil
}

type testIDGenerator struct {
	id string
}

func (g testIDGenerator) NewID() (string, error) {
	return g.id, nil
}

type testCPSpender struct {
	err     error
	command contributionpointapp.SpendCommand
	entry   contributionpointdomain.LedgerEntry
}

func (s *testCPSpender) Spend(ctx context.Context, command contributionpointapp.SpendCommand) (contributionpointdomain.LedgerEntry, error) {
	s.command = command
	if s.err != nil {
		return contributionpointdomain.LedgerEntry{}, s.err
	}
	if s.entry.ID == "" {
		s.entry = contributionpointdomain.LedgerEntry{
			ID:           "cp_ledger_1",
			UserID:       command.UserID,
			PointType:    command.PointType,
			Amount:       -command.Amount,
			Type:         contributionpointdomain.EntryTypeSpend,
			Reason:       command.Reason,
			SourceType:   command.SourceType,
			SourceID:     command.SourceID,
			BalanceAfter: 60,
			CreatedAt:    time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC),
		}
	}

	return s.entry, nil
}

func TestUseCaseListGuilds(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	expected := []guilddomain.Guild{{
		ID:          "guild_typescript",
		Slug:        "typescript",
		Name:        "TypeScript",
		Description: "型の力で支えるギルド。",
		Icon:        "TS",
		Color:       "#3178c6",
		SortOrder:   5,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}
	usecase := NewUseCase(&testRepository{guilds: expected}, testCurrentUserRepository{ok: true}, testIDGenerator{id: "membership_1"})

	guilds, err := usecase.ListGuilds(context.Background())
	if err != nil {
		t.Fatalf("ListGuilds() がエラーを返しました: %v", err)
	}
	if len(guilds) != 1 {
		t.Fatalf("guilds length = %d, 期待値 1", len(guilds))
	}
	if guilds[0].Slug != "typescript" {
		t.Fatalf("Slug = %q, 期待値 typescript", guilds[0].Slug)
	}
}

func TestUseCaseJoinGuild(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	targetGuild := guilddomain.Guild{
		ID:          "guild_go",
		Slug:        "go",
		Name:        "Go",
		Description: "シンプルさと並列処理で前に進むギルド。",
		Icon:        "GO",
		Color:       "#00acd7",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	repository := &testRepository{guilds: []guilddomain.Guild{targetGuild}}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_1"})
	usecase.now = func() time.Time { return now }

	membership, err := usecase.JoinGuild(context.Background(), "session-token", "guild_go")
	if err != nil {
		t.Fatalf("JoinGuild() がエラーを返しました: %v", err)
	}
	if membership.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", membership.Guild.ID)
	}
	if repository.created == nil {
		t.Fatal("CreateMembership() が呼ばれる必要があります")
	}
	if repository.created.UserID != "user_1" {
		t.Fatalf("created user id = %q, 期待値 user_1", repository.created.UserID)
	}
	if !repository.created.JoinedAt.Equal(now) {
		t.Fatalf("joined_at = %v, 期待値 %v", repository.created.JoinedAt, now)
	}
}

func TestUseCaseJoinGuildRejectsAlreadyJoinedUser(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	existing := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    "user_1",
			GuildID:   "guild_go",
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guilddomain.Guild{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	usecase := NewUseCase(&testRepository{activeMembership: &existing}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_2"})

	_, err := usecase.JoinGuild(context.Background(), "session-token", "guild_python")
	if !errors.Is(err, ErrAlreadyJoined) {
		t.Fatalf("JoinGuild() error = %v, 期待値 ErrAlreadyJoined", err)
	}
}

func TestUseCaseJoinGuildRejectsUnknownGuild(t *testing.T) {
	usecase := NewUseCase(&testRepository{}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_1"})

	_, err := usecase.JoinGuild(context.Background(), "session-token", "guild_missing")
	if !errors.Is(err, ErrGuildNotFound) {
		t.Fatalf("JoinGuild() error = %v, 期待値 ErrGuildNotFound", err)
	}
}

func TestUseCaseLeaveMyGuild(t *testing.T) {
	now := time.Date(2026, 5, 16, 13, 0, 0, 0, time.UTC)
	activeMembership := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    "user_1",
			GuildID:   "guild_go",
			JoinedAt:  now.Add(-time.Hour),
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		},
		Guild: guilddomain.Guild{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now.Add(-time.Hour),
			UpdatedAt:   now.Add(-time.Hour),
		},
	}
	repository := &testRepository{activeMembership: &activeMembership}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})
	usecase.now = func() time.Time { return now }

	if err := usecase.LeaveMyGuild(context.Background(), "session-token"); err != nil {
		t.Fatalf("LeaveMyGuild() がエラーを返しました: %v", err)
	}
	if repository.updated == nil {
		t.Fatal("UpdateMembership() が呼ばれる必要があります")
	}
	if repository.updated.LeftAt == nil {
		t.Fatal("left_at が設定されている必要があります")
	}
	if !repository.updated.LeftAt.Equal(now) {
		t.Fatalf("left_at = %v, 期待値 %v", repository.updated.LeftAt, now)
	}
	if !repository.updated.UpdatedAt.Equal(now) {
		t.Fatalf("updated_at = %v, 期待値 %v", repository.updated.UpdatedAt, now)
	}
}

func TestUseCaseLeaveMyGuildRejectsMembershipNotFound(t *testing.T) {
	usecase := NewUseCase(&testRepository{}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	err := usecase.LeaveMyGuild(context.Background(), "session-token")
	if !errors.Is(err, ErrActiveMembershipNotFound) {
		t.Fatalf("LeaveMyGuild() error = %v, 期待値 ErrActiveMembershipNotFound", err)
	}
}

func TestUseCaseContributeCPSpendsCPAndRecordsContribution(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	activeMembership := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    "user_1",
			GuildID:   "guild_go",
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guilddomain.Guild{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	repository := &testRepository{activeMembership: &activeMembership}
	cp := &testCPSpender{}
	usecase := NewUseCaseWithCP(
		repository,
		testCurrentUserRepository{appUser: user.User{ID: "user_1"}, ok: true},
		testIDGenerator{id: "membership_1"},
		testIDGenerator{id: "guild_cp_contribution_1"},
		cp,
	)
	usecase.now = func() time.Time { return now }

	contribution, err := usecase.ContributeCP(context.Background(), "session-token", 40)
	if err != nil {
		t.Fatalf("ContributeCP() がエラーを返しました: %v", err)
	}
	if contribution.ID != "guild_cp_contribution_1" {
		t.Fatalf("contribution id = %q, 期待値 guild_cp_contribution_1", contribution.ID)
	}
	if contribution.GuildID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", contribution.GuildID)
	}
	if contribution.UserID != "user_1" {
		t.Fatalf("user id = %q, 期待値 user_1", contribution.UserID)
	}
	if contribution.PointLedgerID != "cp_ledger_1" {
		t.Fatalf("point ledger id = %q, 期待値 cp_ledger_1", contribution.PointLedgerID)
	}
	if cp.command.Amount != 40 {
		t.Fatalf("CP spend amount = %d, 期待値 40", cp.command.Amount)
	}
	if cp.command.PointType != contributionpointdomain.PointTypeCP {
		t.Fatalf("CP point type = %q, 期待値 CP", cp.command.PointType)
	}
	if cp.command.SourceType != "guild_cp_contribution" {
		t.Fatalf("CP source type = %q, 期待値 guild_cp_contribution", cp.command.SourceType)
	}
	if repository.contribution == nil {
		t.Fatal("CreateCPContribution() が呼ばれる必要があります")
	}
}

func TestUseCaseContributeCPRejectsUserWithoutGuild(t *testing.T) {
	usecase := NewUseCaseWithCP(
		&testRepository{},
		testCurrentUserRepository{appUser: user.User{ID: "user_1"}, ok: true},
		testIDGenerator{id: "membership_1"},
		testIDGenerator{id: "guild_cp_contribution_1"},
		&testCPSpender{},
	)

	_, err := usecase.ContributeCP(context.Background(), "session-token", 40)
	if !errors.Is(err, ErrActiveMembershipNotFound) {
		t.Fatalf("ContributeCP() error = %v, 期待値 ErrActiveMembershipNotFound", err)
	}
}

func TestUseCaseContributeCPRejectsInsufficientBalance(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	activeMembership := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    "user_1",
			GuildID:   "guild_go",
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guilddomain.Guild{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	repository := &testRepository{activeMembership: &activeMembership}
	usecase := NewUseCaseWithCP(
		repository,
		testCurrentUserRepository{appUser: user.User{ID: "user_1"}, ok: true},
		testIDGenerator{id: "membership_1"},
		testIDGenerator{id: "guild_cp_contribution_1"},
		&testCPSpender{err: contributionpointapp.ErrInsufficientBalance},
	)

	_, err := usecase.ContributeCP(context.Background(), "session-token", 40)
	if !errors.Is(err, contributionpointapp.ErrInsufficientBalance) {
		t.Fatalf("ContributeCP() error = %v, 期待値 ErrInsufficientBalance", err)
	}
	if repository.contribution != nil {
		t.Fatal("CP不足時は CreateCPContribution() を呼ばない必要があります")
	}
}

func TestUseCaseListGuildCPContributionsRejectsUnknownGuild(t *testing.T) {
	usecase := NewUseCaseWithCP(
		&testRepository{},
		testCurrentUserRepository{appUser: user.User{ID: "user_1"}, ok: true},
		testIDGenerator{id: "membership_1"},
		testIDGenerator{id: "guild_cp_contribution_1"},
		&testCPSpender{},
	)

	_, err := usecase.ListGuildCPContributions(context.Background(), "guild_missing")
	if !errors.Is(err, ErrGuildNotFound) {
		t.Fatalf("ListGuildCPContributions() error = %v, 期待値 ErrGuildNotFound", err)
	}
}

func TestNewUseCasePanicsWithoutRepository(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("NewUseCase() panic = nil, 期待値 panic")
		}
		if message := fmt.Sprint(recovered); message != "guild repository is required" {
			t.Fatalf("NewUseCase() panic = %q, 期待値 guild repository is required", message)
		}
	}()

	_ = NewUseCase(nil, testCurrentUserRepository{ok: true}, testIDGenerator{id: "membership_1"})
}
