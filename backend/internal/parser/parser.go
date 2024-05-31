package parser

import (
	"CoggersProject/backend/pkg/service/logger"
	"fmt"
	"net"

	"go.uber.org/zap"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

// Сюда передаем адрес сервера в формате string "mc.hypixel.net:25565"
func parseServerInfo(serverAddress string) (string, error) {
	conn, err := net.Dial("udp", serverAddress)
	if err != nil {
		logger.Error("произошла ошибка при подключении к серверу: ", zap.Error(err))
		return "", err
	}
	defer conn.Close()

	// Отправляем Handshake пакет для получения информации о сервере
	packetID := make([]byte, 1)
	packetID[0] = 0xFE // Номер пакета Handshake
	conn.Write(packetID)

	// Получаем данные о сервере
	buffer := make([]byte, 2048)
	_, err = conn.Read(buffer)
	if err != nil {
		logger.Error("произошла ошибка при получении данных от серверу: ", zap.Error(err))
		return "", nil
	}

	// Парсим полученные данные
	serverInfo := string(buffer[3:])
	return serverInfo, nil
}

func (s *Service) GetOnlineInfo(servers map[string]string) (map[string]string, error) {
	// serversInfoMap := make(map[string]string)

	for key, value := range servers {
		serverInfo, err := parseServerInfo(value)
		if err != nil {
			logStr := fmt.Sprintf("не удалось получить данные о сервере %s, ошибка: ", key)
			logger.Error(logStr, zap.Error(err))
		}

		fmt.Printf("Данные, которые были получены с сервера: %s - %s", key, serverInfo)
	}

	return nil, nil
}
