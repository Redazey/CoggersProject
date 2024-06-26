package auth_tests

import (
	"CoggersProject/pkg/db"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

// Добавляет в бд случайного юзера с данным roleId и возвращает его данные
func MockUser(roleId int) (map[string]interface{}, error) {
	UserData := map[string]interface{}{
		"name":      gofakeit.Name(),
		"birthdate": gofakeit.Date(),
		"photourl":  "testimg",
		"push":      gofakeit.Bool(),
		"email":     gofakeit.Email(),
		"password":  gofakeit.Password(true, true, true, true, false, 10),
		"roleId":    roleId,
	}

	// SQL запрос для вставки данных
	sqlStatement := `INSERT INTO Users (name, birthdate, photourl, push, email, password, roleId) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Выполнение SQL запроса
	_, err := db.Conn.Exec(sqlStatement, UserData["name"],
		UserData["birthdate"], UserData["photourl"], UserData["push"],
		UserData["email"], UserData["password"], UserData["roleId"])
	if err != nil {
		return nil, err
	}

	fmt.Println("Данные успешно добавлены в таблицу Users")

	return UserData, nil
}

// Очищает таблицу в бд
func ClearTable(table string) error {
	// SQL запрос для очистки
	sqlStatement := fmt.Sprintf(`DELETE FROM %s`, table)

	// Выполнение SQL запроса
	_, err := db.Conn.Exec(sqlStatement)
	if err != nil {
		return err
	}

	fmt.Println("Таблица очищена")

	return nil
}
