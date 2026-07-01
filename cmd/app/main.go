package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	coreconfig "github.com/egotk/golang-advert-app/internal/core/config"
	coregrpc "github.com/egotk/golang-advert-app/internal/core/grpc"
	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	corepgxpool "github.com/egotk/golang-advert-app/internal/core/postgres/pool/pgx"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	adverthttp "github.com/egotk/golang-advert-app/internal/features/advert/controller/http"
	advertlocal "github.com/egotk/golang-advert-app/internal/features/advert/repo/local"
	advertpostgres "github.com/egotk/golang-advert-app/internal/features/advert/repo/postgres"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	categoryhttp "github.com/egotk/golang-advert-app/internal/features/category/controller"
	categorypostgres "github.com/egotk/golang-advert-app/internal/features/category/repo/postgres"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	usergrpc "github.com/egotk/golang-advert-app/internal/features/user/controller/grpc"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userpostgres "github.com/egotk/golang-advert-app/internal/features/user/repo/postgres"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
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

	logger.Debug("init jwt service")
	jwtConfig := corejwt.NewConfigMust()
	jwtService := corejwt.NewService(jwtConfig)

	logger.Debug("init HTTP server")
	httpServer := corehttp.NewServer(
		corehttp.NewConfigMust(),
		logger,
		corehttp.RequestID(),
		corehttp.Logger(logger),
	)

	logger.Debug("init gRPC server")
	grpcServer := coregrpc.NewServer(
		coregrpc.NewConfigMust(),
		logger,
		coregrpc.RequestID(),
		coregrpc.Logger(logger),
		coregrpc.ErrorHandler(),
		coregrpc.JWToken(
			jwtService,
			userpb.User_List_FullMethodName,
			userpb.User_GetByID_FullMethodName,
			userpb.User_Logout_FullMethodName,
		),
	)
	grpcRegistrar := grpcServer.GetRegistrar()

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

	logger.Debug("init feature: users")
	if err := corevalidator.RegisterValidations(userentity.Validations()...); err != nil {
		logger.Fatal("register user validations", zap.Error(err))
	}
	userRepo := userpostgres.New(pool)
	userUseCase := userusecase.New(userRepo, jwtService)
	userHTTPController := userhttp.New(userUseCase)
	apiVersionRouter.RegisterRoutes(userHTTPController.Routes(jwtService)...)

	userGRPCController := usergrpc.New(userUseCase, logger)
	userpb.RegisterUserServer(grpcRegistrar, userGRPCController)

	logger.Debug("init feature: adverts")
	advertRepo := advertpostgres.New(pool)
	advertStorage := advertlocal.New(config.Root)
	advertUseCase := advertusecase.New(advertRepo, advertStorage)
	advertHTTPController := adverthttp.New(advertUseCase)
	apiVersionRouter.RegisterRoutes(advertHTTPController.Routes(jwtService)...)

	logger.Debug("init feature: categories")
	categoryRepo := categorypostgres.New(pool)
	categoryUseCase := categoryusecase.New(categoryRepo)
	categoryHTTPController := categoryhttp.New(categoryUseCase)
	apiVersionRouter.RegisterRoutes(categoryHTTPController.Routes(jwtService)...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	go func() {
		if err := httpServer.Run(ctx); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	if err := grpcServer.Start(ctx); err != nil {
		logger.Error("gRPC server error", zap.Error(err))
	}
}
