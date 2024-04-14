package app

import (
	"log/slog"
	"time"

	restapp "github.com/rizzmatch/rm_auth/internal/app/rest"
	"github.com/rizzmatch/rm_auth/internal/services/auth"
	"github.com/rizzmatch/rm_auth/internal/storage/postgres"
)

type App struct {
	RESTServer *restapp.App
}

func New(
	log *slog.Logger,
	restPort int,
	host string,
	timeout time.Duration,
) *App {

	postgresStorage, err := postgres.New()
	if err != nil {
		log.Error("failed to connect to db", err)
	}

	authService := auth.New(log, postgresStorage, postgresStorage, timeout)

	restApp := restapp.New(log, *authService, host, restPort, timeout)

	return &App{
		RESTServer: restApp,
	}
}
