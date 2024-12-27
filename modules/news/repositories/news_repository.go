package repositories

import (
	"fmt"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormNewsRepository struct {
  	db *gorm.DB
}

func NewGormNewsRepository(db *gorm.DB) *GormNewsRepository {
	return &GormNewsRepository{db: db}
}

type NewsRepository interface {
	CreateNews(news *entities.News) error
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
}

func (r *GormNewsRepository) CreateNews(news *entities.News) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(news).Error; err != nil {
			return err
		}

		for i := range news.Dialog {
			dialog := &news.Dialog[i]
			if err := tx.Create(dialog).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormNewsRepository) GetNewsByID(id string) (*entities.News, error) {
	var news entities.News
	if err := r.db.Preload("Dialog").First(&news, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &news , nil
}

func (r *GormNewsRepository) GetNewsNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.News{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
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