package mocks

import (
	"mime/multipart"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(user *entities.User, roleName string) (*entities.User, error) {
	args := m.Called(user, roleName)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) Login(email, password string) (string, *entities.User, error) {
	args := m.Called(email, password)
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) LoginAdmin(email, password string) (string, *entities.User, error) {
	args := m.Called(email, password)
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) LoginWithGoogle(user *entities.User) (string, *entities.User, error) {
	args := m.Called(user)
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) ResetPassword(userID, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByID(userID string) (*entities.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUserByID(id string, user entities.User, files *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	args := m.Called(id, user, files, ctx)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserUseCase) VerifyOTP(email, otpCode string) error {
	args := m.Called(email, otpCode)
	return args.Error(0)
}

func (m *MockUserUseCase) ChangedPassword(email, newPassword string) error {
	args := m.Called(email, newPassword)
	return args.Error(0)
}

func (m *MockUserUseCase) CalculateRetirement(userID string) (fiber.Map, error) {
	args := m.Called(userID)
	return args.Get(0).(fiber.Map), args.Error(1)
}

func (m *MockUserUseCase) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserUseCase) UpdateSelectedHouse(userID, nursingHouseID string, transfers []entities.TransferRequest) (*entities.SelectedHouse, error) {
	args := m.Called(userID, nursingHouseID, transfers)
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserUseCase) CreateHistory(history entities.History) (*entities.History, error) {
	args := m.Called(history)
	return args.Get(0).(*entities.History), args.Error(1)
}

func (m *MockUserUseCase) GetHistoryByUserID(userID string) (fiber.Map, error) {
	args := m.Called(userID)
	return args.Get(0).(fiber.Map), args.Error(1)
}

func (m *MockUserUseCase) GetHistoryByMonth(userID string) (map[string]float64, error) {
	args := m.Called(userID)
	return args.Get(0).(map[string]float64), args.Error(1)
}
