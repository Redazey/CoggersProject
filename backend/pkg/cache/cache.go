package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx         = context.Background()
	Rdb         *redis.Client
	CacheEXTime time.Duration
)

func Init(Addr string, Password string, DB int, CacheEx time.Duration) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})

	err := Rdb.Ping(Ctx).Err()
	if err != nil {
		return err
	}

	CacheEXTime = CacheEx

	return nil
}

/*
получение хэш ключа по реквесту
*/
func GetHashKey(request interface{}) (string, error) {
	var hashKey string

	reqJSON, err := json.Marshal(request)
	if err != nil {
		return hashKey, err
	}

	hash := md5.Sum(reqJSON)
	hashKey = hex.EncodeToString(hash[:])

	return hashKey, nil
}

/*
функция для проверки существования таблицы в кэше

принимает:

	table - имя таблицы

возвращает:

	bool - true, если таблица существует, иначе false
	error - ошибка, если возникла
*/
func IsExistInCache(hashKey string) (bool, error) {
	exists, err := Rdb.Exists(Ctx, hashKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

/*
функция для записи данных в кэш, принимает grpc requests
*/
func SaveCache(hashKey string, data interface{}) error {
	cacheData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	Rdb.HSet(Ctx, hashKey, cacheData).Err()

	// Устанавливаем время жизни ключа
	err = Rdb.Expire(Ctx, hashKey, time.Minute*time.Duration(CacheEXTime)).Err()
	if err != nil {
		return err
	}

	return nil
}

/*
Функция для чтения значений по хэш-ключу

возвращает grpc response
*/
func ReadCache(hashKey string) (interface{}, error) {
	response, err := Rdb.Get(Ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}

	return response, nil
}

/*
Функция для удаления значений по хэш-ключу
*/
func DeleteCache(hashKey string) error {
	// Удаляем хэш целиком
	err := Rdb.Del(Ctx, hashKey).Err()
	if err != nil {
		return err
	}
	return nil
}

/*
Функция для удаления значений по шаблону

пример pattern: news_category_*, где * - любое подстановочное значение
*/
func DeleteCacheByPattern(pattern string) error {
	var cursor uint64
	for {
		// Ищем ключи по шаблону
		keys, nextCursor, err := Rdb.Scan(Ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}

		// Удаляем найденные ключи
		if len(keys) > 0 {
			err = Rdb.Del(Ctx, keys...).Err()
			if err != nil {
				return err
			}
		}

		// Обновляем курсор
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

/*
Функция, которая удаляет все протухшие ключ-значения из выбранной таблицы

автоматически применяется при сохранении кэша при помощи функции SaveCache
*/
func DeleteEX(hashKey string) error {
	keys, err := Rdb.HKeys(Ctx, hashKey).Result()
	if err != nil {
		return err
	}

	// удаляем все протухшие ключи из Redis
	for _, key := range keys {
		// Получаем время до истечения срока действия ключа
		ttl := Rdb.TTL(Ctx, key).Val()

		if ttl <= 0 {
			// Если TTL < 0, значит ключ уже истек и можно его удалить
			err := Rdb.Del(Ctx, key).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/*
функция для стирания кэша

нужна в основном для дэбага
*/
func ClearCache(hashKey string) error {
	// Удаление всего кэша из Redis
	err := Rdb.Del(Ctx, hashKey).Err()
	if err != nil {
		return err
	}
	return nil
}
