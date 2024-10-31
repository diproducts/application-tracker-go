package auth

import (
	"context"
	"errors"
	resp "github.com/diproducts/application-tracker-go/internal/lib/api/response"
	"github.com/diproducts/application-tracker-go/internal/lib/auth/tokenutil"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
)

type tokenManager interface {
	ExtractUserIDFromAccessToken(tokenStr string) (int64, error)
}

func NewJWTMiddleware(log *slog.Logger, tm tokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.middleware.auth.NewJWTMiddleware"

			log = log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				msg := "Missing authentication token"
				log.Info(msg)

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error(msg))

				return
			}

			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			userID, err := tm.ExtractUserIDFromAccessToken(tokenStr)
			if err != nil {
				if errors.Is(err, tokenutil.ErrInvalidToken) {
					msg := "Invalid access token"

					log.Info(msg)

					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, resp.Error(msg))

					return
				}

				log.Error("failed to validate access token")

				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("internal server error"))

				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
