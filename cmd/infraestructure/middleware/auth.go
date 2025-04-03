package middlewares

import (
	"apisuario/cmd/domain/entities"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(requiredRole string, jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &entities.Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
            }
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            return
        }

        // Verificar rol si se especificó
        if requiredRole != "" && claims.Tipo != requiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": fmt.Sprintf("Se requiere rol '%s'", requiredRole),
            })
            return
        }

        // Añadir información al contexto
        c.Set("userID", claims.UserID)
        c.Set("userType", claims.Tipo)
        c.Set("deviceID", claims.IdEsp32)

        c.Next()
    }
}