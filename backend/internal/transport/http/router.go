package httptransport

import (
	"aibuddy/internal/auth"
	"log/slog"
	"net/http"
)

type Deps struct {
	Log *slog.Logger
	DB Pinger
	Auth auth.Authenticator
}

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", Healthz)
	mux.Handle("GET /readyz", Readyz(d.DB, d.Log))

	protected := RequireAuth(d.Auth, d.Log)
	mux.Handle("GET /me", protected(http.HandlerFunc(Me)))

	return Chain(mux, Logging(d.Log), Recover(d.Log))
}