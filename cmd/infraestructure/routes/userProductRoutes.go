package routes

import (
	"apisuario/cmd/application/usecases"
	"apisuario/cmd/domain/repositories"
	"apisuario/cmd/infraestructure/controllers"
	middlewares "apisuario/cmd/infraestructure/middleware"

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
    jwtSecret := "tu_clave_secreta_muy_segura" 
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
    createUserUseCase := usecases.NewCreateUserUseCase(userProductRepo)

    // Controladores
    controllerConnect := controllers.NewConnectUserProductController(connectUserProductUseCase)
    controllerVerify := controllers.NewVerifyUserDeviceController(verifyUserDeviceUseCase,jwtSecret)
    controllerCreateUser := controllers.NewCreateUserController(createUserUseCase)

    router.POST("/verify", controllerVerify.VerifyUserDeviceHandler)
    router.POST("/users", controllerCreateUser.CreateUserHandler)

    // Rutas protegidas para admin
    adminRoutes := router.Group("/admin")
    adminRoutes.Use(middlewares.AuthMiddleware("admin", jwtSecret))
    {
        adminRoutes.POST("/connect", controllerConnect.ConnectUserProductHandler)
        // Otras rutas de admin
    }

    // Rutas protegidas para cualquier autenticado
    authRoutes := router.Group("/")
    authRoutes.Use(middlewares.AuthMiddleware("", jwtSecret))
    {
        // Rutas que requieren autenticación
    }

    return router
}