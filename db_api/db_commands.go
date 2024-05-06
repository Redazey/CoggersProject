package db_api

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"fmt"
	"os"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func DBConnect(f func(message map[string]string, db *sqlx.DB) map[string]string) func() map[string]string {
	return func() map[string]string {
		user := os.Getenv("DB_NAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		db_name := os.Getenv("DB_NAME")

		connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", host, user, password, db_name)
		db, err := sqlx.Open("postgres", connStr) // Используйте нужный драйвер для вашей БД

		errCheck(err)
		defer db.Close()

		// Проверка соединения
		err = db.Ping()
		errCheck(err)

		returnData := f(db)
		errCheck(err)

		return returnData
	}
}

func InitiateTables(message map[string]string, db *sqlx.DB) map[string]string {
	SQLqueries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            login TEXT,
            password TEXT,
            token TEXT,
            role INT,
            FOREIGN KEY (role) REFERENCES roles(id)
        )`,
		`CREATE TABLE IF NOT EXISTS roles (
            id SERIAL PRIMARY KEY,
            name TEXT
        )`,
	}

	// Подготовка подключений к бд
	stmt, err := db.Prepare("INSERT INTO roles(name) VALUES(?)")
	errCheck(err)

	for _, query := range SQLqueries {
		db.MustExec(query)
	}

	roles := []string{"user", "employee", "admin"}

	// Проверка, заполнена ли таблица roles
	rows, err := db.Query(`SELECT COUNT(*) FROM roles`)
	errCheck(err)

	var rowsCount int
	err = rows.Scan(&rowsCount)
	errCheck(err)

	if rowsCount == 0 {
		// Вставляем каждую роль в базу данных
		for _, role := range roles {
			_, err = stmt.Exec(role)
			errCheck(err)
		}
	}

	return nil
}

// Добавить кэширование
func GetLoginData(message map[string]string, db *sqlx.DB) map[string]string {
	username := message["username"]
	password := message["password"]
	rows, err := db.Query(
		`SELECT login, password FROM users
        WHERE login = ? AND password = ?`, username, password)
	errCheck(err)

	dbLoginData := make(map[string]string)

	for rows.Next() {
		var key string
		var value string

		err = rows.Scan(&key, &value)
		errCheck(err)

		dbLoginData[key] = value
	}

	return dbLoginData
}

func NewUserRegistration(message map[string]string) {
	username := message["username"]
	password := message["password"]
	getLoginData := DBConnect(GetLoginData)
	dbLoginData = getLoginData(username, password)

	if dbLoginData {
		return errors.New("User already registered")
	} else {
		add_user(username, password, 1)
		return loginFunction
	}
}
