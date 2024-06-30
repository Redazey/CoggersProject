package db

import (
	"CoggersProject/pkg/db"
)

type ServerInfo struct {
	Id        int
	Adress    string
	Name      string
	Version   string
	MaxOnline float64
	Online    float64
}

/*
вовзращает данные о серверах из БД в виде структуры:

	Id        int
	Adress    string
	Name      string
	Version   string
	MaxOnline int NULL
	Online    int NULL
*/
func FetchServersData() (map[string]ServerInfo, error) {
	rows, err := db.Conn.Query("SELECT * FROM servers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	serversDataMap := make(map[string]ServerInfo)

	// Чтение данных из таблицы и добавление их в map
	for rows.Next() {
		var server ServerInfo
		err := rows.Scan(&server.Id, &server.Adress, &server.Name, &server.Version)
		if err != nil {
			return nil, err
		}

		serversDataMap[server.Name] = server
	}

	return serversDataMap, nil
}
