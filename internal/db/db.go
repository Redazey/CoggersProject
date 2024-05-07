package db

import (
	"goRoadMap/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"fmt"
	"os"
)

func dbConnect() (*sqlx.DB, error) {
	err := godotenv.Load("Z:/files/goRoadMap/goRoadMap/.env")
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		host     = os.Getenv("DB_HOST")
		db_name  = os.Getenv("DB_NAME")
	)

	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
		return nil, err
	}

	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		user, db_name, password, host)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		logger.Fatal("ошибка при подключении к БД", zap.Error(err))
		return nil, err
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		logger.Fatal("соединение с БД не прошло проверку: ", zap.Error(err))
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
			logger.Fatal("ошибка при выполнении SQL зарпоса: ", zap.Error(err))
			return err
		}
	}

	stmt, err := db.Prepare("INSERT INTO roles(name) VALUES($1)")
	if err != nil {
		logger.Fatal("ошибка при подготовке SQL зарпоса: ", zap.Error(err))
		return err
	}

	roles := []string{"user", "employee", "admin"}

	rowsCount, err := db.Query(`SELECT COUNT(*) FROM roles`)
	if err != nil {
		logger.Fatal("ошибка при подключении к БД", zap.Error(err))
		return err
	}

	// Проверка, заполнена ли таблица roles
	if rowsCount == nil {
		// Вставляем каждую роль в базу данных
		for _, role := range roles {
			_, err = stmt.Exec(role)
			if err != nil {
				logger.Fatal("ошибка при исполнении SQL запроса: ", zap.Error(err))
				return err
			}
		}
		logger.Info("таблица Roles заполнена")
	} else {
		logger.Info("таблица Roles заполнена")
	}

	logger.Info("таблицы Users и Roles успешно инициализированны")
	return nil
}

// Добавить кэширование
func GetLoginData(message map[string]string) (map[string]string, error) {
	db, err := dbConnect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	username := message["username"]
	password := message["password"]
	rows, err := db.Query(
		`SELECT login, password FROM users
        WHERE login = ? AND password = ?`, username, password)
	if err != nil {
		logger.Fatal("ошибка при исполнении SQL запроса: ", zap.Error(err))
		return nil, err
	}

	dbLoginData := make(map[string]string)

	for rows.Next() {
		var key string
		var value string

		err = rows.Scan(&key, &value)
		if err != nil {
			logger.Fatal("ошибка при получении значение количества строк из запроса: ", zap.Error(err))
			return nil, err
		}

		dbLoginData[key] = value
	}

	return dbLoginData, nil
}

/*
// добавить кэширование
func NewUserRegistration(message map[string]string) (map[string]string, error) {
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
*/
