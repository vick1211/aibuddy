package auth

import (
	"aibuddy/internal/domain"
	"context"
)

type MockAuth struct {
	users map[string]domain.Identity
}

func NewMockAuth(users map[string]domain.Identity) *MockAuth {
	return &MockAuth{users: users}
}

func (m *MockAuth) Authenticate(_ context.Context, token string) (domain.Identity, error) {
	id, ok := m.users[token]
	if !ok {
		return domain.Identity{}, domain.ErrUnauthorized
	}
	return id, nil
}