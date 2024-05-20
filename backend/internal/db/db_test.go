package db_test

import (
	"goRoadMap/backend/config"
	"goRoadMap/backend/internal/cache"
	"goRoadMap/backend/internal/db"
	"goRoadMap/backend/internal/jwtAuth"
	"goRoadMap/backend/pkg/services/logger"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestInitiateTables(t *testing.T) {
	config.Init()
	config := config.GetConfig()
	logger.Init(config.LoggerMode)
	godotenv.Load(config.EnvPath)
	cache.ClearCache()

	testUserData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	expectedUserData := map[string]string{
		"username": "testuser",
		"password": "testpass",
		"roleid":   "1",
	}

	jwtAuth.NewUserRegistration(testUserData)

	t.Run("PullDataFromCache Test", func(t *testing.T) {
		cacheMap, _ := cache.ReadCache("users")

		err := db.PullData("users", cacheMap)
		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
	})

	t.Run("FetchDataFromDB Test", func(t *testing.T) {
		userData, err := db.FetchUserData("testuser")

		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
		assert.Equalf(t, expectedUserData, userData, "Ожидали %v, получили %v", expectedUserData, userData)
	})
}
