package repositories

import (
	"fmt"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormAssRepository struct {
	db *gorm.DB
}

func NewGormAssRepository(db *gorm.DB) *GormAssRepository {
	return &GormAssRepository{db: db}
}


type AssRepository interface {
	CreateAss(asset *entities.Asset) (*entities.Asset, error)
	GetAssByID(id string) (*entities.Asset, error)
	GetAssByUsername(username string) ([]entities.Asset, error)
	GetAssNextID() (string, error)
	UpdateAssByID(asset *entities.Asset) (*entities.Asset, error)
	DeleteAssByID(id string) error
}

func (r *GormAssRepository) CreateAss(asset *entities.Asset) (*entities.Asset, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(asset).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetAssByID(asset.ID)
}

func (r *GormAssRepository) GetAssByID(id string) (*entities.Asset, error) {
	var asset entities.Asset
	if err := r.db.First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}


func (r *GormAssRepository) GetAssByUsername(username string) ([]entities.Asset, error) {
	var assets []entities.Asset
	err := r.db.
		Joins("JOIN users ON users.id = assets.user_id").
		Where("users.username = ?", username).
		Find(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *GormAssRepository) GetAssNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.Asset{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
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

func (r *GormAssRepository) UpdateAssByID(asset *entities.Asset) (*entities.Asset, error) {
	if err := r.db.Save(&asset).Error; err != nil {
		return nil, err
	}

	return r.GetAssByID(asset.ID)
}

func (r *GormAssRepository) DeleteAssByID(id string) error {
    return r.db.Where("id = ?", id).Delete(&entities.Asset{}).Error
}