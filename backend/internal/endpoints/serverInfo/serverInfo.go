package serverinfo

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/pkg/service/logger"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Service interface {
	GetOnlineInfo(map[string]string) (map[string]map[string]interface{}, error)
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

func (e *Endpoint) ServerInfo(ctx echo.Context) error {
	serverMap := make(map[string]string)
	config := config.GetConfig()

	for _, server := range config.Servers {
		serverMap[server.Name] = "https://api.mcsrvstat.us/3/" + server.IP
	}

	returnMsg, err := e.s.GetOnlineInfo(serverMap)
	if err != nil {
		logger.Error("ошибка при попытке получить информацию о серверах: %s", zap.Error(err))
		return err
	}

	return ctx.JSON(http.StatusOK, returnMsg)
}
