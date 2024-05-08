package main

import (
	"encoding/json"
	"net/http"

	"goRoadMap/internal/db"
	"goRoadMap/internal/logger"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

func handler(f func(data map[string]string) (map[string]string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Чтение данных из POST запроса
			var message map[string]string
			err := json.NewDecoder(r.Body).Decode(&message)

			if err != nil {
				logger.Error("ошибка при декодировании json файла: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			returnDataMap, err := f(message)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				logger.Error("ошибка при выполнении функции: ", zap.Error(err))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(returnDataMap)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	logger.Init("info")

	// заворачиваем функции в функцию-декоратор handler
	mux := http.NewServeMux()
	keygen := handler(jwt.Keygen)
	tokenAuth := handler(jwt.TokenAuth)

	mux.HandleFunc("/keygen", keygen)
	mux.HandleFunc("/tokenAuth", tokenAuth)
	mux.HandleFunc("/getLoginData", getLoginData)
	mux.HandleFunc("/newUserRegistration", newUserRegistration)

	logger.Info("Server is running on http://localhost:8080")

	err := db.InitiateTables()
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
	} else {
		logger.Info("Сервер запущен")
	}

	corsHandler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", corsHandler)
}