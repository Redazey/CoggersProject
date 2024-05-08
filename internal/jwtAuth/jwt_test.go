package jwtAuth_test

import (
	"testing"
	"time"

	"goRoadMap/internal/jwtAuth"

	"github.com/golang-jwt/jwt"
)

func TestKeygen(t *testing.T) {
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
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Здесь необходимо указать правильный секретный ключ
	})

	if _, ok := token.Claims.(jwt.MapClaims)["username"]; !ok {
		t.Error("Ошибка в данных токена: отсутствует информация о пользователе")
	}

	exp := token.Claims.(jwt.MapClaims)["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now()) {
		t.Error("Срок действия токена истек")
	}
}
