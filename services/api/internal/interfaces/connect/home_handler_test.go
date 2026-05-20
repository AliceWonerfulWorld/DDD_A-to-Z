package connect

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	homev1 "github.com/jyogi-web/ddd-a-to-z/gen/go/langwar/home/v1"
	homeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/home"
)

type homeUseCaseStub struct {
	data  homeapp.HomeData
	err   error
	token string
}

func (s *homeUseCaseStub) GetHome(_ context.Context, sessionToken string) (homeapp.HomeData, error) {
	s.token = sessionToken
	return s.data, s.err
}

func TestHomeHandlerGetHome(t *testing.T) {
	t.Run("Cookieのセッショントークンでホーム情報を返す", func(t *testing.T) {
		uc := &homeUseCaseStub{data: homeapp.HomeData{
			TotalCP:                  1200,
			TodayCP:                  300,
			PlayerLevel:              6,
			PlayerLevelTotalCP:       2500,
			NextPlayerLevel:          7,
			NextPlayerLevelTotalCP:   3600,
			NextPlayerLevelRemaining: 1100,
			LifetimeTotalEarnedCP:    2500,
		}}
		handler := NewHomeHandler(uc)
		req := connect.NewRequest(&homev1.GetHomeRequest{})
		req.Header().Add("Cookie", (&http.Cookie{Name: sessionCookieName, Value: "session-token"}).String())

		resp, err := handler.GetHome(context.Background(), req)
		if err != nil {
			t.Fatalf("GetHome() がエラーを返しました: %v", err)
		}

		if uc.token != "session-token" {
			t.Fatalf("usecase に渡した token = %q, 期待値 session-token", uc.token)
		}
		if resp.Msg.TotalCp != 1200 {
			t.Fatalf("TotalCp = %d, 期待値 1200", resp.Msg.TotalCp)
		}
		if resp.Msg.TodayCp != 300 {
			t.Fatalf("TodayCp = %d, 期待値 300", resp.Msg.TodayCp)
		}
		if resp.Msg.PlayerLevel != 6 {
			t.Fatalf("PlayerLevel = %d, 期待値 6", resp.Msg.PlayerLevel)
		}
		if resp.Msg.NextPlayerLevelRemaining != 1100 {
			t.Fatalf("NextPlayerLevelRemaining = %d, 期待値 1100", resp.Msg.NextPlayerLevelRemaining)
		}
		if resp.Msg.LifetimeTotalEarnedCp != 2500 {
			t.Fatalf("LifetimeTotalEarnedCp = %d, 期待値 2500", resp.Msg.LifetimeTotalEarnedCp)
		}
	})

	t.Run("Cookieがない場合はUnauthenticatedを返す", func(t *testing.T) {
		handler := NewHomeHandler(&homeUseCaseStub{})

		_, err := handler.GetHome(context.Background(), connect.NewRequest(&homev1.GetHomeRequest{}))

		if connect.CodeOf(err) != connect.CodeUnauthenticated {
			t.Fatalf("connect.CodeOf(error) = %v, 期待値 %v", connect.CodeOf(err), connect.CodeUnauthenticated)
		}
	})

	t.Run("未認証エラーはUnauthenticatedに変換する", func(t *testing.T) {
		handler := NewHomeHandler(&homeUseCaseStub{err: homeapp.ErrUnauthenticated})
		req := connect.NewRequest(&homev1.GetHomeRequest{})
		req.Header().Add("Cookie", (&http.Cookie{Name: sessionCookieName, Value: "expired-token"}).String())

		_, err := handler.GetHome(context.Background(), req)

		if connect.CodeOf(err) != connect.CodeUnauthenticated {
			t.Fatalf("connect.CodeOf(error) = %v, 期待値 %v", connect.CodeOf(err), connect.CodeUnauthenticated)
		}
	})

	t.Run("予期しないエラーはInternalに変換する", func(t *testing.T) {
		handler := NewHomeHandler(&homeUseCaseStub{err: errors.New("database unavailable")})
		req := connect.NewRequest(&homev1.GetHomeRequest{})
		req.Header().Add("Cookie", (&http.Cookie{Name: sessionCookieName, Value: "session-token"}).String())

		_, err := handler.GetHome(context.Background(), req)

		if connect.CodeOf(err) != connect.CodeInternal {
			t.Fatalf("connect.CodeOf(error) = %v, 期待値 %v", connect.CodeOf(err), connect.CodeInternal)
		}
	})
}
