package servParser

import (
	"CoggersProject/config"
	"CoggersProject/internal/app/errorz"
	rec "CoggersProject/internal/app/lib/recovery"
	"CoggersProject/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Service struct {
	Cfg *config.Configuration
}

func New(cfg *config.Configuration) *Service {
	return &Service{
		Cfg: cfg,
	}
}

// Сюда передаем адрес сервера в формате string "https://api.mcsrvstat.us/3/mc.hypixel.net"
func parseServerInfo(serverAddress string) (map[string]interface{}, error) {
	var serverData map[string]interface{}

	response, err := http.Get(fmt.Sprintf("https://api.mcsrvstat.us/3/%s", serverAddress))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return nil, err
	}

	defer response.Body.Close()

	serviceResponse, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("произошла ошибка при получении данных от серверу: ", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(serviceResponse, &serverData)
	if err != nil {
		logger.Error("произошла ошибка при пробразовании из json формата: ", zap.Error(err))

		logStr := fmt.Sprintf("содержание response: %s", serviceResponse)
		logger.Debug(logStr)
		return nil, err
	}

	return serverData, nil
}

func (s *Service) GetServersInfo() (map[string]config.Servers, error) {
	servers := s.Cfg.Servers
	serversInfo := make(map[string]config.Servers)

	defer rec.Recovery()

	for key, server := range servers {
		serverInfo, err := parseServerInfo(server.Adress)
		if err != nil {
			logStr := fmt.Sprintf("не удалось получить данные о сервере %s, ошибка: ", server.Name)
			logger.Error(logStr, zap.Error(err))
			continue
		}

		serverStatus := serverInfo["online"]
		if serverStatus == false {
			logStr := fmt.Sprintf("Сервер %v не отвечает", key)
			logger.Warn(logStr)
			continue
		}

		playersData := serverInfo["players"].(map[string]interface{})

		server.Online = playersData["online"].(float64)
		server.MaxOnline = playersData["max"].(float64)
		serversInfo[server.Name] = server
	}

	if len(servers) == 0 {
		return nil, errorz.ErrServerIsNotResponse
	}

	return serversInfo, nil
}
