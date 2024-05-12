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
	dataMap, err := db.GetData("users")
	if err != nil {
		logger.Error("ошибка при попытке получить данные из БД: ", zap.Error(err))
	}
	cache.SaveCache(dataMap)
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
