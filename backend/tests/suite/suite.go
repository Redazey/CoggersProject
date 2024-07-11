package suite

import (
	"CoggersProject/config"
	"CoggersProject/pkg/db"
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbAuth "CoggersProject/pkg/protos/auth"
	pbServParser "CoggersProject/pkg/protos/servParser"
)

type Suite struct {
	*testing.T
	Cfg *config.Configuration
	Rdb *redis.Client
	Db  *sqlx.DB

	AuthClient       pbAuth.AuthServiceClient
	ServParserClient pbServParser.ServParserServiceClient
}

// New creates new test suite.
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов

	// Читаем конфиг из файла
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("ошибка при инициализации файла конфигурации: %s", err)
	}

	err = db.Init(cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBName)
	if err != nil {
		t.Fatalf("ошибка при инициализации БД: %s", err)
	}

	// Основной родительский контекст
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPCTimeout)

	// Когда тесты пройдут, закрываем контекст
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Создаем кеш
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisAddr + ":" + cfg.Redis.RedisPort,
		Password: cfg.Redis.RedisPassword,
		DB:       cfg.Redis.RedisDBId,
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		t.Fatalf("redis connection failed: %v", err)
	}

	// Создаем клиент
	cc, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	// gRPC-клиент сервера Auth
	authClient := pbAuth.NewAuthServiceClient(cc)
	servParserClient := pbServParser.NewServParserServiceClient(cc)

	return ctx, &Suite{
		T:                t,
		Cfg:              cfg,
		Rdb:              rdb,
		Db:               db.Conn,
		AuthClient:       authClient,
		ServParserClient: servParserClient,
	}
}
