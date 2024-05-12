package db

import (
	"goRoadMap/pkg/services/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"fmt"
	"os"
)

func dbConnect() (*sqlx.DB, error) {
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		host     = os.Getenv("DB_HOST")
		db_name  = os.Getenv("DB_NAME")
	)

	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		user, db_name, password, host)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		logger.Error("ошибка при подключении к БД", zap.Error(err))
		return nil, err
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		logger.Error("соединение с БД не прошло проверку: ", zap.Error(err))
		return nil, err
	}

	logger.Info("соединение с БД успешно установленно")
	return db, nil
}

func InitiateTables() error {
	// обьявляем список выполняемых SQL запросов
	SQLqueries := []string{
		`CREATE TABLE IF NOT EXISTS roles (
            id SERIAL PRIMARY KEY,
            name TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            login TEXT,
            password TEXT,
            token TEXT,
            role INT,
            FOREIGN KEY (role) REFERENCES roles(id)
        )`,
	}
	// подключение к бд
	db, err := dbConnect()
	if err != nil {
		return err
	}

	defer db.Close()

	for _, query := range SQLqueries {
		_, err = db.Exec(query)
		if err != nil {
			logger.Error("ошибка при выполнении SQL зарпоса: ", zap.Error(err))
			return err
		}
	}

	stmt, err := db.Prepare("INSERT INTO roles(name) VALUES($1)")
	if err != nil {
		logger.Error("ошибка при подготовке SQL зарпоса: ", zap.Error(err))
		return err
	}

	roles := []string{"user", "employee", "admin"}

	rowsCount, err := db.Query(`SELECT COUNT(*) FROM roles`)
	if err != nil {
		logger.Error("ошибка при исполнении SQL запроса: ", zap.Error(err))
		return err
	}

	// Проверка, заполнена ли таблица roles
	if rowsCount == nil {
		// Вставляем каждую роль в базу данных
		for _, role := range roles {
			_, err = stmt.Exec(role)
			if err != nil {
				logger.Error("ошибка при исполнении SQL запроса: ", zap.Error(err))
				return err
			}
		}
	}
	logger.Info("таблица Roles заполнена")

	logger.Info("таблицы Users и Roles успешно инициализированны")
	return nil
}

// принимает таблицу как string и возвращает таблицу в виде map
func GetData(table string) (map[string]map[string]string, error) {
	// подключение к бд
	db, err := dbConnect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT * FROM $1`, table)
	if err != nil {
		return nil, err
	}

	// Инициализация map для хранения данных
	usersMap := make(map[string]map[string]string)

	// Чтение данных из таблицы и добавление их в map
	for rows.Next() {
		var name, password, roleId string
		err := rows.Scan(&name, &password, &roleId)
		if err != nil {
			logger.Error("Ошибка при сканировании sql.Rows: ", zap.Error(err))
			return nil, err
		}

		userData := map[string]string{
			"password": password,
			"role":     roleId,
		}
		usersMap[name] = userData
	}

	return usersMap, nil
}
