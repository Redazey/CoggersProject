package db

import (
	"CoggersProject/pkg/db"
)

type ServerInfo struct {
	Id        int
	Ip        string
	Name      string
	Version   string
	MaxOnline int
	Online    int
}

/*
вовзращает данные о серверах из БД в виде структуры:

	Ip        string
	Name      string
	Version   string
	MaxOnline int
	Online    int
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
		err := rows.Scan(&server.Id, &server.Ip, &server.Name, &server.Version, &server.MaxOnline, &server.Online)
		if err != nil {
			return nil, err
		}

		serversDataMap[server.Name] = server
	}

	return serversDataMap, nil
}
