package main

import (
	"build-your-own-reverse-proxy/internal/config"
	"build-your-own-reverse-proxy/internal/runtime"
	"build-your-own-reverse-proxy/internal/server"
	"github.com/lmittmann/tint"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{Level: slog.LevelInfo, TimeFormat: time.Kitchen}))
	slog.SetDefault(logger)

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("Config loaded", "config", cfg)

	rt, err := runtime.Build(cfg, logger)
	if err != nil {
		logger.Error("Failed to create runtime", "error", err)
		os.Exit(1)
	}

	logger.Info("Runtime built successfully")

	srv := server.NewServer(cfg, rt.Proxy)

	logger.Info("Server starting")
	if err := srv.StartServer(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}
