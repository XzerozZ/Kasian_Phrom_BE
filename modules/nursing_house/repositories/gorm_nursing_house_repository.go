package repositories

import (
	"errors"
	
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

func (r *GormNhRepository) GetAllNh() ([]entities.NursingHouse, error) {
	var nhList []entities.NursingHouse
	if err := r.db.Find(&nhList).Error; err != nil {
		return nil, err
	}
	return nhList, nil
}

func (r *GormNhRepository) GetNhByID(id string) (entities.NursingHouse, error) {
	var nh entities.NursingHouse
	if err := r.db.First(&nh, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nh, errors.New("nursing house not found")
		}
		return nh, err
	}
	return nh, nil
}