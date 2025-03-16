package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(_ *testing.T) (*controllers.UserController, *mocks.MockUserUseCase, *fiber.App) {
	mockUseCase := new(mocks.MockUserUseCase)
	controller := controllers.NewUserController(mockUseCase)
	app := fiber.New()
	return controller, mockUseCase, app
}

func TestRegisterHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/register", controller.RegisterHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"uname":    "testuser",
			"email":    "test@example.com",
			"password": "password123",
			"role":     "user",
		}
		jsonBody, _ := json.Marshal(requestBody)

		user := &entities.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUseCase.On("Register", mock.MatchedBy(func(u *entities.User) bool {
			return u.Username == user.Username && u.Email == user.Email && u.Password == user.Password
		}), "user").Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "user created successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Username", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"role":     "user",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "Username is missing", result["message"])
	})

	t.Run("Registration Failure", func(t *testing.T) {
		requestBody := map[string]string{
			"uname":    "testuser",
			"email":    "test@example.com",
			"password": "password123",
			"role":     "user",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("Register", mock.MatchedBy(func(u *entities.User) bool {
			return u.Username == "testuser" && u.Email == "test@example.com" && u.Password == "password123"
		}), "user").Return(&entities.User{}, errors.New("registration failed")).Once()

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrInternalServerError.Code, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})
}

func TestLoginHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/login", controller.LoginHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		user := &entities.User{
			ID:       "user123",
			Username: "testuser",
			Role:     entities.Role{RoleName: "user"},
		}

		mockUseCase.On("Login", "test@example.com", "password123").Return("testtoken", user, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Login successful", result["message"])

		resultData := result["result"].(map[string]interface{})
		assert.Equal(t, "testtoken", resultData["token"])
		assert.Equal(t, "user123", resultData["u_id"])
		assert.Equal(t, "testuser", resultData["uname"])
		assert.Equal(t, "user", resultData["role"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Email", func(t *testing.T) {
		requestBody := map[string]string{
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "Email is missing", result["message"])
	})

	t.Run("Login Failure", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("Login", "test@example.com", "password123").Return("", &entities.User{}, errors.New("login failed")).Once()

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrInternalServerError.Code, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})
}

func TestLoginAdminHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/login-admin", controller.LoginAdminHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "admin@example.com",
			"password": "adminpass",
		}
		jsonBody, _ := json.Marshal(requestBody)

		admin := &entities.User{
			ID:       "admin123",
			Username: "adminuser",
			Role:     entities.Role{RoleName: "admin"},
		}

		mockUseCase.On("LoginAdmin", "admin@example.com", "adminpass").Return("admintoken", admin, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/login-admin", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Login successful", result["message"])

		resultData := result["result"].(map[string]interface{})
		assert.Equal(t, "admintoken", resultData["token"])
		assert.Equal(t, "admin123", resultData["u_id"])
		assert.Equal(t, "adminuser", resultData["uname"])
		assert.Equal(t, "admin", resultData["role"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestLoginWithGoogleHandler(t *testing.T) {
	controller, mockUserUsecase, app := setupTest(t)

	app.Post("/login/google", controller.LoginWithGoogleHandler)

	t.Run("Success", func(t *testing.T) {
		returnedUser := &entities.User{
			ID:       "1",
			Username: "johndoe",
			Role:     entities.Role{RoleName: "user"},
		}
		mockUserUsecase.On("LoginWithGoogle", mock.Anything).Return("token123", returnedUser, nil).Once()

		reqBody := map[string]interface{}{
			"fname":      "John",
			"lname":      "Doe",
			"uname":      "johndoe",
			"email":      "john.doe@example.com",
			"image_link": "https://example.com/profile.jpg",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Login successful", result["message"])

		resultData := result["result"].(map[string]interface{})
		assert.Equal(t, "token123", resultData["token"])
		assert.Equal(t, "1", resultData["u_id"])
		assert.Equal(t, "johndoe", resultData["uname"])
		assert.Equal(t, "user", resultData["role"])
	})

	t.Run("Success with Default Image", func(t *testing.T) {
		mockUserUsecase.ExpectedCalls = nil

		returnedUser := &entities.User{
			ID:       "1",
			Username: "johndoe",
			Role:     entities.Role{RoleName: "user"},
		}
		mockUserUsecase.On("LoginWithGoogle", mock.Anything).Return("token123", returnedUser, nil).Once()

		reqBody := map[string]interface{}{
			"fname":      "John",
			"lname":      "Doe",
			"uname":      "johndoe",
			"email":      "john.doe@example.com",
			"image_link": "",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Login successful", result["message"])
	})

	t.Run("Missing Fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"fname": "John",
			"lname": "Doe",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "Firstname, Lastname, Username, Email or ImageLink is missing", result["message"])
	})

	t.Run("Usecase Error", func(t *testing.T) {
		mockUserUsecase.ExpectedCalls = nil

		mockUserUsecase.On("LoginWithGoogle", mock.Anything).Return("", &entities.User{}, errors.New("database error")).Once()

		reqBody := map[string]interface{}{
			"fname":      "John",
			"lname":      "Doe",
			"uname":      "johndoe",
			"email":      "john.doe@example.com",
			"image_link": "https://example.com/profile.jpg",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Internal Server Error", result["status"])
		assert.Equal(t, "database error", result["message"])
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBuffer([]byte(`{invalid json}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Bad Request", result["status"])
		assert.NotNil(t, result["message"])
	})
}

func TestResetPasswordHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/reset-password", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user123")
		return controller.ResetPasswordHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"old_password": "oldpass",
			"new_password": "newpass",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("ResetPassword", "user123", "oldpass", "newpass").Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Password reset successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Password Fields", func(t *testing.T) {
		requestBody := map[string]string{
			"old_password": "oldpass",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Contains(t, result["message"], "missing")
	})
}

func TestForgotPasswordHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/forgot-password", controller.ForgotPasswordHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("ForgotPassword", "test@example.com").Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Sent OTP successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Email", func(t *testing.T) {
		requestBody := map[string]string{}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "Email is missing", result["message"])
	})
}

func TestVerifyOTPHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/verify-otp", controller.VerifyOTPHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"email": "test@example.com",
			"otp":   "123456",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("VerifyOTP", "test@example.com", "123456").Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "OTP is correct", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Fields", func(t *testing.T) {
		requestBody := map[string]string{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "OTP is missing", result["message"])
	})
}

func TestChangedPasswordHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/change-password", controller.ChangedPasswordHandler)

	t.Run("Success", func(t *testing.T) {
		requestBody := map[string]string{
			"email":       "test@example.com",
			"newPassword": "newpass123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("ChangedPassword", "test@example.com", "newpass123").Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "changed password successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing Fields", func(t *testing.T) {
		requestBody := map[string]string{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "NewPassword is missing", result["message"])
	})
}

func TestLogoutHandler(t *testing.T) {
	controller, _, app := setupTest(t)

	app.Post("/logout", controller.LogoutHandler)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, "Logout successful", result["message"])
	})
}

func TestGetUserByIDHandler(t *testing.T) {
	controller, mockUserUseCase, app := setupTest(t)

	app.Get("/user", controller.GetUserByIDHandler)

	t.Run("GetUserByIDHandler - Success", func(t *testing.T) {
		userID := "123"
		expectedUser := &entities.User{
			ID:        "1",
			Firstname: "John",
			Lastname:  "Doe",
			Username:  "johndoe",
			Email:     "john.doe@example.com",
			Role:      entities.Role{RoleName: "user"},
		}

		mockUserUseCase.On("GetUserByID", userID).Return(expectedUser, nil).Once()
		req := httptest.NewRequest("GET", "/user", nil)
		app := fiber.New()
		app.Get("/user", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetUserByIDHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "User Info retrieved successfully", response["message"])
		assert.NotNil(t, response["result"])
		mockUserUseCase.AssertExpectations(t)
	})

	t.Run("GetUserByIDHandler - No User ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Error", response["status"])
		assert.Contains(t, response["message"], "Unauthorized")
	})

	t.Run("GetUserByIDHandler - Not Found", func(t *testing.T) {
		userID := "123"
		mockUserUseCase.On("GetUserByID", userID).Return(&entities.User{}, errors.New("user not found")).Once()
		req := httptest.NewRequest("GET", "/user", nil)
		app := fiber.New()
		app.Get("/user", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetUserByIDHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Not Found", response["status"])
		assert.Contains(t, response["message"], "user not found")
		mockUserUseCase.AssertExpectations(t)
	})
}

func TestSelectedHouseHandler(t *testing.T) {
	controller, mockUserUseCase, app := setupTest(t)

	app.Get("/selected-house", controller.GetSelectedHouseHandler)

	t.Run("GetSelectedHouseHandler - Success", func(t *testing.T) {
		userID := "user123"
		expectedHouse := &entities.SelectedHouse{
			UserID:         userID,
			NursingHouseID: "1",
			Status:         "In_Progress",
		}

		mockUserUseCase.On("GetSelectedHouse", userID).Return(expectedHouse, nil).Once()
		req := httptest.NewRequest("GET", "/selected-house", nil)
		app := fiber.New()
		app.Get("/selected-house", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetSelectedHouseHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Selected house retrieved successfully", response["message"])
		assert.NotNil(t, response["result"])
		mockUserUseCase.AssertExpectations(t)
	})

	t.Run("GetSelectedHouseHandler - No User ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/selected-house", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Error", response["status"])
		assert.Contains(t, response["message"], "Unauthorized")
	})

	t.Run("GetSelectedHouseHandler - Not Found", func(t *testing.T) {
		userID := "user123"
		mockUserUseCase.On("GetSelectedHouse", userID).Return(&entities.SelectedHouse{}, errors.New("selected house not found")).Once()
		req := httptest.NewRequest("GET", "/selected-house", nil)
		app := fiber.New()
		app.Get("/selected-house", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetSelectedHouseHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Not Found", response["status"])
		assert.Contains(t, response["message"], "selected house not found")
		mockUserUseCase.AssertExpectations(t)
	})
}

func TestUpdateUserByIDHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Put("/user/update", func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-id")
		return controller.UpdateUserByIDHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		userJSON := `{"name":"Test User","email":"test@example.com"}`
		err := writer.WriteField("user", userJSON)
		assert.NoError(t, err)

		part, err := writer.CreateFormFile("images", "test-image.jpg")
		assert.NoError(t, err)
		_, err = part.Write([]byte("dummy image content"))
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/user/update", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		var userData entities.User
		json.Unmarshal([]byte(userJSON), &userData)

		expectedUser := entities.User{
			ID:        "test-user-id",
			Firstname: "Test User",
			Lastname:  "Test LastName",
			Email:     "test@example.com",
		}

		mockUseCase.On("UpdateUserByID", "test-user-id", mock.AnythingOfType("entities.User"), mock.AnythingOfType("*multipart.FileHeader"), mock.Anything).Return(&expectedUser, nil).Once()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "User retrieved successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing User ID", func(t *testing.T) {
		app.Put("/user/update/no-auth", controller.UpdateUserByIDHandler)
		req := httptest.NewRequest(http.MethodPut, "/user/update/no-auth", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
	})

	t.Run("Error - Failed to Parse Body", func(t *testing.T) {
		app.Put("/user/update/bad-body", func(c *fiber.Ctx) error {
			c.Locals("user_id", "test-user-id")
			return controller.UpdateUserByIDHandler(c)
		})

		req := httptest.NewRequest(http.MethodPut, "/user/update/bad-body", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Contains(t, result["message"], "")
	})

	t.Run("Error - Update User Failed", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		userJSON := `{"name":"Test User","email":"test@example.com"}`
		writer.WriteField("user", userJSON)
		writer.Close()

		app.Put("/user/update/error", func(c *fiber.Ctx) error {
			c.Locals("user_id", "test-user-id")
			return controller.UpdateUserByIDHandler(c)
		})

		req := httptest.NewRequest(http.MethodPut, "/user/update/error", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		var userData entities.User
		json.Unmarshal([]byte(userJSON), &userData)

		mockUseCase.On("UpdateUserByID", "test-user-id", mock.AnythingOfType("entities.User"), mock.Anything, mock.Anything).Return(&entities.User{}, errors.New("update error")).Once()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Equal(t, "update error", result["message"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestUpdateSelectedHouseHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Put("/user/houses/:nh_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-id")
		return controller.UpdateSelectedHouseHandler(c)
	})

	t.Run("Success - Regular House", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/user/houses/12345", nil)

		expectedHouse := entities.SelectedHouse{
			UserID:         "test-user-id",
			NursingHouseID: "00001",
		}

		mockUseCase.On("UpdateSelectedHouse", "test-user-id", "12345", []entities.TransferRequest(nil)).Return(&expectedHouse, nil).Once()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "House Updated to user successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Success - Special House with Transfers", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "cash")
		writer.WriteField("name", "Donation")
		writer.WriteField("amount", "1000.50")

		writer.WriteField("type", "goods")
		writer.WriteField("name", "Food")
		writer.WriteField("amount", "500.25")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		expectedTransfers := []entities.TransferRequest{
			{Type: "cash", Name: "Donation", Amount: 1000.50},
			{Type: "goods", Name: "Food", Amount: 500.25},
		}

		expectedHouse := entities.SelectedHouse{
			UserID:         "test-user-id",
			NursingHouseID: "00001",
		}

		mockUseCase.On("UpdateSelectedHouse", "test-user-id", "00001", mock.MatchedBy(func(transfers []entities.TransferRequest) bool {
			if len(transfers) != 2 {
				return false
			}

			return transfers[0].Type == expectedTransfers[0].Type &&
				transfers[1].Type == expectedTransfers[1].Type &&
				transfers[0].Amount == expectedTransfers[0].Amount &&
				transfers[1].Amount == expectedTransfers[1].Amount
		})).Return(&expectedHouse, nil).Once()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "House Updated to user successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing User ID", func(t *testing.T) {
		app.Put("/user/houses/unauthorized/:nh_id", controller.UpdateSelectedHouseHandler)
		req := httptest.NewRequest(http.MethodPut, "/user/houses/unauthorized/12345", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
	})

	t.Run("Bad Request - Missing Nursing House ID", func(t *testing.T) {
		app.Put("/user/houses/", func(c *fiber.Ctx) error {
			c.Locals("user_id", "test-user-id")
			c.Params("nh_id", "")
			return controller.UpdateSelectedHouseHandler(c)
		})

		req := httptest.NewRequest(http.MethodPut, "/user/houses/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Equal(t, "Invalid request: Missing nursing house ID", result["message"])
	})

	t.Run("Bad Request - Invalid Form Data", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Contains(t, result["message"].(string), "Missing required fields: 'type', 'name', or 'amount'")
	})

	t.Run("Bad Request - Missing Required Fields", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "cash")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Equal(t, "Missing required fields: 'type', 'name', or 'amount'", result["message"])
	})

	t.Run("Bad Request - Mismatch Count", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "cash")
		writer.WriteField("type", "goods")
		writer.WriteField("name", "Donation")
		writer.WriteField("amount", "1000")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Equal(t, "Mismatch in count of 'type', 'name', and 'amount'", result["message"])
	})

	t.Run("Bad Request - Invalid Amount Format", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "cash")
		writer.WriteField("name", "Donation")
		writer.WriteField("amount", "not-a-number")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Contains(t, result["message"].(string), "Invalid amount format at index 0")
	})

	t.Run("Bad Request - Negative Amount", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "cash")
		writer.WriteField("name", "Donation")
		writer.WriteField("amount", "-100.50")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/user/houses/00001", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
		assert.Contains(t, result["message"].(string), "Amount cannot be negative at index 0")
	})

	t.Run("Error - Usecase Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/user/houses/12345", nil)

		mockUseCase.On("UpdateSelectedHouse", "test-user-id", "12345", []entities.TransferRequest(nil)).Return(&entities.SelectedHouse{}, errors.New("database error")).Once()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Equal(t, "database error", result["message"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetRetirementPlanHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Get("/user/retirement-plan", func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-id")
		return controller.GetRetirementPlanHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		expectedRetirementPlan := fiber.Map{
			"plan_name":                "My Retirement Plan",
			"allRequiredFund":          float64(3000000),
			"stillneed":                float64(1500000),
			"allretirementfund":        float64(2000000),
			"monthly_expenses":         float64(5000),
			"plan_saving":              float64(300000),
			"all_money":                float64(1500000),
			"saving":                   float64(800000),
			"investment":               float64(700000),
			"all_assets_expense":       float64(2000),
			"nursingHouse_expense":     float64(1500),
			"plan_expense":             float64(1500),
			"annual_savings_return":    float64(3.5),
			"annual_investment_return": float64(7.0),
		}

		mockUseCase.On("CalculateRetirement", "test-user-id").Return(expectedRetirementPlan, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/user/retirement-plan", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "This is user's retirement plan successfully", result["message"])

		resultData := result["result"].(map[string]interface{})

		assert.Equal(t, expectedRetirementPlan["plan_name"], resultData["plan_name"])
		assert.Equal(t, expectedRetirementPlan["allRequiredFund"], resultData["allRequiredFund"])
		assert.Equal(t, expectedRetirementPlan["stillneed"], resultData["stillneed"])
		assert.Equal(t, expectedRetirementPlan["allretirementfund"], resultData["allretirementfund"])
		assert.Equal(t, expectedRetirementPlan["monthly_expenses"], resultData["monthly_expenses"])
		assert.Equal(t, expectedRetirementPlan["plan_saving"], resultData["plan_saving"])
		assert.Equal(t, expectedRetirementPlan["all_money"], resultData["all_money"])
		assert.Equal(t, expectedRetirementPlan["saving"], resultData["saving"])
		assert.Equal(t, expectedRetirementPlan["investment"], resultData["investment"])
		assert.Equal(t, expectedRetirementPlan["all_assets_expense"], resultData["all_assets_expense"])
		assert.Equal(t, expectedRetirementPlan["nursingHouse_expense"], resultData["nursingHouse_expense"])
		assert.Equal(t, expectedRetirementPlan["plan_expense"], resultData["plan_expense"])
		assert.Equal(t, expectedRetirementPlan["annual_savings_return"], resultData["annual_savings_return"])
		assert.Equal(t, expectedRetirementPlan["annual_investment_return"], resultData["annual_investment_return"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing User ID", func(t *testing.T) {
		app.Get("/user/retirement-plan/unauthorized", controller.GetRetirementPlanHandler)
		req := httptest.NewRequest(http.MethodGet, "/user/retirement-plan/unauthorized", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
	})

	t.Run("Error - Calculation Failed", func(t *testing.T) {
		mockUseCase.On("CalculateRetirement", "test-user-id").Return(fiber.Map{}, errors.New("calculation error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/user/retirement-plan", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Equal(t, "calculation error", result["message"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestCreateHistoryHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Post("/histories", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return controller.CreateHistoryHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		requestBody := entities.History{
			Method: "POST",
			Type:   "login",
		}
		jsonBody, _ := json.Marshal(requestBody)

		expectedHistory := entities.History{
			ID:     "history-123",
			Method: "POST",
			Type:   "login",
			UserID: "user-123",
		}

		mockUseCase.On("CreateHistory", mock.MatchedBy(func(h entities.History) bool {
			return h.Method == "POST" && h.Type == "login" && h.UserID == "user-123"
		})).Return(&expectedHistory, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/histories", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "History created successfully", result["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing UserID", func(t *testing.T) {
		appWithoutUser := fiber.New()
		appWithoutUser.Post("/histories", controller.CreateHistoryHandler)

		requestBody := entities.History{
			Method: "POST",
			Type:   "login",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/histories", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := appWithoutUser.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
	})

	t.Run("Invalid Body", func(t *testing.T) {
		invalidJSON := []byte(`{"method": "POST", "type": }`)

		req := httptest.NewRequest(http.MethodPost, "/histories", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrBadRequest.Message, result["status"])
		assert.NotNil(t, result["message"])
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		requestBody := entities.History{}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/histories", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrBadRequest.Message, result["status"])
		assert.Equal(t, "Method or Type is missing.", result["message"])
	})

	t.Run("Usecase Error", func(t *testing.T) {
		requestBody := entities.History{
			Method: "POST",
			Type:   "login",
		}
		jsonBody, _ := json.Marshal(requestBody)

		mockUseCase.On("CreateHistory", mock.Anything).Return(&entities.History{}, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/histories", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrInternalServerError.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrInternalServerError.Message, result["status"])
		assert.Equal(t, "database error", result["message"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetHistoryByUserIDHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Get("/histories", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return controller.GetHistoryByUserIDHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		expectedResponse := fiber.Map{
			"data": []entities.History{
				{
					ID:     "history-123",
					Method: "deposit",
					Type:   "saving",
					Money:  1000.0,
					UserID: "user-123",
				},
				{
					ID:     "history-124",
					Method: "withdraw",
					Type:   "expense",
					Money:  500.0,
					UserID: "user-123",
				},
			},
			"total": 500.0,
		}

		mockUseCase.On("GetHistoryByUserID", "user-123").Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/histories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "Retirement retrieved successfully", result["message"])

		resultData := result["result"].(map[string]interface{})
		assert.NotNil(t, resultData["data"])
		assert.Equal(t, 500.0, resultData["total"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing UserID", func(t *testing.T) {
		appWithoutUser := fiber.New()
		appWithoutUser.Get("/histories", controller.GetHistoryByUserIDHandler)

		req := httptest.NewRequest(http.MethodGet, "/histories", nil)
		resp, _ := appWithoutUser.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
		assert.Nil(t, result["result"])
	})

	t.Run("History Not Found", func(t *testing.T) {
		mockUseCase.On("GetHistoryByUserID", "user-123").Return(fiber.Map{}, errors.New("history not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/histories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Equal(t, "history not found", result["message"])
		assert.Nil(t, result["result"])

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetSummaryHistoryByUserIDHandler(t *testing.T) {
	controller, mockUseCase, app := setupTest(t)

	app.Get("/history/summary", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return controller.GetSummaryHistoryByUserIDHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		expectedResponse := map[string]float64{
			"January":  800.0,
			"February": 950.0,
			"March":    1200.0,
			"April":    500.0,
			"Total":    3450.0,
		}

		mockUseCase.On("GetHistoryByMonth", "user-123").Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/history/summary", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Success", result["status"])
		assert.Equal(t, float64(fiber.StatusOK), result["status_code"])
		assert.Equal(t, "Retirement retrieved successfully", result["message"])

		resultData := result["result"].(map[string]interface{})
		assert.Equal(t, 800.0, resultData["January"])
		assert.Equal(t, 950.0, resultData["February"])
		assert.Equal(t, 1200.0, resultData["March"])
		assert.Equal(t, 500.0, resultData["April"])
		assert.Equal(t, 3450.0, resultData["Total"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing UserID", func(t *testing.T) {
		appWithoutUser := fiber.New()
		appWithoutUser.Get("/history/summary", controller.GetSummaryHistoryByUserIDHandler)

		req := httptest.NewRequest(http.MethodGet, "/history/summary", nil)
		resp, _ := appWithoutUser.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Error", result["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
		assert.Nil(t, result["result"])
	})

	t.Run("History Not Found", func(t *testing.T) {
		mockUseCase.On("GetHistoryByMonth", "user-123").Return(map[string]float64{}, errors.New("history not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/history/summary", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
		assert.Equal(t, "history not found", result["message"])
		assert.Nil(t, result["result"])

		mockUseCase.AssertExpectations(t)
	})
}
