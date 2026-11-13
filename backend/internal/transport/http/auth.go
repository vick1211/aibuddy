package httptransport

import (
	"aibuddy/internal/auth"
	"aibuddy/internal/domain"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func bearerToken(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	writeError(w, http.StatusUnauthorized, "unathorized")
}

func RequireAuth(a auth.Authenticator, log *slog.Logger) Middleware {
	if a == nil {
		panic("httptransport: RequireAuth needs a non-nil Authenticator")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				unauthorized(w)
				return
			}
			id, err := a.Authenticate(r.Context(), token)
			switch {
			case err == nil:
				next.ServeHTTP(w, r.WithContext(auth.WithIdentity(r.Context(), id)))
			case errors.Is(err, domain.ErrUnauthorized):
				unauthorized(w)
			default:
				log.Error("authentication backend failure", "err", err)
				writeError(w, http.StatusServiceUnavailable, "authentication unavailable")
			}
		})
	}
}