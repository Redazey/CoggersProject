package db

import (
	"CoggersProject/pkg/db"
)

type UserData struct {
	Id        int
	Name      string
	Email     string
	Password  string
	RoleId    int
	Birthdate string
	Photourl  string
	Push      bool
}

/*
принимает таблицу как string и возвращает данные о пользователе в виде структуры

	name      string
	email     string
	password  string
	roleId    int
	birthdate string
	photourl  string
	push      bool
*/
func FetchUserData(email string) (UserData, error) {
	rows, err := db.Conn.Query(`SELECT * FROM users
						   		WHERE email = $1`, email)
	if err != nil {
		return UserData{}, err
	}

	defer rows.Close()

	// Чтение данных из таблицы и добавление их в map
	var user UserData
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RoleId, &user.Birthdate, &user.Photourl, &user.Push)
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}
