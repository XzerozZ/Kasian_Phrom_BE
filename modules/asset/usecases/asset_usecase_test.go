package usecases_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAsset_Success(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		TotalCost: 100000.0,
		EndYear:   nextYear,
		UserID:    "user123",
	}

	expectedAsset := asset
	expectedAsset.ID = "ASSET001"
	expectedAsset.Status = "In_Progress"
	expectedAsset.MonthlyExpenses = 100000.0 / float64((currentYear+1-currentYear)*12) // Simple calculation
	expectedAsset.LastCalculatedMonth = int(time.Now().Month())

	mockAssetRepo.On("GetAssetNextID").Return("ASSET001", nil)
	mockAssetRepo.On("CreateAsset", mock.AnythingOfType("*entities.Asset")).Return(&expectedAsset, nil)

	result, err := assetUseCase.CreateAsset(asset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ASSET001", result.ID)
	assert.Equal(t, "In_Progress", result.Status)
	assert.Equal(t, int(time.Now().Month()), result.LastCalculatedMonth)
	mockAssetRepo.AssertExpectations(t)
}

func TestCreateAsset_InvalidTotalCost(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		TotalCost: 0.0,
		EndYear:   nextYear,
		UserID:    "user123",
	}

	result, err := assetUseCase.CreateAsset(asset)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "totalcost must be greater than zero", err.Error())
	mockAssetRepo.AssertNotCalled(t, "CreateAsset")
}

func TestCreateAsset_InvalidEndYear(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	lastYear := strconv.Itoa(currentYear - 1)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		TotalCost: 100000.0,
		EndYear:   lastYear,
		UserID:    "user123",
	}

	result, err := assetUseCase.CreateAsset(asset)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "end year must be greater than or equal to current year", err.Error())
	mockAssetRepo.AssertNotCalled(t, "CreateAsset")
}

func TestGetAssetByID_Success(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	expectedAsset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(expectedAsset, nil)

	result, err := assetUseCase.GetAssetByID("ASSET001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAsset, result)
	mockAssetRepo.AssertExpectations(t)
}

func TestGetAssetByID_NotFound(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	mockAssetRepo.On("GetAssetByID", "ASSET999").Return(nil, errors.New("asset not found"))

	result, err := assetUseCase.GetAssetByID("ASSET999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "asset not found", err.Error())
	mockAssetRepo.AssertExpectations(t)
}

func TestGetAssetByUserID_Success(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)
	currentMonth := int(time.Now().Month())

	assets := []entities.Asset{
		{
			ID:                  "ASSET001",
			Name:                "Test Asset 1",
			Type:                "Property",
			TotalCost:           100000.0,
			CurrentMoney:        10000.0,
			Status:              "In_Progress",
			EndYear:             nextYear,
			LastCalculatedMonth: currentMonth - 1,
			UserID:              "user123",
		},
		{
			ID:                  "ASSET002",
			Name:                "Test Asset 2",
			Type:                "Vehicle",
			TotalCost:           50000.0,
			CurrentMoney:        5000.0,
			Status:              "In_Progress",
			EndYear:             nextYear,
			LastCalculatedMonth: currentMonth,
			UserID:              "user123",
		},
	}

	updatedAsset := assets[0]
	updatedAsset.LastCalculatedMonth = currentMonth
	updatedAsset.MonthlyExpenses = 100000.0 / float64((currentYear+1-currentYear)*12)

	mockAssetRepo.On("GetAssetByUserID", "user123").Return(assets, nil)
	mockAssetRepo.On("UpdateAssetByID", mock.MatchedBy(func(a *entities.Asset) bool {
		return a.ID == "ASSET001"
	})).Return(&updatedAsset, nil)

	results, err := assetUseCase.GetAssetByUserID("user123")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, currentMonth, results[0].LastCalculatedMonth)
	mockAssetRepo.AssertExpectations(t)
}

func TestUpdateAssetByID_Success(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)
	currentMonth := int(time.Now().Month())

	existingAsset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Test Asset",
		Type:                "Property",
		TotalCost:           100000.0,
		CurrentMoney:        10000.0,
		Status:              "In_Progress",
		EndYear:             nextYear,
		MonthlyExpenses:     8333.33,
		LastCalculatedMonth: currentMonth - 1,
		UserID:              "user123",
	}

	updateRequest := entities.Asset{
		Name:      "Updated Asset",
		Type:      "Investment",
		TotalCost: 150000.0,
		EndYear:   nextYear,
		Status:    "In_Progress",
	}

	updatedAsset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Updated Asset",
		Type:                "Investment",
		TotalCost:           150000.0,
		CurrentMoney:        10000.0,
		Status:              "In_Progress",
		EndYear:             nextYear,
		MonthlyExpenses:     12500.0,
		LastCalculatedMonth: currentMonth,
		UserID:              "user123",
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(existingAsset, nil)
	mockAssetRepo.On("UpdateAssetByID", mock.MatchedBy(func(a *entities.Asset) bool {
		return a.ID == "ASSET001" && a.Name == "Updated Asset"
	})).Return(updatedAsset, nil)

	result, err := assetUseCase.UpdateAssetByID("ASSET001", updateRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Asset", result.Name)
	assert.Equal(t, "Investment", result.Type)
	assert.Equal(t, 150000.0, result.TotalCost)
	assert.Equal(t, currentMonth, result.LastCalculatedMonth)
	mockAssetRepo.AssertExpectations(t)
}

func TestUpdateAssetByID_InvalidTotalCost(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	existingAsset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      nextYear,
		UserID:       "user123",
	}

	updateRequest := entities.Asset{
		Name:      "Updated Asset",
		Type:      "Investment",
		TotalCost: 0.0,
		EndYear:   nextYear,
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(existingAsset, nil)

	result, err := assetUseCase.UpdateAssetByID("ASSET001", updateRequest)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "totalcost must be greater than zero", err.Error())
	mockAssetRepo.AssertNotCalled(t, "UpdateAssetByID")
}

func TestUpdateAssetByID_PausedStatus(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	existingAsset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Test Asset",
		Type:                "Property",
		TotalCost:           100000.0,
		CurrentMoney:        10000.0,
		Status:              "In_Progress",
		EndYear:             nextYear,
		MonthlyExpenses:     8333.33,
		LastCalculatedMonth: int(time.Now().Month()),
		UserID:              "user123",
	}

	updateRequest := entities.Asset{
		Name:      "Updated Asset",
		Type:      "Investment",
		TotalCost: 150000.0,
		EndYear:   nextYear,
		Status:    "Paused",
	}

	updatedAsset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Updated Asset",
		Type:                "Investment",
		TotalCost:           150000.0,
		CurrentMoney:        10000.0,
		Status:              "Paused",
		EndYear:             nextYear,
		MonthlyExpenses:     0.0,
		LastCalculatedMonth: 0,
		UserID:              "user123",
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(existingAsset, nil)
	mockAssetRepo.On("UpdateAssetByID", mock.MatchedBy(func(a *entities.Asset) bool {
		return a.ID == "ASSET001" && a.Status == "Paused"
	})).Return(updatedAsset, nil)

	result, err := assetUseCase.UpdateAssetByID("ASSET001", updateRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Paused", result.Status)
	assert.Equal(t, 0.0, result.MonthlyExpenses)
	assert.Equal(t, 0, result.LastCalculatedMonth)
	mockAssetRepo.AssertExpectations(t)
}

func TestDeleteAssetByID_Success(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "retirementplan",
			Name:   "My Retirement",
			Amount: 5000.0,
		},
		{
			Type:   "asset",
			Name:   "Second Asset",
			Amount: 3000.0,
		},
	}

	targetAsset := &entities.Asset{
		ID:           "ASSET002",
		Name:         "Second Asset",
		Type:         "Investment",
		TotalCost:    50000.0,
		CurrentMoney: 2000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	updatedTargetAsset := *targetAsset
	updatedTargetAsset.CurrentMoney = 5000.0

	retirementPlan := &entities.RetirementPlan{
		ID:                     "RET001",
		PlanName:               "My Retirement",
		CurrentSavings:         20000.0,
		UserID:                 "user123",
		CurrentTotalInvestment: 30000.0,
		LastRequiredFunds:      100000.0,
		Status:                 "In_Progress",
	}

	updatedRetirementPlan := *retirementPlan
	updatedRetirementPlan.CurrentSavings = 25000.0

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("CreateHistory", mock.AnythingOfType("*entities.History")).Return(&entities.History{}, nil).Times(3)
	mockAssetRepo.On("FindAssetByNameandUserID", "Second Asset", "user123").Return(targetAsset, nil)
	mockAssetRepo.On("UpdateAssetByID", mock.MatchedBy(func(a *entities.Asset) bool {
		return a.ID == "ASSET002" && a.CurrentMoney == 5000.0
	})).Return(&updatedTargetAsset, nil)
	mockRetirementRepo.On("GetRetirementByUserID", "user123").Return(retirementPlan, nil)
	mockRetirementRepo.On("UpdateRetirementPlan", mock.MatchedBy(func(p *entities.RetirementPlan) bool {
		return p.ID == "RET001" && p.CurrentSavings == 25000.0
	})).Return(&updatedRetirementPlan, nil)
	mockAssetRepo.On("DeleteAssetByID", "ASSET001").Return(nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockRetirementRepo.AssertExpectations(t)
}

func TestDeleteAssetByID_TransferAmountExceedsCurrentMoney(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "retirementplan",
			Name:   "My Retirement",
			Amount: 8000.0,
		},
		{
			Type:   "asset",
			Name:   "Second Asset",
			Amount: 5000.0,
		},
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.Error(t, err)
	assert.Equal(t, "transfer amount exceeds asset's current money", err.Error())
	mockAssetRepo.AssertNotCalled(t, "DeleteAssetByID")
}

func TestDeleteAssetByID_TransferToCompletedAsset(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "asset",
			Name:   "Completed Asset",
			Amount: 5000.0,
		},
	}

	completedAsset := &entities.Asset{
		ID:           "ASSET002",
		Name:         "Completed Asset",
		Type:         "Investment",
		TotalCost:    50000.0,
		CurrentMoney: 50000.0,
		Status:       "Completed",
		EndYear:      "2026",
		UserID:       "user123",
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockAssetRepo.On("FindAssetByNameandUserID", "Completed Asset", "user123").Return(completedAsset, nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.Error(t, err)
	assert.Equal(t, "cannot update completed or paused asset", err.Error())
	mockAssetRepo.AssertNotCalled(t, "DeleteAssetByID")
}

func TestUpdateAssetStatus_InProgress(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 50000.0,
		Status:       "In_Progress",
		EndYear:      nextYear,
		UserID:       "user123",
	}

	err := assetUseCase.UpdateAssetStatus(asset, currentYear)

	assert.NoError(t, err)
	assert.Equal(t, "In_Progress", asset.Status)
}

func TestUpdateAssetStatus_Completed(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 100000.0,
		Status:       "In_Progress",
		EndYear:      nextYear,
		UserID:       "user123",
	}

	err := assetUseCase.UpdateAssetStatus(asset, currentYear)

	assert.NoError(t, err)
	assert.Equal(t, "Completed", asset.Status)
}

func TestUpdateAssetStatus_Paused(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()

	asset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Test Asset",
		Type:                "Property",
		TotalCost:           100000.0,
		CurrentMoney:        50000.0,
		Status:              "In_Progress",
		EndYear:             strconv.Itoa(currentYear),
		MonthlyExpenses:     5000.0,
		LastCalculatedMonth: 5,
		UserID:              "user123",
	}

	mockNotiRepo.On("CreateNotification", mock.AnythingOfType("*entities.Notification")).Return(nil)

	err := assetUseCase.UpdateAssetStatus(asset, currentYear)

	assert.NoError(t, err)
	assert.Equal(t, "Paused", asset.Status)
	assert.Equal(t, 0, asset.LastCalculatedMonth)
	assert.Equal(t, 0.0, asset.MonthlyExpenses)
	mockNotiRepo.AssertExpectations(t)
}

func TestUpdateAssetStatus_AlreadyPaused(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := &entities.Asset{
		ID:                  "ASSET001",
		Name:                "Test Asset",
		Type:                "Property",
		TotalCost:           100000.0,
		CurrentMoney:        50000.0,
		Status:              "Paused",
		EndYear:             nextYear,
		MonthlyExpenses:     5000.0,
		LastCalculatedMonth: 5,
		UserID:              "user123",
	}

	err := assetUseCase.UpdateAssetStatus(asset, currentYear)

	assert.NoError(t, err)
	assert.Equal(t, "Paused", asset.Status)
	assert.Equal(t, 0, asset.LastCalculatedMonth)
	assert.Equal(t, 0.0, asset.MonthlyExpenses)
}

func TestDeleteAssetByID_TransferToHouse(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "house",
			Name:   "Nursing Home",
			Amount: 5000.0,
		},
	}

	selectedHouse := &entities.SelectedHouse{
		NursingHouseID: "NH001",
		UserID:         "user123",
		Status:         "In_Progress",
		CurrentMoney:   20000.0,
		NursingHouse: entities.NursingHouse{
			ID:    "NH001",
			Name:  "Nursing Home",
			Price: 1000.0,
		},
	}

	updatedSelectedHouse := *selectedHouse
	updatedSelectedHouse.CurrentMoney = 25000.0

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("GetSelectedHouse", "user123").Return(selectedHouse, nil)
	mockUserRepo.On("CreateHistory", mock.AnythingOfType("*entities.History")).Return(&entities.History{}, nil).Times(2)
	mockUserRepo.On("UpdateSelectedHouse", mock.MatchedBy(func(h *entities.SelectedHouse) bool {
		return h.CurrentMoney == 25000.0
	})).Return(&updatedSelectedHouse, nil)
	mockAssetRepo.On("DeleteAssetByID", "ASSET001").Return(nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteAssetByID_TransferToCompletedHouse(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "house",
			Name:   "Completed House",
			Amount: 5000.0,
		},
	}

	completedHouse := &entities.SelectedHouse{
		NursingHouseID: "00001",
		UserID:         "user123",
		Status:         "Completed",
		CurrentMoney:   50000.0,
		NursingHouse: entities.NursingHouse{
			ID:    "00001",
			Name:  "Completed House",
			Price: 1000.0,
		},
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("GetSelectedHouse", "user123").Return(completedHouse, nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.Error(t, err)
	assert.Equal(t, "cannot update completed nursing house", err.Error())
	mockAssetRepo.AssertNotCalled(t, "DeleteAssetByID")
}

func TestDeleteAssetByID_HouseCompletedAfterTransfer(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "house",
			Name:   "Almost Complete House",
			Amount: 5000.0,
		},
	}

	selectedHouse := &entities.SelectedHouse{
		NursingHouseID: "NH001",
		UserID:         "user123",
		Status:         "In_Progress",
		CurrentMoney:   295000.0,
		NursingHouse: entities.NursingHouse{
			ID:    "NH001",
			Name:  "Almost Complete House",
			Price: 1000.0, // Monthly price
		},
	}

	updatedSelectedHouse := *selectedHouse
	updatedSelectedHouse.CurrentMoney = 300000.0
	updatedSelectedHouse.Status = "Completed"
	updatedSelectedHouse.MonthlyExpenses = 0
	updatedSelectedHouse.LastCalculatedMonth = 0

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("GetSelectedHouse", "user123").Return(selectedHouse, nil)
	mockUserRepo.On("CreateHistory", mock.AnythingOfType("*entities.History")).Return(&entities.History{}, nil).Times(2)
	mockUserRepo.On("UpdateSelectedHouse", mock.MatchedBy(func(h *entities.SelectedHouse) bool {
		return h.CurrentMoney == 300000.0 && h.Status == "Completed"
	})).Return(&updatedSelectedHouse, nil)
	mockNotiRepo.On("CreateNotification", mock.AnythingOfType("*entities.Notification")).Return(nil)
	mockAssetRepo.On("DeleteAssetByID", "ASSET001").Return(nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockNotiRepo.AssertExpectations(t)
}

func TestDeleteAssetByID_RetirementPlanCompletedAfterTransfer(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "retirementplan",
			Name:   "Almost Complete Retirement",
			Amount: 5000.0,
		},
	}

	retirementPlan := &entities.RetirementPlan{
		ID:                     "RET001",
		PlanName:               "Almost Complete Retirement",
		CurrentSavings:         45000.0,
		CurrentTotalInvestment: 50000.0,
		LastRequiredFunds:      100000.0,
		Status:                 "In_Progress",
		UserID:                 "user123",
	}

	updatedRetirementPlan := *retirementPlan
	updatedRetirementPlan.CurrentSavings = 50000.0
	updatedRetirementPlan.Status = "Completed"
	updatedRetirementPlan.LastMonthlyExpenses = 0

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("CreateHistory", mock.AnythingOfType("*entities.History")).Return(&entities.History{}, nil).Times(2)
	mockRetirementRepo.On("GetRetirementByUserID", "user123").Return(retirementPlan, nil)
	mockRetirementRepo.On("UpdateRetirementPlan", mock.MatchedBy(func(r *entities.RetirementPlan) bool {
		return r.CurrentSavings == 50000.0 && r.Status == "Completed"
	})).Return(&updatedRetirementPlan, nil)
	mockNotiRepo.On("CreateNotification", mock.AnythingOfType("*entities.Notification")).Return(nil)
	mockAssetRepo.On("DeleteAssetByID", "ASSET001").Return(nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
	mockRetirementRepo.AssertExpectations(t)
	mockNotiRepo.AssertExpectations(t)
}

func TestDeleteAssetByID_InvalidTransferType(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	user := &entities.User{
		ID: "user123",
		RetirementPlan: entities.RetirementPlan{
			RetirementAge:  60,
			ExpectLifespan: 85,
		},
	}

	transferRequests := []entities.TransferRequest{
		{
			Type:   "invalid_type",
			Name:   "Some Entity",
			Amount: 5000.0,
		},
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "user123").Return(user, nil)
	mockUserRepo.On("CreateHistory", mock.AnythingOfType("*entities.History")).Return(&entities.History{}, nil)
	mockAssetRepo.On("DeleteAssetByID", "ASSET001").Return(nil)

	err := assetUseCase.DeleteAssetByID("ASSET001", "user123", transferRequests)

	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAssetUseCaseImpl_GetAssetByID_NotFound(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	mockAssetRepo.On("GetAssetByID", "NOTFOUND").Return(nil, errors.New("asset not found"))

	result, err := assetUseCase.GetAssetByID("NOTFOUND")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "asset not found", err.Error())
}

func TestUpdateAssetByID_AssetNotFound(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	updateRequest := entities.Asset{
		Name:      "Updated Asset",
		Type:      "Investment",
		TotalCost: 150000.0,
		EndYear:   "2026",
	}

	mockAssetRepo.On("GetAssetByID", "NOTFOUND").Return(nil, errors.New("asset not found"))

	result, err := assetUseCase.UpdateAssetByID("NOTFOUND", updateRequest)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "asset not found", err.Error())
}

func TestDeleteAssetByID_AssetNotFound(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	mockAssetRepo.On("GetAssetByID", "NOTFOUND").Return(nil, errors.New("asset not found"))

	err := assetUseCase.DeleteAssetByID("NOTFOUND", "user123", []entities.TransferRequest{})

	assert.Error(t, err)
	assert.Equal(t, "asset not found", err.Error())
}

func TestDeleteAssetByID_UserNotFound(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	asset := &entities.Asset{
		ID:           "ASSET001",
		Name:         "Test Asset",
		Type:         "Property",
		TotalCost:    100000.0,
		CurrentMoney: 10000.0,
		Status:       "In_Progress",
		EndYear:      "2026",
		UserID:       "user123",
	}

	mockAssetRepo.On("GetAssetByID", "ASSET001").Return(asset, nil)
	mockUserRepo.On("GetUserByID", "NOTFOUND").Return((*entities.User)(nil), errors.New("user not found"))

	err := assetUseCase.DeleteAssetByID("ASSET001", "NOTFOUND", []entities.TransferRequest{})

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestAssetUseCaseImpl_CreateAsset_FailedToGetNextID(t *testing.T) {
	mockAssetRepo := new(mocks.MockAssetRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNhRepo := new(mocks.MockNhRepository)
	mockRetirementRepo := new(mocks.MockRetirementRepository)
	mockNotiRepo := new(mocks.MockNotiRepository)

	assetUseCase := usecases.NewAssetUseCase(mockAssetRepo, mockUserRepo, mockNhRepo, mockRetirementRepo, mockNotiRepo)

	currentYear := time.Now().Year()
	nextYear := strconv.Itoa(currentYear + 1)

	asset := entities.Asset{
		Name:      "Test Asset",
		Type:      "Property",
		TotalCost: 100000.0,
		EndYear:   nextYear,
		UserID:    "user123",
	}

	mockAssetRepo.On("GetAssetNextID").Return("", errors.New("failed to get next ID"))

	result, err := assetUseCase.CreateAsset(asset)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to get next ID", err.Error())
}
