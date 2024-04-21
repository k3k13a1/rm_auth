package restapp

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	authrest "github.com/rizzmatch/rm_auth/internal/handlers/rest/auth"
	"github.com/rizzmatch/rm_auth/internal/services/auth"
)

type App struct {
	RESTServer *echo.Echo
	host       string
	port       int
	timeout    time.Duration
}

func New(
	authService auth.Auth,
	host string,
	port int,
	timeout time.Duration,
) *App {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/login", func(c echo.Context) error {
		return authrest.Login(c, &authService)
	})

	e.POST("/register", func(c echo.Context) error {
		return authrest.Register(c, &authService)
	})

	return &App{
		RESTServer: e,
		host:       host,
		port:       port,
		timeout:    timeout,
	}
}

func (a *App) Run() error {
	const op = "restapp.Run"

	slog.Info("starting this shit server", slog.String("op", op))

	if err := a.RESTServer.Start(":9241"); err != nil && err != http.ErrServerClosed {
		slog.Error("shutting down the server")
	}
	return nil
}

func (a *App) Stop() error {
	return a.RESTServer.Shutdown(context.TODO())
}
