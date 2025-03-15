package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockFavUseCase struct {
	mock.Mock
}

func (m *MockFavUseCase) CreateFav(fav *entities.Favorite) error {
	args := m.Called(fav)
	return args.Error(0)
}

func (m *MockFavUseCase) GetFavByUserID(userID string) ([]entities.Favorite, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Favorite), args.Error(1)
}

func (m *MockFavUseCase) CheckFav(userID, nursingHouseID string) error {
	args := m.Called(userID, nursingHouseID)
	return args.Error(0)
}

func (m *MockFavUseCase) DeleteFavByID(userID, nursingHouseID string) error {
	args := m.Called(userID, nursingHouseID)
	return args.Error(0)
}
