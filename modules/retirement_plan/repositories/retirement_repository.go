package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormRetirementRepository struct {
	db *gorm.DB
}

func NewGormRetirementRepository(db *gorm.DB) *GormRetirementRepository {
	return &GormRetirementRepository{db: db}
}

type RetirementRepository interface {
	CreateRetirement(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error)
	GetRetirementByID(id string) (*entities.RetirementPlan, error)
	GetRetirementByUserID(userID string) (*entities.RetirementPlan, error)
	UpdateRetirementPlan(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error)
}

func (r *GormRetirementRepository) CreateRetirement(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	if err := r.db.Create(&retirement).Error; err != nil {
		return nil, err
	}

	return r.GetRetirementByID(retirement.ID)
}

func (r *GormRetirementRepository) GetRetirementByID(id string) (*entities.RetirementPlan, error) {
	var retirement entities.RetirementPlan
	if err := r.db.First(&retirement, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &retirement, nil
}

func (r *GormRetirementRepository) GetRetirementByUserID(userID string) (*entities.RetirementPlan, error) {
	var retirement entities.RetirementPlan
	if err := r.db.Where("user_id = ?", userID).Find(&retirement).Error; err != nil {
		return nil, err
	}

	return &retirement, nil
}

func (r *GormRetirementRepository) UpdateRetirementPlan(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	if err := r.db.Save(&retirement).Error; err != nil {
		return nil, err
	}

	return r.GetRetirementByID(retirement.ID)
}
