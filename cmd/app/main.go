package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ARKremlin/go_dnd/internal/config"
	"github.com/ARKremlin/go_dnd/internal/migrations"
	"github.com/ARKremlin/go_dnd/internal/pkg/logger"
	"github.com/ARKremlin/go_dnd/internal/repository/postgres"
	"github.com/ARKremlin/go_dnd/internal/transport/rest"
	"github.com/ARKremlin/go_dnd/internal/transport/rest/handler"
	"github.com/ARKremlin/go_dnd/internal/transport/telegram"
	"github.com/ARKremlin/go_dnd/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	lgr := logger.New(cfg.App.Loglevel)
	slog.SetDefault(lgr)
	lgr.Info("config loaded", "port", cfg.App.Port, "log_level", cfg.App.Loglevel)

	if err := migrations.Run(cfg.DB.DSN()); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	lgr.Info("migrations applied")

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.DB.DSN(), postgres.DefaultPoolConfig())
	if err != nil {
		return fmt.Errorf("pool init: %w", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepo(pool)
	authUC := usecase.NewAuthUseCase(userRepo, []byte(cfg.JWT.Secret), cfg.JWT.TTL)
	authHandler := handler.NewAuthHandler(authUC)

	router := rest.NewRouter(rest.RouterDeps{
		AuthHandler: authHandler,
		JWTSecret:   []byte(cfg.JWT.Secret),
	})

	tgBot, err := telegram.New(cfg.Telegram.BotToken, lgr)
	if err != nil {
		return fmt.Errorf("telegram init: %w", err)
	}

	srv := http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		tgBot.Start(ctx)
	}()

	go func() {
		lgr.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	lgr.Info("shutting down server", "addr", srv.Addr)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErr := srv.Shutdown(shutdownCtx)

	wg.Wait()

	if shutdownErr != nil {
		return fmt.Errorf("server shutdown: %w", shutdownErr)
	}

	lgr.Info("shutdown complete")
	return nil
}
