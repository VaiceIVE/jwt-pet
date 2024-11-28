package http_controllers

import (
	"api/utils/appconfig"
	"api/usecases"
	"net/http"
	"auth-testcase/library/loggerhelper"
	"github.com/labstack/echo/v4"
)

type JWTController struct {
	logger              *loggerhelper.CustomLogger
	echoServer          *echo.Echo
	appConfig           *appconfig.AppConfiguration
	jwtUseCase 			usecases.IJWTUseCase
}

type JWTResponse struct {
	AcessToken string `json:"accessToken" xml:"accessToken"`
	RefreshToken string `json:"refreshToken" xml:"refreshToken"`
}

type IJWTController interface {
	makeRoutes()
}

func NewJWTController(logger *loggerhelper.CustomLogger, echoServer *echo.Echo, jwtUseCase usecases.IJWTUseCase) IJWTController {
	controller := JWTController{
		logger:              logger,
		echoServer:          echoServer,
		jwtUseCase: 		 jwtUseCase,
	}
	controller.makeRoutes()
	return &controller
}

func (c *JWTController) makeRoutes() {
	v1 := c.echoServer.Group("/jwt")

	v1.GET("/access/:id", c.createJWT)
	v1.GET("/refresh", c.refresh)
}

func (c *JWTController) createJWT(ctx echo.Context) error {

	ip := ctx.RealIP()

	c.logger.SugarWithTracing().Infof("Create JWT")

	accessToken, refreshToken, err := c.jwtUseCase.CreateJWTPair(ctx.Param("id"), ip)

	if err != nil{
		c.logger.SugarWithTracing().Error(err.Error())
		return ctx.JSON(http.StatusInternalServerError, "Oops")
	}

	res := &JWTResponse{
		AcessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return ctx.JSON(http.StatusOK, res)	
}

func (c *JWTController) refresh(ctx echo.Context) error {

	c.logger.SugarWithTracing().Infof("Refresh")

	ip := ctx.RealIP()

	var access = ""
	var refresh = ""

	for name, data := range ctx.Request().Header{
		if name == "Authorization"{
			access = data[0]
		}
	}

	cookie, err := ctx.Cookie("Refresh")
	if err != nil {
		return err
	}
	refresh = cookie.Value
	if access == ""{
		return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	if refresh == ""{
		return ctx.JSON(http.StatusBadRequest, "Bad request")
	}

	accessToken, refreshToken, err := c.jwtUseCase.RefreshTokens(access, refresh, ip)

	if err != nil{
		return ctx.JSON(http.StatusBadRequest, "Bad request")
	}

	res := &JWTResponse{
		AcessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return ctx.JSON(http.StatusOK, res)	}
