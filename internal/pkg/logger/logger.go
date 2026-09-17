package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type ctxKey string

const (
	loggerKey    ctxKey = "logger"
	requestIDKey ctxKey = "request_id"
)

func New(lvl string) *slog.Logger {
	var level slog.Level
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	log := FromContext(ctx).With(slog.String("request_id", id))
	ctx = context.WithValue(ctx, loggerKey, log)
	ctx = context.WithValue(ctx, requestIDKey, id)
	return ctx
}

func RequestIDFromContext(ctx context.Context) string {
	raw := ctx.Value(requestIDKey)
	id, ok := raw.(string)
	if ok {
		return id
	}
	return ""
}

func FromContext(ctx context.Context) *slog.Logger {
	raw := ctx.Value(loggerKey)
	log, ok := raw.(*slog.Logger)
	if ok {
		return log
	}
	return slog.Default()
}
