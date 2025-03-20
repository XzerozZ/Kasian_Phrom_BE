package usecases_test

import (
	"errors"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNotiUseCase(t *testing.T) {
	mockRepo := new(mocks.MockNotiRepository)
	useCase := usecases.NewNotiUseCase(mockRepo)

	t.Run("GetNotificationsByUserID - Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "user_123"
		expectedNotifications := []entities.Notification{
			{ID: "1", UserID: userID, Message: "Test notification 1", IsRead: false},
			{ID: "2", UserID: userID, Message: "Test notification 2", IsRead: true},
		}

		mockRepo.On("GetNotificationsByUserID", userID).Return(expectedNotifications, nil).Once()

		notifications, err := useCase.GetNotificationsByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedNotifications, notifications)
		assert.Len(t, notifications, 2)
		assert.Equal(t, "Test notification 1", notifications[0].Message)
		assert.Equal(t, "Test notification 2", notifications[1].Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetNotificationsByUserID - Error", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "user_123"
		expectedError := errors.New("database error")

		mockRepo.On("GetNotificationsByUserID", userID).Return([]entities.Notification{}, expectedError).Once()

		notifications, err := useCase.GetNotificationsByUserID(userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, notifications)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetNotificationsByUserID - Empty Result", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "user_123"
		expectedNotifications := []entities.Notification{}

		mockRepo.On("GetNotificationsByUserID", userID).Return(expectedNotifications, nil).Once()

		notifications, err := useCase.GetNotificationsByUserID(userID)

		assert.NoError(t, err)
		assert.Empty(t, notifications)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkNotificationsAsRead - Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "user_123"

		mockRepo.On("MarkNotificationAsRead", userID).Return(nil).Once()

		err := useCase.MarkNotificationsAsRead(userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkNotificationsAsRead - Error", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "user_123"
		expectedError := errors.New("database error")

		mockRepo.On("MarkNotificationAsRead", userID).Return(expectedError).Once()

		err := useCase.MarkNotificationsAsRead(userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkNotificationsAsRead - User Not Found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		userID := "nonexistent_user"
		expectedError := errors.New("user not found")

		mockRepo.On("MarkNotificationAsRead", userID).Return(expectedError).Once()

		err := useCase.MarkNotificationsAsRead(userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
