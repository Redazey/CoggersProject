package cacher

import (
	"CoggersProject/pkg/cache"
	"fmt"

	"github.com/robfig/cron/v3"
)

func Init(interval string) {
	intervalStr := fmt.Sprintf("%s * * * *", interval)

	c := cron.New()
	_, err := c.AddFunc(intervalStr, func() {
		cache.DeleteEX("*")
	})
	if err != nil {
		return
	}

	c.Start()
}
