package parser_test

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/parser"
	"CoggersProject/backend/pkg/service/cacher"
	"CoggersProject/backend/pkg/service/logger"
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

	s := parser.New()

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)

	serverMap := make(map[string]string)

	for _, server := range config.Servers {
		serverMap[server.Name] = server.IP + ":25565"
	}

	t.Run("Parser Test", func(t *testing.T) {
		_, err := s.GetOnlineInfo(serverMap)

		assert.Nil(t, err, "ожидали Nil, получили ошибку")
	})
}
