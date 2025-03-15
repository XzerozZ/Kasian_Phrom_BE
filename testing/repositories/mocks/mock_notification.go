package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockNotiRepository struct {
	mock.Mock
}

func (m *MockNotiRepository) CreateNotification(notification *entities.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotiRepository) GetNotificationsByUserID(userID string) ([]entities.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Notification), args.Error(1)
}

func (m *MockNotiRepository) MarkNotificationAsRead(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
