package main

import (
	"build-your-own-reverse-proxy/src/internal/config"
	"build-your-own-reverse-proxy/src/internal/runtime"
	"build-your-own-reverse-proxy/src/internal/server"
	"github.com/lmittmann/tint"
	"log/slog"
	"net/http"
	"os"
	"time"
	"os/signal"
	"context"
)

func main() {
	logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen}))
	slog.SetDefault(logger)

	cfg, err := config.Load("configs/config_advanced.yaml")
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

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<- sigChan

		logger.Info("Shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.StopServer(ctx); err != nil {
			logger.Error("Failed to stop server", "error", err)
		} else {
			logger.Info("Server stopped gracefully")
		}
	}()

	logger.Info("Server starting")
	if err := srv.StartServer(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}
