package usecases

import (
	"api/dataproviders/main_db_provider"
	"api/utils/appconfig"
	"auth-testcase/library/loggerhelper"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTUseCase struct {
	logger          *loggerhelper.CustomLogger
	jwtRepo        	main_db_provider.IJWTRepository
	jwtConfig       appconfig.JWTConfiguration
	mailerUseCase	IMailUseCase
}

type IJWTUseCase interface {
	CreateJWTPair(guid string, ip string) (string, string, error)
	RefreshTokens(access string, refresh string, ip string) (string, string, error)
}

func NewJwtUseCase(logger *loggerhelper.CustomLogger,  jwtRepo main_db_provider.IJWTRepository, jwtConfig appconfig.JWTConfiguration, mailerUseCase IMailUseCase) IJWTUseCase {
	uc := JWTUseCase{
		logger:          logger,
		jwtRepo:         jwtRepo,
		jwtConfig:       jwtConfig,
		mailerUseCase: 	 mailerUseCase,
	}
	return &uc
}

func (uc *JWTUseCase) CreateJWTPair(guid string, ip string) (string, string, error) {

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": guid,
		"ip": ip,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": guid,
		"ip": ip,
		"exp": time.Now().Add(time.Hour * 24 * 5).Unix(),
		"iat": time.Now().Unix(),
	})

	accessToken, err := at.SignedString([]byte(uc.jwtConfig.JWTAccessSecret))

    if err != nil {
        return "", "", err
    }

	refreshToken, err := rt.SignedString([]byte(uc.jwtConfig.JWTRefreshSecret))

    if err != nil {
        return "", "", err
    }

	hashedAccessToken, err := bcrypt.GenerateFromPassword([]byte(strings.Split(accessToken, ".")[2][:71]), 14)

	if err != nil {
		uc.logger.SugarWithTracing().Info(err.Error())
        return "", "", err
    }
	
	_, err = uc.jwtRepo.CreateAccess(guid, string(hashedAccessToken), time.Now().Add(time.Hour * 24 * 5).Unix(), ip)

	if err != nil{
		return "", "", err
	}

	return accessToken, refreshToken, nil
	
}

func (uc *JWTUseCase) RefreshTokens(access string, refresh string, ip string) (string, string, error) {
	
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.jwtConfig.JWTAccessSecret), nil
    })

    if err != nil{
		uc.logger.SugarWithTracing().Error(err.Error())
		return "", "", err
    }

	claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid{
		return "", "", errors.New("invalid token")
	}
	guid := claims["guid"].(string)

	jwtData, err := uc.jwtRepo.GetDataByGuid(guid)

	if err != nil{
		return "", "", err
	}

	invalid := bcrypt.CompareHashAndPassword([]byte(jwtData.Hash), []byte(strings.Split(access, ".")[2][:71]))
	
	if invalid != nil {
		return "", "", errors.New("invalid access token")
	}

	if ip != jwtData.LastLoginIp{
		uc.mailerUseCase.SendNotificationMail(guid, ip)
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": guid,
		"ip": ip,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	accessToken, err := at.SignedString([]byte(uc.jwtConfig.JWTAccessSecret))

    if err != nil {
        return "", "", err
    }

	uc.jwtRepo.UpdateAccess(guid, time.Now().Add(time.Hour * 24 * 5).Unix(), ip)

	return accessToken, refresh, nil
}
