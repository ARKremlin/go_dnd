package usecase

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/google/uuid"
)

type mockPlayerRepo struct {
	domain.UserRepository
	existing      *domain.User
	afterConflict *domain.User
	getErr        error
	createErrs    []error
	updateErr     error
	gets          int
	creates       int
	updates       int
}

func (m *mockPlayerRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	m.gets++
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.gets > 1 && m.afterConflict != nil {
		return m.afterConflict, nil
	}
	if m.existing == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.existing, nil
}

func (m *mockPlayerRepo) Create(ctx context.Context, u *domain.User) error {
	m.creates++
	if m.creates <= len(m.createErrs) {
		return m.createErrs[m.creates-1]
	}
	return nil
}

func (m *mockPlayerRepo) Update(ctx context.Context, u *domain.User) error {
	m.updates++
	return m.updateErr
}

func playerWithUsername(username string) *domain.User {
	tgID := int64(42)
	u := &domain.User{ID: uuid.New(), TelegramID: &tgID, Role: domain.RolePlayer}
	if username != "" {
		u.TelegramUsername = &username
	}
	return u
}

func usernameOf(u *domain.User) string {
	if u.TelegramUsername == nil {
		return ""
	}
	return *u.TelegramUsername
}

func TestPlayerUsecase_GetOrCreatePlayer(t *testing.T) {
	dbErr := errors.New("connection refused")

	test := []struct {
		name          string
		inUsername    string
		existing      *domain.User
		getErr        error
		createErrs    []error
		updateErr     error
		afterConflict *domain.User
		wantErr       error
		wantUsername  string
		wantCreates   int
		wantUpdates   int
	}{
		{
			name:         "new player with username",
			inUsername:   "@Petr",
			wantUsername: "petr",
			wantCreates:  1,
		},
		{
			name:        "new player without username",
			inUsername:  "",
			wantCreates: 1,
		},
		{
			name:         "existing player, same username",
			inUsername:   "@Petr",
			existing:     playerWithUsername("petr"),
			wantUsername: "petr",
		},
		{
			name:         "existing player, username changed",
			inUsername:   "@New_petr",
			existing:     playerWithUsername("old_petr"),
			wantUsername: "new_petr",
			wantUpdates:  1,
		},
		{
			name:         "existing player, empty username",
			inUsername:   "",
			existing:     playerWithUsername("petr"),
			wantUsername: "petr",
		},
		{
			name:         "existing player without username gets one",
			inUsername:   "@Petr",
			existing:     playerWithUsername(""),
			wantUsername: "petr",
			wantUpdates:  1,
		},
		{
			name:    "get failure",
			getErr:  dbErr,
			wantErr: dbErr,
		},
		{
			name:        "create failure",
			createErrs:  []error{dbErr},
			wantErr:     dbErr,
			wantCreates: 1,
		},
		{
			name:         "username taken on update keeps old one",
			inUsername:   "taken",
			existing:     playerWithUsername("old_petr"),
			updateErr:    domain.ErrUserAlreadyExists,
			wantUsername: "old_petr",
			wantUpdates:  1,
		},
		{
			name:        "update failure",
			inUsername:  "new_petr",
			existing:    playerWithUsername("old_petr"),
			updateErr:   dbErr,
			wantErr:     dbErr,
			wantUpdates: 1,
		},
		{
			name:          "parallel start already created the player",
			inUsername:    "petr",
			createErrs:    []error{domain.ErrUserAlreadyExists},
			afterConflict: playerWithUsername("petr"),
			wantUsername:  "petr",
			wantCreates:   1,
		},
		{
			name:         "username taken on create makes player without it",
			inUsername:   "taken",
			createErrs:   []error{domain.ErrUserAlreadyExists, nil},
			wantUsername: "",
			wantCreates:  2,
		},
		{
			name:        "conflict without username is an error",
			inUsername:  "",
			createErrs:  []error{domain.ErrUserAlreadyExists},
			wantErr:     domain.ErrUserAlreadyExists,
			wantCreates: 1,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPlayerRepo{
				existing:      tt.existing,
				afterConflict: tt.afterConflict,
				getErr:        tt.getErr,
				createErrs:    tt.createErrs,
				updateErr:     tt.updateErr,
			}
			
			uc := NewPlayerUseCase(repo, slog.New(slog.DiscardHandler))

			got, err := uc.GetOrCreatePlayer(context.Background(), 42, tt.inUsername)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("GetOrCreatePlayer(): got error %v, want %v", err, tt.wantErr)
			}
			if repo.creates != tt.wantCreates {
				t.Errorf("Create calls: got %d, want %d", repo.creates, tt.wantCreates)
			}
			if repo.updates != tt.wantUpdates {
				t.Errorf("Update calls: got %d, want %d", repo.updates, tt.wantUpdates)
			}
			if tt.wantErr != nil {
				return
			}

			if got.Role != domain.RolePlayer {
				t.Errorf("role: got %q, want %q", got.Role, domain.RolePlayer)
			}
			if got.TelegramID == nil || *got.TelegramID != 42 {
				t.Errorf("telegram id: got %v, want 42", got.TelegramID)
			}
			if usernameOf(got) != tt.wantUsername {
				t.Errorf("telegram username: got %q, want %q", usernameOf(got), tt.wantUsername)
			}
		})
	}
}
