package jwtAuth_test

import (
	"goRoadMap/internal/errorz"
	"goRoadMap/internal/jwtAuth"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func TestKeygen(t *testing.T) {
	err := godotenv.Load("Z:/files/goRoadMap/goRoadMap/.env")
	if err != nil {
		t.Errorf("Error loading .env file:  %v", err)
	}

	secret := os.Getenv("JWT_KEY")
	data := map[string]string{
		"username": "testuser",
	}

	// Вызов тестируемой функции
	tokenMap, err := jwtAuth.Keygen(data)

	// Проверка на ошибку
	if err != nil {
		t.Errorf("Не ожидали ошибку, получили: %v", err)
	}

	// Проверка наличия сообщения возвращаемого токена
	if _, ok := tokenMap["message"]; !ok {
		t.Error("Ожидали наличие сообщения в токене")
	}

	// Проверка на генерацию корректного JWT-токена
	tokenString := tokenMap["message"]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Errorf("Неверный секретный ключ: %v", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims)["username"]; !ok {
		t.Error("Ошибка в данных токена: отсутствует информация о пользователе")
	}

	exp := token.Claims.(jwt.MapClaims)["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now()) {
		t.Error("Срок действия токена истек")
	}
}

func TestTokenAuth(t *testing.T) {
	err := godotenv.Load("Z:/files/goRoadMap/goRoadMap/.env")
	if err != nil {
		t.Errorf("Error loading .env file:  %v", err)
	}

	secret := os.Getenv("JWT_KEY")

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
	_, err = jwtAuth.TokenAuth(data)

	// Проверяем наличие ошибок
	if err != nil {
		t.Errorf("Не ожидали ошибку, получили: %v", err)
	}

	// В данном тесте не проверяем содержимое returnDataMap, так как функция TokenAuth не модифицирует его

	// Проверяем корректность обработки истекшего токена
	claims := jwt.MapClaims{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Errorf("Неверный секретный ключ: %v", err)
	}

	_ = token.Claims.Valid()

	now := time.Now()
	claims["exp"] = now.Unix()

	tokenString, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	dataExpired := map[string]string{
		"token": tokenString,
	}

	_, errExpired := jwtAuth.TokenAuth(dataExpired)
	if errExpired != errorz.TokenExpired {
		t.Errorf("Ожидали ошибку TokenExpired, получили другую ошибку: %v", errExpired)
	}

}
