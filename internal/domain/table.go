package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Table struct {
	ID        uuid.UUID
	DMID      uuid.UUID
	Name      string
	CreatedAt time.Time
}

type TableRepository interface {
	Create(ctx context.Context, t *Table) error
}
