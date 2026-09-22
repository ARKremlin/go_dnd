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

	//// SMOKE
	//
	//ctx := context.Background()
	//pool, err := postgres.NewPool(ctx, cfg.DB.DSN(), postgres.DefaultPoolConfig())
	//if err != nil {
	//	lgr.Error("pool init:", err)
	//	os.Exit(1)
	//}
	//defer pool.Close()
	//
	//repo := postgres.NewUserRepo(pool)
	//
	//uname := "dm_" + uuid.New().String()[:8]
	//pwdHash, err := hasher.HashPassword("secret")
	//
	//if err != nil {
	//	lgr.Error("hashing password:", err)
	//	os.Exit(1)
	//}
	//u := domain.User{
	//	Username:     &uname,
	//	PasswordHash: &pwdHash,
	//	Role:         domain.RoleDM,
	//}
	//if err := repo.Create(ctx, &u); err != nil {
	//	lgr.Error("create user:", err)
	//	os.Exit(1)
	//}
	//lgr.Info("created", "id", u.ID, "username", *u.Username, "created_at", u.CreatedAt)
	//
	//got, err := repo.GetByUsername(ctx, uname)
	//if err != nil {
	//	lgr.Error("get user by username:", err)
	//	os.Exit(1)
	//}
	//if got.ID != u.ID {
	//	lgr.Error("id mismatch", "got", got.ID, "created", u.ID)
	//	os.Exit(1)
	//}
	//lgr.Info("found", "id", got.ID, "role", got.Role)
	//dup := &domain.User{
	//	Username:     &uname,
	//	PasswordHash: &pwdHash,
	//	Role:         domain.RoleDM,
	//}
	//err = repo.Create(ctx, dup)
	//if !errors.Is(err, domain.ErrUserAlreadyExists) {
	//	lgr.Error("expected ErrUserAlreadyExists", "got", err)
	//	os.Exit(1)
	//}
	//lgr.Info("duplicate rejected as expected")
	//// end of SMOKE

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
