package jwt_api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Handler(f func(data map[string]string) (bool, string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Чтение данных из POST запроса
			var message map[string]string
			err := json.NewDecoder(r.Body).Decode(&message)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			boolMsg, stringMsg, err := f(message)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			returnDataMap := map[string]interface{}{
				"status":  boolMsg,
				"message": stringMsg,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(returnDataMap)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func Keygen(data map[string]string) (bool, string, error) {

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
		return false, "JWT token signature error", err
	}

	return true, tokenString, nil
}

func TokenAuth(data map[string]string) (bool, string, error) {
	tokenString := data["token"]
	secretKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, "Error was occured while parsing jwtToken", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

		if time.Now().After(expirationTime) {
			return false, "Error was occured while authentification process", errors.New("token expired")
		}
	} else {
		return false, "Error was occured while validation process", jwt.ValidationError{}
	}

	return true, "Validation successfull", nil
}
