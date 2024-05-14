package cacher

import (
	"fmt"
	"goRoadMap/internal/cache"
	"goRoadMap/internal/db"
	"goRoadMap/pkg/services/logger"

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
	/*
	   в случае, если мы будем выгружать всю таблицу в кэш, что очень сомнительно
	   	// теперь обновляем кэш
	   	for table := range cacheTables {
	   		dataMap, err := db.GetData(cacheTables[table])
	   		if err != nil {
	   			logger.Error("ошибка при обновлении кэша ", zap.Error(err))
	   			return
	   		}

	   		err = cache.SaveCache(cacheTables[table], cache.ConvertMap(dataMap))
	   		if err != nil {
	   			logger.Error("ошибка при обновлении кэша ", zap.Error(err))
	   			return
	   		}

	   		err = cache.DeleteEX(cacheTables[table])
	   		if err != nil {
	   			logger.Error("ошибка при обновлении кэша ", zap.Error(err))
	   			return
	   		}
	   	}
	*/
}

func Init(interval string) {
	intervalStr := fmt.Sprintf("0 */%s * * * *", interval)

	c := cron.New()
	_, err := c.AddFunc(intervalStr, cacheUpdate)
	if err != nil {
		logger.Error("ошибка при инициализации кэша: ", zap.Error(err))
		return
	}

	c.Start()
}
