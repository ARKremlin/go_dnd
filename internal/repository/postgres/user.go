package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

const userColumns = "id, telegram_id, telegram_username, username, password_hash, role, created_at"

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	const q = `
INSERT INTO users (telegram_id, telegram_username, username, password_hash, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, q,
		u.TelegramID,
		u.TelegramUsername,
		u.Username,
		u.PasswordHash,
		string(u.Role),
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return mapPgError(err)
	}
	return nil
}
func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = "SELECT " + userColumns + " FROM users WHERE id = $1"
	return r.scanOne(ctx, q, id)
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const q = "SELECT " + userColumns + " FROM users WHERE username = $1"
	return r.scanOne(ctx, q, username)
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	const q = "SELECT " + userColumns + " FROM users WHERE telegram_id = $1"
	return r.scanOne(ctx, q, telegramID)
}

func (r *UserRepo) GetByTelegramUsername(ctx context.Context, telegramUsername string) (*domain.User, error) {
	const q = "SELECT " + userColumns + " FROM users WHERE telegram_username = $1"
	return r.scanOne(ctx, q, telegramUsername)
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	const q = `
			UPDATE users
				SET telegram_username = $1,
					username = $2,
					password_hash = $3
				WHERE id = $4`
	tag, err := r.pool.Exec(ctx, q,
		u.TelegramUsername,
		u.Username,
		u.PasswordHash,
		u.ID)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) scanOne(ctx context.Context, q string, args ...any) (*domain.User, error) {
	var (
		u    domain.User
		role string
	)
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&u.ID,
		&u.TelegramID,
		&u.TelegramUsername,
		&u.Username,
		&u.PasswordHash,
		&role,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, mapPgError(err)
	}
	u.Role = domain.Role(role)
	return &u, nil
}

func mapPgError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrUserAlreadyExists
	}
	return err
}
