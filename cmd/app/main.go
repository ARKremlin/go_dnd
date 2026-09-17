package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ARKremlin/go_dnd/internal/config"
	"github.com/ARKremlin/go_dnd/internal/pkg/logger"
	"github.com/ARKremlin/go_dnd/internal/transport/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config load:", err)
		os.Exit(1)
	}
	lgr := logger.New(cfg.App.Loglevel)
	slog.SetDefault(lgr)
	lgr.Info("config loaded", "port", cfg.App.Port, "log_level", cfg.App.Loglevel)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt)
	defer stop()
	go func() {
		lgr.Info("http server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.Error("http server error", "err", err)
		}
	}()
	<-ctx.Done()
	lgr.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		lgr.Error("http server shutdown error", "err", err)
	}

	lgr.Info("http server stopped")
}
