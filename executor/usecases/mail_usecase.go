package usecases

import (
	"api/utils/appconfig"
	"auth-testcase/library/loggerhelper"
	"net/smtp"

)

type MailUseCase struct {
	logger          *loggerhelper.CustomLogger
	mailerConfig    appconfig.MailerConfig
}

type IMailUseCase interface {
	SendNotificationMail(guid string, ip string) (error)
}

func NewMailUseCase(logger *loggerhelper.CustomLogger,  mailerConfig appconfig.MailerConfig) IMailUseCase {
	uc := MailUseCase{
		logger:          	logger,
		mailerConfig:       mailerConfig,
	}
	return &uc
}

func (uc *MailUseCase) SendNotificationMail(guid string, ip string) (error) {
	//Send email based on user GUID that some new ip had an access to account
	message := []byte("This is a test email message." + ip + " Had access to your profile")

	toMock := []string{"vinakich@gmail.com"}

	auth := smtp.PlainAuth("", uc.mailerConfig.From, uc.mailerConfig.Password, uc.mailerConfig.SmtpHost)

	err := smtp.SendMail(uc.mailerConfig.SmtpHost+":"+uc.mailerConfig.SmtpPort, auth, uc.mailerConfig.From, toMock, message)
	if err != nil {
		uc.logger.SugarNoTracing().Error(err.Error())

		return err
	}
	return nil
}