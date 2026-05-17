package guild

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type testRepository struct {
	guilds           []guilddomain.Guild
	activeMembership *guilddomain.MembershipWithGuild
	created          *guilddomain.Membership
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
