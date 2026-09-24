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
	"github.com/ARKremlin/go_dnd/internal/migrations"
	"github.com/ARKremlin/go_dnd/internal/pkg/logger"
	"github.com/ARKremlin/go_dnd/internal/repository/postgres"
	"github.com/ARKremlin/go_dnd/internal/transport/rest"
	"github.com/ARKremlin/go_dnd/internal/transport/rest/handler"
	"github.com/ARKremlin/go_dnd/internal/usecase"
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

	if err := migrations.Run(cfg.DB.DSN()); err != nil {
		lgr.Error("migrations failed", "err", err)
		os.Exit(1)
	}
	lgr.Info("migrations applied")

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.DB.DSN(), postgres.DefaultPoolConfig())
	if err != nil {
		lgr.Error("pool initlgr.Info(\"server stooped\")", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepo(pool)
	authUC := usecase.NewAuthUseCase(userRepo, []byte(cfg.JWT.Secret), cfg.JWT.TTL)
	authHandler := handler.NewAuthHandler(authUC)

	router := rest.NewRouter(rest.RouterDeps{
		AuthHandler: authHandler,
		JWTSecret:   []byte(cfg.JWT.Secret),
	})

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		lgr.Error("server shutdown:", "err", err)
	}

	lgr.Info("server stopped")
}
