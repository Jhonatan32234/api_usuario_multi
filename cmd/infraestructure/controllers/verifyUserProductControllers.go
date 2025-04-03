package controllers

import (
	"apisuario/cmd/application/usecases"
	"apisuario/cmd/domain/entities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type VerifyUserDeviceController struct {
	usecase *usecases.VerifyUserDeviceUseCase
	jwtSecret string
}

func NewVerifyUserDeviceController(usecase *usecases.VerifyUserDeviceUseCase, jwtSecret string) *VerifyUserDeviceController {	
	return &VerifyUserDeviceController{
		usecase: usecase,
		jwtSecret: jwtSecret,
	}
}
func (c *VerifyUserDeviceController) VerifyUserDeviceHandler(ctx *gin.Context) {
    var data entities.UserVeryfication
    if err := ctx.ShouldBindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de solicitud inválidos"})
        return
    }

    user, device, err := c.usecase.VerifyCredentials(data)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Generar token JWT
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &entities.Claims{
        UserID:  user.ID,
        Tipo:    user.Tipo,
        IdEsp32: device.IdEsp32,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "api-usuario",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(c.jwtSecret))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
        return
    }

    // Enviar token en header y datos en cuerpo
    ctx.Header("Authorization", "Bearer "+tokenString)
    ctx.JSON(http.StatusOK, gin.H{
        "tipo":     user.Tipo,
        "id_esp32": device.IdEsp32,
    })
}