package cache_test

import (
	"CoggersProject/pkg/cache"
	"CoggersProject/pkg/cache/tests/suite"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveCache(t *testing.T) {
	suite.New(t)

	testResp := "testResponse"
	hashKey := fmt.Sprintf("test_%s", testResp)

	cache.ClearCache(hashKey)
	err := cache.SaveCache(hashKey, testResp)

	assert.Nil(t, err, "ошибка при сохранении кэша: %s", err)
}
