package handler

import (
	"encoding/json"
	"goRoadMap/pkg/services/logger"
	"net/http"

	"go.uber.org/zap"
)

func Handler(f func(data map[string]string) (map[string]string, error)) func(w http.ResponseWriter, r *http.Request) {
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
