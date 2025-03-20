package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) CreateAsset(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByID(id string) (*entities.Asset, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) UpdateAssetByID(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) DeleteAssetByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAssetRepository) GetAssetNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAssetRepository) FindAssetByNameandUserID(name string, userID string) (*entities.Asset, error) {
	args := m.Called(name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}
