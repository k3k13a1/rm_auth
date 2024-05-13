package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/rizzmatch/rm_auth/internal/handlers/grpc/auth"
	"github.com/rizzmatch/rm_auth/internal/services/auth"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(
	authService auth.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, &authService)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := slog.With(slog.String("op", op))
	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := slog.With(slog.String("op", op))
	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
