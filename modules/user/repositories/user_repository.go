package repositories

import (
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	CreateSelectedHouse(selectedHouse *entities.SelectedHouse) error
	FindUserByEmail(email string) (entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	GetRoleByName(name string) (entities.Role, error)
	UpdateUserByID(user *entities.User) (*entities.User, error)
	CreateOTP(otp *entities.OTP) error
	GetOTPByUserID(userID string) (*entities.OTP, error)
	DeleteOTP(userID string) error

	GetSelectedHouse(userID string) (*entities.SelectedHouse, error)
	UpdateSelectedHouse(selectedHouse *entities.SelectedHouse) (*entities.SelectedHouse, error)

	CreateHistory(history *entities.History) (*entities.History, error)
	GetHistoryByUserID(userID string) ([]entities.History, error)
	GetHistoryInRange(userID string, startDate, endDate time.Time) ([]entities.History, error)
	GetUserDepositsInRange(userID string, startDate, endDate time.Time) ([]entities.History, error)
	GetUserHistoryByMonth(userID string) (map[string]float64, error)
}

func (r *GormUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		selectedHouse := &entities.SelectedHouse{
			UserID:         user.ID,
			NursingHouseID: "00001",
			CurrentMoney:   0.0,
			Status:         "Completed",
		}

		if err := tx.Create(&selectedHouse).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) CreateSelectedHouse(selectedHouse *entities.SelectedHouse) error {
	return r.db.Create(&selectedHouse).Error
}

func (r *GormUserRepository) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *GormUserRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	err := r.db.Preload("Quiz.Risk").Preload("Role").Preload("Assets").Preload("RetirementPlan").Preload("House.NursingHouse.Images").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	var selectedHouse entities.SelectedHouse
	err := r.db.Preload("NursingHouse.Images").Where("user_id = ?", userID).First(&selectedHouse).Error
	if err != nil {
		return nil, err
	}

	return &selectedHouse, nil
}

func (r *GormUserRepository) GetRoleByName(name string) (entities.Role, error) {
	var role entities.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return entities.Role{}, err
	}

	return role, nil
}

func (r *GormUserRepository) UpdateUserByID(user *entities.User) (*entities.User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) UpdateSelectedHouse(selectedHouse *entities.SelectedHouse) (*entities.SelectedHouse, error) {
	if err := r.db.Model(&entities.SelectedHouse{}).Where("user_id = ?", selectedHouse.UserID).Updates(map[string]interface{}{
		"nursing_house_id": selectedHouse.NursingHouseID,
		"current_money":    selectedHouse.CurrentMoney,
		"status":           selectedHouse.Status,
	}).Error; err != nil {
		return nil, err
	}

	return r.GetSelectedHouse(selectedHouse.UserID)
}

func (r *GormUserRepository) CreateOTP(otp *entities.OTP) error {
	if err := r.db.Create(otp).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
	var otp entities.OTP
	if err := r.db.Where("user_id = ?", userID).First(&otp).Error; err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *GormUserRepository) DeleteOTP(userID string) error {
	if err := r.db.Delete(&entities.OTP{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) CreateHistory(history *entities.History) (*entities.History, error) {
	if err := r.db.Create(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

func (r *GormUserRepository) GetHistoryByUserID(userID string) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormUserRepository) GetHistoryInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ? AND track_date BETWEEN ? AND ?", userID, startDate, endDate).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormUserRepository) GetUserDepositsInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("user_id = ? AND method = ? AND track_date BETWEEN ? AND ?", userID, "deposit", startDate, endDate).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormUserRepository) GetUserHistoryByMonth(userID string) (map[string]float64, error) {
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
