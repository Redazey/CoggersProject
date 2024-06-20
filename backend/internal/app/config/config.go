package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
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
	EXTime         time.Duration `env:"EXTime" envDefault:"15m"`
	UpdateInterval string        `env:"updateTime" envDefault:"15"`
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

	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&config.Redis)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&config.DB)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfig() *Configuration {
	return &Configuration{}
}
