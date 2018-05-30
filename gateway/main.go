package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"arctron.cn/arctron/arcplus/gateway/conf"
	"arctron.cn/arctron/arcplus/gateway/ctrl"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	ctrl.Logger = e.Logger

	defer ctrl.CloseConn()
	// Connect to micro services
	if err := ctrl.ConnSrv(); err != nil {
		e.Logger.Fatal("exits due to service connect error: %v", err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{Generator: func() string {
		id := uuid.New()
		return id.String()
	}}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "hello, world",
		})
	})

	// 验证码接口
	e.POST("/sendsms", ctrl.SendSms)
	e.POST("/checkcode", ctrl.CheckCode)

	// 用户部分接口
	e.POST("/regist", ctrl.Regist)

	// 电子名片部分接口
	e.POST("/expert", ctrl.Expert)

	// Restricted group
	r := e.Group("/auth")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "token",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
		SigningKey:    []byte(conf.SignKey),
	}))

	// 获取七牛云上传token
	r.POST("/mediatoken", ctrl.MediaToken)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
