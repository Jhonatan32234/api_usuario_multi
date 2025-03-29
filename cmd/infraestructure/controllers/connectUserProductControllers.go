package controllers

import (
    "apisuario/cmd/domain/entities"
    "apisuario/cmd/application/usecases"
    "github.com/gin-gonic/gin"
    "net/http"
)

type ConnectUserProductController struct {
    useCase *usecases.ConnectUserProductUseCase
}

func NewConnectUserProductController(useCase *usecases.ConnectUserProductUseCase) *ConnectUserProductController {
    return &ConnectUserProductController{useCase: useCase}
}

func (c *ConnectUserProductController) ConnectUserProductHandler(ctx *gin.Context) {
    var data entities.UserProduct
    if err := ctx.ShouldBindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if err := c.useCase.Execute(data); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"message": "User connected with device"})
}


