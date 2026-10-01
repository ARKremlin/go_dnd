package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/ARKremlin/go_dnd/internal/pkg/token"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/ARKremlin/go_dnd/internal/pkg/hasher"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	usernameMinLen = 3
	usernameMaxLen = 32
	passwordMinLen = 8
	passwordMaxLen = 72
)

type AuthUseCase struct {
	users  domain.UserRepository
	secret []byte
	ttl    time.Duration
}

func NewAuthUseCase(users domain.UserRepository, secret []byte, ttl time.Duration) *AuthUseCase {
	return &AuthUseCase{users: users, secret: secret, ttl: ttl}
}

func (uc *AuthUseCase) RegisterDM(ctx context.Context, username, password string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	hash, err := hasher.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &domain.User{
		Username:     &username,
		PasswordHash: &hash,
		Role:         domain.RoleDM,
	}
	if err := uc.users.Create(ctx, u); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return u, nil
}

func (uc *AuthUseCase) LoginDM(ctx context.Context, username, password string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	u, err := uc.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to get user by username: %w", err)
	}
	if u.PasswordHash == nil {
		return "", ErrInvalidCredentials
	}
	if !hasher.CheckPassword(*u.PasswordHash, password) {
		return "", ErrInvalidCredentials
	}

	tok, err := token.GenerateToken(u.ID, string(u.Role), uc.secret, uc.ttl)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tok, nil
}

func (uc *AuthUseCase) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	u, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return u, nil
}

func validateUsername(s string) error {
	if s == "" {
		return fmt.Errorf("%w: empty username", ErrInvalidInput)
	}
	n := utf8.RuneCountInString(s)
	if n < usernameMinLen {
		return fmt.Errorf("%w: username too short (min %d)", ErrInvalidInput, usernameMinLen)
	}
	if n > usernameMaxLen {
		return fmt.Errorf("%w: username too long (max %d)", ErrInvalidInput, usernameMaxLen)
	}

	for _, r := range s {
		if r == '_' {
			continue
		}
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		return fmt.Errorf("%w: invalid character in username %q", ErrInvalidInput, r)
	}
	return nil
}

func validatePassword(s string) error {
	if len(s) < passwordMinLen {
		return fmt.Errorf("%w: password too short (min %d)", ErrInvalidInput, passwordMinLen)
	}
	if len(s) > passwordMaxLen {
		return fmt.Errorf("%w: password too long (max %d)", ErrInvalidInput, passwordMaxLen)
	}
	return nil
}
