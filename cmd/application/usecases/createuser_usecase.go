package usecases

import (
	"apisuario/cmd/domain/entities"
	"apisuario/cmd/domain/repositories"
)

type CreateUserUseCase struct {
	repo repositories.UserProductRepository
}

func NewCreateUserUseCase(repo repositories.UserProductRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(data entities.UserVeryfication) error {
	return uc.repo.CreateUser(data)
}