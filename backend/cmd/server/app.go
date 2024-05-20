package main

import (
	"log"
	"net/http"

	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/db"
	"goRoadMap/backend/internal/jwtAuth"
	"goRoadMap/backend/pkg/handler"
	"goRoadMap/backend/pkg/services/cacher"
	"goRoadMap/backend/pkg/services/logger"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	// инициализируем конфиг, .env, логгер и кэш
	config.Init()
	config := config.GetConfig()

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
		return
	}

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)

	// заворачиваем функции в функцию-декоратор handler
	mux := http.NewServeMux()
	keygen := handler.Handler(jwtAuth.Keygen)
	tokenAuth := handler.Handler(jwtAuth.TokenAuth)
	getLoginData := handler.Handler(jwtAuth.GetLoginData)
	newUserRegistration := handler.Handler(jwtAuth.NewUserRegistration)

	mux.HandleFunc("/keygen", keygen)
	mux.HandleFunc("/tokenAuth", tokenAuth)
	mux.HandleFunc("/getLoginData", getLoginData)
	mux.HandleFunc("/newUserRegistration", newUserRegistration)

	logger.Info("Server is running on http://localhost:8080")

	err = db.InitiateTables()
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
	} else {
		logger.Info("Сервер запущен")
	}

	corsHandler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", corsHandler)
}
