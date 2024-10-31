package app

import (
	"context"
	"errors"
	"github.com/diproducts/application-tracker-go/internal/config"
	"github.com/diproducts/application-tracker-go/internal/lib/auth/password_hasher"
	"github.com/diproducts/application-tracker-go/internal/lib/auth/tokenutil"
	"github.com/diproducts/application-tracker-go/internal/lib/logger/sl"
	"github.com/diproducts/application-tracker-go/internal/repository/storage/postgresql"
	"github.com/diproducts/application-tracker-go/internal/transport/http/routers"
	"github.com/diproducts/application-tracker-go/internal/usecase"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(cfg *config.Config) {
	log := setupLogger(cfg.Env)

	db, err := postgresql.InitDB(&cfg.DB)
	if err != nil {
		log.Error("failed to init db connection", sl.Err(err))
		return
	}
	defer db.Close()

	passwordHasher := password_hasher.NewBcryptPasswordHasher()
	userRepository := postgresql.NewUserRepository(db)

	tokenManager := tokenutil.NewJWTTokenManager(
		cfg.AccessSecret,
		cfg.RefreshSecret,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)

	userUsecase := usecase.NewUserUsecase(passwordHasher, userRepository, tokenManager, log)

	ctx := context.TODO()
	router := chi.NewRouter()
	router.Mount("/api", routers.NewAPIRouter(ctx, log, userUsecase))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server starting", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error starting server", sl.Err(err))
	}

	log.Error("server stopped")

	// TODO: graceful shutdown
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
