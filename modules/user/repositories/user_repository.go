package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"github.com/google/uuid"
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
	FindUserByEmail(email string) (entities.User, error)
	FindUserByID(id uuid.UUID) (*entities.User, error)
	GetRoleByName(name string) (entities.Role, error)
}

func (r *GormUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user , nil
}

func (r *GormUserRepository) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *GormUserRepository) FindUserByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetRoleByName(name string) (entities.Role, error) {
	var role entities.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return entities.Role{}, err
	}
	return role, nil
}