package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(user *entities.User, roleName string) (*entities.User, error) {
	args := m.Called(user, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) Login(email, password string) (string, *entities.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) LoginAdmin(email, password string) (string, *entities.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) ResetPassword(userID, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserUseCase) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserUseCase) VerifyOTP(email, otp string) error {
	args := m.Called(email, otp)
	return args.Error(0)
}

func (m *MockUserUseCase) ChangedPassword(userID, newPassword string) error {
	args := m.Called(userID, newPassword)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByID(id string) (*entities.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserUseCase) UpdateUserByID(id string, user entities.User, files *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	args := m.Called(id, user, files, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateSelectedHouse(userID, nursingHouseID string) (*entities.SelectedHouse, error) {
	args := m.Called(userID, nursingHouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SelectedHouse), args.Error(1)
}

func (m *MockUserUseCase) CalculateRetirement(userID string) (fiber.Map, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(fiber.Map), args.Error(1)
}

func setupApp(controller *controllers.UserController) *fiber.App {
	app := fiber.New()
	app.Post("/register", controller.RegisterHandler)
	app.Post("/login", controller.LoginHandler)
	app.Post("/admin/login", controller.LoginAdminHandler)
	app.Post("/forgot-password", controller.ForgotPasswordHandler)
	app.Post("/verify-otp", controller.VerifyOTPHandler)
	app.Post("/resetpassword", controller.ResetPasswordHandler)
	app.Post("/logout", controller.LogoutHandler)
	app.Get("/users/:id", controller.GetUserByIDHandler)
	app.Get("/selected-house", controller.GetSelectedHouseHandler)
	app.Put("/users/:id", controller.UpdateUserByIDHandler)
	app.Put("/selected-house/:nh_id", controller.UpdateSelectedHouseHandler)
	app.Get("/retirement-plan", controller.GetRetirementPlanHandler)
	return app
}

func TestRegisterHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	tests := []struct {
		name           string
		reqBody        map[string]string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "successful registration",
			reqBody: map[string]string{
				"uname":    "testuser",
				"email":    "test@test.com",
				"password": "password123",
				"role":     "User",
			},
			expectedStatus: fiber.StatusCreated,
			mockSetup: func() {
				expectedUser := &entities.User{
					Username: "testuser",
					Email:    "test@test.com",
				}
				mockUseCase.On("Register", mock.AnythingOfType("*entities.User"), "User").Return(expectedUser, nil)
			},
		},
		{
			name: "missing username",
			reqBody: map[string]string{
				"email":    "test@test.com",
				"password": "password123",
				"role":     "User",
			},
			expectedStatus: fiber.StatusBadRequest,
			mockSetup:      func() {},
		},
		{
			name: "registration error",
			reqBody: map[string]string{
				"uname":    "testuser",
				"email":    "test@test.com",
				"password": "password123",
				"role":     "User",
			},
			expectedStatus: fiber.StatusInternalServerError,
			mockSetup: func() {
				mockUseCase.On("Register", mock.AnythingOfType("*entities.User"), "User").
					Return(nil, errors.New("registration failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			jsonBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	tests := []struct {
		name           string
		reqBody        map[string]string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "successful login",
			reqBody: map[string]string{
				"email":    "test@test.com",
				"password": "password123",
			},
			expectedStatus: fiber.StatusOK,
			mockSetup: func() {
				expectedUser := &entities.User{
					Username: "testuser",
					Role:     entities.Role{RoleName: "User"},
				}
				mockUseCase.On("Login", "test@test.com", "password123").
					Return("jwt.token.here", expectedUser, nil)
			},
		},
		{
			name: "missing email",
			reqBody: map[string]string{
				"password": "password123",
			},
			expectedStatus: fiber.StatusBadRequest,
			mockSetup:      func() {},
		},
		{
			name: "invalid credentials",
			reqBody: map[string]string{
				"email":    "test@test.com",
				"password": "wrongpass",
			},
			expectedStatus: fiber.StatusInternalServerError,
			mockSetup: func() {
				mockUseCase.On("Login", "test@test.com", "wrongpass").
					Return("", nil, errors.New("invalid credentials"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			jsonBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestForgotPasswordHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	tests := []struct {
		name           string
		reqBody        map[string]string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "successful forgot password",
			reqBody: map[string]string{
				"email": "test@test.com",
			},
			expectedStatus: fiber.StatusOK,
			mockSetup: func() {
				mockUseCase.On("ForgotPassword", "test@test.com").Return(nil)
			},
		},
		{
			name:           "missing email",
			reqBody:        map[string]string{},
			expectedStatus: fiber.StatusBadRequest,
			mockSetup:      func() {},
		},
		{
			name: "forgot password error",
			reqBody: map[string]string{
				"email": "test@test.com",
			},
			expectedStatus: fiber.StatusInternalServerError,
			mockSetup: func() {
				mockUseCase.On("ForgotPassword", "test@test.com").
					Return(errors.New("forgot password failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			jsonBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/forgot-password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestVerifyOTPHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	tests := []struct {
		name           string
		reqBody        map[string]string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "successful OTP verification",
			reqBody: map[string]string{
				"email": "test@test.com",
				"otp":   "123456",
			},
			expectedStatus: fiber.StatusOK,
			mockSetup: func() {
				mockUseCase.On("VerifyOTP", "test@test.com", "123456").Return(nil)
			},
		},
		{
			name: "missing email",
			reqBody: map[string]string{
				"otp": "123456",
			},
			expectedStatus: fiber.StatusBadRequest,
			mockSetup:      func() {},
		},
		{
			name: "invalid OTP",
			reqBody: map[string]string{
				"email": "test@test.com",
				"otp":   "123456",
			},
			expectedStatus: fiber.StatusInternalServerError,
			mockSetup: func() {
				mockUseCase.On("VerifyOTP", "test@test.com", "123456").
					Return(errors.New("invalid OTP"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			jsonBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/verify-otp", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	req := httptest.NewRequest("POST", "/logout", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Success", result["status"])
	assert.Equal(t, "Logout successful", result["message"])
}

func TestGetUserByIDHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "successful get user",
			userID:         "123",
			expectedStatus: fiber.StatusOK,
			mockSetup: func() {
				expectedUser := &entities.User{
					Username: "testuser",
					Email:    "test@test.com",
				}
				mockUseCase.On("GetUserByID", "123").Return(expectedUser, nil)
			},
		},
		{
			name:           "user not found",
			userID:         "999",
			expectedStatus: fiber.StatusNotFound,
			mockSetup: func() {
				mockUseCase.On("GetUserByID", "999").Return(nil, errors.New("user not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			fctx := &fasthttp.RequestCtx{}
			ctx := app.AcquireCtx(fctx)
			defer app.ReleaseCtx(ctx)
			ctx.Locals("user_id", tt.userID)
			err := controller.GetUserByIDHandler(ctx)
			if err != nil {
				assert.Equal(t, tt.expectedStatus, err.(*fiber.Error).Code)
			} else {
				assert.Equal(t, tt.expectedStatus, ctx.Response().StatusCode())
			}
		})
	}
}

func TestGetRetirementPlanHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := fiber.New()

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		setupAuth      bool
		mockSetup      func()
	}{
		{
			name:           "successful retirement calculation",
			userID:         "123",
			expectedStatus: fiber.StatusOK,
			setupAuth:      true,
			mockSetup: func() {
				expectedResult := fiber.Map{
					"required_funds":  1000000.0,
					"monthly_savings": 5000.0,
				}
				mockUseCase.On("CalculateRetirement", "123").Return(expectedResult, nil)
			},
		},
		{
			name:           "calculation error",
			userID:         "999",
			expectedStatus: fiber.StatusNotFound,
			setupAuth:      true,
			mockSetup: func() {
				mockUseCase.On("CalculateRetirement", "999").Return(fiber.Map{}, errors.New("calculation error"))
			},
		},
		{
			name:           "unauthorized access",
			userID:         "",
			expectedStatus: fiber.StatusUnauthorized,
			setupAuth:      false,
			mockSetup:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			fctx := &fasthttp.RequestCtx{}
			ctx := app.AcquireCtx(fctx)
			defer app.ReleaseCtx(ctx)
			if tt.setupAuth {
				ctx.Locals("user_id", tt.userID)
			}

			err := controller.GetRetirementPlanHandler(ctx)
			if err != nil {
				assert.Equal(t, tt.expectedStatus, err.(*fiber.Error).Code)
			} else {
				assert.Equal(t, tt.expectedStatus, ctx.Response().StatusCode())
			}
		})
	}
}
