package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormHistoryRepository struct {
	db *gorm.DB
}

func NewGormHistoryRepository(db *gorm.DB) *GormHistoryRepository {
	return &GormHistoryRepository{db: db}
}

type HistoryRepository interface {
	CreateHistory(history *entities.History) (*entities.History, error)
	GetHistoryByUserID(userID string) ([]entities.History, error)
}

func (r *GormHistoryRepository) CreateHistory(history *entities.History) (*entities.History, error) {
	if err := r.db.Create(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

func (r *GormHistoryRepository) GetHistoryByUserID(userID string) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}
