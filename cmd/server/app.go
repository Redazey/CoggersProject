package main

import (
	"log"
	"net/http"

	"goRoadMap/config"
	"goRoadMap/internal/db"
	"goRoadMap/internal/jwtAuth"
	"goRoadMap/pkg/handler"
	"goRoadMap/pkg/services/cacher"
	"goRoadMap/pkg/services/logger"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	// инициализируем конфиг, кэш, .env и логгер
	config.Init()
	config := config.GetConfig()

	cacher.Init(config.Cache.UpdateInterval)

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
		return
	}

	logger.Init(config.LoggerMode)

	// заворачиваем функции в функцию-декоратор handler
	mux := http.NewServeMux()
	keygen := handler.Handler(jwtAuth.Keygen)
	tokenAuth := handler.Handler(jwtAuth.TokenAuth)
	getLoginData := handler.Handler(db.GetLoginData)
	newUserRegistration := handler.Handler(db.NewUserRegistration)

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
