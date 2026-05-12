package http

import (
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	t.Run("ヘルスチェックを返す", func(t *testing.T) {
		router := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))
		request := httptest.NewRequest(stdhttp.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != stdhttp.StatusOK {
			t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
		}

		if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Fatalf("Content-Type = %q, 期待値 application/json*", got)
		}

		if got := strings.TrimSpace(response.Body.String()); got != `{"status":"ok"}` {
			t.Fatalf("レスポンスボディ = %q, 期待値 health status", got)
		}
	})
}
