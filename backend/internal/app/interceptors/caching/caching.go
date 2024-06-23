package caching

import (
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/logger"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		hashKey, err := cache.GetHashKey(req)
		if err != nil {
			logger.Error("возникла ошибка при получении хэш-ключа", zap.Error(err))
		}

		response, err := cache.ReadCache(hashKey)
		if err != nil {
			logger.Error("ошибка при получении значения из кэша", zap.Error(err))
		}

		if response != nil {
			return response, nil // Возвращаем результат из кэша, если он есть
		}

		// Если в кэше нет значения, вызываем сам эндпоинт
		res, err := handler(ctx, req)
		if err == nil {
			cache.SaveCache(hashKey, res) // Сохраняем результат в кэше
		}

		return res, err
	}
}

func CacheMiddleware(next grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		hashKey, err := cache.GetHashKey(req)
		if err != nil {
			logger.Error("возникла ошибка при получении хэш-ключа", zap.Error(err))
		}

		response, err := cache.ReadCache(hashKey)
		if err != nil {
			logger.Error("ошибка при получении значения из кэша", zap.Error(err))
		}

		if response != nil {
			return response, nil // Возвращаем результат из кэша, если он есть
		}

		// Если в кэше нет значения, вызываем сам эндпоинт
		res, err := next(ctx, req)
		if err == nil {
			cache.SaveCache(hashKey, res) // Сохраняем результат в кэше
		}

		return res, err
	}
}
