// Package http は HTTP ルーティングとリクエスト・レスポンスのアダプターを管理する。
package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	stdhttp "net/http"
)

type RouteRegistrar interface {
	RegisterRoutes(mux *stdhttp.ServeMux)
}

func NewRouter(logger *slog.Logger, registrars ...RouteRegistrar) *stdhttp.ServeMux {
	mux := stdhttp.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		var body bytes.Buffer

		if err := json.NewEncoder(&body).Encode(map[string]string{"status": "ok"}); err != nil {
			logger.Error("failed to encode health response", "error", err)
			stdhttp.Error(w, "Internal Server Error", stdhttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusOK)

		if _, err := w.Write(body.Bytes()); err != nil {
			logger.Error("failed to write health response", "error", err)
		}
	})

	for _, registrar := range registrars {
		registrar.RegisterRoutes(mux)
	}

	mux.Handle("/api/", stdhttp.StripPrefix("/api", mux))

	return mux
}
