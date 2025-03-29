package controllers

import (
	"apisuario/cmd/application/usecases"
	"apisuario/cmd/domain/entities"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


type VerifyUserDeviceController struct {
	usecase *usecases.VerifyUserDeviceUseCase
}


func NewVerifyUserDeviceController(usecase *usecases.VerifyUserDeviceUseCase) *VerifyUserDeviceController {	
	return &VerifyUserDeviceController{usecase: usecase}
}

func (c *VerifyUserDeviceController) VerifyUserDeviceHandler(ctx *gin.Context) {
    log.Print("VerifyUserDeviceHandler called")
    var data entities.UserVeryfication 
    if err := ctx.ShouldBindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    deviceID, err := c.usecase.Execute(data)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "No relationship found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"id_esp32": deviceID})
}
