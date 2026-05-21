package home

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeCurrentUserRepository struct {
	appUser user.User
	found   bool
	err     error
	token   string
	now     time.Time
}

func (r *fakeCurrentUserRepository) FindUserBySessionToken(_ context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	r.token = sessionToken
	r.now = now
	return r.appUser, r.found, r.err
}

type fakeContributionPointReader struct {
	balance     int64
	todayEarned int64
	totalEarned int64
	err         error
	userIDs     []user.ID
}

func (r *fakeContributionPointReader) GetBalance(_ context.Context, userID user.ID) (int64, error) {
	r.userIDs = append(r.userIDs, userID)
	return r.balance, r.err
}

func (r *fakeContributionPointReader) GetTodayEarned(_ context.Context, userID user.ID) (int64, error) {
	r.userIDs = append(r.userIDs, userID)
	return r.todayEarned, r.err
}

func (r *fakeContributionPointReader) GetTotalEarned(_ context.Context, userID user.ID) (int64, error) {
	r.userIDs = append(r.userIDs, userID)
	return r.totalEarned, r.err
}

func TestUseCaseGetHome(t *testing.T) {
	t.Run("セッショントークンからホーム表示用CPとレベル進捗を返す", func(t *testing.T) {
		auth := &fakeCurrentUserRepository{
			appUser: user.User{ID: "user_1"},
			found:   true,
		}
		cp := &fakeContributionPointReader{
			balance:     1200,
			todayEarned: 300,
			totalEarned: 2500,
		}
		usecase := NewUseCase(auth, cp)
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC)
		}

		data, err := usecase.GetHome(context.Background(), "session-token")
		if err != nil {
			t.Fatalf("GetHome() がエラーを返しました: %v", err)
		}

		if auth.token != "session-token" {
			t.Fatalf("認証へ渡した token = %q, 期待値 session-token", auth.token)
		}
		if auth.now != usecase.now() {
			t.Fatalf("認証へ渡した now = %v, 期待値 %v", auth.now, usecase.now())
		}
		if len(cp.userIDs) != 3 {
			t.Fatalf("CP reader 呼び出し回数 = %d, 期待値 3", len(cp.userIDs))
		}
		for _, got := range cp.userIDs {
			if got != "user_1" {
				t.Fatalf("CP reader に渡した userID = %q, 期待値 user_1", got)
			}
		}
		if data.TotalCP != 1200 {
			t.Fatalf("TotalCP = %d, 期待値 1200", data.TotalCP)
		}
		if data.TodayCP != 300 {
			t.Fatalf("TodayCP = %d, 期待値 300", data.TodayCP)
		}
		if data.LifetimeTotalEarnedCP != 2500 {
			t.Fatalf("LifetimeTotalEarnedCP = %d, 期待値 2500", data.LifetimeTotalEarnedCP)
		}
		if data.PlayerLevel != 6 {
			t.Fatalf("PlayerLevel = %d, 期待値 6", data.PlayerLevel)
		}
		if data.PlayerLevelTotalCP != 2500 {
			t.Fatalf("PlayerLevelTotalCP = %d, 期待値 2500", data.PlayerLevelTotalCP)
		}
		if data.NextPlayerLevel != 7 {
			t.Fatalf("NextPlayerLevel = %d, 期待値 7", data.NextPlayerLevel)
		}
		if data.NextPlayerLevelTotalCP != 3600 {
			t.Fatalf("NextPlayerLevelTotalCP = %d, 期待値 3600", data.NextPlayerLevelTotalCP)
		}
		if data.NextPlayerLevelRemaining != 1100 {
			t.Fatalf("NextPlayerLevelRemaining = %d, 期待値 1100", data.NextPlayerLevelRemaining)
		}
	})

	t.Run("空のセッショントークンは未認証エラーを返す", func(t *testing.T) {
		usecase := NewUseCase(&fakeCurrentUserRepository{}, &fakeContributionPointReader{})

		_, err := usecase.GetHome(context.Background(), "")

		if !errors.Is(err, ErrUnauthenticated) {
			t.Fatalf("GetHome() error = %v, 期待値 ErrUnauthenticated", err)
		}
	})

	t.Run("セッションに紐づくユーザーがない場合は未認証エラーを返す", func(t *testing.T) {
		usecase := NewUseCase(
			&fakeCurrentUserRepository{found: false},
			&fakeContributionPointReader{},
		)

		_, err := usecase.GetHome(context.Background(), "missing-session")

		if !errors.Is(err, ErrUnauthenticated) {
			t.Fatalf("GetHome() error = %v, 期待値 ErrUnauthenticated", err)
		}
	})

	t.Run("CP取得エラーは呼び出し元に返す", func(t *testing.T) {
		cpErr := errors.New("cp unavailable")
		usecase := NewUseCase(
			&fakeCurrentUserRepository{appUser: user.User{ID: "user_1"}, found: true},
			&fakeContributionPointReader{err: cpErr},
		)

		_, err := usecase.GetHome(context.Background(), "session-token")

		if !errors.Is(err, cpErr) {
			t.Fatalf("GetHome() error = %v, 期待値 %v", err, cpErr)
		}
	})
}
