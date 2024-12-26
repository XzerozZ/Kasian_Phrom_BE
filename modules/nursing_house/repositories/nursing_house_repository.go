package repositories

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type GormNhRepository struct {
  	db *gorm.DB
}

func NewGormNhRepository(db *gorm.DB) *GormNhRepository {
	return &GormNhRepository{db: db}
}

type NhRepository interface {
	CreateNh(nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetActiveNh() ([]entities.NursingHouse, error)
	GetInactiveNh() ([]entities.NursingHouse, error)
	GetNhByID(id string) (entities.NursingHouse, error)
	GetNhNextID() (string, error)
	UpdateNhByID(nursingHouse entities.NursingHouse) (entities.NursingHouse, error)
	DeleteNhByID(id string) error
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

func (r *GormNhRepository) GetActiveNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Where("status = ?", "Active").Find(&nursingHouses).Error; err != nil {
		return nil, err
	}
	return nursingHouses, nil
}

func (r *GormNhRepository) GetInactiveNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Where("status = ?", "Inactive").Find(&nursingHouses).Error; err != nil {
		return nil, err
	}
	return nursingHouses, nil
}

func (r *GormNhRepository) GetNhByID(id string) (entities.NursingHouse, error) {
	var nursingHouse entities.NursingHouse
	if err := r.db.First(&nursingHouse, id).Error; err != nil {
		return entities.NursingHouse{}, err
	}
	return nursingHouse , nil
}

func (r *GormNhRepository) GetNhNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.NursingHouse{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
		return "", err
	}
	
	maxIDInt := 0
	if maxID != "" {
		maxIDInt, _ = strconv.Atoi(maxID)
	}
	nextID := maxIDInt + 1
	formattedID := fmt.Sprintf("%05d", nextID)
	return formattedID, nil
}

func (r *GormNhRepository) UpdateNhByID(nursingHouse entities.NursingHouse) (entities.NursingHouse, error) {
	if err := r.db.Save(&nursingHouse).Error; err != nil {
		return entities.NursingHouse{}, err
	}
	return nursingHouse, nil
}

func (r *GormNhRepository) DeleteNhByID(id string) error {
	if err := r.db.Delete(&entities.NursingHouse{}, id).Error; err != nil {
		return err
	}
	return nil
}