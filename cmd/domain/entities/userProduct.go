package entities

type UserProduct struct {
	ID       int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`	
	IdEsp32 string `json:"id_esp32"`
}

