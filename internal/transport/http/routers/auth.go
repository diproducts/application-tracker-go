package routers

import (
	"context"
	"github.com/diproducts/application-tracker-go/internal/transport/http/handlers/user/create"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func NewAuthRoutes(ctx context.Context, log *slog.Logger, userCreator create.UserCreator) chi.Router {
	r := chi.NewRouter()
	r.Post("/register", create.New(ctx, log, userCreator))
	return r
}
