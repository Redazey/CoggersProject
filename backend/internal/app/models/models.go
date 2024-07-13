package models

type (
	ServerInfo struct {
		Id        int
		Adress    string
		Name      string
		Version   string
		MaxOnline float64
		Online    float64
	}
)
