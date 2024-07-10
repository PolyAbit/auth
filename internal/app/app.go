package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/PolyAbit/auth/internal/app/grpc"
	"github.com/PolyAbit/auth/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	tokenSecret string,
) *App {
	// TODO: init storage and service

	authService := auth.New(log, nil, tokenTTL, tokenSecret)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
