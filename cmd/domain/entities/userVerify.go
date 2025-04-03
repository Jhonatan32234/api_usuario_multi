package entities



type UserVeryfication struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Tipo     string `json:"tipo,omitempty"` // omitempty para que sea opcional
}
