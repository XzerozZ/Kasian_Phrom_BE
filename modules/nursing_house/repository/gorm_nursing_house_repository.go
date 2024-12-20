package adapters

import (
	"gorm.io/gorm"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type GormNhRepository struct {
  	db *gorm.DB
}

func NewGormNhRepository(db *gorm.DB) entities.NhRepository {
 	return &GormNhRepository{
		db: db
	}
}

func (r *GormNhRepository) CreateNh(req *Nursing_House) (*Nursing_House, error) {
	if err := r.db.Create(&req).Error; err != nil {
		return nil, err
	}
	return req, nil
}