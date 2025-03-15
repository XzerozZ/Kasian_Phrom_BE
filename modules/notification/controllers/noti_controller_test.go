package controllers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNotiHandlers(t *testing.T) {
	mockUseCase := new(mocks.MockNotiUseCase)
	controller := controllers.NewNotiController(mockUseCase)

	t.Run("GetNotificationsByUserIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Get("/notifications", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.GetNotificationsByUserIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockNotifications := []entities.Notification{
			{ID: "1", UserID: "user_123", Message: "Test notification 1", IsRead: false},
			{ID: "2", UserID: "user_123", Message: "Test notification 2", IsRead: true},
		}

		mockUseCase.On("GetNotificationsByUserID", "user_123").Return(mockNotifications, nil).Once()
		req := httptest.NewRequest("GET", "/notifications", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Notification retrieved successfully", responseMap["message"])
		assert.NotNil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNotificationsByUserIDHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Get("/notifications", controller.GetNotificationsByUserIDHandler)
		req := httptest.NewRequest("GET", "/notifications", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Error", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), responseMap["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", responseMap["message"])
		assert.Nil(t, responseMap["result"])
	})

	t.Run("GetNotificationsByUserIDHandler - NotFound", func(t *testing.T) {
		app := fiber.New()
		app.Get("/notifications", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.GetNotificationsByUserIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("GetNotificationsByUserID", "user_123").Return([]entities.Notification{}, errors.New("no notifications found")).Once()
		req := httptest.NewRequest("GET", "/notifications", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "No Notification found for this user", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("MarkAsReadHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Put("/notifications/mark-read", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.MarkAsReadHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("MarkNotificationsAsRead", "user_123").Return(nil).Once()
		req := httptest.NewRequest("PUT", "/notifications/mark-read", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Read notification successfully", responseMap["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("MarkAsReadHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Put("/notifications/mark-read", controller.MarkAsReadHandler)
		req := httptest.NewRequest("PUT", "/notifications/mark-read", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Error", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusUnauthorized), responseMap["status_code"])
		assert.Equal(t, "Unauthorized: Missing user ID", responseMap["message"])
		assert.Nil(t, responseMap["result"])
	})

	t.Run("MarkAsReadHandler - NotFound", func(t *testing.T) {
		app := fiber.New()
		app.Put("/notifications/mark-read", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.MarkAsReadHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("MarkNotificationsAsRead", "user_123").Return(errors.New("notifications not found")).Once()
		req := httptest.NewRequest("PUT", "/notifications/mark-read", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, fiber.ErrNotFound.Message, responseMap["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), responseMap["status_code"])
		assert.Equal(t, "notifications not found", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})
}
