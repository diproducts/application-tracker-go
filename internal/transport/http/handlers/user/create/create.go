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

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name,omitempty"`
}

type Response struct {
	resp.Response
	Id int64 `json:"id,omitempty"`
}

type UserCreator interface {
	CreateUser(ctx context.Context, email, password, name string) (int64, error)
}

func New(ctx context.Context, log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.user.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			msg := "failed to decode request"
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := userCreator.CreateUser(ctx, req.Email, req.Password, req.Name)
		if err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				msg := "user with this email already exists"
				log.Info(msg)

				render.JSON(w, r, resp.Error(msg))

				return
			}
			msg := "failed to create user"
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))

			return
		}

		log.Info("new user created")

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
