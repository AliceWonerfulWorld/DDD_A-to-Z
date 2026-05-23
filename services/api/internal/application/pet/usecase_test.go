package pet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
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

type stubPetTrainingRepository struct {
	petWithGuild petapp.PetWithGuild
	found        bool
	updated      *petdomain.Pet
}

func (s *stubPetTrainingRepository) FindPetByIDForUser(_ context.Context, petID petdomain.ID, userID user.ID) (petapp.PetWithGuild, bool, error) {
	if !s.found || s.petWithGuild.Pet.ID != petID || s.petWithGuild.Pet.UserID != userID {
		return petapp.PetWithGuild{}, false, nil
	}
	return s.petWithGuild, true, nil
}

func (s *stubPetTrainingRepository) UpdatePet(_ context.Context, pet petdomain.Pet) error {
	s.updated = &pet
	return nil
}

type stubCPSpender struct {
	balance int64
	spent   *contributionpointapp.SpendCommand
}

func (s *stubCPSpender) Spend(_ context.Context, command contributionpointapp.SpendCommand) (contributionpointdomain.LedgerEntry, error) {
	if s.balance < command.Amount {
		return contributionpointdomain.LedgerEntry{}, contributionpointapp.ErrInsufficientBalance
	}
	s.spent = &command
	s.balance -= command.Amount
	return contributionpointdomain.LedgerEntry{
		ID:           "point_ledger_1",
		UserID:       command.UserID,
		PointType:    command.PointType,
		Amount:       -command.Amount,
		Type:         contributionpointdomain.EntryTypeSpend,
		Reason:       command.Reason,
		SourceType:   command.SourceType,
		SourceID:     command.SourceID,
		BalanceAfter: s.balance,
		CreatedAt:    time.Date(2026, 5, 23, 1, 0, 0, 0, time.UTC),
	}, nil
}

type stubIDGenerator struct {
	id string
}

func (g stubIDGenerator) NewID() (string, error) {
	return g.id, nil
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

func TestTrainPetSuccess(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	goGuild := mustGuild(t, "guild_go", "go", "Go", now)

	tests := []struct {
		stat            string
		spentCP         int64
		expectedBalance int64
		reason          string
		assertResult    func(t *testing.T, result petapp.TrainPetResult)
		assertUpdated   func(t *testing.T, pet petdomain.Pet)
	}{
		{
			stat:            "power",
			spentCP:         10,
			expectedBalance: 110,
			reason:          "pet_training_power",
			assertResult: func(t *testing.T, result petapp.TrainPetResult) {
				t.Helper()
				if result.Pet.Power != 7 {
					t.Fatalf("Power = %d, 期待値 7", result.Pet.Power)
				}
			},
			assertUpdated: func(t *testing.T, pet petdomain.Pet) {
				t.Helper()
				if pet.Stats.Strength != 8 {
					t.Fatalf("updated Strength = %d, 期待値 8", pet.Stats.Strength)
				}
			},
		},
		{
			stat:            "hp",
			spentCP:         20,
			expectedBalance: 100,
			reason:          "pet_training_hp",
			assertResult: func(t *testing.T, result petapp.TrainPetResult) {
				t.Helper()
				if result.Pet.MaxHP != 40 {
					t.Fatalf("MaxHP = %d, 期待値 40", result.Pet.MaxHP)
				}
			},
			assertUpdated: func(t *testing.T, pet petdomain.Pet) {
				t.Helper()
				if pet.Stats.Vitality != 7 {
					t.Fatalf("updated Vitality = %d, 期待値 7", pet.Stats.Vitality)
				}
			},
		},
		{
			stat:            "guard",
			spentCP:         10,
			expectedBalance: 110,
			reason:          "pet_training_guard",
			assertResult: func(t *testing.T, result petapp.TrainPetResult) {
				t.Helper()
				if result.Pet.Guard != 6 {
					t.Fatalf("Guard = %d, 期待値 6", result.Pet.Guard)
				}
			},
			assertUpdated: func(t *testing.T, pet petdomain.Pet) {
				t.Helper()
				if pet.Stats.Vitality != 7 {
					t.Fatalf("updated Vitality = %d, 期待値 7", pet.Stats.Vitality)
				}
			},
		},
		{
			stat:            "speed",
			spentCP:         10,
			expectedBalance: 110,
			reason:          "pet_training_speed",
			assertResult: func(t *testing.T, result petapp.TrainPetResult) {
				t.Helper()
				if result.Pet.Speed != 7 {
					t.Fatalf("Speed = %d, 期待値 7", result.Pet.Speed)
				}
			},
			assertUpdated: func(t *testing.T, pet petdomain.Pet) {
				t.Helper()
				if pet.Stats.Agility != 8 {
					t.Fatalf("updated Agility = %d, 期待値 8", pet.Stats.Agility)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.stat, func(t *testing.T) {
			goPet := mustPet(t, "pet_go", testUser.ID, "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7}, now)
			pets := &stubPetTrainingRepository{petWithGuild: petapp.PetWithGuild{Pet: goPet, Guild: goGuild}, found: true}
			cp := &stubCPSpender{balance: 120}
			uc := petapp.NewUseCaseWithTraining(
				&stubCurrentUser{user: testUser, found: true},
				&stubCPBalanceReader{},
				&stubPetReader{},
				&stubCurrentGuildReader{},
				pets,
				cp,
				stubIDGenerator{id: "pet_training_1"},
				nil,
			)

			result, err := uc.TrainPet(context.Background(), petapp.TrainPetCommand{
				SessionToken: "valid-token",
				PetID:        "pet_go",
				Stat:         tt.stat,
			})
			if err != nil {
				t.Fatalf("TrainPet() error = %v", err)
			}
			if result.SpentCP != tt.spentCP || result.CPBalance != tt.expectedBalance {
				t.Fatalf("result = %+v, 期待値 spent %d balance %d", result, tt.spentCP, tt.expectedBalance)
			}
			tt.assertResult(t, result)
			if pets.updated == nil {
				t.Fatal("updated pet = nil")
			}
			tt.assertUpdated(t, *pets.updated)
			if cp.spent == nil {
				t.Fatal("spent command = nil")
			}
			if cp.spent.Amount != tt.spentCP || cp.spent.Reason != tt.reason || cp.spent.SourceID != "pet_training_1" {
				t.Fatalf("spent command = %+v", cp.spent)
			}
		})
	}
}

func TestTrainPetInvalidStat(t *testing.T) {
	uc := petapp.NewUseCaseWithTraining(&stubCurrentUser{found: true}, &stubCPBalanceReader{}, &stubPetReader{}, &stubCurrentGuildReader{}, &stubPetTrainingRepository{}, &stubCPSpender{}, stubIDGenerator{id: "pet_training_1"}, nil)

	_, err := uc.TrainPet(context.Background(), petapp.TrainPetCommand{SessionToken: "valid-token", PetID: "pet_go", Stat: "luck"})
	if !errors.Is(err, petapp.ErrInvalidTrainStat) {
		t.Fatalf("error = %v, 期待値 ErrInvalidTrainStat", err)
	}
}

func TestTrainPetInsufficientCPDoesNotUpdatePet(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	goGuild := mustGuild(t, "guild_go", "go", "Go", now)
	goPet := mustPet(t, "pet_go", testUser.ID, "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7}, now)
	pets := &stubPetTrainingRepository{petWithGuild: petapp.PetWithGuild{Pet: goPet, Guild: goGuild}, found: true}
	uc := petapp.NewUseCaseWithTraining(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPBalanceReader{},
		&stubPetReader{},
		&stubCurrentGuildReader{},
		pets,
		&stubCPSpender{balance: 5},
		stubIDGenerator{id: "pet_training_1"},
		nil,
	)

	_, err := uc.TrainPet(context.Background(), petapp.TrainPetCommand{SessionToken: "valid-token", PetID: "pet_go", Stat: "power"})
	if !errors.Is(err, petapp.ErrInsufficientCP) {
		t.Fatalf("error = %v, 期待値 ErrInsufficientCP", err)
	}
	if pets.updated != nil {
		t.Fatalf("updated pet = %+v, 期待値 nil", pets.updated)
	}
}

func TestTrainPetCannotTrainOtherUsersPet(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	otherPet := mustPet(t, "pet_go", "user_2", "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7}, now)
	pets := &stubPetTrainingRepository{petWithGuild: petapp.PetWithGuild{Pet: otherPet, Guild: mustGuild(t, "guild_go", "go", "Go", now)}, found: true}
	cp := &stubCPSpender{balance: 120}
	uc := petapp.NewUseCaseWithTraining(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPBalanceReader{},
		&stubPetReader{},
		&stubCurrentGuildReader{},
		pets,
		cp,
		stubIDGenerator{id: "pet_training_1"},
		nil,
	)

	_, err := uc.TrainPet(context.Background(), petapp.TrainPetCommand{SessionToken: "valid-token", PetID: "pet_go", Stat: "power"})
	if !errors.Is(err, petapp.ErrPetNotFound) {
		t.Fatalf("error = %v, 期待値 ErrPetNotFound", err)
	}
	if cp.spent != nil || cp.balance != 120 {
		t.Fatalf("cp spender = %+v balance %d, 期待値 未消費 balance 120", cp.spent, cp.balance)
	}
	if pets.updated != nil {
		t.Fatalf("updated pet = %+v, 期待値 nil", pets.updated)
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
