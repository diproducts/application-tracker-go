package routers

import (
	"context"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/diproducts/application-tracker-go/internal/transport/http/handlers/user/create"
	"github.com/diproducts/application-tracker-go/internal/transport/http/handlers/user/login"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type userManager interface {
	CreateUser(ctx context.Context, email, password, name string) (int64, error)
	Login(ctx context.Context, email, password string) (models.Tokens, error)
}

func NewAuthRoutes(ctx context.Context, log *slog.Logger, userManager userManager) chi.Router {
	r := chi.NewRouter()
	r.Post("/register", create.New(ctx, log, userManager))
	r.Post("/login", login.New(ctx, log, userManager))
	return r
}
