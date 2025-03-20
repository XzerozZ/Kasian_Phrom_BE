package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockTransRepository struct {
	mock.Mock
}

func (m *MockTransRepository) CreateTransaction(transaction *entities.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransRepository) GetTransactionByID(id string) (*entities.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Transaction), args.Error(1)
}

func (m *MockTransRepository) GetTransactionByUserID(userID string) ([]map[string]interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockTransRepository) GetTransactionByLoanIDs(loanIDs []string) ([]entities.Transaction, error) {
	args := m.Called(loanIDs)
	return args.Get(0).([]entities.Transaction), args.Error(1)
}

func (m *MockTransRepository) GetLatestTransactionByLoanID(loanID string) (*entities.Transaction, error) {
	args := m.Called(loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Transaction), args.Error(1)
}

func (m *MockTransRepository) UpdateTransaction(transaction *entities.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransRepository) DeleteTransaction(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransRepository) DeleteTransactionsByLoanID(loanID string) error {
	args := m.Called(loanID)
	return args.Error(0)
}

func (m *MockTransRepository) CountTransactionsByLoanID(loanID string) (int, error) {
	args := m.Called(loanID)
	return args.Int(0), args.Error(1)
}
