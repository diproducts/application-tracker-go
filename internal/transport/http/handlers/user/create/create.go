package create

import (
	"context"
	"errors"
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
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name,omitempty"`
}

type response struct {
	resp.Response
	Id int64 `json:"id,omitempty"`
}

type userCreator interface {
	CreateUser(ctx context.Context, email, password, name string) (int64, error)
}

func New(ctx context.Context, log *slog.Logger, userCreator userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.user.create"

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

		id, err := userCreator.CreateUser(ctx, req.Email, req.Password, req.Name)
		if err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				msg := "user with this email already exists"
				log.Info(msg)

				render.Status(r, http.StatusConflict)
				render.JSON(w, r, resp.Error(msg))

				return
			}
			msg := "failed to create user"
			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(msg))

			return
		}

		log.Info("new user created")

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
