package cache_test

import (
	"CoggersProject/pkg/cache"
	"CoggersProject/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var _, st = suite.New(t)

	err := cache.Init(st.Cfg.Redis.RedisAddr+":"+st.Cfg.Redis.RedisPort, st.Cfg.Redis.RedisUsername, st.Cfg.Redis.RedisPassword, 0, st.Cfg.Cache.EXTime)
	assert.Nil(t, err, "Expected no error during initialization")
}

func TestGetHashKey(t *testing.T) {
	testReq := map[string]interface{}{
		"key": "value",
	}
	expectedHash := "a7353f7cddce808de0032747a0b7be50"

	hashKey, err := cache.GetHashKey(testReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hashKey)
}

func TestIsExistInCache(t *testing.T) {
	hashKey := "mockHashKey"

	isExists, err := cache.IsExistInCache(hashKey)

	assert.NoError(t, err)
	assert.Equal(t, false, isExists)
}

func TestSaveCache(t *testing.T) {
	hashKey := "mockHashKey"
	mockResponse := "mockResponse"

	err := cache.SaveCache(hashKey, mockResponse)

	assert.NoError(t, err)
	// Дополнительные проверки, если нужны
}

// Продолжай в том же духе для остальных функций
