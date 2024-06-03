package serverinfo_test

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/pkg/service/cacher"
	"CoggersProject/backend/pkg/service/logger"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func TestServerInfo(t *testing.T) {
	config.Init()
	config := config.GetConfig()

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
	}

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)
	client := &http.Client{}

	t.Run("ServerInfo Test", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/ServerInfo", nil)
		if err != nil {
			logger.Error("Ошибка при отправке запроса: ", zap.Error(err))
			return
		}

		res, err := client.Do(req)
		if err != nil {
			logger.Error("Ошибка при отправке запроса: ", zap.Error(err))
			return
		}

		receivedData, _ := io.ReadAll(res.Body)

		fmt.Printf("Получены данные от сервера: %s", receivedData)
	})
}
