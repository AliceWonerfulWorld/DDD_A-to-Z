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
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type testRepository struct {
	guilds           []guilddomain.Guild
	activeMembership *guilddomain.MembershipWithGuild
	created          *guilddomain.Membership
	updated          *guilddomain.Membership
	pet              *petdomain.Pet
	createdPet       *petdomain.Pet
	createPetErr     error
	contribution     *guilddomain.CPContribution
	contributions    []guilddomain.CPContribution
	activityLogs     []guilddomain.ActivityLog
	membersByGuild   map[guilddomain.ID][]guilddomain.MemberContribution
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

func (r testRepository) ListActiveMembersByGuild(ctx context.Context, guildID guilddomain.ID) ([]guilddomain.MemberContribution, error) {
	return r.membersByGuild[guildID], nil
}

func (r testRepository) ListActivityLogsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.ActivityLog, error) {
	if limit > 0 && len(r.activityLogs) > limit {
		return r.activityLogs[:limit], nil
	}

	return r.activityLogs, nil
}

func (r *testRepository) CreateMembership(ctx context.Context, membership guilddomain.Membership) error {
	r.created = &membership
	return nil
}

func (r *testRepository) UpdateMembership(ctx context.Context, membership guilddomain.Membership) error {
	r.updated = &membership
	return nil
}

func (r testRepository) FindPetByUserAndGuild(ctx context.Context, userID user.ID, guildID guilddomain.ID) (petdomain.Pet, bool, error) {
	if r.pet == nil {
		return petdomain.Pet{}, false, nil
	}
	if r.pet.UserID != userID || r.pet.GuildID != guildID {
		return petdomain.Pet{}, false, nil
	}

	return *r.pet, true, nil
}

func (r *testRepository) CreatePet(ctx context.Context, pet petdomain.Pet) error {
	if r.createPetErr != nil {
		return r.createPetErr
	}
	r.createdPet = &pet
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

	result, err := usecase.JoinGuild(context.Background(), "session-token", "guild_go")
	if err != nil {
		t.Fatalf("JoinGuild() がエラーを返しました: %v", err)
	}
	if result.Membership.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", result.Membership.Guild.ID)
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
	if repository.createdPet == nil {
		t.Fatal("CreatePet() が呼ばれる必要があります")
	}
	if repository.createdPet.UserID != "user_1" || repository.createdPet.GuildID != "guild_go" {
		t.Fatalf("created pet owner = %q/%q, 期待値 user_1/guild_go", repository.createdPet.UserID, repository.createdPet.GuildID)
	}
	if result.GrantedPet == nil {
		t.Fatal("GrantedPet が設定されている必要があります")
	}
	if result.PetAlreadyOwned {
		t.Fatal("PetAlreadyOwned = true, 期待値 false")
	}
}

func TestUseCaseJoinGuildDoesNotGrantDuplicatePet(t *testing.T) {
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
	existingPet := petdomain.Pet{
		ID:        "pet_existing",
		UserID:    "user_1",
		GuildID:   "guild_go",
		Attribute: petdomain.AttributeGo,
		Stats:     petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7},
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(-time.Hour),
	}
	repository := &testRepository{
		guilds: []guilddomain.Guild{targetGuild},
		pet:    &existingPet,
	}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_1"})
	usecase.now = func() time.Time { return now }

	result, err := usecase.JoinGuild(context.Background(), "session-token", "guild_go")
	if err != nil {
		t.Fatalf("JoinGuild() がエラーを返しました: %v", err)
	}
	if repository.createdPet != nil {
		t.Fatal("CreatePet() は呼ばれない必要があります")
	}
	if result.GrantedPet != nil {
		t.Fatalf("GrantedPet = %#v, 期待値 nil", result.GrantedPet)
	}
	if !result.PetAlreadyOwned {
		t.Fatal("PetAlreadyOwned = false, 期待値 true")
	}
}

func TestUseCaseJoinGuildTreatsDuplicatePetInsertAsAlreadyOwned(t *testing.T) {
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
	repository := &testRepository{
		guilds:       []guilddomain.Guild{targetGuild},
		createPetErr: ErrPetAlreadyOwned,
	}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_1"})
	usecase.now = func() time.Time { return now }

	result, err := usecase.JoinGuild(context.Background(), "session-token", "guild_go")
	if err != nil {
		t.Fatalf("JoinGuild() がエラーを返しました: %v", err)
	}
	if result.GrantedPet != nil {
		t.Fatalf("GrantedPet = %#v, 期待値 nil", result.GrantedPet)
	}
	if !result.PetAlreadyOwned {
		t.Fatal("PetAlreadyOwned = false, 期待値 true")
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

func TestUseCaseJoinGuildGrantsMissingPetForAlreadyJoinedUser(t *testing.T) {
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
	repository := &testRepository{activeMembership: &existing}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "pet_repair_1"})
	usecase.now = func() time.Time { return now }

	result, err := usecase.JoinGuild(context.Background(), "session-token", "guild_go")
	if !errors.Is(err, ErrAlreadyJoined) {
		t.Fatalf("JoinGuild() error = %v, 期待値 ErrAlreadyJoined", err)
	}
	if repository.createdPet == nil {
		t.Fatal("CreatePet() が呼ばれる必要があります")
	}
	if repository.createdPet.UserID != "user_1" || repository.createdPet.GuildID != "guild_go" {
		t.Fatalf("created pet owner = %q/%q, 期待値 user_1/guild_go", repository.createdPet.UserID, repository.createdPet.GuildID)
	}
	if result.GrantedPet == nil {
		t.Fatal("GrantedPet が設定されている必要があります")
	}
	if result.PetAlreadyOwned {
		t.Fatal("PetAlreadyOwned = true, 期待値 false")
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

func TestUseCaseGetMyGuildDetailsListsMembersForActiveGuild(t *testing.T) {
	now := time.Date(2026, 5, 19, 9, 0, 0, 0, time.UTC)
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
	usecase := NewUseCase(&testRepository{
		activeMembership: &activeMembership,
		membersByGuild: map[guilddomain.ID][]guilddomain.MemberContribution{
			"guild_go": {{
				UserID:        "user_1",
				Name:          "Alice",
				TotalEarnedCP: 120,
				JoinedAt:      now,
			}},
			"guild_python": {{
				UserID:        "user_2",
				Name:          "Bob",
				TotalEarnedCP: 80,
				JoinedAt:      now,
			}},
		},
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	details, ok, err := usecase.GetMyGuildDetails(context.Background(), "session-token")
	if err != nil {
		t.Fatalf("GetMyGuildDetails() がエラーを返しました: %v", err)
	}
	if !ok {
		t.Fatal("GetMyGuildDetails() ok = false, 期待値 true")
	}
	if details.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", details.Guild.ID)
	}
	if len(details.Members) != 1 {
		t.Fatalf("members length = %d, 期待値 1", len(details.Members))
	}
	if details.Members[0].UserID != "user_1" {
		t.Fatalf("members[0].UserID = %q, 期待値 user_1", details.Members[0].UserID)
	}
}

func TestUseCaseGetGuildDashboardRequiresOwnGuild(t *testing.T) {
	now := time.Date(2026, 5, 19, 9, 0, 0, 0, time.UTC)
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
			ID:                 "guild_go",
			Slug:               "go",
			Name:               "Go",
			Description:        "シンプルさと並列処理で前に進むギルド。",
			Icon:               "GO",
			Color:              "#00acd7",
			SortOrder:          1,
			MemberCount:        2,
			TotalContributedCP: 160,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
	usecase := NewUseCase(&testRepository{
		activeMembership: &activeMembership,
		membersByGuild: map[guilddomain.ID][]guilddomain.MemberContribution{
			"guild_go": {{
				UserID:        "user_1",
				Name:          "Alice",
				TotalEarnedCP: 120,
				JoinedAt:      now,
			}},
		},
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	dashboard, err := usecase.GetGuildDashboard(context.Background(), "session-token", "guild_go")
	if err != nil {
		t.Fatalf("GetGuildDashboard() がエラーを返しました: %v", err)
	}
	if dashboard.Guild.TotalContributedCP != 160 {
		t.Fatalf("total contributed cp = %d, 期待値 160", dashboard.Guild.TotalContributedCP)
	}
	if len(dashboard.Members) != 1 {
		t.Fatalf("members length = %d, 期待値 1", len(dashboard.Members))
	}
}

func TestUseCaseGetGuildDashboardRejectsOtherGuild(t *testing.T) {
	now := time.Date(2026, 5, 19, 9, 0, 0, 0, time.UTC)
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
	usecase := NewUseCase(&testRepository{
		activeMembership: &activeMembership,
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	_, err := usecase.GetGuildDashboard(context.Background(), "session-token", "guild_python")
	if !errors.Is(err, ErrGuildAccessDenied) {
		t.Fatalf("GetGuildDashboard() error = %v, 期待値 ErrGuildAccessDenied", err)
	}
}

func TestUseCaseGetGuildDashboardRejectsUserWithoutGuild(t *testing.T) {
	usecase := NewUseCase(&testRepository{}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	_, err := usecase.GetGuildDashboard(context.Background(), "session-token", "guild_go")
	if !errors.Is(err, ErrActiveMembershipNotFound) {
		t.Fatalf("GetGuildDashboard() error = %v, 期待値 ErrActiveMembershipNotFound", err)
	}
}

func TestUseCaseListGuildActivityLogs(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
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
			Description: "Go guild",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	usecase := NewUseCase(&testRepository{
		activeMembership: &activeMembership,
		activityLogs: []guilddomain.ActivityLog{{
			ID:         "user_1:commit:repo:sha",
			UserID:     "user_1",
			Player:     "Alice",
			Type:       "commit",
			Repo:       "jyogi-web/DDD_A-to-Z",
			Message:    "Add activity log",
			Language:   "Go",
			CP:         1,
			OccurredAt: now,
		}},
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	logs, err := usecase.ListGuildActivityLogs(context.Background(), "session-token", "guild_go", 20)
	if err != nil {
		t.Fatalf("ListGuildActivityLogs() がエラーを返しました: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("logs length = %d, 期待値 1", len(logs))
	}
	if logs[0].Message != "Add activity log" {
		t.Fatalf("message = %q, 期待値 Add activity log", logs[0].Message)
	}
}

func TestUseCaseListGuildActivityLogsRejectsOtherGuild(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
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
			Description: "Go guild",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	usecase := NewUseCase(&testRepository{
		activeMembership: &activeMembership,
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testIDGenerator{id: "membership_unused"})

	_, err := usecase.ListGuildActivityLogs(context.Background(), "session-token", "guild_python", 20)
	if !errors.Is(err, ErrGuildAccessDenied) {
		t.Fatalf("ListGuildActivityLogs() error = %v, 期待値 ErrGuildAccessDenied", err)
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
