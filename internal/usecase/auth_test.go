package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	domain.UserRepository
	getByID func(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return m.getByID(ctx, id)
}

func TestAuthUseCase_GetMe(t *testing.T) {
	id := uuid.New()
	username := "master"
	existing := &domain.User{ID: id, Username: &username, Role: domain.RoleDM}
	dbErr := errors.New("connection refused")

	tests := []struct {
		name     string
		repo     func(ctx context.Context, id uuid.UUID) (*domain.User, error)
		wantUser *domain.User
		wantErr  error
	}{
		{
			name: "user found",
			repo: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return existing, nil
			},
			wantUser: existing,
		},
		{
			name: "user not found",
			repo: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name: "repository failure",
			repo: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, dbErr
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewAuthUseCase(&mockUserRepo{getByID: tt.repo}, []byte("secret"), time.Hour)

			got, err := uc.GetMe(context.Background(), id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("uc.GetMe(): got error %v, want %v", err, tt.wantErr)
			}

			if got != tt.wantUser {
				t.Fatalf("uc.GetMe(): got %v, want %v", got, tt.wantUser)
			}
		})
	}
}
