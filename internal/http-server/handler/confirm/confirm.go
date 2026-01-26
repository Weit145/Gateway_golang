package confirm

import (
	"log/slog"
	"net/http"

	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	TokenEmail string
}

type GRPCConfirmUser interface {
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.confirm.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		token := chi.URLParam(r, "token")
		log.Info("Token param", slog.String("token", token))

		if token == "" {
			log.Error("Token param missing")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("token is required"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Confirm token received", slog.String("token", req.TokenEmail))
	}
}
