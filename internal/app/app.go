package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/PolyAbit/auth/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage and service

	// TODO: replace nil to real service
	grpcApp := grpcapp.New(log, nil, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
