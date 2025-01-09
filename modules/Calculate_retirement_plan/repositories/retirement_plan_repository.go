package repositories

import (
	"fmt"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormRetRepository struct {
	db *gorm.DB
}

func NewGormRetRepository(db *gorm.DB) *GormRetRepository {
	return &GormRetRepository{db: db}
}

type RetRepository interface {
	CreateRet(retirementPlan *entities.RetirementPlan) (*entities.RetirementPlan, error)
	GetRetByID(id string) (*entities.RetirementPlan, error)
	GetRetNextID() (string, error)
}

func (r *GormRetRepository) CreateRet(retirementPlan *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	fmt.Println("เข้ามาใน repo")
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(retirementPlan).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetRetByID(retirementPlan.ID)
}

func (r *GormRetRepository) GetRetByID(id string) (*entities.RetirementPlan, error) {
	var retirementPlan entities.RetirementPlan
	if err := r.db.First(&retirementPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &retirementPlan, nil
}

func (r *GormRetRepository) GetRetNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.RetirementPlan{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
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
