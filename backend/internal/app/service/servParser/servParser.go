package servParser

import (
	"CoggersProject/config"
	"CoggersProject/internal/app/errorz"
	rec "CoggersProject/internal/app/lib/recovery"
	"CoggersProject/pkg/logger"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Tnze/go-mc/bot"
	"go.uber.org/zap"
)

type parsedData struct {
	Players struct {
		MaxOnline int `json:"max"`
		Online    int `json:"online"`
	}
}

type Service struct {
	Cfg *config.Configuration
}

func New(cfg *config.Configuration) *Service {
	return &Service{
		Cfg: cfg,
	}
}

func parseServerInfo(serverAddress string) (*parsedData, error) {
	response, _, err := bot.PingAndList(serverAddress)
	if err != nil {
		fmt.Println("Ошибка при пинге сервера:", err)
		return nil, err
	}
	var serverData *parsedData
	json.Unmarshal(response, &serverData)

	return serverData, nil
}

func (s *Service) GetServersInfo() ([]config.Servers, error) {
	servers := s.Cfg.Servers
	var wg sync.WaitGroup
	serversInfo := make([]config.Servers, len(servers))

	defer rec.Recovery()

	for i, server := range servers {
		wg.Add(1)
		go func(i int, server config.Servers) {
			defer wg.Done()
			serverInfo, err := parseServerInfo(server.Adress)
			if err != nil {
				logStr := fmt.Sprintf("не удалось получить данные о сервере %s, ошибка: ", server.Name)
				logger.Error(logStr, zap.Error(err))
				return
			}

			server.Online = serverInfo.Players.Online
			server.MaxOnline = serverInfo.Players.MaxOnline
			serversInfo[i] = server
		}(i, server)
	}
	wg.Wait()

	if len(servers) == 0 {
		return nil, errorz.ErrServerIsNotResponse
	}

	return serversInfo, nil
}
