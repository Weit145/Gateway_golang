package logout

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	TokenAsset string
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=GRPCLogOutUser
type GRPCLogOutUser interface {
	LogOutUser(ctx context.Context, token string) error
}

func New(log *slog.Logger, grpcLogoutUser GRPCLogOutUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.logout.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		req.TokenAsset = r.Header.Get("Authorization")
		if !strings.HasPrefix(req.TokenAsset, "Bearer ") {
			log.Error("Failed to decode handler")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"invalid handler: ",
			))
			return
		}

		err := grpcLogoutUser.LogOutUser(r.Context(), req.TokenAsset)
		if err != nil {
			log.Error("Failed logout")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(
				"Failed logout",
			))
			return
		}

		token := strings.TrimPrefix(req.TokenAsset, "Bearer ")
		token = strings.TrimSpace(token)

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Authorization", "")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Logout user")
	}
}
