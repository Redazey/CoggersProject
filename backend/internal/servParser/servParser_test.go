package servParser_test

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/servParser"
	"CoggersProject/backend/pkg/service/cacher"
	"CoggersProject/backend/pkg/service/logger"
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	config.Init()
	config := config.GetConfig()

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
	}

	s := servParser.New()

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)

	serverMap := make(map[string]string)

	for _, server := range config.Servers {
		serverMap[server.Name] = "https://api.mcsrvstat.us/3/" + server.IP
	}

	t.Run("Parser Test", func(t *testing.T) {
		expectedServInfo := make(map[string]map[string]interface{})
		ServersInfo, err := s.GetOnlineInfo(serverMap)

		logStr := fmt.Sprintf("Результат выполнения функции: %v", ServersInfo)
		logger.Info(logStr)
		assert.IsType(t, expectedServInfo, ServersInfo, "ожидали: %T, получили: %T", expectedServInfo, ServersInfo)
		assert.Nil(t, err, "ожидали Nil, получили ошибку")
	})
}
