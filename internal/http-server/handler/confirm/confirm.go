package confirm

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Weit145/GATEWAY_golang/internal/lib/logger"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	GRPCauth "github.com/Weit145/proto-repo/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	TokenEmail string
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=GRPCConfirmUser
type GRPCConfirmUser interface {
	RegistrationUser(ctx context.Context, token string) (*GRPCauth.CookieResponse, error)
}

type Response struct {
	TokenAsset string
}

func New(log *slog.Logger, grpcConfirm GRPCConfirmUser) http.HandlerFunc {
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

		resp, err := grpcConfirm.RegistrationUser(r.Context(), token)
		if err != nil {
			log.Error("Failed confirm user", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"Failed confirm user",
			))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     resp.Cookie.Key,
			Value:    resp.Cookie.Value,
			Path:     "/", // доступна всему сайту
			Expires:  time.Now().Add(time.Duration(resp.Cookie.MaxAge) * time.Hour),
			HttpOnly: resp.Cookie.Httponly, // теперь можно видеть в браузере через JS
			Secure:   resp.Cookie.Secure,   // локально без HTTPS
			SameSite: http.SameSiteLaxMode, // Lax позволит отправлять cookie на GET запросы
		})

		var out Response
		out.TokenAsset = resp.AccessToken

		w.Header().Set("Authorization", "Bearer "+out.TokenAsset)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Confirm token received", slog.String("token", req.TokenEmail))
	}
}
