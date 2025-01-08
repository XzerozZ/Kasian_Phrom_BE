package repositories

import (
	"fmt"
	"strconv"
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
	CreateFav(fav *entities.Favorite) (*entities.Favorite, error)
	GetFavByUserID(userID string) ([]entities.Favorite, error)
	CheckFav(userID string, nursingHouseID string) (*entities.Favorite, error)
	DeleteFavByID(userID string, nursingHouseID string) error
}

func (r *GormFavRepository) CreateFav(fav *entities.Favorite) (*entities.Favorite, error) {
	if err := r.db.Create(fav).Error; err != nil {
		return nil, err
	}
	return fav, nil
}

func (r *GormFavRepository) GetFavByUserID(userID string) ([]entities.Favorite, error) {
	var fav entities.Favorite
	if err := r.db.Where("user_id = ?", id).First(&fav).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fav, nil
}

