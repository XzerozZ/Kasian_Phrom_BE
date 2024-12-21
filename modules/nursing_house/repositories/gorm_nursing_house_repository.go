package repositories

import (
	"gorm.io/gorm"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type GormNhRepository struct {
  	db *gorm.DB
}

type NhRepository interface {
	Create(nursingHouse *entities.NursingHouse) error
}

func NewGormNhRepository(db *gorm.DB) *GormNhRepository {
 	return &GormNhRepository{db: db}
}

func (r *GormNhRepository) Create(nursingHouse *entities.NursingHouse) error {
	return r.db.Create(nursingHouse).Error
}