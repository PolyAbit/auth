package main

import (
	"log/slog"

	"github.com/PolyAbit/auth/internal/app"
	"github.com/PolyAbit/auth/internal/config"
	"github.com/PolyAbit/auth/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	log.Info("Init config and logger")
	log.Info("Current config", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL, cfg.TokenSecret)

	application.GRPCServer.MustRun()
}
