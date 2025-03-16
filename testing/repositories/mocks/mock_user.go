package mocks

import (
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) CreateSelectedHouse(selectedHouse *entities.SelectedHouse) error {
	args := m.Called(selectedHouse)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*entities.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (*entities.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetRoleByName(name string) (entities.Role, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Role), args.Error(1)
}

func (m *MockUserRepository) UpdateUserByID(user *entities.User) (*entities.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) CreateOTP(otp *entities.OTP) error {
	args := m.Called(otp)
	return args.Error(0)
}

func (m *MockUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.OTP), args.Error(1)
}

func (m *MockUserRepository) DeleteOTP(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserRepository) UpdateSelectedHouse(selectedHouse *entities.SelectedHouse) (*entities.SelectedHouse, error) {
	args := m.Called(selectedHouse)
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserRepository) CreateHistory(history *entities.History) (*entities.History, error) {
	args := m.Called(history)
	return args.Get(0).(*entities.History), args.Error(1)
}

func (m *MockUserRepository) GetHistoryByUserID(userID string) ([]entities.History, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.History), args.Error(1)
}

func (m *MockUserRepository) GetHistoryInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).([]entities.History), args.Error(1)
}

func (m *MockUserRepository) GetUserDepositsInRange(userID string, startDate, endDate time.Time) ([]entities.History, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).([]entities.History), args.Error(1)
}

func (m *MockUserRepository) GetUserHistoryByMonth(userID string) (map[string]float64, error) {
	args := m.Called(userID)
	return args.Get(0).(map[string]float64), args.Error(1)
}
