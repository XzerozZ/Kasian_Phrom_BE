package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createValidRetirementPlan() entities.RetirementPlan {
	birthDate := time.Now().AddDate(-30, 0, 0).Format("02-01-2006")
	return entities.RetirementPlan{
		UserID:                  "test-user-id",
		BirthDate:               birthDate,
		ExpectLifespan:          80,
		RetirementAge:           60,
		PlanName:                "Test Plan",
		MonthlyIncome:           50000,
		MonthlyExpenses:         30000,
		CurrentSavings:          500000,
		CurrentSavingsReturns:   3.0,
		CurrentTotalInvestment:  1000000,
		InvestmentReturn:        7.0,
		ExpectedInflation:       2.5,
		ExpectedMonthlyExpenses: 40000,
		AnnualExpenseIncrease:   3.0,
		AnnualSavingsReturn:     3.0,
		AnnualInvestmentReturn:  7.0,
	}
}

func TestCreateRetirement_Success(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()

	expectedReturnPlan := retirementPlan
	expectedReturnPlan.ID = "generated-uuid"
	expectedReturnPlan.Status = "In_Progress"

	mockRepo.On("CreateRetirement", mock.AnythingOfType("*entities.RetirementPlan")).Return(&expectedReturnPlan, nil)

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 30, age)
	assert.Equal(t, "In_Progress", result.Status)
	mockRepo.AssertExpectations(t)
}

func TestCreateRetirement_NegativeCurrentSavings(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	retirementPlan.CurrentSavings = -1

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currentSavings must be greater than or equal to zero")
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertNotCalled(t, "CreateRetirement")
}

func TestCreateRetirement_NegativeMonthlyIncome(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	retirementPlan.MonthlyIncome = -1

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "monthlyIncome must be greater than or equal to zero")
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertNotCalled(t, "CreateRetirement")
}

func TestCreateRetirement_ZeroCurrentSavingsReturns(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	retirementPlan.CurrentSavingsReturns = 0

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "monthlyIncome must be greater than zero")
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertNotCalled(t, "CreateRetirement")
}

func TestCreateRetirement_AgeOlderThanRetirementAge(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	retirementPlan.BirthDate = time.Now().AddDate(-65, 0, 0).Format("02-01-2006")

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "age must be less than RetirementAge")
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertNotCalled(t, "CreateRetirement")
}

func TestCreateRetirement_RetirementAgeHigherThanLifespan(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	retirementPlan.RetirementAge = 85
	retirementPlan.ExpectLifespan = 80

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retirementAge must be less than ExpectLifespan")
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertNotCalled(t, "CreateRetirement")
}

func TestCreateRetirement_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	retirementPlan := createValidRetirementPlan()
	expectedError := errors.New("database error")

	mockRepo.On("CreateRetirement", mock.AnythingOfType("*entities.RetirementPlan")).Return(&entities.RetirementPlan{}, expectedError)

	result, age, err := useCase.CreateRetirement(retirementPlan)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, age)
	mockRepo.AssertExpectations(t)
}

func TestGetRetirementByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	expectedPlan := createValidRetirementPlan()
	expectedPlan.ID = "test-id"

	mockRepo.On("GetRetirementByID", "test-id").Return(&expectedPlan, nil)

	result, err := useCase.GetRetirementByID("test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPlan.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetRetirementByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	expectedError := errors.New("record not found")

	mockRepo.On("GetRetirementByID", "non-existent-id").Return(nil, expectedError)

	result, err := useCase.GetRetirementByID("non-existent-id")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetRetirementByUserID_Success(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	expectedPlan := createValidRetirementPlan()
	expectedPlan.UserID = "test-user-id"

	mockRepo.On("GetRetirementByUserID", "test-user-id").Return(&expectedPlan, nil)

	result, err := useCase.GetRetirementByUserID("test-user-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPlan.UserID, result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestGetRetirementByUserID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	expectedError := errors.New("record not found")

	mockRepo.On("GetRetirementByUserID", "non-existent-user-id").Return(nil, expectedError)

	result, err := useCase.GetRetirementByUserID("non-existent-user-id")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRetirementByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	userID := "test-user-id"
	existingPlan := createValidRetirementPlan()
	existingPlan.ID = "test-id"
	existingPlan.UserID = userID
	existingPlan.LastCalculatedMonth = int(time.Now().Month())
	existingPlan.LastRequiredFunds = 10000000
	existingPlan.LastMonthlyExpenses = 35000

	updatedPlan := createValidRetirementPlan()
	updatedPlan.ID = "test-id"
	updatedPlan.UserID = userID
	updatedPlan.MonthlyIncome = 60000
	updatedPlan.ExpectedMonthlyExpenses = 45000

	mockRepo.On("GetRetirementByUserID", userID).Return(&existingPlan, nil)
	mockRepo.On("UpdateRetirementPlan", mock.AnythingOfType("*entities.RetirementPlan")).Return(&updatedPlan, nil)

	result, err := useCase.UpdateRetirementByID(userID, updatedPlan)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRetirementByID_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	userID := "non-existent-user-id"
	updatedPlan := createValidRetirementPlan()
	expectedError := errors.New("user not found")

	mockRepo.On("GetRetirementByUserID", userID).Return(nil, expectedError)

	result, err := useCase.UpdateRetirementByID(userID, updatedPlan)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertNotCalled(t, "UpdateRetirementPlan")
}

func TestUpdateRetirementByID_NegativeMonthlyIncome(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	userID := "test-user-id"
	existingPlan := createValidRetirementPlan()

	updatedPlan := createValidRetirementPlan()
	updatedPlan.MonthlyIncome = -1

	mockRepo.On("GetRetirementByUserID", userID).Return(&existingPlan, nil)

	result, err := useCase.UpdateRetirementByID(userID, updatedPlan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "monthlyIncome must be greater than or equal to zero")
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "UpdateRetirementPlan")
}

func TestUpdateRetirementByID_StatusCompletedWhenFundsSufficient(t *testing.T) {
	mockRepo := new(mocks.MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)

	userID := "test-user-id"
	existingPlan := createValidRetirementPlan()
	existingPlan.ID = "test-id"
	existingPlan.UserID = userID
	existingPlan.LastCalculatedMonth = (int(time.Now().Month()) + 1) % 12
	existingPlan.LastRequiredFunds = 1000000

	updatedPlan := createValidRetirementPlan()
	updatedPlan.ID = "test-id"
	updatedPlan.UserID = userID
	updatedPlan.CurrentSavings = 500000
	updatedPlan.CurrentTotalInvestment = 1000000

	expectedUpdatedPlan := updatedPlan
	expectedUpdatedPlan.Status = "Completed"
	expectedUpdatedPlan.LastMonthlyExpenses = 0

	mockRepo.On("GetRetirementByUserID", userID).Return(&existingPlan, nil)
	mockRepo.On("UpdateRetirementPlan", mock.AnythingOfType("*entities.RetirementPlan")).Return(&expectedUpdatedPlan, nil)

	result, err := useCase.UpdateRetirementByID(userID, updatedPlan)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Completed", result.Status)
	assert.Equal(t, float64(0), result.LastMonthlyExpenses)
	mockRepo.AssertExpectations(t)
}
