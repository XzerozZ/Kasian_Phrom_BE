package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockRetirementUseCase struct {
	mock.Mock
}

func (m *MockRetirementUseCase) CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, int, error) {
	args := m.Called(retirement)
	if result := args.Get(0); result != nil {
		return result.(*entities.RetirementPlan), args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockRetirementUseCase) GetRetirementByID(id string) (*entities.RetirementPlan, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entities.RetirementPlan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRetirementUseCase) GetRetirementByUserID(userID string) (*entities.RetirementPlan, error) {
	args := m.Called(userID)
	if result := args.Get(0); result != nil {
		return result.(*entities.RetirementPlan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRetirementUseCase) UpdateRetirementByID(userID string, retirement entities.RetirementPlan) (*entities.RetirementPlan, error) {
	args := m.Called(userID, retirement)
	if result := args.Get(0); result != nil {
		return result.(*entities.RetirementPlan), args.Error(1)
	}
	return nil, args.Error(1)
}
