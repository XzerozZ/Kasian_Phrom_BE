package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRetirementRepository struct {
	mock.Mock
}

func (m *MockRetirementRepository) CreateRetirement(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	args := m.Called(retirement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RetirementPlan), args.Error(1)
}

func (m *MockRetirementRepository) GetRetirementByID(id string) (*entities.RetirementPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RetirementPlan), args.Error(1)
}

func (m *MockRetirementRepository) GetRetirementByUserID(userID string) (*entities.RetirementPlan, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RetirementPlan), args.Error(1)
}

func (m *MockRetirementRepository) UpdateRetirementPlan(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	args := m.Called(retirement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RetirementPlan), args.Error(1)
}

func setupTest() (*usecases.RetirementUseCaseImpl, *MockRetirementRepository) {
	mockRepo := new(MockRetirementRepository)
	useCase := usecases.NewRetirementUseCase(mockRepo)
	return useCase, mockRepo
}

func TestCreateRetirement(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful creation", func(t *testing.T) {
		validRetirement := entities.RetirementPlan{
			BirthDate:              time.Now().AddDate(-30, 0, 0).Format("02-01-2006"),
			RetirementAge:          60,
			ExpectLifespan:         80,
			CurrentSavings:         100000,
			MonthlyIncome:          5000,
			MonthlyExpenses:        3000,
			CurrentSavingsReturns:  5,
			CurrentTotalInvestment: 50000,
			InvestmentReturn:       7,
			ExpectedInflation:      3,
			AnnualExpenseIncrease:  2,
			AnnualSavingsReturn:    4,
			AnnualInvestmentReturn: 6,
		}

		mockRepo.On("CreateRetirement", mock.AnythingOfType("*entities.RetirementPlan")).Return(&validRetirement, nil)

		result, age, err := useCase.CreateRetirement(validRetirement)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 30, age)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid current savings", func(t *testing.T) {
		invalidRetirement := entities.RetirementPlan{
			BirthDate:      time.Now().AddDate(-30, 0, 0).Format("02-01-2000"),
			RetirementAge:  60,
			ExpectLifespan: 80,
			CurrentSavings: -1000,
		}

		result, age, err := useCase.CreateRetirement(invalidRetirement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, age)
		assert.Equal(t, "currentSavings must be greater than or equal to zero", err.Error())
	})

	t.Run("invalid age", func(t *testing.T) {
		invalidRetirement := entities.RetirementPlan{
			BirthDate:              time.Now().AddDate(-30, 0, 0).Format("02-01-2000"),
			RetirementAge:          60,
			ExpectLifespan:         80,
			CurrentSavings:         100000,
			MonthlyIncome:          5000,
			MonthlyExpenses:        3000,
			CurrentSavingsReturns:  5,
			CurrentTotalInvestment: 50000,
			InvestmentReturn:       7,
			ExpectedInflation:      3,
		}

		result, age, err := useCase.CreateRetirement(invalidRetirement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, age)
		assert.Equal(t, "age must be less than RetirementAge", err.Error())
	})

	t.Run("invalid retirement age", func(t *testing.T) {
		invalidRetirement := entities.RetirementPlan{
			BirthDate:              time.Now().AddDate(-30, 0, 0).Format("02-01-2000"),
			RetirementAge:          85,
			ExpectLifespan:         80,
			CurrentSavings:         100000,
			MonthlyIncome:          5000,
			MonthlyExpenses:        3000,
			CurrentSavingsReturns:  5,
			CurrentTotalInvestment: 50000,
			InvestmentReturn:       7,
			ExpectedInflation:      3,
		}

		result, age, err := useCase.CreateRetirement(invalidRetirement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, age)
		assert.Equal(t, "retirementAge must be less than ExpectLifespan", err.Error())
	})
}

func TestGetRetirementByID(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedRetirement := &entities.RetirementPlan{
			ID:            "test-id",
			RetirementAge: 60,
		}

		mockRepo.On("GetRetirementByID", "test-id").Return(expectedRetirement, nil)

		result, err := useCase.GetRetirementByID("test-id")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRetirementByID", "non-existent").Return(nil, errors.New("retirement plan not found"))

		result, err := useCase.GetRetirementByID("non-existent")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateRetirementByID(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful update", func(t *testing.T) {
		existingPlan := &entities.RetirementPlan{
			ID:            "test-id",
			RetirementAge: 60,
		}

		updatedPlan := entities.RetirementPlan{
			BirthDate:              time.Now().AddDate(-30, 0, 0).Format("02-01-2000"),
			RetirementAge:          65,
			ExpectLifespan:         85,
			CurrentSavings:         150000,
			MonthlyIncome:          6000,
			MonthlyExpenses:        3500,
			CurrentSavingsReturns:  5.5,
			CurrentTotalInvestment: 60000,
			InvestmentReturn:       7.5,
			ExpectedInflation:      3.2,
			AnnualExpenseIncrease:  2.5,
			AnnualSavingsReturn:    4.5,
			AnnualInvestmentReturn: 6.5,
		}

		mockRepo.On("GetRetirementByUserID", "test-id").Return(existingPlan, nil)
		mockRepo.On("UpdateRetirementPlan", mock.AnythingOfType("*entities.RetirementPlan")).Return(&updatedPlan, nil)

		result, err := useCase.UpdateRetirementByID("test-id", updatedPlan)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 65, result.RetirementAge)
		mockRepo.AssertExpectations(t)
	})

	t.Run("plan not found", func(t *testing.T) {
		updatedPlan := entities.RetirementPlan{
			RetirementAge: 65,
		}

		mockRepo.On("GetRetirementByUserID", "non-existent").Return(nil, errors.New("retirement plan not found"))

		result, err := useCase.UpdateRetirementByID("non-existent", updatedPlan)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid monthly income", func(t *testing.T) {
		existingPlan := &entities.RetirementPlan{
			ID:            "test-id",
			RetirementAge: 60,
		}

		invalidPlan := entities.RetirementPlan{
			BirthDate:      time.Now().AddDate(-30, 0, 0).Format("2006-01-02"),
			RetirementAge:  65,
			ExpectLifespan: 85,
			MonthlyIncome:  -1000,
		}

		mockRepo.On("GetRetirementByUserID", "test-id").Return(existingPlan, nil)

		result, err := useCase.UpdateRetirementByID("test-id", invalidPlan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "monthlyIncome must be greater than or equal to zero", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRetirementByUserID(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedRetirement := &entities.RetirementPlan{
			ID:            "test-id",
			RetirementAge: 60,
		}

		mockRepo.On("GetRetirementByUserID", "user-id").Return(expectedRetirement, nil)

		result, err := useCase.GetRetirementByUserID("user-id")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRetirementByUserID", "non-existent").Return(nil, errors.New("retirement plan not found"))

		result, err := useCase.GetRetirementByUserID("non-existent")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
