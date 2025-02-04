package repositories

import (
	"fmt"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormAssetRepository struct {
	db *gorm.DB
}

func NewGormAssetRepository(db *gorm.DB) *GormAssetRepository {
	return &GormAssetRepository{db: db}
}

type AssetRepository interface {
	CreateAsset(asset *entities.Asset) (*entities.Asset, error)
	GetAssetByID(id string) (*entities.Asset, error)
	GetAssetByUserID(userID string) ([]entities.Asset, error)
	GetAssetNextID() (string, error)
	UpdateAssetByID(asset *entities.Asset) (*entities.Asset, error)
	DeleteAssetByID(id string) error
}

func (r *GormAssetRepository) CreateAsset(asset *entities.Asset) (*entities.Asset, error) {
	if err := r.db.Create(&asset).Error; err != nil {
		return nil, err
	}

	return r.GetAssetByID(asset.ID)
}

func (r *GormAssetRepository) GetAssetByID(id string) (*entities.Asset, error) {
	var asset entities.Asset
	if err := r.db.First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *GormAssetRepository) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	var assets []entities.Asset
	if err := r.db.Where("user_id = ?", userID).Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *GormAssetRepository) GetAssetNextID() (string, error) {
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

func (r *GormAssetRepository) UpdateAssetByID(asset *entities.Asset) (*entities.Asset, error) {
	if err := r.db.Save(&asset).Error; err != nil {
		return nil, err
	}

	return r.GetAssetByID(asset.ID)
}

func (r *GormAssetRepository) DeleteAssetByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Asset{}).Error
}
