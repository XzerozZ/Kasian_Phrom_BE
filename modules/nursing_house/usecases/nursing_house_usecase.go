package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
)

type NhUseCase interface {
	CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetNhByID(id int) (entities.NursingHouse, error)
}

type NhUseCaseImpl struct {
	nhrepo repositories.NhRepository
}

func NewNhUseCase(nhrepo repositories.NhRepository) *NhUseCaseImpl {
	return &NhUseCaseImpl{nhrepo: nhrepo}
}

func (u *NhUseCaseImpl) CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error) {
	return u.nhrepo.CreateNh(nursingHouse)
}

func (u *NhUseCaseImpl) GetAllNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetAllNh()
}

func (u *NhUseCaseImpl) GetNhByID(id int) (entities.NursingHouse, error) {
	return u.nhrepo.GetNhByID(id)
}