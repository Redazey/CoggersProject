package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"goRoadMap/jwt_api"

	"github.com/rs/cors"
)

func handler(f func(data map[string]string) (bool, string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Чтение данных из POST запроса
			var message map[string]string
			err := json.NewDecoder(r.Body).Decode(&message)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			stringMsg, err := f(message)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			returnDataMap := map[string]interface{}{
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

func main() {
	mux := http.NewServeMux()
	keygen := handler(jwt_api.Keygen)
	tokenAuth := handler(jwt_api.TokenAuth)

	mux.HandleFunc("/keygen", keygen)
	mux.HandleFunc("/tokenAuth", tokenAuth)

	fmt.Println("Server is running on http://localhost:8080")

	corsHandler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", corsHandler)
}
