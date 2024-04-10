package main

import (
	"github.com/rizzmatch/rm_auth/internal/app"
	"github.com/rizzmatch/rm_auth/internal/config"
	"github.com/rizzmatch/rm_auth/internal/logger"
)

func main() {
	cfg := config.SetupConfig()

	log := logger.SetupLogger(cfg.Env)
	log.Info("zapili")
	log.Debug("zapili")

	application := app.New(log, cfg.REST.Port, cfg.REST.Host, cfg.REST.Timeout)

	go application.RESTServer.Run()

	// TODO: graceful shutdown
}
