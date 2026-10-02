package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/ARKremlin/go_dnd/internal/domain"
)

type PlayerUseCase struct {
	users domain.UserRepository
}

func NewPlayerUseCase(users domain.UserRepository) *PlayerUseCase {
	return &PlayerUseCase{users: users}
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

	u.TelegramUsername = &username
	if err := uc.users.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update telegram username: %w", err)
	}
	return u, nil
}

func (uc *PlayerUseCase) createPlayer(ctx context.Context, telegramID int64, username string) (*domain.User, error) {
	u := &domain.User{
		TelegramID: &telegramID,
		Role:       domain.RolePlayer,
	}
	if username != "" {
		u.TelegramUsername = &username
	}

	if err := uc.users.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return u, nil
}
