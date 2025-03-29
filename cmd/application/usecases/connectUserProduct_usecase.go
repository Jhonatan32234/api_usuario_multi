package usecases

import (
    "apisuario/cmd/domain/entities"
    "apisuario/cmd/domain/repositories"
)

type ConnectUserProductUseCase struct {
    repo repositories.UserProductRepository
}

func NewConnectUserProductUseCase(repo repositories.UserProductRepository) *ConnectUserProductUseCase {
    return &ConnectUserProductUseCase{repo: repo}
}

func (uc *ConnectUserProductUseCase) Execute(data entities.UserProduct) error {
    return uc.repo.ConnectUserProduct(data)
}
