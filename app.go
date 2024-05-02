package main

import (
	"fmt"
	"net/http"

	"goRoadMap/jwt_api"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	keygen := jwt_api.Handler(jwt_api.Keygen)
	tokenAuth := jwt_api.Handler(jwt_api.TokenAuth)

	mux.HandleFunc("/keygen", keygen)
	mux.HandleFunc("/tokenAuth", tokenAuth)

	fmt.Println("Server is running on http://localhost:8080")

	corsHandler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", corsHandler)
}