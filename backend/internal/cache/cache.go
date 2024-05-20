package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/errorz"
	"goRoadMap/backend/pkg/services/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

/*
функция для подготовки мапы к записи в кэш
columns - поля, по которым будет сгенерирован md5hash

формат входящей мапы:

	map[string]string = {
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

формат выходящей мапы:

	map[string]map[string]interface{} = {
		"md5hash": {
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}
	+ md5hash key
*/
func ConvertMap(inputMap map[string]string, columns ...string) (map[string]map[string]interface{}, string) {
	var mainKey string

	hash := md5.Sum([]byte(strings.Join(columns, "")))
	mainKey = hex.EncodeToString(hash[:])

	outputMap := make(map[string]map[string]interface{})
	outputMap[mainKey] = make(map[string]interface{})

	for key, value := range inputMap {
		outputMap[mainKey][key] = value
	}

	return outputMap, mainKey
}

func redisConnect() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := conn.Ping(ctx).Err()
	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis: ", zap.Error(err))
	}
	return conn
}

/*
поиск данных в кэше по md5 хэш-ключу

запрашиваем поиск по ключу input, и что должно вернуться по ключу output
*/
func IsDataInCache(table string, input string, output string) (interface{}, error) {
	cacheMap, err := ReadCache(table)
	if cacheMap[input] != nil && err == nil {
		return cacheMap[input][output], nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

/*
функция для записи данных в кэш, принимает мапы, после конвертации функцией ConvertMap

Вид входящей мапы:

	map[string]map[string]interface{}{
		"md5hash": {
			"username": "exampleUser",
			"password": "examplePass",
			"roleid":   "exampleRoleid",
		},
	}
*/
func SaveCache(table string, cacheMap map[string]map[string]interface{}) error {
	config.Init()
	config := config.GetConfig()

	conn := redisConnect()
	defer conn.Close()

	if cacheMap == nil {
		return errorz.ErrNilCacheData
	}

	for key, args := range cacheMap {
		// Устанавливаем значение в хэш-таблицу
		jsonMap, err := json.Marshal(args)
		if err != nil {
			logger.Error("Ошибка при преобразовании кэша в json: ", zap.Error(err))
			return err
		}

		err = conn.HSet(ctx, table, key, jsonMap).Err()
		if err != nil {
			logger.Error("Ошибка при сохранении кэша в Redis: ", zap.Error(err))
			return err
		}

		// Устанавливаем время жизни ключа
		err = conn.Expire(ctx, key, time.Minute*time.Duration(config.Cache.EXTime)).Err()
		if err != nil {
			logger.Error("Ошибка при установки срока жизни кэша: ", zap.Error(err))
			return err
		}
	}

	logger.Info("Кэш успешно сохранён в Redis")
	// удаляем устаревшие данные
	err := DeleteEX(table)
	if err != nil {
		logger.Error("Ошибка при удалении устаревшего кэша: ", zap.Error(err))
		return err
	}
	return nil
}

/*
Функция для чтения значений по хэш-ключу

возвращает map вида:

	map[string]map[string]interface{} = {
		"md5hash": {
			"username": "exampleUser",
			"password": "examplePass",
			"roleid":   "exampleRoleid",
		},
	}
*/
func ReadCache(table string) (map[string]map[string]interface{}, error) {
	conn := redisConnect()
	defer conn.Close()

	cacheMap := make(map[string]map[string]interface{})

	// Получаем все поля и значения из хэша
	result, err := conn.HGetAll(ctx, table).Result()
	if err != nil {
		logger.Error("Ошибка при извлечении кэша из Redis: ", zap.Error(err))
		return nil, err
	}

	// Преобразуем результат в map[string]interface{}
	for key, value := range result {
		var tempMap map[string]interface{}
		err := json.Unmarshal([]byte(value), &tempMap)
		if err != nil {
			logger.Error("Ошибка при декодировании кэша из JSON: ", zap.Error(err))
			return nil, err
		}
		cacheMap[key] = tempMap
	}

	return cacheMap, nil
}

/*
Функция, которая удаляет все протухшие ключ-значения из выбранной таблицы

автоматически применяется при сохранении кэша при помощи функции SaveCache
*/
func DeleteEX(table string) error {
	conn := redisConnect()
	defer conn.Close()

	keys, err := conn.HKeys(ctx, table).Result()
	if err != nil {
		logger.Error("ошибка при получении ключей из Redis: ", zap.Error(err))
		return err
	}

	// удаляем все протухшие ключи из Redis
	for _, key := range keys {
		// Получаем время до истечения срока действия ключа
		ttl := conn.TTL(ctx, key).Val()

		if ttl <= 0 {
			// Если TTL < 0, значит ключ уже истек и можно его удалить
			err := conn.Del(ctx, key).Err()
			if err != nil {
				logger.Error("Ошибка при удалении протухших значений из Redis: ", zap.Error(err))
				return err
			}
		}
	}

	logger.Info("Протухшие ключи удалены из кэша")
	return nil
}

/*
функция для стирания кэша

нужна в основном для дэбага
*/
func ClearCache() {
	conn := redisConnect()
	defer conn.Close()

	// Удаление всего кэша из Redis
	err := conn.FlushAll(ctx).Err()
	if err != nil {
		logger.Error("Ошибка при удалении кэша из Redis:", zap.Error(err))
		return
	}

	logger.Info("Протухшие ключи удалены из кэша")
}
