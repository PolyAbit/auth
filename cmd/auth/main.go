package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/PolyAbit/auth/internal/app"
	"github.com/PolyAbit/auth/internal/config"
	"github.com/PolyAbit/auth/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	log.Info("Init config and logger")
	log.Info("Current config", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.GRPC.GatewayPort, cfg.StoragePath, cfg.TokenTTL, cfg.TokenSecret)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("code", sign.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}
