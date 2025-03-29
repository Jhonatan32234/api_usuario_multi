package routes

import (
	"apisuario/cmd/application/usecases"
	"apisuario/cmd/domain/repositories"
	"apisuario/cmd/infraestructure/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()

    userProductRepo := repositories.NewUserRepositoryDB(db)

    connectUserProductUseCase := usecases.NewConnectUserProductUseCase(userProductRepo)
    verifyUserDeviceUseCase := usecases.NewVerifyUserDeviceUseCase(userProductRepo)

    controllerConnect := controllers.NewConnectUserProductController(connectUserProductUseCase)
    controllerVerify := controllers.NewVerifyUserDeviceController(verifyUserDeviceUseCase)

    router.POST("/connect", controllerConnect.ConnectUserProductHandler)
    router.POST("/verify", controllerVerify.VerifyUserDeviceHandler)

    return router
}
