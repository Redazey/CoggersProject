package suite

import (
	"CoggersProject/internal/app/config"
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

type Suite struct {
	*testing.T
	Cfg *config.Configuration
	Rdb *redis.Client
}

// New creates new test suite.
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов

	// Читаем конфиг из файла
	cfg, err := config.NewConfig("../../../.env")
	if err != nil {
		t.Fatalf("ошибка при инициализации файла конфигурации: %s", err)
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

	return ctx, &Suite{
		T:   t,
		Cfg: cfg,
		Rdb: rdb,
	}
}
