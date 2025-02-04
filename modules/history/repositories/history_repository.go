package repositories

import (
	"time"

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
	GetHistoryInRange(userID string, startDate, endDate time.Time) ([]entities.History, error)
	GetUserDepositsInRange(userID string, startDate, endDate time.Time) ([]entities.History, error)
	GetUserHistoryByMonth(userID string) (map[string]float64, error)
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

func (r *GormHistoryRepository) GetHistoryInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ? AND track_date BETWEEN ? AND ?", userID, startDate, endDate).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) GetUserDepositsInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ? AND method = ? AND track_date BETWEEN ? AND ?", userID, "deposit", startDate, endDate).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) GetUserHistoryByMonth(userID string) (map[string]float64, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, err
	}

	historyByMonth := make(map[string]float64)
	for _, history := range histories {
		monthKey := history.TrackDate.Format("2006-01")
		if history.Method == "deposit" {
			historyByMonth[monthKey] += history.Money
		} else if history.Method == "withdraw" {
			historyByMonth[monthKey] -= history.Money
		}
	}

	return historyByMonth, nil
}
