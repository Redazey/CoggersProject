package app

import (
	"CoggersProject/config"
	pbAuth "CoggersProject/gen/go/auth"
	"CoggersProject/internal/app/endpoint/grpcAuth"
	"CoggersProject/internal/app/lib/cacher"
	"CoggersProject/internal/app/service/auth"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/db"
	"CoggersProject/pkg/logger"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	auth *auth.Service

	server *grpc.Server
}

func New() (*App, error) {
	// инициализируем конфиг, логгер и кэш
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	logger.Init(cfg.LoggerLevel)
	cacher.Init(cfg.Cache.UpdateInterval)

	a := &App{}

	a.server = grpc.NewServer()

	// обьявляем сервисы
	a.auth = auth.New(cfg)

	// регистрируем эндпоинты
	serviceAuth := &grpcAuth.Endpoint{
		Auth: a.auth,
	}
	pbAuth.RegisterAuthServiceServer(a.server, serviceAuth)

	err = cache.Init(cfg.Redis.RedisAddr+":"+cfg.Redis.RedisPort, cfg.Redis.RedisUsername, cfg.Redis.RedisPassword, cfg.Redis.RedisDBId, cfg.Cache.EXTime)
	if err != nil {
		logger.Error("ошибка при инициализации кэша: ", zap.Error(err))
		return nil, err
	}

	err = db.Init(cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBName)
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Fatal("Ошибка при открытии listener: ", zap.Error(err))
	}

	err = a.server.Serve(lis)
	if err != nil {
		logger.Fatal("Ошибка при инициализации сервера: ", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Stop() {
	logger.Info("закрытие gRPC сервера")

	a.server.GracefulStop()
}
