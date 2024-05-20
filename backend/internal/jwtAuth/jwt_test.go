package jwtAuth_test

import (
	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/errorz"
	"goRoadMap/backend/internal/jwtAuth"
	"goRoadMap/backend/pkg/services/logger"

	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	config.Init()
	config := config.GetConfig()
	logger.Init(config.LoggerMode)
	godotenv.Load(config.EnvPath)

	err := godotenv.Load(config.EnvPath)
	assert.Nil(t, err, "Ошибка при открытии env файла: %v", err)

	secret := os.Getenv("JWT_KEY")
	data := map[string]string{
		"username": "testuser",
	}

	t.Run("GenerateToken", func(t *testing.T) {
		// Вызов тестируемой функции
		tokenMap, err := jwtAuth.Keygen(data)
		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)

		// Проверка наличия сообщения возвращаемого токена
		assert.NotNil(t, tokenMap["message"], "Ожидаем наличие сообщения в токене")

		// Проверка на генерацию корректного JWT-токена
		tokenString := tokenMap["message"]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.Nil(t, err, "Неверный секретный ключ: %v", err)

		assert.NotNil(t, token.Claims.(jwt.MapClaims)["username"],
			"Ошибка в данных токена: отсутствует информация о пользователе")

		exp := token.Claims.(jwt.MapClaims)["exp"].(float64)
		expTime := time.Unix(int64(exp), 0)
		assert.False(t, expTime.Before(time.Now()), "Срок действия токена истек")
	})
}

func TestTokenAuth(t *testing.T) {
	secret := os.Getenv("JWT_KEY")

	t.Run("TokenAuthentication", func(t *testing.T) {
		// Генерируем JWT-токен для теста
		tokenData := map[string]string{
			"username": "testuser",
		}
		tokenMap, _ := jwtAuth.Keygen(tokenData)
		tokenString := tokenMap["message"]

		// Формируем входные данные для функции TokenAuth
		data := map[string]string{
			"token": tokenString,
		}

		// Вызываем тестируемую функцию
		_, err := jwtAuth.TokenAuth(data)
		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)

		// В данном тесте не проверяем содержимое returnDataMap, так как функция TokenAuth не модифицирует его

		// Проверяем корректность обработки истекшего токена
		claims := jwt.MapClaims{}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.Nil(t, err, "Неверный секретный ключ: %v", err)

		_ = token.Claims.Valid()

		now := time.Now()
		claims["exp"] = now.Unix()

		tokenString, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		dataExpired := map[string]string{
			"token": tokenString,
		}

		_, errExpired := jwtAuth.TokenAuth(dataExpired)
		assert.EqualError(t, errExpired, errorz.ErrTokenExpired.Error(),
			"Ожидаем ошибку TokenExpired, получили другую ошибку: %v", errExpired)
	})
}
