package config

import (
	"encoding/json"
	"log"
	"os"
)

type Cache struct {
	EXTime         int    `json:"EXTime"`
	UpdateInterval string `json:"updateInterval"`
}

type Config struct {
	EnvPath    string `json:"envPath"`
	LoggerMode string `json:"loggerMode"`
	Cache      Cache
}

var config Config

func Init() {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Ошибка при попытке прочитать файл конфигурации:", err)
		return
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Ошибка при распаковывании файла конфигурации:", err)
		return
	}
}

func GetConfig() *Config {
	return &config
}
