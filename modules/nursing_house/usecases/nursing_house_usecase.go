package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
)

type NhUseCase interface {
	CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetActiveNh() ([]entities.NursingHouse, error)
	GetNhByID(id int) (entities.NursingHouse, error)
	UpdateNhByID(id int,nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	DeleteNhByID(id int) error
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

func (u *NhUseCaseImpl) GetActiveNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetActiveNh()
}

func (u *NhUseCaseImpl) GetNhByID(id int) (entities.NursingHouse, error) {
	return u.nhrepo.GetNhByID(id)
}

func (u *NhUseCaseImpl) UpdateNhByID(id int,nursingHouse entities.NursingHouse) (entities.NursingHouse, error) {
	existingNh, err := u.nhrepo.GetNhByID(id)
	if err != nil {
		return entities.NursingHouse{}, err
	}
	
	existingNh.Name = nursingHouse.Name
	existingNh.Province = nursingHouse.Province
	existingNh.Address = nursingHouse.Address
	existingNh.Price = nursingHouse.Price
	existingNh.Google_map = nursingHouse.Google_map
	existingNh.Phone_number = nursingHouse.Phone_number
	existingNh.Web_site = nursingHouse.Web_site
	existingNh.Time = nursingHouse.Time
	existingNh.Status = nursingHouse.Status

	updatedNh, err := u.nhrepo.UpdateNhByID(existingNh)
	if err != nil {
		return entities.NursingHouse{}, err
	}
	return updatedNh, nil
}  		

func (u *NhUseCaseImpl) DeleteNhByID(id int) error {
	return u.nhrepo.DeleteNhByID(id)
}