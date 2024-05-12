package cache

import (
	"encoding/json"
	"goRoadMap/config"
	"goRoadMap/internal/errorz"

	"github.com/gomodule/redigo/redis"
)

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		//logger.Fatal("Ошибка при подключении к Redis: ", zap.Error(err))
		return nil
	} else {
		//logger.Info("Подключение к Redis установлено")
		return conn
	}
}

// запрашиваем поиск по ключу input, и что должно вернуться по ключу output
func IsDataInCache(table string, input string, output string) (interface{}, error) {
	cacheMap, err := ReadCache(table, input)
	if cacheMap != nil && err == nil {
		return cacheMap[output], nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

/*
функция принимает map вида:

	map[string]map[string]string{
		"username": {
			"password": "exampass",
			"roleid":   "1",
		}
	}
*/
func SaveCache(cacheMap map[string]map[string]string, table string) error {
	config.Init()
	config := config.GetConfig()

	conn := redisConnect()

	defer conn.Close()

	if cacheMap == nil {
		return errorz.ErrNilCacheData
	} else {
		for key, value := range cacheMap {
			jsonMap, err := json.Marshal(value)
			if err != nil {
				//logger.Error("Ошибка при преобразовании строк в map: ", zap.Error(err))
				return err
			}
			// Устанавливаем значение в хэш-таблицу
			_, err = conn.Do("HSET", table, key, jsonMap)

			if err != nil {
				//logger.Fatal("Ошибка при сохранении кэша в Redis: ", zap.Error(err))
				return err
			}
			// Устанавливаем TTL для всей хэш-таблицы
			_, err = conn.Do("EXPIRE", table, config.Cache.EXTime)

			if err != nil {
				//logger.Fatal("Ошибка при сохранении кэша в Redis: ", zap.Error(err))
				return err
			}
		}
	}

	//logger.Info("Кэш успешно сохранён в Redis")
	return nil
}

/*
table - таблица в которой ищем, target - цель нашего поиска

функция возвращает map вида:

	map[string]string = {
		"password": "examplePass",
		"roleid":   "exampleRoleid",
	}
*/
func ReadCache(table string, target string) (map[string]string, error) {
	conn := redisConnect()
	defer conn.Close()

	dataBytes, err := redis.Bytes(conn.Do("HGET", table, target))
	if err != nil {
		//logger.Error("Ошибка при преобразовании строк в map: ", zap.Error(err))
		return nil, err
	}
	// Преобразуем данные []uint8 в map[string]string
	var result map[string]string
	err = json.Unmarshal(dataBytes, &result)
	if err != nil {
		//logger.Error("Ошибка при преобразовании строк в map: ", zap.Error(err))
		return nil, err
	}
	return result, nil
}
