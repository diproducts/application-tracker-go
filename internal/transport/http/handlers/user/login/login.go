package login

import (
	"context"
	"errors"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	resp "github.com/diproducts/application-tracker-go/internal/lib/api/response"
	"github.com/diproducts/application-tracker-go/internal/lib/logger/sl"
	"github.com/diproducts/application-tracker-go/internal/usecase"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type response struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type loginProvider interface {
	Login(ctx context.Context, email, password string) (models.Tokens, error)
}

func New(ctx context.Context, log *slog.Logger, loginProvider loginProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.user.login"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			msg := "failed to decode request"
			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(msg))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Info("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		tokens, err := loginProvider.Login(ctx, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, usecase.ErrInvalidCredentials) {
				log.Info("invalid credentials")

				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("invalid email or password"))

				return
			}
			log.Error("failed to login user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("something went wrong"))

			return
		}

		log.Info("user logged in")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response{
			Response:     resp.OK(),
			AccessToken:  tokens.Access,
			RefreshToken: tokens.Refresh,
		})
	}
}
