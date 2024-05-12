package db

import (
	"goRoadMap/internal/cache"
	"goRoadMap/internal/errorz"
	"goRoadMap/internal/jwtAuth"
	"goRoadMap/pkg/services/logger"

	"go.uber.org/zap"
)

// передаем в эту функцию username и password
func GetLoginData(message map[string]string) (map[string]string, error) {
	username := message["username"]
	password := message["password"]

	cachePwd, err := cache.IsDataInCache("users", username, "password")
	if err != nil {
		logger.Error("ошибка при поиске данных в кэше Redis: ", zap.Error(err))
		return nil, err
	}

	if cachePwd != nil && cachePwd == password {
		return jwtAuth.Keygen(message)
	} else if cachePwd != nil && cachePwd != password {
		return nil, errorz.ErrUserNotFound
	}

	return nil, errorz.ErrUserNotFound
}

// передаем в эту функцию username и password
func NewUserRegistration(message map[string]string) (map[string]string, error) {
	cacheData, err := cache.ReadCache("users", message["username"])
	if err != nil {
		return nil, err
	}
	if cacheData != nil {
		return nil, errorz.ErrUserExists
	}

	formatedMessage := map[string]map[string]string{
		message["username"]: {
			"password": message["password"],
			"roleid":   "1",
		},
	}

	err = cache.SaveCache(formatedMessage, "users")
	if err != nil {
		logger.Error("ошибка при регистрации нового пользователя: ", zap.Error(err))
		return nil, err
	}

	return jwtAuth.Keygen(message)
}
