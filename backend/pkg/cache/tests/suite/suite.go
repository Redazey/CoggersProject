package suite

import (
	"CoggersProject/config"
	"CoggersProject/pkg/cache"
	"context"
	"testing"
)

type Suite struct {
	*testing.T
	env *config.Enviroment
}

// New creates new test suite.
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов

	// Читаем конфиг из файла
	env, err := config.NewEnv("../../../.env")
	if err != nil {
		t.Fatalf("ошибка при инициализации файла конфигурации: %s", err)
	}

	// Основной родительский контекст
	ctx, cancelCtx := context.WithTimeout(context.Background(), env.GRPCTimeout)

	// Когда тесты пройдут, закрываем контекст
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Инициализируем кеш
	cache.Init(
		env.Redis.RedisAddr+":"+env.Redis.RedisPort,
		env.Redis.RedisPassword,
		env.Redis.RedisDBId,
		env.Cache.EXTime,
	)

	return ctx, &Suite{
		T:   t,
		env: env,
	}
}
