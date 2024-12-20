package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type NhUse struct {
	NhRepo entities.NhRepository
}

func NewNhUseCase(nhRepo entities.NhRepository) entities.NhUsecase {
	return &NhUse{
		NhRepo: nhRepo
	}
}

func (s *NhUse) CreateNh(req *Nursing_House) (*Nursing_House, error) {
	return s.NhRepo.CreateNh(nursing_house)
}