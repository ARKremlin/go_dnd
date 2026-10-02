package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ARKremlin/go_dnd/internal/domain"
)

type PlayerUseCase struct {
	users domain.UserRepository
	log   *slog.Logger
}

func NewPlayerUseCase(users domain.UserRepository, log *slog.Logger) *PlayerUseCase {
	return &PlayerUseCase{users: users, log: log}
}

func (uc *PlayerUseCase) GetOrCreatePlayer(ctx context.Context, telegramID int64, telegramUsername string) (*domain.User, error) {
	username := domain.NormalizeTelegramUsername(telegramUsername)

	u, err := uc.users.GetByTelegramID(ctx, telegramID)
	switch {
	case err == nil:
		return uc.syncUsername(ctx, u, username)
	case errors.Is(err, domain.ErrUserNotFound):
		return uc.createPlayer(ctx, telegramID, username)
	default:
		return nil, fmt.Errorf("failed to get user by telegram id(%d): %w", telegramID, err)
	}
}

func (uc *PlayerUseCase) syncUsername(ctx context.Context, u *domain.User, username string) (*domain.User, error) {
	if username == "" {
		return u, nil
	}
	if u.TelegramUsername != nil && *u.TelegramUsername == username {
		return u, nil
	}

	old := u.TelegramUsername
	u.TelegramUsername = &username

	err := uc.users.Update(ctx, u)
	switch {
	case err == nil:
		return u, nil
	case errors.Is(err, domain.ErrUserAlreadyExists):
		u.TelegramUsername = old
		uc.log.Warn("telegram username is taken by another account, keeping the old one",
			"user_id", u.ID)
		return u, nil
	default:
		u.TelegramUsername = old
		return nil, fmt.Errorf("failed to update telegram username: %w", err)
	}
}

func (uc *PlayerUseCase) createPlayer(ctx context.Context, telegramID int64, username string) (*domain.User, error) {
	u := &domain.User{
		TelegramID: &telegramID,
		Role:       domain.RolePlayer,
	}
	if username != "" {
		u.TelegramUsername = &username
	}

	err := uc.users.Create(ctx, u)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	existing, getErr := uc.users.GetByTelegramID(ctx, telegramID)
	switch {
	case getErr == nil:
		return uc.syncUsername(ctx, existing, username)
	case errors.Is(getErr, domain.ErrUserNotFound):

	default:
		return nil, fmt.Errorf("failed to get player after conflict: %w", getErr)
	}

	if username == "" {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	uc.log.Warn("telegram username is taken by another account, creating player without it",
		"telegram_id", telegramID)
	u.TelegramUsername = nil
	if err := uc.users.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create player without username: %w", err)
	}
	return u, nil
}
