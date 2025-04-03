package entities

import "github.com/golang-jwt/jwt/v4"

type AuthResponse struct {
    Token   string `json:"token"`
    Tipo    string `json:"tipo"`
    IdEsp32 string `json:"id_esp32"`
}


type Claims struct {
	UserID  int    `json:"user_id"`
	Tipo    string `json:"tipo"`
	IdEsp32 string `json:"id_esp32"`
	jwt.RegisteredClaims
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Tipo     string `json:"tipo"`
}

type Device struct {
    ID      int    `json:"id"`
    IdEsp32 string `json:"id_esp32"`
}