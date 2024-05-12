package cache_test

import (
	"goRoadMap/internal/cache"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDummyData(t *testing.T) {
	//logger.Init("info")
	t.Run("FillingWithData", func(t *testing.T) {
		dummydata := map[string]map[string]string{
			"testuser": {
				"password": "exampass",
				"roleid":   "1",
			},
		}
		err := cache.SaveCache(dummydata, "test")
		assert.NoError(t, err, "ошибка при заполнении Redis: %v", err)
	})
}
func TestIsDataInCache(t *testing.T) {
	// Тестирование кейса, когда данные есть в кэше
	t.Run("DataInCache", func(t *testing.T) {
		result, err := cache.IsDataInCache("test", "testuser", "password")

		assert.NoError(t, err, "ошибка: %v", err)
		assert.Equal(t, "exampass", result, "ожидаемое значение - \"exampass\", вышло: %v", result)
	})

	// Тестирование кейса, когда данных нет в кэше
	t.Run("NoDataInCache", func(t *testing.T) {
		result, err := cache.IsDataInCache("test", "ghostuser", "password")

		assert.Error(t, err, "ошибка")
		assert.Nil(t, result, "ожидалось nil, вышло: %v", result)
	})
}

func TestReadCache(t *testing.T) {
	// Тестирование чтения кэша с определенным ключом
	t.Run("ReadSpecificKey", func(t *testing.T) {
		expectedValue := map[string]string{
			"password": "exampass",
			"roleid":   "1",
		}
		result, err := cache.ReadCache("test", "testuser")

		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, expectedValue, result, "Expected value %s", expectedValue)
	})

	// Допишите другие тесты, если необходимо
}
