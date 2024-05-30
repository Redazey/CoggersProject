package db

import (
	"CoggersProject/backend/pkg/service/logger"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"fmt"
	"os"
)

func dbConnect() (*sqlx.DB, error) {
	defer func() {
		if r := recover(); r != nil {
			// Обработка паники
			fmt.Println("Паника обработана:", r)
		}
	}()
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
	defer func() {
		if r := recover(); r != nil {
			// Обработка паники
			fmt.Println("Паника обработана:", r)
		}
	}()
	// обьявляем список выполняемых SQL запросов
	SQLqueries := []string{
		`CREATE TABLE IF NOT EXISTS roles (
            id SERIAL PRIMARY KEY,
            name TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT,
            password TEXT,
            roleid INT,
            FOREIGN KEY (roleid) REFERENCES roles(id)
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

	rowsCount, err := db.Query("SELECT COUNT(*) FROM roles")
	if err != nil {
		logger.Error("ошибка при исполнении SQL запроса: ", zap.Error(err))
		return err
	}

	defer rowsCount.Close()

	var count int
	if rowsCount.Next() {
		err := rowsCount.Scan(&count)
		if err != nil {
			logger.Error("ошибка при сканировании результата запроса: ", zap.Error(err))
			return err
		}
	}
	// Проверка, заполнена ли таблица roles
	if count == 0 {
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
func FetchUserData(username string) (map[string]string, error) {
	// подключение к бд
	db, err := dbConnect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT username, password, roleid 
						   FROM users
						   WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Получение информации о столбцах
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Инициализация именованного массива, который содержит структуру для сканирования
	values := make([]interface{}, len(columns))
	for i := range columns {
		values[i] = new(interface{})
	}

	// Инициализация мапы для хранения данных
	dataMap := make(map[string]string)

	// Чтение данных из таблицы и добавление их в map
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i, colName := range columns {
			val := *(values[i].(*interface{}))
			if val == nil {
				dataMap[colName] = ""
			} else {
				dataMap[colName] = fmt.Sprintf("%v", val)
			}
		}
	}

	return dataMap, nil
}

// принимает map - значения, которые нужно внести в БД и string - таблицу, в которую будем вносить значения
func PullData(table string, data map[string]map[string]interface{}) error {
	db, err := dbConnect()
	if err != nil {
		return err
	}

	for _, keyData := range data {
		var (
			columns []string
			values  []string
		)

		for key, value := range keyData {
			columns = append(columns, key)
			values = append(values, fmt.Sprintf("%s", value))
		}
		cmdStr := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?)`, table, strings.Join(columns, ", "))
		query, args, err := sqlx.In(cmdStr, values)

		if err != nil {
			logger.Error("Ошибка при выполнении SQL-запроса: ", zap.Error(err))
			return err
		}

		query = db.Rebind(query)
		_, err = db.Query(query, args...)
		if err != nil {
			logger.Error("Ошибка при выполнении SQL-запроса: ", zap.Error(err))
			return err
		}
	}

	logger.Info("Данные в БД были успешно обновлены")
	return nil
}
