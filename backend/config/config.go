package config

import (
	"encoding/json"
	"log"
	"os"
)

type Cache struct {
	EXTime         int    `json:"EXTime"`
	UpdateInterval string `json:"updateTime"`
}

type Servers struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type Config struct {
	EnvPath    string `json:"envPath"`
	LoggerMode string `json:"loggerMode"`
	Cache      Cache
	Servers    []Servers `json:"servers"`
}

var config Config

func Init() {
	configFile, err := os.ReadFile("C:/CoggersProject/backend/config/config.json")
	if err != nil {
		log.Fatal("Ошибка при попытке прочитать файл конфигурации: ", err)
		return
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Ошибка при распаковывании файла конфигурации: ", err)
		return
	}
}

func GetConfig() *Config {
	return &config
}
