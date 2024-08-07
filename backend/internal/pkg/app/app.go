package app

import (
	"CoggersProject/config"
	"CoggersProject/internal/app/endpoint/grpcAuth"
	"CoggersProject/internal/app/endpoint/grpcServParser"
	"CoggersProject/internal/app/lib/cacher"
	"CoggersProject/internal/app/service/auth"
	"CoggersProject/internal/app/service/servParser"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/db"
	"CoggersProject/pkg/logger"
	pbAuth "CoggersProject/pkg/protos/auth"
	pbServParser "CoggersProject/pkg/protos/servParser"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	auth       *auth.Service
	servParser *servParser.Service

	server *grpc.Server
}

func New() (*App, error) {
	// инициализируем конфиг, логгер и кэш
	env, err := config.NewEnv()
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	cfg, err := config.NewConfig("./config/servers.json")
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	logger.Init(env.LoggerLevel)
	cacher.Init(env.Cache.UpdateInterval)

	a := &App{}

	a.server = grpc.NewServer()

	// инициализируем сервисы
	a.auth = auth.New(env)
	a.servParser = servParser.New(cfg)

	// инициализируем эндпоинты
	serviceAuth := &grpcAuth.Endpoint{
		Auth:      a.auth,
		SecretKey: env.JwtSecret,
	}

	serviceServParser := &grpcServParser.Endpoint{
		Parser: a.servParser,
	}

	pbAuth.RegisterAuthServiceServer(a.server, serviceAuth)
	pbServParser.RegisterServParserServiceServer(a.server, serviceServParser)

	err = cache.Init(env.Redis.RedisAddr+":"+env.Redis.RedisPort, env.Redis.RedisPassword, env.Redis.RedisDBId, env.Cache.EXTime)
	if err != nil {
		logger.Error("ошибка при инициализации кэша: ", zap.Error(err))
		return nil, err
	}

	err = db.Init(env.DB.DBUser, env.DB.DBPassword, env.DB.DBHost, env.DB.DBName)
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", ":8080")
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
