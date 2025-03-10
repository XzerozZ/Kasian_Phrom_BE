package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuizUseCase struct {
	mock.Mock
}

func (m *MockQuizUseCase) CreateQuiz(userID string, weights []int) (*entities.Quiz, error) {
	args := m.Called(userID, weights)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}

func (m *MockQuizUseCase) GetQuizByUserID(userID string) (*entities.Quiz, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Quiz), args.Error(1)
}

func setupTestApp() (*fiber.App, *MockQuizUseCase) {
	app := fiber.New()
	mockUseCase := new(MockQuizUseCase)
	controller := controllers.NewQuizController(mockUseCase)

	app.Post("/quiz", func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-id")
		return controller.CreateQuizHandler(c)
	})

	app.Post("/quiz/no-auth", controller.CreateQuizHandler)

	app.Get("/quiz", func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-id")
		return controller.GetQuizByUserIDHandler(c)
	})

	app.Get("/quiz/no-auth", controller.GetQuizByUserIDHandler)

	return app, mockUseCase
}

func createMultipartForm(weights []string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, weight := range weights {
		if err := writer.WriteField("weight", weight); err != nil {
			return nil, "", err
		}
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, contentType, nil
}

func TestCreateQuizHandler_Success(t *testing.T) {
	app, mockUseCase := setupTestApp()
	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	expectedQuiz := &entities.Quiz{
		UserID: "test-user-id",
	}

	mockUseCase.On("CreateQuiz", "test-user-id", weights).Return(expectedQuiz, nil)
	weightStrings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	body, contentType, err := createMultipartForm(weightStrings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/quiz", body)
	req.Header.Set("Content-Type", contentType)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Success", result["status"])
	assert.Equal(t, float64(200), result["status_code"])
	assert.Equal(t, "Quiz created successfully", result["message"])

	mockUseCase.AssertExpectations(t)
}

func TestCreateQuizHandler_Unauthorized(t *testing.T) {
	app, _ := setupTestApp()

	weightStrings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	body, contentType, err := createMultipartForm(weightStrings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/quiz/no-auth", body)
	req.Header.Set("Content-Type", contentType)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
	assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
}

func TestCreateQuizHandler_InvalidFormData(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/quiz", strings.NewReader("Invalid form data"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
	assert.Equal(t, "Failed to parse form data", result["message"])
}

func TestCreateQuizHandler_MissingWeights(t *testing.T) {
	app, _ := setupTestApp()

	weightStrings := []string{"1", "2", "3", "4", "5"}
	body, contentType, err := createMultipartForm(weightStrings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/quiz", body)
	req.Header.Set("Content-Type", contentType)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
	assert.Equal(t, "Must answer 12 quiz", result["message"])
}

func TestCreateQuizHandler_InvalidWeightValue(t *testing.T) {
	app, _ := setupTestApp()

	weightStrings := []string{"1", "2", "3", "4", "not-a-number", "6", "7", "8", "9", "10"}
	body, contentType, err := createMultipartForm(weightStrings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/quiz", body)
	req.Header.Set("Content-Type", contentType)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusBadRequest), result["status_code"])
	assert.Equal(t, "Invalid weight value at position 5", result["message"])
}

func TestCreateQuizHandler_UseCaseError(t *testing.T) {
	app, mockUseCase := setupTestApp()
	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	expectedError := errors.New("database error")
	mockUseCase.On("CreateQuiz", "test-user-id", weights).Return(nil, expectedError)

	weightStrings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	body, contentType, err := createMultipartForm(weightStrings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/quiz", body)
	req.Header.Set("Content-Type", contentType)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusInternalServerError), result["status_code"])
	assert.Equal(t, "database error", result["message"])

	mockUseCase.AssertExpectations(t)
}

func TestGetQuizByUserIDHandler_Success(t *testing.T) {
	app, mockUseCase := setupTestApp()

	expectedQuiz := &entities.Quiz{
		UserID: "test-user-id",
	}

	mockUseCase.On("GetQuizByUserID", "test-user-id").Return(expectedQuiz, nil)

	req := httptest.NewRequest(http.MethodGet, "/quiz", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Success", result["status"])
	assert.Equal(t, float64(200), result["status_code"])
	assert.Equal(t, "Quiz retrieved successfully", result["message"])

	mockUseCase.AssertExpectations(t)
}

func TestGetQuizByUserIDHandler_Unauthorized(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/quiz/no-auth", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Error", result["status"])
	assert.Equal(t, float64(fiber.StatusUnauthorized), result["status_code"])
	assert.Equal(t, "Unauthorized: Missing user ID", result["message"])
}

func TestGetQuizByUserIDHandler_NotFound(t *testing.T) {
	app, mockUseCase := setupTestApp()

	expectedError := errors.New("quiz not found")
	mockUseCase.On("GetQuizByUserID", "test-user-id").Return(nil, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/quiz", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, fiber.ErrNotFound.Message, result["status"])
	assert.Equal(t, float64(fiber.ErrNotFound.Code), result["status_code"])
	assert.Equal(t, "This user has not answered quiz yet.", result["message"])

	mockUseCase.AssertExpectations(t)
}
