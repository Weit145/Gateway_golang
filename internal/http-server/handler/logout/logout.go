package logout

import (
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

type GRPCLogOutUser interface {
}

func New(log *slog.Logger) http.HandlerFunc {
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
		token := strings.TrimPrefix(req.TokenAsset, "Bearer ")
		token = strings.TrimSpace(token)

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token", // та же кука, которую нужно удалить
			Value:    "",              // очищаем значение
			Path:     "/",             // путь должен совпадать с кукой, чтобы браузер её удалил
			MaxAge:   -1,              // говорит браузеру удалить куку сразу
			HttpOnly: true,            // как было, можно оставить true
			Secure:   false,           // для локальной разработки
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Authorization", "")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Success())

		log.Info("Logout user")
	}
}
