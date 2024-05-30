package cacher

import (
	"CoggersProject/backend/pkg/cache"
	"CoggersProject/backend/pkg/db"
	"CoggersProject/backend/pkg/service/logger"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func cacheUpdate() {
	// здесь настраиваем названия таблиц, которые будут сохранятся в кэше
	cacheTables := []string{"users"}

	// сначала загружаем весь кэш в БД
	for table := range cacheTables {
		// для этого получаем весь кэш
		cacheMap, err := cache.ReadCache(cacheTables[table])
		if err != nil {
			return
		}
		// затем загружаем в БД
		db.PullData(cacheTables[table], cacheMap)
	}
}

func Init(interval string) {
	intervalStr := fmt.Sprintf("%s * * * *", interval)

	c := cron.New()
	_, err := c.AddFunc(intervalStr, cacheUpdate)
	if err != nil {
		logger.Error("ошибка при инициализации кэша: ", zap.Error(err))
		return
	}

	c.Start()
}
