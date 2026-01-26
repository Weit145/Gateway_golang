package main

import (
	"log/slog"
	"os"

	"github.com/Weit145/GATEWAY_golang/internal/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	//Init config
	cfg := config.MustLoad()

	//Init logger
	log := setupLogger(cfg.Env)
	log.Info("Start GATEWAY")

	//Init router
	router := chi.NewRouter()
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
