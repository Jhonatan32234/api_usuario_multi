package usecases

import (
	"apisuario/cmd/domain/entities"
	"apisuario/cmd/domain/repositories"
)

type VerifyUserDeviceUseCase struct {
	repo repositories.UserProductRepository
}


func NewVerifyUserDeviceUseCase(repo repositories.UserProductRepository) *VerifyUserDeviceUseCase {
	return &VerifyUserDeviceUseCase{repo: repo}
}


func (uc *VerifyUserDeviceUseCase) Execute(data entities.UserVeryfication) (string, error) {
	return uc.repo.VerifyUserDevice(data)
}