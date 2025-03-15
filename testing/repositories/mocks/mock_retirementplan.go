package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockRetirementRepository struct {
	mock.Mock
}

func (m *MockRetirementRepository) CreateRetirement(retirement *entities.RetirementPlan) (*entities.RetirementPlan, error) {
	args := m.Called(retirement)
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
	return args.Get(0).(*entities.RetirementPlan), args.Error(1)
}
