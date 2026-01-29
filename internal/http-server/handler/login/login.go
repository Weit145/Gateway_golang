package login

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	GRPCauth "github.com/Weit145/proto-repo/auth"

	"github.com/Weit145/GATEWAY_golang/internal/lib/logger"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	name     string
	password string
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=GRPCLoginUser
type GRPCLoginUser interface {
	Authenticate(ctx context.Context, login, password string) (*GRPCauth.CookieResponse, error)
}

type Response struct {
	TokenAsset string
}

func New(log *slog.Logger, grpcLoginUser GRPCLoginUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		var ok bool

		req.name, req.password, ok = r.BasicAuth()
		if len(req.name) < 4 || len(req.password) < 6 || !ok || req.name == "" || req.password == "" {
			log.Error("Failed to decode handler")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"invalid handler: ",
			))
			return
		}

		resp, err := grpcLoginUser.Authenticate(r.Context(), req.name, req.password)
		if err != nil {
			log.Error("Failed login user", logger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(
				"Failed login user",
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

		log.Info("Login user")

	}
}
