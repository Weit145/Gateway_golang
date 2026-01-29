package refresh

import (
	"context"
	"log/slog"
	"net/http"

	GRPCauth "github.com/Weit145/proto-repo/auth"

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

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=GRPCRefreshToken
type GRPCRefreshToken interface {
	RefreshToken(ctx context.Context, refreshToken string) (*GRPCauth.AccessTokenResponse, error)
}

func New(log *slog.Logger, grpcRefreshToken GRPCRefreshToken) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.refresh.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// var req Request

		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Error("Failed cookie")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed cookie"))
			return
		}

		if cookie.Value == "" {
			log.Error("Failed cookie value")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed cookie value"))
			return
		}

		resp, err := grpcRefreshToken.RefreshToken(r.Context(), cookie.Value)
		if err != nil {
			log.Error("Failed refresh token", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"Failed refresh token",
			))
			return
		}

		var out Response
		out.TokenAsset = resp.AccessToken

		w.Header().Set("Authorization", "Bearer "+out.TokenAsset)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

	}
}
