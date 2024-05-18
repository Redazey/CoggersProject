package jwtAuth

import (
	"database/sql"
	"goRoadMap/internal/cache"
	"goRoadMap/internal/db"
	"goRoadMap/internal/errorz"
	"goRoadMap/pkg/services/logger"

	"go.uber.org/zap"
)

// передаем в эту функцию username и password
func GetLoginData(message map[string]string) (map[string]string, error) {
	userData, hashKey := cache.ConvertMap(message, "username", "password")

	cachePwd, err := cache.IsDataInCache("users", hashKey, "password")
	if err != nil {
		logger.Error("ошибка при поиске данных в кэше Redis: ", zap.Error(err))
		return nil, err
	}

	if cachePwd != nil && cachePwd == message["password"] {
		return Keygen(message)
	} else if cachePwd == nil {
		// написать функцию для получения юзера из бд здесь
		dbMap, err := db.FetchUserData(message["username"])
		if err != nil {
			return nil, err
		}
		if dbMap != nil && dbMap["password"] == message["password"] {
			// сохраняем залогиненого юзера в кэш
			cache.SaveCache("users", userData)

			// авторизуем его
			return Keygen(message)
		}
	}

	return nil, errorz.ErrUserNotFound
}

// передаем в эту функцию username и password
func NewUserRegistration(message map[string]string) (map[string]string, error) {
	message["roleid"] = "1"
	userData, hashKey := cache.ConvertMap(message, "username", "password")

	cachePwd, err := cache.IsDataInCache("users", hashKey, "password")
	if err != nil {
		return nil, err
	}

	// если пароль у юзера есть, значит и юзер существует
	if cachePwd != "" {
		dbMap, err := db.FetchUserData(message["username"])
		if err != sql.ErrNoRows && err != nil {
			return nil, err
		}

		if len(dbMap) != 0 {
			logger.Info("такой пользователь уже существует")
			return nil, errorz.ErrUserExists
		}
	}

	err = cache.SaveCache("users", userData)
	if err != nil {
		logger.Error("ошибка при регистрации пользователя: %s", zap.Error(err))
		return nil, err
	}

	return Keygen(message)
}
