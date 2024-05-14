package jwtAuth

import (
	"goRoadMap/internal/cache"
	"goRoadMap/internal/errorz"
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
		return Keygen(message)
	} else if cachePwd == nil {
		// написать функцию для получения юзера из бд здесь
	}

	return nil, errorz.ErrUserNotFound
}

// передаем в эту функцию username и password
func NewUserRegistration(message map[string]string) (map[string]string, error) {
	cachePwd, err := cache.IsDataInCache("users", message["username"], "password")
	if err != nil {
		return nil, err
	}
	if cachePwd != nil {
		// написать функцию для получения юзера из бд здесь
	}

	err = cache.SaveCache("users", cache.ConvertMap(message))
	if err != nil {
		logger.Error("ошибка при регистрации нового пользователя: ", zap.Error(err))
		return nil, err
	}

	return Keygen(message)
}
