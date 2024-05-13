package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rizzmatch/rm_auth/internal/app"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/logger"
)

func main() {
	cfg := config.SetupConfig()

	logger.SetupLogger(cfg.AppEnv)
	slog.Info("starting auth service")

	application := app.New(*cfg)

	go application.RESTSrv.Run()
	go application.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.RESTSrv.Stop()
	application.GRPCSrv.Stop()
	slog.Info("Gracefully stoped")
}
