package cache_test

import (
	pbAuth "CoggersProject/gen/go/auth"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/cache/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistInCache(t *testing.T) {
	_, st := suite.New(t)
	cache.Rdb = st.Rdb

	testReq := pbAuth.LoginRequest{
		Email:    "test@mail.test",
		Password: "testpwd",
	}

	hashKey, err := cache.GetHashKey(testReq)

	err = cache.SaveCache(hashKey, testReq)
	assert.NoError(t, err, "ошибка при создании тестового кэша: %s", err)

	isExists, err := cache.IsExistInCache(hashKey)

	assert.NoError(t, err)
	assert.Equal(t, false, isExists)
}
