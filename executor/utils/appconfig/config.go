package appconfig

import (
	"auth-testcase/library/viperr"
	"github.com/spf13/viper"
	"api/dataproviders/main_db_provider"
)

type IAppConfiguration interface {
	GetDatabaseConfig() DatabaseConfiguration
	GetHttpServerConfig() HttpServerConfiguration
}

func NewAppConfigurationFromEnvFile(pathToEnvFile string) *AppConfiguration {
	viperhelper.ReadFromEnv(pathToEnvFile)
	appConfig := AppConfiguration{}
	appConfig.database = DatabaseConfiguration{
		Host:     viper.GetString("DATABASE_HOST"),
		Username: viper.GetString("DATABASE_USER"),
		Password: viper.GetString("DATABASE_PASSWORD"),
		DatabaseName: viper.GetString("DATABASE_NAME"),
	}
	appConfig.httpServer = HttpServerConfiguration{
		Host: viper.GetString("SERVER_HTTP_HOST"),
		Port: viper.GetInt("SERVER_HTTP_PORT"),
	}
	appConfig.jwt = JWTConfiguration{
		JWTAccessSecret: viper.GetString("ACCESS_SECRET"),
		JWTRefreshSecret: viper.GetString("REFRESH_SECRET"),
	}
	appConfig.mail = MailerConfig{
		From: viper.GetString("MAILER_FROM"),
		Password: viper.GetString("MAILER_PASSWORD"),
		SmtpHost: viper.GetString("MAILER_SMTP_HOST"),
		SmtpPort: viper.GetString("MAILER_SMTP_PORT"),
	}
	return &appConfig
}

func (c *AppConfiguration) GetDatabaseConfig() DatabaseConfiguration {
	return c.database
}
func (c *AppConfiguration) GetDatabaseConfigForDbProvider() main_db_provider.DatabaseConfiguration {
	return main_db_provider.DatabaseConfiguration{
		Username: c.database.Username,
		Password: c.database.Password,
		Host:     c.database.Host,
		DatabaseName: c.database.DatabaseName,
		PoolSize: 10,
	}
}

func (c *AppConfiguration) GetHttpServerConfig() HttpServerConfiguration {
	return c.httpServer
}

func (c *AppConfiguration) GetJWTConfig() JWTConfiguration {
	return c.jwt
}

func (c *AppConfiguration) GetMailerConfig() MailerConfig {
	return c.mail
}