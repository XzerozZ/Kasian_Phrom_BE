package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockAssetUseCase struct {
	mock.Mock
}

func (m *MockAssetUseCase) CreateAsset(asset entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	if result := args.Get(0); result != nil {
		return result.(*entities.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAssetUseCase) GetAssetByID(id string) (*entities.Asset, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entities.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAssetUseCase) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	args := m.Called(userID)
	if result := args.Get(0); result != nil {
		return result.([]entities.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAssetUseCase) UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error) {
	args := m.Called(id, asset)
	if result := args.Get(0); result != nil {
		return result.(*entities.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAssetUseCase) DeleteAssetByID(id string, userID string, transfers []entities.TransferRequest) error {
	args := m.Called(id, userID, transfers)
	return args.Error(0)
}
