package usecases_test

import (
	"errors"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/usecases"
	"github.com/stretchr/testify/assert"
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

func TestCreateQuiz_NewUser(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	mockRepo.On("GetQuizByUserID", "new-user-id").Return(nil, errors.New("quiz not found"))
	expectedQuiz := &entities.Quiz{
		UserID: "new-user-id",
		RiskID: 2,
	}

	mockRepo.On("CreateQuiz", mock.MatchedBy(func(q *entities.Quiz) bool {
		return q.UserID == "new-user-id" && q.RiskID == 2
	})).Return(expectedQuiz, nil)

	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	weights := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	result, err := quizUseCase.CreateQuiz("new-user-id", weights)

	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateQuiz_ExistingUser(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	existingQuiz := &entities.Quiz{
		UserID: "existing-user-id",
		RiskID: 1,
	}

	mockRepo.On("GetQuizByUserID", "existing-user-id").Return(existingQuiz, nil)
	mockRepo.On("DeleteQuiz", "existing-user-id").Return(nil)

	expectedQuiz := &entities.Quiz{
		UserID: "existing-user-id",
		RiskID: 3,
	}

	mockRepo.On("CreateQuiz", mock.MatchedBy(func(q *entities.Quiz) bool {
		return q.UserID == "existing-user-id" && q.RiskID == 3
	})).Return(expectedQuiz, nil)

	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	weights := []int{2, 3, 2, 3, 2, 3, 2, 3, 2, 3}
	result, err := quizUseCase.CreateQuiz("existing-user-id", weights)

	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateQuiz_DeleteFails(t *testing.T) {
	mockRepo := new(MockQuizRepository)

	existingQuiz := &entities.Quiz{
		UserID: "existing-user-id",
		RiskID: 1,
	}

	mockRepo.On("GetQuizByUserID", "existing-user-id").Return(existingQuiz, nil)
	mockRepo.On("DeleteQuiz", "existing-user-id").Return(errors.New("delete failed"))
	quizUseCase := usecases.NewQuizUseCase(mockRepo)

	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result, err := quizUseCase.CreateQuiz("existing-user-id", weights)

	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateQuiz_ZeroWeights(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	mockRepo.On("GetQuizByUserID", "user-id").Return(nil, errors.New("quiz not found"))

	expectedQuiz := &entities.Quiz{
		UserID: "user-id",
		RiskID: 1,
	}

	mockRepo.On("CreateQuiz", mock.MatchedBy(func(q *entities.Quiz) bool {
		return q.UserID == "user-id" && q.RiskID == 1
	})).Return(expectedQuiz, nil)

	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	weights := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	result, err := quizUseCase.CreateQuiz("user-id", weights)

	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateQuiz_HighWeights(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	mockRepo.On("GetQuizByUserID", "user-id").Return(nil, errors.New("quiz not found"))

	expectedQuiz := &entities.Quiz{
		UserID: "user-id",
		RiskID: 5,
	}

	mockRepo.On("CreateQuiz", mock.MatchedBy(func(q *entities.Quiz) bool {
		return q.UserID == "user-id" && q.RiskID == 5
	})).Return(expectedQuiz, nil)

	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	weights := []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
	result, err := quizUseCase.CreateQuiz("user-id", weights)

	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateQuiz_CreateFails(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	mockRepo.On("GetQuizByUserID", "user-id").Return(nil, errors.New("quiz not found"))

	mockRepo.On("CreateQuiz", mock.MatchedBy(func(q *entities.Quiz) bool {
		return q.UserID == "user-id" && q.RiskID == 2
	})).Return(nil, errors.New("create failed"))

	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	weights := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	result, err := quizUseCase.CreateQuiz("user-id", weights)

	assert.Error(t, err)
	assert.Equal(t, "create failed", err.Error())
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetQuizByUserID_Success(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	expectedQuiz := &entities.Quiz{
		UserID: "test-user-id",
		RiskID: 3,
	}

	mockRepo.On("GetQuizByUserID", "test-user-id").Return(expectedQuiz, nil)
	quizUseCase := usecases.NewQuizUseCase(mockRepo)
	result, err := quizUseCase.GetQuizByUserID("test-user-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedQuiz, result)
	mockRepo.AssertExpectations(t)
}

func TestGetQuizByUserID_NotFound(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	mockRepo.On("GetQuizByUserID", "non-existent-user").Return(nil, errors.New("quiz not found"))
	quizUseCase := usecases.NewQuizUseCase(mockRepo)

	result, err := quizUseCase.GetQuizByUserID("non-existent-user")

	assert.Error(t, err)
	assert.Equal(t, "quiz not found", err.Error())
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
