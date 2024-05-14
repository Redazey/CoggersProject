package cache_test

import (
	"goRoadMap/internal/cache"
	"goRoadMap/pkg/services/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDummyData(t *testing.T) {
	logger.Init("info")
	t.Run("FillingWithData", func(t *testing.T) {
		dummydata := map[string]string{
			"username": "testuser",
			"password": "exampass",
			"roleid":   "1",
		}
		err := cache.SaveCache("test", cache.ConvertMap(dummydata))
		assert.NoError(t, err, "ошибка при заполнении Redis: %v", err)
	})
}
func TestIsDataInCache(t *testing.T) {
	// Тестирование кейса, когда данные есть в кэше
	t.Run("DataInCache", func(t *testing.T) {
		result, err := cache.IsDataInCache("test", "testuser", "password")

		assert.NoError(t, err, "ошибок при выполнении не найдено")
		assert.Equal(t, "exampass", result, "ожидаемое значение - \"exampass\", вышло: %v", result)
	})

	// Тестирование кейса, когда данных нет в кэше
	t.Run("NoDataInCache", func(t *testing.T) {
		result, err := cache.IsDataInCache("test", "ghostuser", "password")
		assert.Nil(t, err, "неожиданная ошибка: %v", err)
		assert.Nil(t, result, "ожидалось nil, вышло: %v", result)
	})
}

func TestReadCache(t *testing.T) {
	// Тестирование чтения кэша с определенным ключом
	t.Run("ReadSpecificKey", func(t *testing.T) {
		expectedValue := map[string]interface{}{
			"password": "exampass",
			"roleid":   "1",
		}
		result, err := cache.ReadCache("test")

		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, expectedValue, result["testuser"], "Expected value %s", expectedValue)
	})
}
