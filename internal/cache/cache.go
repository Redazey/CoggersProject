package cache

import (
	"goRoadMap/internal/errorz"
	"goRoadMap/internal/logger"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func redisConnect() redis.Conn {
	// Создаем новый клиент Redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")

	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis: ", zap.Error(err))
		return nil
	} else {
		logger.Info("Подключение к Redis установлено")
		return conn
	}
}

func SaveCache(data map[string]string) error {
	conn := redisConnect()

	defer conn.Close()

	if data == nil {
		return errorz.NilCacheData
	} else {

		for key, value := range data {
			_, err := conn.Do("SET", key, value)
			if err != nil {
				logger.Fatal("Ошибка при сохранении кэша в Redis: ", zap.Error(err))
				return err
			}
		}
	}

	logger.Info("Кэш успешно сохранён в Redis")
	return nil
}

func ReadCache() map[string]string {
	cacheMap := make(map[string]string)
	conn := redisConnect()

	defer conn.Close()

	cursor := 0
	for cursor != 0 {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", "*", "COUNT", 1000))
		if err != nil {
			logger.Error("Ошибка при сканировании Redis:", zap.Error(err))
			break
		}

		cursor, _ = redis.Int(values[0], nil)
		keys, _ := redis.Strings(values[1], nil)

		for _, key := range keys {
			val, _ := redis.String(conn.Do("GET", key))
			cacheMap[key] = val
		}
	}

	return cacheMap
}
