package app

import (
	"github.com/labstack/gommon/log"
	grpcapp "github.com/rizzmatch/rm_auth/internal/app/grpc"
	restapp "github.com/rizzmatch/rm_auth/internal/app/rest"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/services/auth"
	"github.com/rizzmatch/rm_auth/internal/storage/postgres"
)

type App struct {
	RESTSrv *restapp.App
	GRPCSrv *grpcapp.App
}

func New(
	cfg config.Config,
) *App {

	postgresStorage, err := postgres.New(cfg.PSQLUser, cfg.PSQLUser, cfg.PSQLDB, cfg.PSQLHost, cfg.PSQLPort, cfg.PSQLSSLMode)
	if err != nil {
		log.Error("failed to connect to db", err)
	}

	authService := auth.New(postgresStorage, cfg.RESTTimeout)

	grpcApp := grpcapp.New(*authService, cfg.GRPCPort)
	restApp := restapp.New(*authService, cfg.RESTHost, cfg.RESTPort, cfg.RESTTimeout)

	return &App{
		RESTSrv: restApp,
		GRPCSrv: grpcApp,
	}
}
