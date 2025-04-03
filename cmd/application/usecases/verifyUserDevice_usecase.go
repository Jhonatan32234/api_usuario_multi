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

func (uc *VerifyUserDeviceUseCase) VerifyCredentials(data entities.UserVeryfication) (*entities.User, *entities.Device, error) {
    user, device, err := uc.repo.VerifyUserAndDevice(data) // El método se llama VerifyUserAndDevice, no VerifyUserDevice
    if err != nil {
        return nil, nil, err
    }
    return user, device, nil
}