package refresh

import (
	"log/slog"
	"net/http"

	"github.com/Weit145/GATEWAY_golang/internal/lib/logger"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	TokenRefresh string
}

type Response struct {
	TokenAsset string
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

		// var req Request

		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Error("Failed cookie", logger.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed cookie"))
			return
		}

		if cookie.Value == "" {
			log.Error("Failed cookie value", logger.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed cookie value"))
			return
		}

		var resp Response
		resp.TokenAsset = "abc123"

		w.Header().Set("Authorization", "Bearer "+resp.TokenAsset)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

	}
}
