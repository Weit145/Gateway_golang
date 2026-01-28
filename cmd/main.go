package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Weit145/GATEWAY_golang/internal/config"
	"github.com/Weit145/GATEWAY_golang/internal/grpc/auth"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/confirm"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/login"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/logout"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/refresh"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/registration"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	//Init config
	cfg := config.MustLoad()

	//Init logger
	log := setupLogger(cfg.Env)
	log.Info("Start GATEWAY")

	//Init grpc client

	client, err := auth.New("auth-service:50051")
	if err != nil {
		log.Error("failed to create auth client:", err)
	}
	//Init router
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/registration", func(r chi.Router) {
		r.Post("/", registration.New(log, client))
		r.Get("/confirm/{token}", confirm.New(log, client))
	})

	router.Get("/refresh", refresh.New(log))
	router.Post("/login", login.New(log))
	router.Get("/logout", logout.New(log))

	srv := &http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start HTTP server", logger.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
