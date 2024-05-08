package jwtAuth

import (
	"os"
	"time"

	"github.com/joho/godotenv"

	"goRoadMap/internal/errorz"
	"goRoadMap/internal/logger"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func Keygen(data map[string]string) (map[string]string, error) {
	err := godotenv.Load("Z:/files/goRoadMap/goRoadMap/.env")
	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
		return nil, err
	}

	returnDataMap := make(map[string]string)
	username := data["username"]

	// Создаем новый JWT токен
	token := jwt.New(jwt.SigningMethodHS256)

	// Указываем параметры для токена
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	secretKey := []byte(os.Getenv("JWT_KEY"))

	// Подписываем токен с помощью указанного секретного ключа
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		logger.Fatal("ошибка при подписи токена: ", zap.Error(err))
		return nil, err
	} else {
		returnDataMap["message"] = tokenString
	}

	return returnDataMap, nil
}

func TokenAuth(data map[string]string) (map[string]string, error) {
	returnDataMap := make(map[string]string)
	tokenString := data["token"]
	secretKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		logger.Fatal("ошибка при чтении токена: ", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

		if time.Now().After(expirationTime) {
			return nil, errorz.TokenExpired
		}
	} else {
		logger.Fatal("невалидный токен: ", zap.Error(jwt.ValidationError{}))
		return nil, jwt.ValidationError{}
	}

	return returnDataMap, nil
}