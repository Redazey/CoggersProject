package jwt

import (
	"time"

	"CoggersProject/internal/app/errorz"

	"github.com/golang-jwt/jwt"
)

func Keygen(email string, pwd string, secretKey string) (string, error) {
	// Создаем новый JWT токен
	token := jwt.New(jwt.SigningMethodHS256)

	// Указываем параметры для токена
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["password"] = pwd
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Подписываем токен с помощью указанного секретного ключа
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TokenAuth(tokenString string, secretKey string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

		if time.Now().After(expirationTime) {
			return false, errorz.ErrTokenExpired
		}

	} else {
		return false, errorz.ErrValidation
	}

	return true, nil
}
