package main

import (
	"context"
	"os"
	"os/signal"
	"pvz/configs"
	"pvz/internal"
	"pvz/internal/bootstrap"
	"pvz/internal/logger"
	"syscall"
	"time"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.Log.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	logger.Init(cfg.LogLvl)
	logger.Log.Info("logger initialized", "level", cfg.LogLvl)

	deps, err := bootstrap.InitDeps()
	if err != nil {
		logger.Log.Error("failed to initialize dependencies", "err", err)
		os.Exit(1)
	}

	server := internal.NewServer(cfg, deps)

	errCh := make(chan error, 1)
	go func() {
		logger.Log.Info("starting server")
		if servErr := server.Start(); servErr != nil {
			errCh <- servErr
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalCh:
		logger.Log.Warn("shutdown signal received", "signal", sig)
	case err := <-errCh:
		logger.Log.Error("server encountered an error", "err", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Log.Info("shutting down server")
	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
		logger.Log.Error("graceful shutdown failed", "err", shutdownErr)
		os.Exit(1)
	}

	logger.Log.Info("server gracefully stopped")
}
