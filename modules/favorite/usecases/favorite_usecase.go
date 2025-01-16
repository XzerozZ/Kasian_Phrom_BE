package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/repositories"
)

type FavUseCase interface {
	CreateFav(fav *entities.Favorite) error
	GetFavByUserID(userID string) ([]entities.Favorite, error)
	CheckFav(userID string, nursingHouseID string) error
	DeleteFavByID(userID string, nursingHouseID string) error
}

type FavUseCaseImpl struct {
	favrepo repositories.FavRepository
}

func NewFavUseCase(favrepo repositories.FavRepository) *FavUseCaseImpl {
	return &FavUseCaseImpl{favrepo: favrepo}
}

func (u *FavUseCaseImpl) CreateFav(fav *entities.Favorite) error {
	return u.favrepo.CreateFav(fav)
}

func (u *FavUseCaseImpl) GetFavByUserID(userID string) ([]entities.Favorite, error) {
	return u.favrepo.GetFavByUserID(userID)
}

func (u *FavUseCaseImpl) CheckFav(userID string, nursingHouseID string) error {
	return u.favrepo.CheckFav(userID, nursingHouseID)
}

func (u *FavUseCaseImpl) DeleteFavByID(userID string, nursingHouseID string) error {
	return u.favrepo.DeleteFavByID(userID, nursingHouseID)
}
