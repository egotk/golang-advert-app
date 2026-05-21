package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	coreconfig "github.com/egotk/golang-advert-app/internal/core/config"
	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	corelogger "github.com/egotk/golang-advert-app/internal/core/logger"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
)

func main() {
	config := coreconfig.NewMust()
	time.Local = config.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	fmt.Println("advert app starting")
	logger, err := corezaplogger.New(corezaplogger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("app time zone", corelogger.Any("zone", time.Local))

	logger.Debug("init HTTP server")
	httpServer := corehttp.NewServer(
		corehttp.NewConfigMust(),
		logger,
		corehttp.RequestID(),
		corehttp.Logger(logger),
	)

	apiVersionRouter := corehttp.NewAPIVersionRouter(corehttp.ApiV1)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server error", corelogger.Error(err))
	}
}
