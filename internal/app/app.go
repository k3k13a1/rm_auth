package app

import (
	"log/slog"

	restapp "github.com/rizzmatch/rm_auth/internal/app/rest"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/services/auth"
	"github.com/rizzmatch/rm_auth/internal/storage/postgres"
)

type App struct {
	RESTServer *restapp.App
}

func New(
	log *slog.Logger,
	cfg config.Config,
) *App {

	postgresStorage, err := postgres.New(cfg.PSQLUser, cfg.PSQLUser, cfg.PSQLDB, cfg.PSQLHost, cfg.PSQLPort, cfg.PSQLSSLMode)
	if err != nil {
		log.Error("failed to connect to db", err)
	}

	authService := auth.New(log, postgresStorage, postgresStorage, cfg.RESTTimeout)

	restApp := restapp.New(log, *authService, cfg.RESTHost, cfg.RESTPort, cfg.RESTTimeout)

	return &App{
		RESTServer: restApp,
	}
}
