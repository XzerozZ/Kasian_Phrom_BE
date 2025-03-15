package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) CreateLoan(loan *entities.Loan) (*entities.Loan, error) {
	args := m.Called(loan)
	return args.Get(0).(*entities.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetLoanByID(id string) (*entities.Loan, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Loan), args.Get(1).(map[string]interface{}), args.Error(2)
}

func (m *MockLoanRepository) GetAllLoansByStatus(statuses []string) ([]entities.Loan, error) {
	args := m.Called(statuses)
	return args.Get(0).([]entities.Loan), args.Error(1)
}

func (m *MockLoanRepository) UpdateLoanByID(loan *entities.Loan) (*entities.Loan, error) {
	args := m.Called(loan)
	return args.Get(0).(*entities.Loan), args.Error(1)
}

func (m *MockLoanRepository) DeleteLoanByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
