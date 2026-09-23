package middleware

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const (
	ctxKeyUserID ctxKey = iota
	ctxKeyRole
)

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, id)
}
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	raw := ctx.Value(ctxKeyUserID)
	id, ok := raw.(uuid.UUID)
	return id, ok
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ctxKeyRole, role)
}
func RoleFromContext(ctx context.Context) (string, bool) {
	raw := ctx.Value(ctxKeyRole)
	role, ok := raw.(string)
	return role, ok
}
