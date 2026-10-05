package postgres

import (
	"context"
	"fmt"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.TableRepository = (*TableRepo)(nil)

type TableRepo struct {
	pool *pgxpool.Pool
}

func NewTableRepo(pool *pgxpool.Pool) *TableRepo {
	return &TableRepo{pool: pool}
}

func (r *TableRepo) Create(ctx context.Context, t *domain.Table) error {
	const q = `
INSERT INTO game_tables (dm_id, name)
VALUES ($1, $2)
RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, q, t.DMID, t.Name).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert table: %w", err)
	}
	return nil
}
