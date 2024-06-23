package auth

import (
	"CoggersProject/config"
	"CoggersProject/internal/app/errorz"
	"CoggersProject/internal/app/lib/db"
	"CoggersProject/internal/app/lib/jwt"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/logger"
	"database/sql"

	"go.uber.org/zap"
)

type Service struct {
	Cfg *config.Configuration
}

func New(cfg *config.Configuration) *Service {
	return &Service{
		Cfg: cfg,
	}
}

func (s *Service) UserLogin(email string, password string) (string, error) {
	msg := map[string]string{
		"email":    email,
		"password": password,
	}

	hashKey, err := cache.GetHashKey(msg)
	if err != nil {
		logger.Error("ошибка при генерации хэш-ключа: ", zap.Error(err))
		return "", err
	}

	cacheData, err := cache.ReadCache(hashKey)
	if err != nil {
		logger.Error("ошибка при поиске данных в кэше Redis: ", zap.Error(err))
		return "", err
	}

	cachePwd, ok := cacheData.(map[string]string)["password"]
	if !ok {
		return "", errorz.ErrTypeConversion
	}

	if cachePwd != "" && cachePwd != password {
		return "", errorz.ErrUserNotFound
	} else if cachePwd == "" {
		dbData, err := db.FetchUserData(email)
		if err != nil {
			return "", err
		}

		if dbData.Password != password {
			return "", err
		}

		// сохраняем залогиненого юзера в кэш
		err = cache.SaveCache(hashKey, dbData)
		if err != nil {
			return "", err
		}
	}

	// генерируем jwt токен и данных юзера для использования в дальнейшем
	key, err := jwt.Keygen(email, password, s.Cfg.JwtSecret)
	if err != nil {
		logger.Error("ошибка при генерации токена: ", zap.Error(err))
		return "", err
	}

	// авторизуем его
	return key, nil
}

func (s *Service) NewUserRegistration(email string, password string) (string, error) {
	msg := map[string]string{
		"email":    email,
		"password": password,
	}

	hashKey, err := cache.GetHashKey(msg)
	if err != nil {
		logger.Error("ошибка при генерации хэш-ключа: ", zap.Error(err))
		return "", err
	}

	exists, err := cache.IsExistInCache(hashKey)
	if err != nil {
		logger.Error("ошибка при поиске значения в кэше: ", zap.Error(err))
		return "", err
	}

	if exists {
		return "", errorz.ErrUserExists
	}

	dbMap, err := db.FetchUserData(email)
	if err != sql.ErrNoRows && err != nil {
		return "", err
	}

	err = cache.SaveCache(hashKey, dbMap)
	if err != nil {
		logger.Error("ошибка при регистрации пользователя: ", zap.Error(err))
		return "", err
	}

	key, err := jwt.Keygen(email, password, s.Cfg.JwtSecret)
	if err != nil {
		logger.Error("ошибка при генерации токена: ", zap.Error(err))
		return "", err
	}

	return key, nil
}

func (s *Service) IsAdmin(tokenString string) (bool, error) {
	userData, err := jwt.UserDataFromJwt(tokenString, s.Cfg.JwtSecret)
	if err != nil {
		logger.Error("ошибка при распаковки jwt-токена: ", zap.Error(err))
		return false, err
	}

	hashKey, err := cache.GetHashKey(userData)
	if err != nil {
		logger.Error("ошибка при генерации хэш-ключа: ", zap.Error(err))
		return false, err
	}

	cacheUserData, err := cache.ReadCache(hashKey)
	if err != nil {
		logger.Error("ошибка при поиске данных в кэше Redis: ", zap.Error(err))
		return false, err
	}

	userDataMap, ok := cacheUserData.(map[string]string)
	if !ok {
		return false, errorz.ErrTypeConversion
	}
	roleId := userDataMap["roleId"]

	if roleId != "" {
		if roleId == "1" {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
}
