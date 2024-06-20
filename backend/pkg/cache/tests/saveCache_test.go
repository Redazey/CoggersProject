package cache_test

import (
	pbAuth "CoggersProject/gen/go/auth"
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/cache/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveCache(t *testing.T) {
	_, st := suite.New(t)
	cache.Rdb = st.Rdb

	testResp := pbAuth.AuthResponse{
		Key: "testKey",
	}

	hashKey, err := cache.GetHashKey(testResp)

	cache.ClearCache(hashKey)
	err = cache.SaveCache(hashKey, testResp)

	assert.Nil(t, err, "ошибка при сохранении кэша: %s", err)
}
