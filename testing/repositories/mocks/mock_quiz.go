package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) CreateQuiz(quiz *entities.Quiz) (*entities.Quiz, error) {
	args := m.Called(quiz)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetQuizByUserID(userID string) (*entities.Quiz, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}

func (m *MockQuizRepository) DeleteQuiz(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
