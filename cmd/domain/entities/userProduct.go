package entities

type UserProduct struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`	
	IdEsp32  string `json:"id_esp32"`
	Tipo     string `json:"tipo,omitempty"` // omitempty para que sea opcional
}
