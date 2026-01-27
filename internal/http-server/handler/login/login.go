package login

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	name     string
	password string
}

type GRPCLoginUser interface {
}

type Response struct {
	TokenAsset string
}

func New(log *slog.Logger) http.HandlerFunc {
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

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "abc123",
			Path:     "/", // доступна всему сайту
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,                // теперь можно видеть в браузере через JS
			Secure:   false,                // локально без HTTPS
			SameSite: http.SameSiteLaxMode, // Lax позволит отправлять cookie на GET запросы
		})

		var resp Response
		resp.TokenAsset = "abc123"

		w.Header().Set("Authorization", "Bearer "+resp.TokenAsset)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Login user")

	}
}
