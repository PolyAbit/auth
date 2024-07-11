package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/PolyAbit/auth/internal/app/grpc"
	"github.com/PolyAbit/auth/internal/services/auth"
	"github.com/PolyAbit/auth/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	httpPort int,
	storagePath string,
	tokenTTL time.Duration,
	tokenSecret string,
) *App {
	storage, err := sqlite.New(storagePath)

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL, tokenSecret)

	grpcApp := grpcapp.New(log, authService, grpcPort, httpPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
