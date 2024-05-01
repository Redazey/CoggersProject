package jwt_api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Handler(f func(data map[string]interface{}) (string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Чтение данных из POST запроса
			var message map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			returned_message, err := f(message)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			returnDataMap := map[string]interface{}{
				"message": returned_message,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(returnDataMap)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func Keygen(data map[string]interface{}) (string, error) {

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
		return "JWT token signature error", err
	}

	return tokenString, nil
}
