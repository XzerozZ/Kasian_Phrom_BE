package usecases_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) GetAssetNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAssetRepository) CreateAsset(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByID(id string) (*entities.Asset, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) UpdateAssetByID(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) DeleteAssetByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateAsset(t *testing.T) {
	mockRepo := new(MockAssetRepository)
	useCase := usecases.NewAssetUseCase(mockRepo)

	mockRepo.On("GetAssetNextID").Return("1", nil)
	mockRepo.On("CreateAsset", mock.Anything).Return(&entities.Asset{ID: "1", Name: "Test Asset", Type: "Property", EndYear: strconv.Itoa(time.Now().Year() + 1), TotalCost: 1000}, nil)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		EndYear:   strconv.Itoa(time.Now().Year() + 1),
		TotalCost: 1000,
	}

	createdAsset, err := useCase.CreateAsset(asset)

	assert.NoError(t, err)
	assert.NotNil(t, createdAsset)
	assert.Equal(t, "1", createdAsset.ID)
}

func TestCreateAsset_InvalidTotalCost(t *testing.T) {
	mockRepo := new(MockAssetRepository)
	useCase := usecases.NewAssetUseCase(mockRepo)
	mockRepo.On("GetAssetNextID").Return("1", nil)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		EndYear:   strconv.Itoa(time.Now().Year() + 1),
		TotalCost: 0,
	}

	createdAsset, err := useCase.CreateAsset(asset)

	assert.Error(t, err)
	assert.Nil(t, createdAsset)
	assert.Equal(t, "totalcost must be greater than zero", err.Error())
}

func TestGetAssetByID(t *testing.T) {
	mockRepo := new(MockAssetRepository)
	useCase := usecases.NewAssetUseCase(mockRepo)

	mockRepo.On("GetAssetByID", "1").Return(&entities.Asset{ID: "1", Name: "Test Asset"}, nil)

	asset, err := useCase.GetAssetByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "1", asset.ID)
}

func TestUpdateAssetByID(t *testing.T) {
	mockRepo := new(MockAssetRepository)
	useCase := usecases.NewAssetUseCase(mockRepo)

	existingAsset := entities.Asset{ID: "1", Name: "Old Asset", Type: "Property", EndYear: strconv.Itoa(time.Now().Year() + 1), TotalCost: 500}
	updatedAsset := entities.Asset{ID: "1", Name: "Updated Asset", Type: "Property", EndYear: strconv.Itoa(time.Now().Year() + 2), TotalCost: 1000}

	mockRepo.On("GetAssetByID", "1").Return(&existingAsset, nil)
	mockRepo.On("UpdateAssetByID", mock.Anything).Return(&updatedAsset, nil)

	result, err := useCase.UpdateAssetByID("1", updatedAsset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Asset", result.Name)
}

func TestDeleteAssetByID(t *testing.T) {
	mockRepo := new(MockAssetRepository)
	useCase := usecases.NewAssetUseCase(mockRepo)

	mockRepo.On("GetAssetByID", "1").Return(&entities.Asset{ID: "1"}, nil)
	mockRepo.On("DeleteAssetByID", "1").Return(nil)

	err := useCase.DeleteAssetByID("1")

	assert.NoError(t, err)
}
