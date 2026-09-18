package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RolePlayer Role = "player"
	RoleDM     Role = "dm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID               uuid.UUID
	TelegramID       *int64
	TelegramUsername *string
	Username         *string
	PasswordHash     *string
	Role             Role
	CreatedAt        time.Time
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	GetByTelegramUsername(ctx context.Context, telegramUsername string) (*User, error)
	Update(ctx context.Context, u *User) error
}
