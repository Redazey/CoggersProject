package app

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/endpoints/auth"
	"CoggersProject/backend/internal/mw"
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
	e    *auth.Endpoint
	s    *jwtAuth.Service
	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.s = jwtAuth.New()

	a.e = auth.New(a.s)

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

	a.echo.GET("/UserLogin", a.e.UserLogin)
	a.echo.GET("/NewUserRegistration", a.e.NewUserRegistration)

	err = db.InitiateTables()
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
		return nil, err
	}
	logger.Error("zxc")

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
