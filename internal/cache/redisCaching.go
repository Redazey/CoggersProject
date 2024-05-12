package cache

import (
	"goRoadMap/config"
	"goRoadMap/internal/errorz"
	"goRoadMap/pkg/services/logger"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")

	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis: ", zap.Error(err))
		return nil
	} else {
		logger.Info("Подключение к Redis установлено")
		return conn
	}
}

// запрашиваем поиск по ключу input, и что должно вернуться по ключу output
func IsDataInCache(input string, output string) (interface{}, error) {
	cacheMap, err := ReadCache(input)
	if cacheMap != nil && err == nil {
		return cacheMap[output], nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

func SaveCache(cacheMap map[string]string) error {
	config.Init()
	config := config.GetConfig()

	conn := redisConnect()

	defer conn.Close()

	if cacheMap == nil {
		return errorz.ErrNilCacheData
	} else {
		for key, value := range cacheMap {
			_, err := conn.Do("SET", key, value, "EX", config.Cache.EXTime)

			if err != nil {
				logger.Fatal("Ошибка при сохранении кэша в Redis: ", zap.Error(err))
				return err
			}
		}
	}

	logger.Info("Кэш успешно сохранён в Redis")
	return nil
}

// target - цель нашего поиска, * - если хотим полностью отсканировать кэш
func ReadCache(target string) (map[string]string, error) {
	cacheMap := make(map[string]string)
	conn := redisConnect()

	defer conn.Close()

	cursor := 0
	for cursor != 0 {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", target, "COUNT", 1000))
		if err != nil {
			logger.Error("Ошибка при сканировании Redis:", zap.Error(err))
			return nil, err
		}

		cursor, _ = redis.Int(values[0], nil)
		keys, _ := redis.Strings(values[1], nil)

		for _, key := range keys {
			val, _ := redis.String(conn.Do("GET", key))
			cacheMap[key] = val
		}
	}

	return cacheMap, nil
}
