package app

import (
	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/endpoints/auth"
	"goRoadMap/backend/pkg/db"
	"goRoadMap/backend/pkg/jwtAuth"
	"goRoadMap/backend/pkg/services/cacher"
	"goRoadMap/backend/pkg/services/logger"
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

	// здесь используем MW (логгер, рековер, т.д.)
	// a.echo.Use()

	a.echo.GET("/tokenAuth", a.e.GetLoginData)
	a.echo.GET("/newUserRegistration", a.e.NewUserRegistration)

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
	} else {
		logger.Info("Server is running on http://localhost:8080")
	}

	return nil
}
