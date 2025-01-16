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
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) LoginAdmin(email, password string) (string, *entities.User, error) {
	args := m.Called(email, password)
	return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUseCase) UpdateUserByID(id string, user entities.User, files *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	args := m.Called(id, user, files, ctx)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) CalculateRetirement(userID string) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

func setupApp(controller *controllers.UserController) *fiber.App {
	app := fiber.New()
	app.Post("/register", controller.RegisterHandler)
	app.Post("/login", controller.LoginHandler)
	app.Post("/admin/login", controller.LoginAdminHandler)
	app.Post("/forgot-password", controller.ForgotPasswordHandler)
	app.Post("/verify-otp", controller.VerifyOTPHandler)
	app.Post("/reset-password", controller.ResetPasswordHandler)
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

	t.Run("successful registration", func(t *testing.T) {
		reqBody := map[string]string{
			"uname":    "testuser",
			"email":    "test@test.com",
			"password": "password123",
			"role":     "User",
		}

		jsonBody, _ := json.Marshal(reqBody)
		expectedUser := &entities.User{
			Username: "testuser",
			Email:    "test@test.com",
		}

		mockUseCase.On("Register", mock.AnythingOfType("*entities.User"), "User").Return(expectedUser, nil)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestLoginHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)
	t.Run("successful login", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@test.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedToken := "jwt.token.here"
		expectedUser := &entities.User{
			Username: "testuser",
			Role:     entities.Role{RoleName: "User"},
		}

		mockUseCase.On("Login", "test@test.com", "password123").Return(expectedToken, expectedUser, nil)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestLoginAdminHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)
	t.Run("successful admin login", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "admin@test.com",
			"password": "admin123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedToken := "admin.jwt.token"
		expectedUser := &entities.User{
			Username: "admin",
			Role:     entities.Role{RoleName: "Admin"},
		}

		mockUseCase.On("LoginAdmin", "admin@test.com", "admin123").Return(expectedToken, expectedUser, nil)

		req := httptest.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUpdateUserHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)
	t.Run("successful user update", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		userFields := map[string]string{
			"username":  "updateduser",
			"email":     "updated@test.com",
			"firstname": "John",
			"lastname":  "Doe",
		}

		for key, value := range userFields {
			_ = writer.WriteField(key, value)
		}
		writer.Close()

		expectedUser := &entities.User{
			Username:  "updateduser",
			Email:     "updated@test.com",
			Firstname: "John",
			Lastname:  "Doe",
		}

		mockUseCase.On("UpdateUserByID",
			"123",
			mock.MatchedBy(func(u entities.User) bool {
				return u.Username == userFields["username"] &&
					u.Email == userFields["email"] &&
					u.Firstname == userFields["firstname"] &&
					u.Lastname == userFields["lastname"]
			}),
			(*multipart.FileHeader)(nil),
			mock.AnythingOfType("*fiber.Ctx"),
		).Return(expectedUser, nil)

		req := httptest.NewRequest("PUT", "/users/123", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetRetirementPlanHandler(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := setupApp(controller)

	t.Run("successful retirement calculation", func(t *testing.T) {
		expectedResult := fiber.Map{
			"required_funds":  1000000.0,
			"monthly_savings": 5000.0,
		}

		mockUseCase.On("CalculateRetirement", "user123").Return(expectedResult, nil)

		req := httptest.NewRequest("GET", "/retirement-plan", nil)
		ctx := app.AcquireCtx(req)
		ctx.Locals("user_id", "user123")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/retirement-plan", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("calculation error", func(t *testing.T) {
		mockUseCase.On("CalculateRetirement", "user123").Return(fiber.Map{}, errors.New("calculation error"))

		req := httptest.NewRequest("GET", "/retirement-plan", nil)
		ctx := app.AcquireCtx(req)
		ctx.Locals("user_id", "user123")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
