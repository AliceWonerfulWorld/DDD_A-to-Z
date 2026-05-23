package pet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type stubCurrentUser struct {
	user  user.User
	found bool
	err   error
}

func (s *stubCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return s.user, s.found, s.err
}

type stubCPBalanceReader struct {
	balance int64
	err     error
}

func (s *stubCPBalanceReader) GetBalance(_ context.Context, _ user.ID) (int64, error) {
	return s.balance, s.err
}

type stubPetReader struct {
	pets []petapp.PetWithGuild
	err  error
}

func (s *stubPetReader) ListPetsByUser(_ context.Context, _ user.ID) ([]petapp.PetWithGuild, error) {
	return s.pets, s.err
}

type stubCurrentGuildReader struct {
	membership guilddomain.MembershipWithGuild
	found      bool
	err        error
}

func (s *stubCurrentGuildReader) FindActiveMembershipByUserID(_ context.Context, _ user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	return s.membership, s.found, s.err
}

func TestGetMyPetsEmptyToken(t *testing.T) {
	uc := petapp.NewUseCase(&stubCurrentUser{}, &stubCPBalanceReader{}, &stubPetReader{}, &stubCurrentGuildReader{})

	_, err := uc.GetMyPets(context.Background(), "")
	if !errors.Is(err, petapp.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetMyPetsSessionNotFound(t *testing.T) {
	uc := petapp.NewUseCase(&stubCurrentUser{found: false}, &stubCPBalanceReader{}, &stubPetReader{}, &stubCurrentGuildReader{})

	_, err := uc.GetMyPets(context.Background(), "invalid-token")
	if !errors.Is(err, petapp.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetMyPetsSuccess(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	goGuild := mustGuild(t, "guild_go", "go", "Go", now)
	rustGuild := mustGuild(t, "guild_rust", "rust", "Rust", now)
	goPet := mustPet(t, "pet_go", testUser.ID, "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7}, now)
	rustPet := mustPet(t, "pet_rust", testUser.ID, "guild_rust", petdomain.AttributeRust, petdomain.Stats{Vitality: 8, Strength: 8, Agility: 4}, now.Add(-time.Hour))

	uc := petapp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPBalanceReader{balance: 120},
		&stubPetReader{pets: []petapp.PetWithGuild{
			{Pet: goPet, Guild: goGuild},
			{Pet: rustPet, Guild: rustGuild},
		}},
		&stubCurrentGuildReader{membership: guilddomain.MembershipWithGuild{Guild: goGuild}, found: true},
	)

	result, err := uc.GetMyPets(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.CPBalance != 120 {
		t.Fatalf("CPBalance = %d, 期待値 120", result.CPBalance)
	}
	if len(result.Pets) != 2 {
		t.Fatalf("pets length = %d, 期待値 2", len(result.Pets))
	}
	if result.CurrentGuildPet == nil {
		t.Fatal("CurrentGuildPet = nil, 期待値 Go pet")
	}
	if result.CurrentGuildPet.ID != "pet_go" {
		t.Fatalf("CurrentGuildPet.ID = %q, 期待値 pet_go", result.CurrentGuildPet.ID)
	}
	if result.CurrentGuildPet.Name != "Gopher" || result.CurrentGuildPet.Species != "gopher" {
		t.Fatalf("unexpected current guild pet display fields: %+v", result.CurrentGuildPet)
	}
	if result.CurrentGuildPet.MaxHP != 35 || result.CurrentGuildPet.Power != 6 || result.CurrentGuildPet.Guard != 5 || result.CurrentGuildPet.Speed != 6 {
		t.Fatalf("unexpected current guild pet stats: %+v", result.CurrentGuildPet)
	}
}

func TestGetMyPetsNoGuildNoPets(t *testing.T) {
	testUser := user.User{ID: "user_1"}
	uc := petapp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPBalanceReader{balance: 120},
		&stubPetReader{pets: []petapp.PetWithGuild{}},
		&stubCurrentGuildReader{found: false},
	)

	result, err := uc.GetMyPets(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.CurrentGuildPet != nil {
		t.Fatalf("CurrentGuildPet = %+v, 期待値 nil", result.CurrentGuildPet)
	}
	if len(result.Pets) != 0 {
		t.Fatalf("pets length = %d, 期待値 0", len(result.Pets))
	}
}

func TestGetMyPetsPetReaderError(t *testing.T) {
	testUser := user.User{ID: "user_1"}
	petErr := errors.New("pet query failed")
	uc := petapp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPBalanceReader{balance: 120},
		&stubPetReader{err: petErr},
		&stubCurrentGuildReader{},
	)

	_, err := uc.GetMyPets(context.Background(), "valid-token")
	if !errors.Is(err, petErr) {
		t.Fatalf("expected petErr, got %v", err)
	}
}

func mustPet(t *testing.T, id petdomain.ID, userID user.ID, guildID guilddomain.ID, attribute petdomain.Attribute, stats petdomain.Stats, now time.Time) petdomain.Pet {
	t.Helper()
	foundPet, err := petdomain.NewPet(petdomain.Pet{
		ID:        id,
		UserID:    userID,
		GuildID:   guildID,
		Attribute: attribute,
		Stats:     stats,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to build pet: %v", err)
	}
	return foundPet
}

func mustGuild(t *testing.T, id guilddomain.ID, slug, name string, now time.Time) guilddomain.Guild {
	t.Helper()
	foundGuild, err := guilddomain.NewGuild(guilddomain.Guild{
		ID:          id,
		Slug:        slug,
		Name:        name,
		Description: name + " guild",
		Icon:        name,
		Color:       "#123456",
		SortOrder:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("failed to build guild: %v", err)
	}
	return foundGuild
}
