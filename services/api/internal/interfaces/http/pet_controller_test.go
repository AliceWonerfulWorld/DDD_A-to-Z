package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
)

type stubPetCurrentUser struct {
	u     user.User
	found bool
	err   error
}

func (s *stubPetCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return s.u, s.found, s.err
}

type stubPetCPBalanceReader struct {
	balance int64
	err     error
}

func (s *stubPetCPBalanceReader) GetBalance(_ context.Context, _ user.ID) (int64, error) {
	return s.balance, s.err
}

type stubPetReader struct {
	pets []petapp.PetWithGuild
	err  error
}

func (s *stubPetReader) ListPetsByUser(_ context.Context, _ user.ID) ([]petapp.PetWithGuild, error) {
	return s.pets, s.err
}

type stubPetCurrentGuildReader struct {
	membership guilddomain.MembershipWithGuild
	found      bool
	err        error
}

func (s *stubPetCurrentGuildReader) FindActiveMembershipByUserID(_ context.Context, _ user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	return s.membership, s.found, s.err
}

type stubHTTPPetBattleReader struct {
	opponents     []petapp.PetWithGuild
	attacker      petapp.PetWithGuild
	attackerFound bool
	defender      petapp.PetWithGuild
	defenderFound bool
}

func (s *stubHTTPPetBattleReader) ListOpponentPets(_ context.Context, userID user.ID) ([]petapp.PetWithGuild, error) {
	opponents := make([]petapp.PetWithGuild, 0, len(s.opponents))
	for _, opponent := range s.opponents {
		if opponent.Pet.UserID != userID {
			opponents = append(opponents, opponent)
		}
	}
	return opponents, nil
}

func (s *stubHTTPPetBattleReader) FindPetByIDForUser(_ context.Context, petID petdomain.ID, userID user.ID) (petapp.PetWithGuild, bool, error) {
	if !s.attackerFound || s.attacker.Pet.ID != petID || s.attacker.Pet.UserID != userID {
		return petapp.PetWithGuild{}, false, nil
	}
	return s.attacker, true, nil
}

func (s *stubHTTPPetBattleReader) FindOpponentPetByID(_ context.Context, petID petdomain.ID, userID user.ID) (petapp.PetWithGuild, bool, error) {
	if !s.defenderFound || s.defender.Pet.ID != petID || s.defender.Pet.UserID == userID {
		return petapp.PetWithGuild{}, false, nil
	}
	return s.defender, true, nil
}

func newPetMux(currentUser *stubPetCurrentUser, cp *stubPetCPBalanceReader, pets *stubPetReader, guild *stubPetCurrentGuildReader) *http.ServeMux {
	uc := petapp.NewUseCase(currentUser, cp, pets, guild)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	controller := httpapi.NewPetController(uc, logger)
	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)
	return mux
}

func newPetBattleMux(currentUser *stubPetCurrentUser, battlePets *stubHTTPPetBattleReader) *http.ServeMux {
	uc := petapp.NewUseCaseWithTrainingAndBattle(
		currentUser,
		&stubPetCPBalanceReader{},
		&stubPetReader{},
		&stubPetCurrentGuildReader{},
		nil,
		nil,
		nil,
		nil,
		battlePets,
	)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	controller := httpapi.NewPetController(uc, logger)
	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)
	return mux
}

func TestPetControllerNoCookie(t *testing.T) {
	mux := newPetMux(&stubPetCurrentUser{}, &stubPetCPBalanceReader{}, &stubPetReader{}, &stubPetCurrentGuildReader{})

	req := httptest.NewRequest("GET", "/pets/me", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestPetControllerSuccess(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	goGuild := mustHTTPPetGuild(t, "guild_go", "go", "Go", now)
	goPet := mustHTTPPet(t, "pet_go", testUser.ID, "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 7, Agility: 7}, now)
	mux := newPetMux(
		&stubPetCurrentUser{u: testUser, found: true},
		&stubPetCPBalanceReader{balance: 120},
		&stubPetReader{pets: []petapp.PetWithGuild{{Pet: goPet, Guild: goGuild}}},
		&stubPetCurrentGuildReader{membership: guilddomain.MembershipWithGuild{Guild: goGuild}, found: true},
	)

	req := httptest.NewRequest("GET", "/pets/me", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["cpBalance"].(float64) != 120 {
		t.Fatalf("cpBalance = %v, 期待値 120", body["cpBalance"])
	}
	currentGuildPet, ok := body["currentGuildPet"].(map[string]any)
	if !ok {
		t.Fatal("expected currentGuildPet object")
	}
	if currentGuildPet["guildId"] != "guild_go" {
		t.Fatalf("guildId = %v, 期待値 guild_go", currentGuildPet["guildId"])
	}
	if currentGuildPet["maxHp"].(float64) != 35 {
		t.Fatalf("maxHp = %v, 期待値 35", currentGuildPet["maxHp"])
	}
	if _, ok := currentGuildPet["acquiredAt"].(string); !ok {
		t.Fatal("expected acquiredAt string")
	}
	pets, ok := body["pets"].([]any)
	if !ok {
		t.Fatal("expected pets array")
	}
	if len(pets) != 1 {
		t.Fatalf("pets length = %d, 期待値 1", len(pets))
	}
}

func TestPetControllerNoGuildNoPets(t *testing.T) {
	testUser := user.User{ID: "user_1"}
	mux := newPetMux(
		&stubPetCurrentUser{u: testUser, found: true},
		&stubPetCPBalanceReader{balance: 120},
		&stubPetReader{pets: []petapp.PetWithGuild{}},
		&stubPetCurrentGuildReader{found: false},
	)

	req := httptest.NewRequest("GET", "/pets/me", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["currentGuildPet"] != nil {
		t.Fatalf("currentGuildPet = %v, 期待値 nil", body["currentGuildPet"])
	}
	pets, ok := body["pets"].([]any)
	if !ok {
		t.Fatal("expected pets array")
	}
	if len(pets) != 0 {
		t.Fatalf("pets length = %d, 期待値 0", len(pets))
	}
}

func TestPetControllerInternalError(t *testing.T) {
	testUser := user.User{ID: "user_1"}
	mux := newPetMux(
		&stubPetCurrentUser{u: testUser, found: true},
		&stubPetCPBalanceReader{err: errors.New("db error")},
		&stubPetReader{},
		&stubPetCurrentGuildReader{},
	)

	req := httptest.NewRequest("GET", "/pets/me", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestPetControllerListBattleOpponents(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	otherUser := user.User{ID: "user_2"}
	rustGuild := mustHTTPPetGuild(t, "guild_rust", "rust", "Rust", now)
	rustPet := mustHTTPPet(t, "pet_rust", otherUser.ID, "guild_rust", petdomain.AttributeRust, petdomain.Stats{Vitality: 8, Strength: 8, Agility: 4}, now)
	mux := newPetBattleMux(
		&stubPetCurrentUser{u: testUser, found: true},
		&stubHTTPPetBattleReader{opponents: []petapp.PetWithGuild{{Pet: rustPet, Guild: rustGuild}}},
	)

	req := httptest.NewRequest("GET", "/pets/battle/opponents", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	opponents, ok := body["opponents"].([]any)
	if !ok || len(opponents) != 1 {
		t.Fatalf("opponents = %v, 期待値 length 1", body["opponents"])
	}
	opponent := opponents[0].(map[string]any)
	if opponent["petId"] != "pet_rust" || opponent["ownerUserId"] != "user_2" {
		t.Fatalf("opponent = %+v", opponent)
	}
	if opponent["maxHp"].(float64) != 45 || opponent["power"].(float64) != 7 {
		t.Fatalf("opponent stats = %+v", opponent)
	}
}

func TestPetControllerBattlePet(t *testing.T) {
	now := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	testUser := user.User{ID: "user_1"}
	otherUser := user.User{ID: "user_2"}
	goGuild := mustHTTPPetGuild(t, "guild_go", "go", "Go", now)
	rustGuild := mustHTTPPetGuild(t, "guild_rust", "rust", "Rust", now)
	attacker := petapp.PetWithGuild{
		Pet:   mustHTTPPet(t, "pet_go", testUser.ID, "guild_go", petdomain.AttributeGo, petdomain.Stats{Vitality: 6, Strength: 20, Agility: 7}, now),
		Guild: goGuild,
	}
	defender := petapp.PetWithGuild{
		Pet:   mustHTTPPet(t, "pet_rust", otherUser.ID, "guild_rust", petdomain.AttributeRust, petdomain.Stats{Vitality: 4, Strength: 8, Agility: 4}, now),
		Guild: rustGuild,
	}
	mux := newPetBattleMux(
		&stubPetCurrentUser{u: testUser, found: true},
		&stubHTTPPetBattleReader{attacker: attacker, attackerFound: true, defender: defender, defenderFound: true},
	)

	req := httptest.NewRequest("POST", "/pets/pet_go/battle", bytes.NewBufferString(`{"opponentPetId":"pet_rust"}`))
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["result"] != "win" || body["winnerPetId"] != "pet_go" {
		t.Fatalf("body = %+v", body)
	}
	turns := body["turns"].([]any)
	if len(turns) == 0 {
		t.Fatal("turns length = 0")
	}
	firstTurn := turns[0].(map[string]any)
	if firstTurn["actorPetId"] != "pet_go" || firstTurn["message"] != "Gopher attacks Ferris for 18 damage." {
		t.Fatalf("first turn = %+v", firstTurn)
	}
	attackerBody := body["attacker"].(map[string]any)
	if attackerBody["petId"] != "pet_go" || attackerBody["name"] != "Gopher" {
		t.Fatalf("attacker = %+v", attackerBody)
	}
}

func TestPetControllerBattleRejectsOwnPetAsOpponent(t *testing.T) {
	mux := newPetBattleMux(
		&stubPetCurrentUser{u: user.User{ID: "user_1"}, found: true},
		&stubHTTPPetBattleReader{},
	)

	req := httptest.NewRequest("POST", "/pets/pet_go/battle", bytes.NewBufferString(`{"opponentPetId":"pet_go"}`))
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func mustHTTPPet(t *testing.T, id petdomain.ID, userID user.ID, guildID guilddomain.ID, attribute petdomain.Attribute, stats petdomain.Stats, now time.Time) petdomain.Pet {
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

func mustHTTPPetGuild(t *testing.T, id guilddomain.ID, slug, name string, now time.Time) guilddomain.Guild {
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
