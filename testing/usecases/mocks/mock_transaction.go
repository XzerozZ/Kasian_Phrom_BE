package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockTransactionUseCase struct {
	mock.Mock
}

func (m *MockTransactionUseCase) CreateTransactionsForAllUsers() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransactionUseCase) MarkTransactiontoPaid(transactionID, userID string) error {
	args := m.Called(transactionID, userID)
	return args.Error(0)
}

func (m *MockTransactionUseCase) GetTransactionByUserID(userID string) ([]map[string]interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}
