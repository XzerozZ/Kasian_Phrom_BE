package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockQuizUseCase struct {
	mock.Mock
}

func (m *MockQuizUseCase) CreateQuiz(userID string, weights []int) (*entities.Quiz, error) {
	args := m.Called(userID, weights)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}

func (m *MockQuizUseCase) GetQuizByUserID(userID string) (*entities.Quiz, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}
