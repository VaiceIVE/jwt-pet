package main

import (
	httpcontroller "api/entrypoints/http_controllers"
	"api/usecases"
	"api/utils/appconfig"
	"fmt"
	"auth-testcase/library/loggerhelper"
	"api/dataproviders/main_db_provider"
	echoPrometheus "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func MakeHttpServer(logger *loggerhelper.CustomLogger, appConfig *appconfig.AppConfiguration) {
	echoServer := echo.New()
	setupEchoMiddlewares(echoServer, logger)

	_ = httpcontroller.NewJWTController(logger, echoServer,
		usecases.NewJwtUseCase(logger, main_db_provider.NewJWTRepository(logger, appConfig.GetDatabaseConfigForDbProvider()), appConfig.GetJWTConfig(), usecases.NewMailUseCase(logger, appConfig.GetMailerConfig())))

	zap.S().Infof("Server start HTTP Server %s:%d", appConfig.GetHttpServerConfig().Host, appConfig.GetHttpServerConfig().Port)
	err := echoServer.Start(fmt.Sprintf(":%d", appConfig.GetHttpServerConfig().Port))
	if err != nil {
		zap.S().Fatalf("HTTP Server error %v", err)
	}
}

func setupEchoMiddlewares(echoServer *echo.Echo, logger *loggerhelper.CustomLogger) {
	echoServer.Use(loggerhelper.EchoCustomLogger(logger))
	prometheus := echoPrometheus.NewPrometheus("http_app_server", nil)
	prometheus.Use(echoServer)
}