package routes

import (
    "apisuario/cmd/application/usecases"
    "apisuario/cmd/domain/repositories"
    "apisuario/cmd/infraestructure/controllers"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()

    // Middleware para manejar solicitudes OPTIONS y descartarlas con un 200
    router.Use(func(c *gin.Context) {
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        c.Next()
    })

    // Configuración de CORS personalizada
    router.Use(cors.New(cors.Config{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Repositorios y casos de uso
    userProductRepo := repositories.NewUserRepositoryDB(db)
    connectUserProductUseCase := usecases.NewConnectUserProductUseCase(userProductRepo)
    verifyUserDeviceUseCase := usecases.NewVerifyUserDeviceUseCase(userProductRepo)

    // Controladores
    controllerConnect := controllers.NewConnectUserProductController(connectUserProductUseCase)
    controllerVerify := controllers.NewVerifyUserDeviceController(verifyUserDeviceUseCase)

    // Rutas
    router.POST("/connect", controllerConnect.ConnectUserProductHandler)
    router.POST("/verify", controllerVerify.VerifyUserDeviceHandler)

    return router
}
