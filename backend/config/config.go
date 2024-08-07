package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Файл переменных окружения
type Enviroment struct {
	LoggerLevel string        `env:"loggerMode" envDefault:"debug"`
	GRPCTimeout time.Duration `env:"GRPC_TIMEOUT" envDefault:"10h"`
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
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDBId     int    `env:"REDIS_DB_ID,required"`
}

type Cache struct {
	EXTime         time.Duration `json:"EXTime"`
	UpdateInterval string        `json:"updateInterval"`
}

// Файл конфигурации
type Configuration struct {
	Servers []Servers `json:"servers"`
}

type Servers struct {
	Adress    string `json:"Adress"`
	Name      string `json:"Name"`
	Version   string `json:"Version"`
	Online    int
	MaxOnline int
}

var config Configuration
var enviroment Enviroment

/*
Структура env файла

	-------GENERAL------
	LoggerLevel string
	GRPCTimeout time.Duration
	JwtSecret   string
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
func NewEnv(envPath ...string) (*Enviroment, error) {
	err := godotenv.Load(envPath...)
	if err != nil {
		log.Fatalf("Файл .env не найден: %s", err)
	}

	err = env.Parse(&enviroment)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&enviroment.Redis)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&enviroment.DB)
	if err != nil {
		return nil, err
	}

	return &enviroment, nil
}

func NewConfig(cfgPath string) (*Configuration, error) {
	configFile, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal("Ошибка при попытке прочитать файл конфигурации:", err)
		return nil, err
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Ошибка при распаковывании файла конфигурации:", err)
		return nil, err
	}

	return &config, nil
}

func GetEnv() *Enviroment {
	return &Enviroment{}
}

func GetConfig() *Configuration {
	return &Configuration{}
}
