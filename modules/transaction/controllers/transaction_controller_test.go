package controllers_test

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionsForAllUsersHandler(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	controller := controllers.NewTransactionController(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("CreateTransactionsForAllUsers").Return(nil).Once()

		app := fiber.New()
		app.Post("/transactions/create-all", controller.CreateTransactionsForAllUsersHandler)

		req := httptest.NewRequest("POST", "/transactions/create-all", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		mockUseCase.On("CreateTransactionsForAllUsers").Return(errors.New("no loans found for transaction creation")).Once()

		app := fiber.New()
		app.Post("/transactions/create-all", controller.CreateTransactionsForAllUsersHandler)

		req := httptest.NewRequest("POST", "/transactions/create-all", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMarkTransactiontoPaidHandler(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	controller := controllers.NewTransactionController(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("MarkTransactiontoPaid", "trans-123", "user-123").Return(nil).Once()

		app := fiber.New()
		app.Put("/transactions/:id/mark-paid", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return controller.MarkTransactiontoPaidHandler(c)
		})

		req := httptest.NewRequest("PUT", "/transactions/trans-123/mark-paid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing User ID", func(t *testing.T) {
		app := fiber.New()
		app.Put("/transactions/:id/mark-paid", controller.MarkTransactiontoPaidHandler)

		req := httptest.NewRequest("PUT", "/transactions/trans-123/mark-paid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		mockUseCase.On("MarkTransactiontoPaid", "trans-123", "user-123").Return(errors.New("transaction is not in a payable state")).Once()

		app := fiber.New()
		app.Put("/transactions/:id/mark-paid", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return controller.MarkTransactiontoPaidHandler(c)
		})

		req := httptest.NewRequest("PUT", "/transactions/trans-123/mark-paid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetTransactionByUserIDHandler(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	controller := controllers.NewTransactionController(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockTransactions := []map[string]interface{}{
			{
				"trans_id":   "trans-123",
				"status":     "ชำระ",
				"created_at": time.Now(),
				"loan": map[string]interface{}{
					"id":   "loan-123",
					"name": "Home Loan",
				},
			},
			{
				"trans_id":   "trans-456",
				"status":     "ชำระแล้ว",
				"created_at": time.Now(),
				"loan": map[string]interface{}{
					"id":   "loan-456",
					"name": "Car Loan",
				},
			},
		}

		mockUseCase.On("GetTransactionByUserID", "user-123").Return(mockTransactions, nil).Once()

		app := fiber.New()
		app.Get("/transactions", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return controller.GetTransactionByUserIDHandler(c)
		})

		req := httptest.NewRequest("GET", "/transactions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized - Missing User ID", func(t *testing.T) {
		app := fiber.New()
		app.Get("/transactions", controller.GetTransactionByUserIDHandler)

		req := httptest.NewRequest("GET", "/transactions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		mockUseCase.On("GetTransactionByUserID", "user-123").Return([]map[string]interface{}{}, errors.New("not found")).Once()

		app := fiber.New()
		app.Get("/transactions", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return controller.GetTransactionByUserIDHandler(c)
		})

		req := httptest.NewRequest("GET", "/transactions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
