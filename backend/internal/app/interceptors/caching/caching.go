package caching

import (
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/logger"
	"context"
	"reflect"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		hashKey, err := cache.GetHashKey(req)
		if err != nil {
			logger.Error("возникла ошибка при получении хэш-ключа", zap.Error(err))
		}

		responseStruct := reflect.TypeOf(handler)
		m := responseStruct.Out(0).(proto.Message)

		response, err := cache.ReadProtoCache(hashKey, m)
		if err != nil {
			logger.Error("ошибка при получении значения из кэша", zap.Error(err))
		}

		if response != nil {
			return response, nil // Возвращаем результат из кэша, если он есть
		}

		// Если в кэше нет значения, вызываем сам эндпоинт
		res, err := handler(ctx, req)
		if err == nil {
			cache.SaveProtoCache(hashKey, res.(proto.Message)) // Сохраняем результат в кэше
		}

		return res, err
	}
}
