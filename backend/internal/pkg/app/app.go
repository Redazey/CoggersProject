package app

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/endpoints/auth"
	serverinfo "CoggersProject/backend/internal/endpoints/serverInfo"
	"CoggersProject/backend/internal/mw"
	"CoggersProject/backend/internal/servParser"
	"CoggersProject/backend/pkg/db"
	"CoggersProject/backend/pkg/jwtAuth"
	"CoggersProject/backend/pkg/service/cacher"
	"CoggersProject/backend/pkg/service/logger"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type App struct {
	auth       *auth.Endpoint
	serverinfo *serverinfo.Endpoint
	jwt        *jwtAuth.Service
	servParser *servParser.Service
	echo       *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.jwt = jwtAuth.New()
	a.servParser = servParser.New()

	a.auth = auth.New(a.jwt)
	a.serverinfo = serverinfo.New(a.servParser)

	// инициализируем конфиг, .env, логгер и кэш
	config.Init()
	config := config.GetConfig()

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
		return nil, err
	}

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)

	a.echo = echo.New()
	a.echo.Use(mw.Recovery)

	a.echo.GET("/UserLogin", a.auth.UserLogin)
	a.echo.GET("/NewUserRegistration", a.auth.NewUserRegistration)
	a.echo.GET("/ServerInfo", a.serverinfo.ServerInfo)

	err = db.InitiateTables()
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	err := a.echo.Start(":8080")
	if err != nil {
		logger.Fatal("Ошибка при инициализации сервера: ", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Stop() error {
	err := a.echo.Close()
	if err != nil {
		logger.Fatal("Ошибка при инициализации сервера: ", zap.Error(err))
		return err
	}

	return nil
}
