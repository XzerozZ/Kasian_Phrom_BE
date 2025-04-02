package repositories

import (
	"errors"
	"fmt"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormFavRepository struct {
	db *gorm.DB
}

func NewGormFavRepository(db *gorm.DB) *GormFavRepository {
	return &GormFavRepository{db: db}
}

type FavRepository interface {
	CreateFav(fav *entities.Favorite) error
	GetFavByUserID(userID string) ([]entities.Favorite, error)
	CheckFav(userID string, nursingHouseID string) error
	DeleteFavByID(userID string, nursingHouseID string) error
}

func (r *GormFavRepository) CreateFav(fav *entities.Favorite) error {
	var user entities.User
	var nursingHouse entities.NursingHouse
	if err := r.db.First(&user, "id = ?", fav.UserID).Error; err != nil {
		return fmt.Errorf("user_id not found: %v", err)
	}

	if err := r.db.First(&nursingHouse, "id = ?", fav.NursingHouseID).Error; err != nil {
		return fmt.Errorf("nursing_house_id not found: %v", err)
	}

	if err := r.db.Create(&fav).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormFavRepository) GetFavByUserID(userID string) ([]entities.Favorite, error) {
	var favs []entities.Favorite
	if err := r.db.Preload("NursingHouse.Images").Where("user_id = ?", userID).Find(&favs).Error; err != nil {
		return nil, err
	}

	return favs, nil
}

func (r *GormFavRepository) CheckFav(userID string, nursingHouseID string) error {
	var fav entities.Favorite
	if err := r.db.Where("user_id = ? AND nursing_house_id = ?", userID, nursingHouseID).First(&fav).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not favorited nursing house")
		}

		return err
	}

	return nil
}

func (r *GormFavRepository) DeleteFavByID(userID string, nursingHouseID string) error {
	if err := r.db.Where("user_id = ? AND nursing_house_id = ?", userID, nursingHouseID).Delete(&entities.Favorite{}).Error; err != nil {
		return err
	}

	return nil
}
