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

    // Aplica CORS antes de definir las rutas
    router.Use(cors.Default()) // Esto habilita CORS con la configuración predeterminada (permite todos los orígenes)

    // Si deseas más configuraciones personalizadas, puedes usar cors.New(), por ejemplo:
    // router.Use(cors.New(cors.Config{
    //     AllowOrigins: []string{"http://localhost:3000"},  // Solo permite ciertos orígenes
    //     AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    //     AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    // }))

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
