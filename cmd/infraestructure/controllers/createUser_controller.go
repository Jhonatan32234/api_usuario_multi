package controllers

import (
	"apisuario/cmd/application/usecases"
	"apisuario/cmd/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserController struct {
	usecase *usecases.CreateUserUseCase
}

func NewCreateUserController(usecase *usecases.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{usecase: usecase}
}

func (c *CreateUserController) CreateUserHandler(ctx *gin.Context) {
	var data entities.UserVeryfication
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.usecase.Execute(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}