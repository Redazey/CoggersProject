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
	Env *config.Enviroment
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
	env, err := config.NewEnv("../../.env")
	if err != nil {
		t.Fatalf("ошибка при инициализации файла конфигурации: %s", err)
	}

	err = db.Init(env.DB.DBUser, env.DB.DBPassword, env.DB.DBHost, env.DB.DBName)
	if err != nil {
		t.Fatalf("ошибка при инициализации БД: %s", err)
	}

	// Основной родительский контекст
	ctx, cancelCtx := context.WithTimeout(context.Background(), env.GRPCTimeout)

	// Когда тесты пройдут, закрываем контекст
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Создаем кеш
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.Redis.RedisAddr + ":" + env.Redis.RedisPort,
		Password: env.Redis.RedisPassword,
		DB:       env.Redis.RedisDBId,
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
		Env:              env,
		Rdb:              rdb,
		Db:               db.Conn,
		AuthClient:       authClient,
		ServParserClient: servParserClient,
	}
}
