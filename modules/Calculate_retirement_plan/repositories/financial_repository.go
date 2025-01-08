package repositories

import (
	"fmt"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormFinRepository struct {
	db *gorm.DB
}

func NewGormFinRepository(db *gorm.DB) *GormFinRepository {
	return &GormFinRepository{db: db}
}

type FinRepository interface {
	CreateFin(financial *entities.Financial) (*entities.Financial, error)
	GetFinByID(id string) (*entities.Financial, error)
	GetFinByUserID(id string) (*entities.Financial, error)
	GetFinNextID() (string, error)
}

func (r *GormFinRepository) CreateFin(financial *entities.Financial) (*entities.Financial, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(financial).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetFinByID(financial.ID)
}

func (r *GormFinRepository) GetFinByID(financial_id string) (*entities.Financial, error) {
	var financial entities.Financial
	if err := r.db.First(&financial, "financial_id = ?", financial_id).Error; err != nil {
		return nil, err
	}

	return &financial, nil
}

func (r *GormFinRepository) GetFinByUserID(user_id string) (*entities.Financial, error) {
	var financial entities.Financial
	if err := r.db.First(&financial, "user_id = ?", user_id).Error; err != nil {
		return nil, err
	}

	return &financial, nil
}

func (r *GormFinRepository) GetFinNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.Financial{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
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
