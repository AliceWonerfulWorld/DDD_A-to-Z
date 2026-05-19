package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type guildTestRepository struct {
	guilds           []guilddomain.Guild
	activeMembership *guilddomain.MembershipWithGuild
	updated          *guilddomain.Membership
	contributions    []guilddomain.CPContribution
	members          []guilddomain.MemberContribution
	expectedGuildID  guilddomain.ID
}

func (r guildTestRepository) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	return r.guilds, nil
}

func (r guildTestRepository) FindGuildByID(ctx context.Context, guildID guilddomain.ID) (guilddomain.Guild, bool, error) {
	for _, guild := range r.guilds {
		if guild.ID == guildID {
			return guild, true, nil
		}
	}

	return guilddomain.Guild{}, false, nil
}

func (r guildTestRepository) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	if r.activeMembership == nil {
		return guilddomain.MembershipWithGuild{}, false, nil
	}

	return *r.activeMembership, true, nil
}

func (r guildTestRepository) ListActiveMembersByGuild(ctx context.Context, guildID guilddomain.ID) ([]guilddomain.MemberContribution, error) {
	if r.expectedGuildID != "" && guildID != r.expectedGuildID {
		return []guilddomain.MemberContribution{}, nil
	}

	return r.members, nil
}

func (r guildTestRepository) CreateMembership(ctx context.Context, membership guilddomain.Membership) error {
	if r.activeMembership != nil {
		return guildapp.ErrAlreadyJoined
	}

	return nil
}

func (r *guildTestRepository) UpdateMembership(ctx context.Context, membership guilddomain.Membership) error {
	if r.activeMembership == nil {
		return guildapp.ErrActiveMembershipNotFound
	}

	r.updated = &membership
	return nil
}

func (r guildTestRepository) CreateCPContribution(ctx context.Context, contribution guilddomain.CPContribution) error {
	return nil
}

func (r guildTestRepository) ListCPContributionsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.CPContribution, error) {
	var contributions []guilddomain.CPContribution
	for _, contribution := range r.contributions {
		if contribution.GuildID == guildID {
			contributions = append(contributions, contribution)
		}
	}

	return contributions, nil
}

func (r guildTestRepository) ListCPContributionsByUser(ctx context.Context, userID user.ID, limit int) ([]guilddomain.CPContribution, error) {
	var contributions []guilddomain.CPContribution
	for _, contribution := range r.contributions {
		if contribution.UserID == userID {
			contributions = append(contributions, contribution)
		}
	}

	return contributions, nil
}

type guildTestCurrentUserRepository struct {
	appUser user.User
	ok      bool
}

func (r guildTestCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return r.appUser, r.ok, nil
}

type guildTestIDGenerator struct{}

func (g guildTestIDGenerator) NewID() (string, error) {
	return "membership_1", nil
}

type guildContributionTestIDGenerator struct{}

func (g guildContributionTestIDGenerator) NewID() (string, error) {
	return "guild_cp_contribution_1", nil
}

type guildControllerTestCPSpender struct {
	err error
}

func (s guildControllerTestCPSpender) Spend(ctx context.Context, command contributionpointapp.SpendCommand) (contributionpointdomain.LedgerEntry, error) {
	if s.err != nil {
		return contributionpointdomain.LedgerEntry{}, s.err
	}

	return contributionpointdomain.LedgerEntry{
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
	}, nil
}

func TestGuildControllerListGuilds(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{
		guilds: []guilddomain.Guild{{
			ID:                 "guild_typescript",
			Slug:               "typescript",
			Name:               "TypeScript",
			Description:        "型の力で支えるギルド。",
			Icon:               "📘",
			Color:              "#3178c6",
			SortOrder:          5,
			MemberCount:        12,
			TotalContributedCP: 120,
			CreatedAt:          now,
			UpdatedAt:          now,
		}},
	}, guildTestCurrentUserRepository{ok: true}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		Guilds []struct {
			ID                 string `json:"id"`
			Slug               string `json:"slug"`
			Name               string `json:"name"`
			Description        string `json:"description"`
			Icon               string `json:"icon"`
			Color              string `json:"color"`
			MemberCount        int64  `json:"member_count"`
			TotalContributedCP int64  `json:"total_contributed_cp"`
		} `json:"guilds"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if len(body.Guilds) != 1 {
		t.Fatalf("guilds length = %d, 期待値 1", len(body.Guilds))
	}
	if body.Guilds[0].ID != "guild_typescript" {
		t.Fatalf("id = %q, 期待値 guild_typescript", body.Guilds[0].ID)
	}
	if body.Guilds[0].MemberCount != 12 {
		t.Fatalf("member_count = %d, 期待値 12", body.Guilds[0].MemberCount)
	}
	if body.Guilds[0].TotalContributedCP != 120 {
		t.Fatalf("total_contributed_cp = %d, 期待値 120", body.Guilds[0].TotalContributedCP)
	}
}

func TestGuildControllerJoinGuild(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{
		guilds: []guilddomain.Guild{{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			MemberCount: 1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}},
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/guilds/guild_go/join", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusCreated {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusCreated)
	}

	var body struct {
		Guild struct {
			ID string `json:"id"`
		} `json:"guild"`
		Membership struct {
			ID       string `json:"id"`
			JoinedAt string `json:"joined_at"`
		} `json:"membership"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", body.Guild.ID)
	}
	if body.Membership.JoinedAt == "" {
		t.Fatal("joined_at が設定されている必要があります")
	}
}

func TestGuildControllerGetMyGuild(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
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
			MemberCount:        1,
			TotalContributedCP: 80,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{
		activeMembership: &activeMembership,
		expectedGuildID:  "guild_go",
		members: []guilddomain.MemberContribution{{
			UserID:        "user_1",
			Name:          "Alice",
			TotalEarnedCP: 120,
			JoinedAt:      now.Add(-time.Hour),
		}, {
			UserID:        "user_2",
			Name:          "Bob",
			TotalEarnedCP: 80,
			JoinedAt:      now,
		}},
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/me/guild", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		Guild struct {
			ID                 string `json:"id"`
			TotalContributedCP int64  `json:"total_contributed_cp"`
		} `json:"guild"`
		Members []struct {
			UserID        string `json:"user_id"`
			Name          string `json:"name"`
			TotalEarnedCP int64  `json:"total_earned_cp"`
			JoinedAt      string `json:"joined_at"`
		} `json:"members"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", body.Guild.ID)
	}
	if body.Guild.TotalContributedCP != 80 {
		t.Fatalf("total_contributed_cp = %d, 期待値 80", body.Guild.TotalContributedCP)
	}
	if len(body.Members) != 2 {
		t.Fatalf("members length = %d, 期待値 2", len(body.Members))
	}
	if body.Members[0].Name != "Alice" {
		t.Fatalf("members[0].name = %q, 期待値 Alice", body.Members[0].Name)
	}
	if body.Members[0].TotalEarnedCP != 120 {
		t.Fatalf("members[0].total_earned_cp = %d, 期待値 120", body.Members[0].TotalEarnedCP)
	}
}

func TestGuildControllerGetGuildDashboard(t *testing.T) {
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
			TotalContributedCP: 180,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{
		activeMembership: &activeMembership,
		expectedGuildID:  "guild_go",
		members: []guilddomain.MemberContribution{{
			UserID:        "user_1",
			Name:          "Alice",
			TotalEarnedCP: 120,
			JoinedAt:      now,
		}},
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds/guild_go/dashboard", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		State string `json:"state"`
		Guild struct {
			ID                 string `json:"id"`
			MemberCount        int64  `json:"member_count"`
			TotalContributedCP int64  `json:"total_contributed_cp"`
		} `json:"guild"`
		Members []struct {
			UserID string `json:"user_id"`
		} `json:"members"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.State != "joined" {
		t.Fatalf("state = %q, 期待値 joined", body.State)
	}
	if body.Guild.ID != "guild_go" {
		t.Fatalf("guild id = %q, 期待値 guild_go", body.Guild.ID)
	}
	if body.Guild.MemberCount != 2 {
		t.Fatalf("member_count = %d, 期待値 2", body.Guild.MemberCount)
	}
	if body.Guild.TotalContributedCP != 180 {
		t.Fatalf("total_contributed_cp = %d, 期待値 180", body.Guild.TotalContributedCP)
	}
	if len(body.Members) != 1 {
		t.Fatalf("members length = %d, 期待値 1", len(body.Members))
	}
}

func TestGuildControllerGetGuildDashboardRejectsOtherGuild(t *testing.T) {
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
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{
		activeMembership: &activeMembership,
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds/guild_python/dashboard", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusForbidden {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusForbidden)
	}
}

func TestGuildControllerGetGuildDashboardRejectsMembershipNotFound(t *testing.T) {
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds/guild_go/dashboard", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusNotFound {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusNotFound)
	}
}

func TestGuildControllerLeaveMyGuild(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
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
			MemberCount: 1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	repository := &guildTestRepository{
		activeMembership: &activeMembership,
	}
	controller := NewGuildController(guildapp.NewUseCase(repository, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodDelete, "/me/guild", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusNoContent {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusNoContent)
	}
	if repository.updated == nil {
		t.Fatal("UpdateMembership() が呼ばれる必要があります")
	}
	if repository.updated.LeftAt == nil {
		t.Fatal("left_at が設定されている必要があります")
	}
}

func TestGuildControllerLeaveMyGuildRejectsMembershipNotFound(t *testing.T) {
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodDelete, "/me/guild", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusNotFound {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusNotFound)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.Error.Code != "guild_membership_not_found" {
		t.Fatalf("error.code = %q, 期待値 guild_membership_not_found", body.Error.Code)
	}
	if body.Error.Message != "active guild membership not found" {
		t.Fatalf("error.message = %q, 期待値 active guild membership not found", body.Error.Message)
	}
}

func TestGuildControllerContributeCP(t *testing.T) {
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
			MemberCount: 1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	controller := NewGuildController(guildapp.NewUseCaseWithCP(&guildTestRepository{
		activeMembership: &activeMembership,
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}, guildContributionTestIDGenerator{}, guildControllerTestCPSpender{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/me/guild/cp-contributions", bytes.NewBufferString(`{"amount":40}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusCreated {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusCreated)
	}

	var body struct {
		Contribution struct {
			ID            string `json:"id"`
			GuildID       string `json:"guild_id"`
			UserID        string `json:"user_id"`
			PointLedgerID string `json:"point_ledger_id"`
			Amount        int64  `json:"amount"`
			CreatedAt     string `json:"created_at"`
		} `json:"contribution"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.Contribution.Amount != 40 {
		t.Fatalf("amount = %d, 期待値 40", body.Contribution.Amount)
	}
	if body.Contribution.GuildID != "guild_go" {
		t.Fatalf("guild_id = %q, 期待値 guild_go", body.Contribution.GuildID)
	}
	if body.Contribution.PointLedgerID != "cp_ledger_1" {
		t.Fatalf("point_ledger_id = %q, 期待値 cp_ledger_1", body.Contribution.PointLedgerID)
	}
}

func TestGuildControllerContributeCPRejectsInsufficientBalance(t *testing.T) {
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
			MemberCount: 1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	controller := NewGuildController(guildapp.NewUseCaseWithCP(&guildTestRepository{
		activeMembership: &activeMembership,
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}, guildContributionTestIDGenerator{}, guildControllerTestCPSpender{
		err: contributionpointapp.ErrInsufficientBalance,
	}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/me/guild/cp-contributions", bytes.NewBufferString(`{"amount":40}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusConflict {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusConflict)
	}
}

func TestGuildControllerListCPContributions(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	controller := NewGuildController(guildapp.NewUseCaseWithCP(&guildTestRepository{
		guilds: []guilddomain.Guild{{
			ID:          "guild_go",
			Slug:        "go",
			Name:        "Go",
			Description: "シンプルさと並列処理で前に進むギルド。",
			Icon:        "GO",
			Color:       "#00acd7",
			SortOrder:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}},
		contributions: []guilddomain.CPContribution{{
			ID:            "guild_cp_contribution_1",
			GuildID:       "guild_go",
			UserID:        "user_1",
			PointLedgerID: "cp_ledger_1",
			Amount:        40,
			CreatedAt:     now,
		}},
	}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}, guildContributionTestIDGenerator{}, guildControllerTestCPSpender{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds/guild_go/cp-contributions", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		Contributions []struct {
			ID     string `json:"id"`
			Amount int64  `json:"amount"`
		} `json:"contributions"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if len(body.Contributions) != 1 {
		t.Fatalf("contributions length = %d, 期待値 1", len(body.Contributions))
	}
	if body.Contributions[0].Amount != 40 {
		t.Fatalf("amount = %d, 期待値 40", body.Contributions[0].Amount)
	}
}

func TestGuildControllerListCPContributionsRejectsUnknownGuild(t *testing.T) {
	controller := NewGuildController(guildapp.NewUseCaseWithCP(&guildTestRepository{}, guildTestCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, guildTestIDGenerator{}, guildContributionTestIDGenerator{}, guildControllerTestCPSpender{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds/guild_missing/cp-contributions", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusNotFound {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusNotFound)
	}
}

func TestGuildControllerJoinGuildRejectsUnauthenticated(t *testing.T) {
	controller := NewGuildController(guildapp.NewUseCase(&guildTestRepository{}, guildTestCurrentUserRepository{}, guildTestIDGenerator{}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/guilds/guild_go/join", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusUnauthorized)
	}
}
