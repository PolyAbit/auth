package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	grpcauth "github.com/PolyAbit/auth/internal/grpc/auth"
	"github.com/PolyAbit/auth/internal/lib/logger/sl"
	authv1 "github.com/PolyAbit/protos/gen/go/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "google.golang.org/grpc/metadata"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	grpcPort   int
	httpPort   int
}

func New(log *slog.Logger, authService grpcauth.Auth, gRPCPort int, httpPort int) *App {
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	grpcauth.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		grpcPort:   gRPCPort,
		httpPort:   httpPort,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	// const op = "grpcapp.Run"

	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := a.startGrpcServer(); err != nil {
			a.log.Error("failed start grpc", sl.Err(err))
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.startHttpServer(ctx); err != nil {
			a.log.Error("failed start http", sl.Err(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) startGrpcServer() error {
	const op = "app.grpc.startGrpcServer"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			h.ServeHTTP(w, r)
	})
}
func (a *App) startHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := authv1.RegisterAuthHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", a.grpcPort), opts)
	if err != nil {
		return err
	}

	a.log.Info("gateway server started", slog.String("addr", fmt.Sprintf(":%d", a.httpPort)))

	return http.ListenAndServe(fmt.Sprintf("localhost:%d", a.httpPort), enableCors(mux))
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.grpcPort))

	a.gRPCServer.GracefulStop()
}
