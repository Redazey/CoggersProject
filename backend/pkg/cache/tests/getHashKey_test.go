package cache_test

import (
	pbAuth "CoggersProject/gen/go/auth"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/cache/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashKey(t *testing.T) {
	suite.New(t)

	testReq := pbAuth.LoginRequest{
		Email:    "test@mail.test",
		Password: "testpwd",
	}
	expectedHash := "f661e41cd98daee327c1d220f90c7397"

	hashKey, err := cache.GetHashKey(testReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hashKey)
}

/*
func TestIsExistInCache(t *testing.T) {
	testReq := pbAuth.LoginRequest{
		Email:    "test@mail.test",
		Password: "testpwd",
	}

	hashKey, err := cache.GetHashKey(testReq)

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
*/
