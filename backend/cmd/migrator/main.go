package main

import (
	"CoggersProject/config"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath string
	// Путь до папки с миграциями.
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	env, err := config.NewEnv()
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфига: %s", err)
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", env.DB.DBUser, env.DB.DBPassword, env.DB.DBHost, env.DB.DBName)

	m, err := migrate.New(
		"file://"+migrationsPath,
		connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to apply")
		} else {
			log.Fatal(err)
		}
	} else {
		log.Println("Migrations applied successfully")
	}
}
