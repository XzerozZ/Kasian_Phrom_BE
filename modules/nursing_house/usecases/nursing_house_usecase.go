package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
)

type NhUsecase interface {
	CreateNh(nursingHouse *entities.NursingHouse) error
}

type NhUsecaseImpl struct {
	nhrepo repositories.NhRepository
}

type NHUsecase interface {
	CreateNursingHouse(nursingHouse *entities.NursingHouse) error
}

func NewNhUseCase(nhrepo repositories.NhRepository) *NhUsecaseImpl {
	return &NhUsecaseImpl{nhrepo: nhrepo}
}

func (u *NhUsecaseImpl) CreateNh(nursingHouse *entities.NursingHouse) error {
	return u.nhrepo.Create(nursingHouse)
}

func (u *NhUsecaseImpl) GetAllNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetAllNh()
}

func (u *NhUseCaseImpl) GetNhByID(id string) (entities.NursingHouse, error) {
	return u.nhrepo.GetNhByID(id)
}