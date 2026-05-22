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
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	corepgxpool "github.com/egotk/golang-advert-app/internal/core/postgres/pool/pgx"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/http"
	userpostgres "github.com/egotk/golang-advert-app/internal/features/user/repo/postgres"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"go.uber.org/zap"
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

	logger.Debug("app time zone", zap.Any("zone", time.Local))

	logger.Debug("init HTTP server")
	httpServer := corehttp.NewServer(
		corehttp.NewConfigMust(),
		logger,
		corehttp.RequestID(),
		corehttp.Logger(logger),
	)

	logger.Debug("init postgres connection pool")
	pool, err := corepgxpool.NewPool(
		ctx,
		corepgxpool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	apiVersionRouter := corehttp.NewAPIVersionRouter(corehttp.ApiV1)
	userRepo := userpostgres.New(pool)
	userUseCase := userusecase.New(userRepo)
	userHTTPController := userhttp.New(userUseCase)
	apiVersionRouter.RegisterRoutes(userHTTPController.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server error", zap.Error(err))
	}
}
