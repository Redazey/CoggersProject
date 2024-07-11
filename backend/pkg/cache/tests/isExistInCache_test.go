package cache_test

import (
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/cache/tests/suite"
	pbAuth "CoggersProject/pkg/protos/auth"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistInCache(t *testing.T) {
	suite.New(t)

	testReq := &pbAuth.LoginRequest{
		Email:    "test@mail.test",
		Password: "testpwd",
	}

	hashKey := fmt.Sprintf("test_%s_%s", testReq.Email, testReq.Password)

	err := cache.SaveProtoCache(hashKey, testReq)
	assert.NoError(t, err, "ошибка при создании тестового кэша: %s", err)

	isExists, err := cache.IsExistInCache(hashKey)

	assert.NoError(t, err)
	assert.Equal(t, true, isExists)
}
