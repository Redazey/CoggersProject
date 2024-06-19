package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	LoggerLevel string        `json:"loggerMode" jsonDefault:"debug"`
	GRPCTimeout time.Duration `json:"GRPC_TIMEOUT" jsonDefault:"10h"`
	Servers     []Servers     `json:"servers"`
	JwtSecret   string        `env:"JWT_SECRET,required"`
	DB          DB
	Redis       Redis
	Cache       Cache
}

type DB struct {
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBHost     string `env:"DB_HOST,required"`
}

type Redis struct {
	RedisAddr     string `env:"REDIS_ADDR,required"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisUsername string `env:"REDIS_USERNAME,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDBId     int    `env:"REDIS_DB_ID,required"`
}

type Cache struct {
	EXTime         time.Duration `json:"EXTime" jsonDefault:"15m"`
	UpdateInterval string        `json:"updateTime" jsonDefault:"15"`
}

type Servers struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

var config Configuration

/*
Структура файла конфигурации

	-------GENERAL------
	Port          string
	gRPCTimeout   time.Duration
	LoggerLevel   string
	---------DB---------
	DBUser        string
	DBPassword    string
	DBName        string
	DBHost        string
	-------REDIS--------
	RedisAddr     string
	RedisPort     string
	RedisPassword string
	RedisDBid     int
	-------CACHE--------
	CacheInterval string
	CacheEXTime   int
*/
func NewConfig(files ...string) (*Configuration, error) {
	err := godotenv.Load(files...)
	if err != nil {
		log.Fatalf("Файл .env не найден: %s", err)
	}

	configFile, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal("Ошибка при попытке прочитать файл конфигурации: ", err)
		return nil, err
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Ошибка при распаковывании файла конфигурации: ", err)
		return nil, err
	}

	cfg := Configuration{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg.Redis)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg.DB)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetConfig() *Configuration {
	return &Configuration{}
}
