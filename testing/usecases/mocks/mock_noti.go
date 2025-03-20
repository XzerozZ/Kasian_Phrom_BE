package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockNotiUseCase struct {
	mock.Mock
}

func (m *MockNotiUseCase) GetNotificationsByUserID(userID string) ([]entities.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Notification), args.Error(1)
}

func (m *MockNotiUseCase) MarkNotificationsAsRead(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
