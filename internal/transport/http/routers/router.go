package routers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewAPIRouter(
	ctx context.Context,
	log *slog.Logger,
	userManager userManager,
) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Mount("/auth", NewAuthRoutes(ctx, log, userManager))

	return r
}
