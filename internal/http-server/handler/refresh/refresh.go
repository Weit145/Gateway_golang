package refresh

import (
	"log/slog"
	"net/http"

	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
)

type Request struct {
	TokenRefresh string
}

type Response struct {
	TokenAsset string
	response.Response
}

type GRPCRefreshToken interface {
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.refresh.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		r.Cookie("refresh_token")
	}
}
