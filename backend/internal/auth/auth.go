package auth

import (
	"aibuddy/internal/domain"
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (domain.Identity, error)
}

type ctxKey struct{}

func WithIdentity(ctx context.Context, id domain.Identity) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func IdentityFrom(ctx context.Context) (domain.Identity, bool) {
	id, ok := ctx.Value(ctxKey{}).(domain.Identity)
	return id, ok
}