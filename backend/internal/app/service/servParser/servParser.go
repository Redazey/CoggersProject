package servParser

import (
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
}

func New() *Service {
	return &Service{}
}

// Сюда передаем адрес сервера в формате string "https://api.mcsrvstat.us/3/mc.hypixel.net"
func parseServerInfo(serverAddress string) (map[string]interface{}, error) {
	var serverData map[string]interface{}

	response, err := http.Get(serverAddress)
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
		return nil, err
	}

	return serverData, nil
}

func (s *Service) GetOnlineInfo(servers map[string]string) (map[string]map[string]interface{}, error) {
	serversInfo := make(map[string]map[string]interface{})
	defer rec.Recovery()

	for key, value := range servers {
		serverInfo, err := parseServerInfo(value)
		if err != nil {
			logStr := fmt.Sprintf("не удалось получить данные о сервере %s, ошибка: ", key)
			logger.Error(logStr, zap.Error(err))
		}

		serverStatus := serverInfo["online"]
		if serverStatus == false {
			logStr := fmt.Sprintf("Сервер %v не отвечает", key)
			logger.Info(logStr)
			continue
		}

		serversInfo[key] = make(map[string]interface{})

		playersData := serverInfo["players"].(map[string]interface{})

		serversInfo[key]["online"] = playersData["online"]
		serversInfo[key]["max"] = playersData["max"]
		fmt.Printf("%v: %v\n", key, playersData)
	}

	if len(serversInfo) == 0 {
		return nil, errorz.ErrServerIsNotResponse
	}

	return serversInfo, nil
}
