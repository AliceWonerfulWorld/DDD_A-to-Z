package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	homeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/home"
)

type homeUseCaseStub struct {
	data homeapp.HomeData
	err  error
}

func (s homeUseCaseStub) GetHome(_ context.Context, _ string) (homeapp.HomeData, error) {
	return s.data, s.err
}

func TestHomeControllerGetHome(t *testing.T) {
	controller := NewHomeController(
		homeUseCaseStub{data: homeapp.HomeData{
			TotalCP:                  1200,
			TodayCP:                  300,
			PlayerLevel:              6,
			PlayerLevelTotalCP:       2500,
			NextPlayerLevel:          7,
			NextPlayerLevelTotalCP:   3600,
			NextPlayerLevelRemaining: 1100,
			LifetimeTotalEarnedCP:    2500,
		}},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	request := httptest.NewRequest(stdhttp.MethodGet, "/home", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	response := httptest.NewRecorder()

	controller.getHome(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		TotalCP                  int64 `json:"total_cp"`
		TodayCP                  int64 `json:"today_cp"`
		PlayerLevel              int   `json:"player_level"`
		PlayerLevelTotalCP       int64 `json:"player_level_total_cp"`
		NextPlayerLevel          int   `json:"next_player_level"`
		NextPlayerLevelTotalCP   int64 `json:"next_player_level_total_cp"`
		NextPlayerLevelRemaining int64 `json:"next_player_level_remaining"`
		LifetimeTotalEarnedCP    int64 `json:"lifetime_total_earned_cp"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンス JSON の decode に失敗しました: %v", err)
	}

	if body.TotalCP != 1200 {
		t.Fatalf("total_cp = %d, 期待値 1200", body.TotalCP)
	}
	if body.TodayCP != 300 {
		t.Fatalf("today_cp = %d, 期待値 300", body.TodayCP)
	}
	if body.PlayerLevel != 6 {
		t.Fatalf("player_level = %d, 期待値 6", body.PlayerLevel)
	}
	if body.PlayerLevelTotalCP != 2500 {
		t.Fatalf("player_level_total_cp = %d, 期待値 2500", body.PlayerLevelTotalCP)
	}
	if body.NextPlayerLevel != 7 {
		t.Fatalf("next_player_level = %d, 期待値 7", body.NextPlayerLevel)
	}
	if body.NextPlayerLevelTotalCP != 3600 {
		t.Fatalf("next_player_level_total_cp = %d, 期待値 3600", body.NextPlayerLevelTotalCP)
	}
	if body.NextPlayerLevelRemaining != 1100 {
		t.Fatalf("next_player_level_remaining = %d, 期待値 1100", body.NextPlayerLevelRemaining)
	}
	if body.LifetimeTotalEarnedCP != 2500 {
		t.Fatalf("lifetime_total_earned_cp = %d, 期待値 2500", body.LifetimeTotalEarnedCP)
	}
}

func TestHomeControllerGetHomeRejectsUnauthenticated(t *testing.T) {
	controller := NewHomeController(
		homeUseCaseStub{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	request := httptest.NewRequest(stdhttp.MethodGet, "/home", nil)
	response := httptest.NewRecorder()

	controller.getHome(response, request)

	if response.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("status = %d, 期待値 %d", response.Code, stdhttp.StatusUnauthorized)
	}
}

func TestHomeControllerGetHomeReturnsCPError(t *testing.T) {
	controller := NewHomeController(
		homeUseCaseStub{err: errors.New("cp unavailable")},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	request := httptest.NewRequest(stdhttp.MethodGet, "/home", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	response := httptest.NewRecorder()

	controller.getHome(response, request)

	if response.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("status = %d, 期待値 %d", response.Code, stdhttp.StatusInternalServerError)
	}
}

func TestHomeControllerGetHomeReturnsUnauthenticatedFromUseCase(t *testing.T) {
	controller := NewHomeController(
		homeUseCaseStub{err: homeapp.ErrUnauthenticated},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	request := httptest.NewRequest(stdhttp.MethodGet, "/home", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})
	response := httptest.NewRecorder()

	controller.getHome(response, request)

	if response.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("status = %d, 期待値 %d", response.Code, stdhttp.StatusUnauthorized)
	}
}
