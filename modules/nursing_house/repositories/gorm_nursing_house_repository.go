package repositories

import (
	"gorm.io/gorm"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type GormNhRepository struct {
  	db *gorm.DB
}

type NhRepository interface {
	CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetNhByID(id int) (entities.NursingHouse, error)
}

func NewGormNhRepository(db *gorm.DB) *GormNhRepository {
 	return &GormNhRepository{db: db}
}

func (r *GormNhRepository) CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error) {
	if err := r.db.Create(&nursingHouse).Error; err != nil {
		return entities.NursingHouse{}, err
	}
	return nursingHouse , nil
}

func (r *GormNhRepository) GetAllNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Find(&nursingHouses).Error; err != nil {
		return nil, err
	}
	return nursingHouses , nil
}

func (r *GormNhRepository) GetNhByID(id int) (entities.NursingHouse, error) {
	var nursingHouse entities.NursingHouse
	if err := r.db.First(&nursingHouse, id).Error; err != nil {
		return entities.NursingHouse{}, err
	}
	return nursingHouse , nil
}