package auth_test

import (
	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/errorz"
	"goRoadMap/backend/pkg/cache"
	"goRoadMap/backend/pkg/db"
	"goRoadMap/backend/pkg/jwtAuth"
	"goRoadMap/backend/pkg/services/logger"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	config.Init()
	config := config.GetConfig()
	logger.Init(config.LoggerMode)
	godotenv.Load(config.EnvPath)
	db.InitiateTables()

	testUserData := map[string]string{
		"username": "testuserreg",
		"password": "testpassreg",
	}
	_, hashKey := cache.ConvertMap(testUserData, "username", "password")

	// регистрируем пользователя в БД
	t.Run("NewUserRegistration Test", func(t *testing.T) {

		// работу keygen мы уже проверяем в другом тесте, так что здесь ошибку мы не берем
		expectedKey, _ := jwtAuth.Keygen(testUserData)
		JWTkey, err := jwtAuth.NewUserRegistration(testUserData)
		if err == nil {
			assert.Equalf(t, expectedKey, JWTkey, "Ожидалось получить %v, получили: %v", expectedKey, JWTkey)
		} else if err == errorz.ErrUserExists {
			logger.Info("Пользователь уже существует")
		} else {
			assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
		}
		// проверяем, сохранила ли функция нашего нового пользователя в кэш
		cachePwd, _ := cache.IsDataInCache("users", hashKey, "password")
		assert.Equalf(t, testUserData["password"], cachePwd, "Ожидалось получить testpass, получили: %v", cachePwd)
	})

	// теперь пробуем найти в БД нашего зарегистрированного пользователя и авторизовать его
	t.Run("GetLoginData Test", func(t *testing.T) {
		// работу keygen мы уже проверяем в другом тесте, так что здесь ошибку мы не берем
		expectedKey, _ := jwtAuth.Keygen(testUserData)

		JWTkey, err := jwtAuth.GetLoginData(testUserData)
		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
		assert.Equalf(t, expectedKey, JWTkey, "Ожидалось получить %v, получили: %v", expectedKey, JWTkey)
	})
}
