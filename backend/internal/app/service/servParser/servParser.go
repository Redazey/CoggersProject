package servParser

import (
	"CoggersProject/internal/app/errorz"
	"CoggersProject/internal/app/lib/db"
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

func (s *Service) GetServersInfo() (map[string]db.ServerInfo, error) {
	servers, err := db.FetchServersData()
	if err != nil {
		logger.Error("ошибка при получении данных о серверах из БД: ", zap.Error(err))
		return nil, err
	}

	defer rec.Recovery()

	for key, value := range servers {
		serverInfo, err := parseServerInfo(value.Adress)
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

		playersData := serverInfo["players"].(map[string]interface{})

		value.Online = playersData["online"].(int)
		value.MaxOnline = playersData["max"].(int)
		fmt.Printf("%v: %v\n", key, playersData)
	}

	if len(servers) == 0 {
		return nil, errorz.ErrServerIsNotResponse
	}

	return servers, nil
}
