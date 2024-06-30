package suite

import (
	"CoggersProject/config"
	"CoggersProject/pkg/cache"
	"context"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg *config.Configuration
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

	// Инициализируем кеш
	cache.Init(
		cfg.Redis.RedisAddr+":"+cfg.Redis.RedisPort,
		cfg.Redis.RedisPassword,
		cfg.Redis.RedisDBId,
		cfg.Cache.EXTime,
	)

	return ctx, &Suite{
		T:   t,
		Cfg: cfg,
	}
}
