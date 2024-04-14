package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rizzmatch/rm_auth/internal/app"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/logger"
)

func main() {
	cfg := config.SetupConfig()

	log := logger.SetupLogger(cfg.Env)
	log.Info("zapili")
	log.Debug("zapili")

	application := app.New(log, *cfg)

	application.RESTServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.RESTServer.Stop()
	log.Info("Gracefully stoped")
}
