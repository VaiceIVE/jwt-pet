package appconfig

type AppConfiguration struct {
	database   DatabaseConfiguration
	httpServer HttpServerConfiguration
	jwt        JWTConfiguration
	mail 	   MailerConfig
}

type DatabaseConfiguration struct {
	Username 	 string
	Password 	 string
	Host     	 string
	DatabaseName string
}

type JWTConfiguration struct{
	JWTAccessSecret    string
    JWTRefreshSecret   string
}

type HttpServerConfiguration struct {
	Host string
	Port int
}

type MailerConfig struct{
	From string
	Password string
	SmtpHost string
	SmtpPort string
}