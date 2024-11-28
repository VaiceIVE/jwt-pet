package main

import (
	"api/utils/appconfig"
	"auth-testcase/library/loggerhelper"
)

func main() {
	logger := loggerhelper.NewCustomLogger()
	appConfig := appconfig.NewAppConfigurationFromEnvFile(".env")

	logger.SugarNoTracing().Infof("Run executor")

	MakeHttpServer(logger, appConfig)
}