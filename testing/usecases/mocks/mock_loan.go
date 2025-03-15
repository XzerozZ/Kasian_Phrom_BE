package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockLoanUseCase struct {
	mock.Mock
}

func (m *MockLoanUseCase) CreateLoan(loan entities.Loan) (*entities.Loan, error) {
	args := m.Called(loan)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUseCase) GetLoanByID(id string) (*entities.Loan, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUseCase) GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Loan), args.Get(1).(map[string]interface{}), args.Error(2)
}

func (m *MockLoanUseCase) UpdateLoanStatusByID(id string, loan entities.Loan) (*entities.Loan, error) {
	args := m.Called(id, loan)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUseCase) DeleteLoanByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
