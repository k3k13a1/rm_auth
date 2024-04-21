package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/gommon/log"
	"github.com/rizzmatch/rm_auth/internal/app"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/logger"
)

func main() {
	cfg := config.SetupConfig()

	logger.SetupLogger(cfg.AppEnv)
	slog.Info("zapili")
	slog.Debug("zapili")

	application := app.New(*cfg)

	application.RESTServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.RESTServer.Stop()
	log.Info("Gracefully stoped")
}
