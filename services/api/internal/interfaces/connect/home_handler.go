package connect

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	homev1 "github.com/jyogi-web/ddd-a-to-z/gen/go/langwar/home/v1"
	"github.com/jyogi-web/ddd-a-to-z/gen/go/langwar/home/v1/homev1connect"
	homeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/home"
)

const sessionCookieName = "lang_war_session"

type homeUseCase interface {
	GetHome(ctx context.Context, sessionToken string) (homeapp.HomeData, error)
}

type HomeHandler struct {
	uc homeUseCase
}

var _ homev1connect.HomeServiceHandler = (*HomeHandler)(nil)

func NewHomeHandler(uc homeUseCase) *HomeHandler {
	return &HomeHandler{uc: uc}
}

func (h *HomeHandler) GetHome(
	ctx context.Context,
	req *connect.Request[homev1.GetHomeRequest],
) (*connect.Response[homev1.GetHomeResponse], error) {
	cookie, err := extractCookie(req.Header(), sessionCookieName)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	data, err := h.uc.GetHome(ctx, cookie)
	if errors.Is(err, homeapp.ErrUnauthenticated) {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&homev1.GetHomeResponse{
		TotalCp:                  data.TotalCP,
		TodayCp:                  data.TodayCP,
		PlayerLevel:              data.PlayerLevel,
		PlayerLevelTotalCp:       data.PlayerLevelTotalCP,
		NextPlayerLevel:          data.NextPlayerLevel,
		NextPlayerLevelTotalCp:   data.NextPlayerLevelTotalCP,
		NextPlayerLevelRemaining: data.NextPlayerLevelRemaining,
		LifetimeTotalEarnedCp:    data.LifetimeTotalEarnedCP,
	}), nil
}

func extractCookie(headers http.Header, name string) (string, error) {
	req := &http.Request{Header: headers}
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
