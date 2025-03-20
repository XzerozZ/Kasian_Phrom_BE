package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockFavRepository struct {
	mock.Mock
}

func (m *MockFavRepository) CreateFav(fav *entities.Favorite) error {
	args := m.Called(fav)
	return args.Error(0)
}

func (m *MockFavRepository) GetFavByUserID(userID string) ([]entities.Favorite, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Favorite), args.Error(1)
}

func (m *MockFavRepository) CheckFav(userID string, nursingHouseID string) error {
	args := m.Called(userID, nursingHouseID)
	return args.Error(0)
}

func (m *MockFavRepository) DeleteFavByID(userID string, nursingHouseID string) error {
	args := m.Called(userID, nursingHouseID)
	return args.Error(0)
}
