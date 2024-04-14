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

	postgresStorage, err := postgres.New(cfg.POSTGRES.User, cfg.POSTGRES.Password, cfg.POSTGRES.Database, cfg.POSTGRES.Host, cfg.POSTGRES.Port, cfg.POSTGRES.SSLMode)
	if err != nil {
		log.Error("failed to connect to db", err)
	}

	authService := auth.New(log, postgresStorage, postgresStorage, cfg.REST.Timeout)

	restApp := restapp.New(log, *authService, cfg.REST.Host, cfg.REST.Port, cfg.REST.Timeout)

	return &App{
		RESTServer: restApp,
	}
}
