package db_api

import (
	"github.com/jmoiron/sqlx"

	"fmt"
	"os"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func dbConnect() {

	user := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", host, user, password, db_name)
	db, err := sqlx.Open("postgres", connStr) // Используйте нужный драйвер для вашей БД
	errCheck(err)

	// Проверка соединения
	err = db.Ping()
	errCheck(err)

	defer db.Close()

}
