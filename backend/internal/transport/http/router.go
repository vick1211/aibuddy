package httptransport

import (
	"log/slog"
	"net/http"
)

type Deps struct {
	Log *slog.Logger
	DB Pinger
}

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", Healthz)
	mux.Handle("GET /readyz", Readyz(d.DB, d.Log))

	return Chain(mux, Logging(d.Log), Recover(d.Log))
}