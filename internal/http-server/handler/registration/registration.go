package registration

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	grpc "github.com/Weit145/proto-repo/auth"

	"github.com/Weit145/GATEWAY_golang/internal/lib/logger"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Email string `json:"email" validate:"required,email"`
	Login string `json:"login" validate:"required,min=4"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=GRPCRegistrationUser
type GRPCRegistrationUser interface {
	CreateUser(ctx context.Context, login, email, password, username string) (*grpc.Okey, error)
}

func New(log *slog.Logger, grpcRegistration GRPCRegistrationUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.registration.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		username, password, ok := r.BasicAuth()

		if len(username) < 4 || len(password) < 6 || !ok || username == "" || password == "" {
			log.Error("Failed to decode handler")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"invalid handler: ",
			))
			return
		}

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"empty request",
			))

			return
		}

		if err != nil {
			log.Error("Failed to decode request", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"invalid request body: ",
			))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			log.Error("Validation error", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"validation error: "+err.Error(),
			))
			return
		}

		_, err = grpcRegistration.CreateUser(r.Context(), req.Login, req.Email, password, username)
		if err != nil {
			log.Error("Failed create user", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"Failed create user",
			))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Registration: ",
			slog.String("Username: ", username),
		)
	}
}
