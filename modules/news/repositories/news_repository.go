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
	CreateNews(news *entities.News) (*entities.News, error)
	GetAllNews() ([]entities.News, error)
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
	UpdateNewsByID(news *entities.News) (*entities.News, error)
	DeleteDialog(id string) error
	DeleteNewsByID(id string) error
}

func (r *GormNewsRepository) CreateNews(news *entities.News) (*entities.News, error) {
	if err := r.db.Create(&news).Error; err != nil {
		return nil, err
	}

	return r.GetNewsByID(news.ID)
}

func(r *GormNewsRepository) GetAllNews() ([]entities.News, error){
	var news []entities.News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}

	return news , nil
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

func (r *GormNewsRepository) UpdateNewsByID(news *entities.News) (*entities.News, error) {
	if err := r.db.Save(&news).Error; err != nil {
		return nil, err
	}

	return r.GetNewsByID(news.ID)
}

func (r *GormNewsRepository) DeleteDialog(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Dialog{}).Error
}

func (r *GormNewsRepository) DeleteNewsByID(id string) error {
    return r.db.Where("id = ?", id).Delete(&entities.News{}).Error
}
