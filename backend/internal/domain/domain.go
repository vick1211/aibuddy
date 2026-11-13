package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Role string

const (
	RoleIntern  Role = "intern"
	RoleMentor  Role = "mentor"
	RoleManager Role = "manager"
	RoleHRBP    Role = "hr_bp"
)

type Identity struct {
	UserID uuid.UUID
	Role   Role
	Team   string
}

var (
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidArgs  = errors.New("invalid arguments")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)
