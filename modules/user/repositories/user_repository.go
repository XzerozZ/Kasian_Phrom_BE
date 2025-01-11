package repositories

import (
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
	GetSelectedHouse(userID string) (*entities.SelectedHouse, error)
	GetRoleByName(name string) (entities.Role, error)
	UpdateUserByID(user *entities.User) (*entities.User, error)
	UpdateSelectedHouse(selectedHouse *entities.SelectedHouse) (*entities.SelectedHouse, error)
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
	err := r.db.Preload("Role").Preload("Assets").Preload("RetirementPlan").Preload("House.NursingHouse").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	var selectedHouse entities.SelectedHouse
	err := r.db.Preload("NursingHouse").Where("user_id = ?", userID).First(&selectedHouse).Error
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
    if err := r.db.Save(&selectedHouse).Error; err != nil {
        return nil, err
    }

    return r.GetSelectedHouse(selectedHouse.UserID)
}