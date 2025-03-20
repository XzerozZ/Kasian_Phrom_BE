package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	controller "github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanHandlers(t *testing.T) {
	mockLoanUseCase := new(mocks.MockLoanUseCase)
	mockTransUseCase := new(mocks.MockTransactionUseCase)
	controller := controller.NewLoanController(mockLoanUseCase, mockTransUseCase)

	t.Run("CreateLoanHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Post("/loans", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateLoanHandler(c)
		})

		mockLoan := &entities.Loan{
			ID:     "loan123",
			Name:   "Home Loan",
			Type:   "Mortgage",
			UserID: "user123",
		}

		mockLoanUseCase.On("CreateLoan", mock.MatchedBy(func(l entities.Loan) bool {
			return l.Name == "Home Loan" && l.Type == "Mortgage" && l.UserID == "user123"
		})).Return(mockLoan, nil).Once()

		loanJSON, _ := json.Marshal(fiber.Map{
			"name": "Home Loan",
			"type": "Mortgage",
		})

		req := httptest.NewRequest("POST", "/loans", bytes.NewBuffer(loanJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Loan created successfully", responseMap["message"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("CreateLoanHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Post("/loans-no-auth", controller.CreateLoanHandler)

		loanJSON, _ := json.Marshal(fiber.Map{
			"name": "Home Loan",
			"type": "Mortgage",
		})

		req := httptest.NewRequest("POST", "/loans-no-auth", bytes.NewBuffer(loanJSON))
		req.Header.Set("Content-Type", "application/json")
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

	t.Run("CreateLoanHandler - Invalid Data", func(t *testing.T) {
		app := fiber.New()
		app.Post("/loans-invalid", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateLoanHandler(c)
		})

		invalidLoanJSON, _ := json.Marshal(fiber.Map{
			"type": "Mortgage",
		})

		req := httptest.NewRequest("POST", "/loans-invalid", bytes.NewBuffer(invalidLoanJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, fiber.ErrBadRequest.Message, responseMap["status"])
		assert.Equal(t, float64(fiber.ErrBadRequest.Code), responseMap["status_code"])
		assert.Equal(t, "Name or Type is missing.", responseMap["message"])
		assert.Nil(t, responseMap["result"])
	})

	t.Run("CreateLoanHandler - Database Error", func(t *testing.T) {
		app := fiber.New()
		app.Post("/loans-error", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateLoanHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("CreateLoan", mock.MatchedBy(func(l entities.Loan) bool {
			return l.Name == "Failed Loan" && l.Type == "Error" && l.UserID == "user123"
		})).Return(nil, errors.New("database error")).Once()

		errorLoanJSON, _ := json.Marshal(fiber.Map{
			"name": "Failed Loan",
			"type": "Error",
		})

		req := httptest.NewRequest("POST", "/loans-error", bytes.NewBuffer(errorLoanJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Internal Server Error", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusInternalServerError), responseMap["status_code"])
		assert.Equal(t, "database error", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("GetLoanByIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Get("/loans/:id", controller.GetLoanByIDHandler)

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoan := &entities.Loan{
			ID:     "loan123",
			Name:   "Home Loan",
			Type:   "Mortgage",
			UserID: "user123",
		}

		mockLoanUseCase.On("GetLoanByID", "loan123").Return(mockLoan, nil).Once()

		req := httptest.NewRequest("GET", "/loans/loan123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Asset retrieved successfully", responseMap["message"])
		assert.NotNil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("GetLoanByIDHandler - Not Found", func(t *testing.T) {
		app := fiber.New()
		app.Get("/loans/:id", controller.GetLoanByIDHandler)

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("GetLoanByID", "nonexistent").Return(nil, errors.New("loan not found")).Once()

		req := httptest.NewRequest("GET", "/loans/nonexistent", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "loan not found", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("GetLoanByUserIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Get("/user/loans", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.GetLoanByUserIDHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoans := []entities.Loan{
			{ID: "loan123", Name: "Home Loan", Type: "Mortgage", UserID: "user123"},
			{ID: "loan456", Name: "Car Loan", Type: "Auto", UserID: "user123"},
		}

		loanSummary := map[string]interface{}{
			"total_loan":               float64(2),
			"total_amount":             18000.0,
			"total_transaction_amount": 1500.0,
		}

		mockLoanUseCase.On("GetLoanByUserID", "user123").Return(mockLoans, loanSummary, nil).Once()

		req := httptest.NewRequest("GET", "/user/loans", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Loans retrieved successfully", responseMap["message"])
		assert.NotNil(t, responseMap["result"])

		result := responseMap["result"].(map[string]interface{})
		loansList, ok := result["loans"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, loansList, 2)

		for _, item := range loansList {
			loan, ok := item.(map[string]interface{})
			assert.True(t, ok)
			assert.Contains(t, loan, "id")
			assert.Contains(t, loan, "name")
			assert.Contains(t, loan, "type")
			assert.Contains(t, loan, "userID")
		}

		summary, ok := result["summary"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, loanSummary, summary)

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("GetLoanByUserIDHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Get("/user/loans-no-auth", controller.GetLoanByUserIDHandler)

		req := httptest.NewRequest("GET", "/user/loans-no-auth", nil)
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

	t.Run("GetLoanByUserIDHandler - Not Found", func(t *testing.T) {
		app := fiber.New()
		app.Get("/user/loans-empty", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user456")
			return controller.GetLoanByUserIDHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("GetLoanByUserID", mock.Anything).Return([]entities.Loan{}, map[string]interface{}{}, errors.New("no loans found")).Once()

		req := httptest.NewRequest("GET", "/user/loans-empty", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "no loans found", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("UpdateLoanStatusByIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Put("/loans/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.UpdateLoanStatusByIDHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockUpdatedLoan := &entities.Loan{
			ID:     "loan123",
			Name:   "Home Loan",
			Type:   "Mortgage",
			Status: "Approved",
			UserID: "user123",
		}

		mockLoanUseCase.On("UpdateLoanStatusByID", "loan123", mock.MatchedBy(func(l entities.Loan) bool {
			return l.Status == "Approved" && l.UserID == "user123"
		})).Return(mockUpdatedLoan, nil).Once()

		updateJSON, _ := json.Marshal(fiber.Map{
			"status": "Approved",
		})

		req := httptest.NewRequest("PUT", "/loans/loan123", bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Asset update successfully", responseMap["message"])
		assert.NotNil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("UpdateLoanStatusByIDHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Put("/loans/:id", controller.UpdateLoanStatusByIDHandler)

		updateJSON, _ := json.Marshal(fiber.Map{
			"status": "Approved",
		})

		req := httptest.NewRequest("PUT", "/loans/loan123", bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
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

	t.Run("UpdateLoanStatusByIDHandler - Invalid Body", func(t *testing.T) {
		app := fiber.New()
		app.Put("/loans/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.UpdateLoanStatusByIDHandler(c)
		})

		invalidJSON := []byte(`{"status": Invalid JSON}`)

		req := httptest.NewRequest("PUT", "/loans/loan123", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, fiber.ErrNotFound.Message, responseMap["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), responseMap["status_code"])
		assert.NotEqual(t, "", responseMap["message"])
		assert.Nil(t, responseMap["result"])
	})

	t.Run("UpdateLoanStatusByIDHandler - Not Found", func(t *testing.T) {
		app := fiber.New()
		app.Put("/loans/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.UpdateLoanStatusByIDHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("UpdateLoanStatusByID", "nonexistent", mock.MatchedBy(func(l entities.Loan) bool {
			return l.Status == "Approved" && l.UserID == "user123"
		})).Return(nil, errors.New("loan not found")).Once()

		updateJSON, _ := json.Marshal(fiber.Map{
			"status": "Approved",
		})

		req := httptest.NewRequest("PUT", "/loans/nonexistent", bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.ErrNotFound.Code, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, fiber.ErrNotFound.Message, responseMap["status"])
		assert.Equal(t, float64(fiber.ErrNotFound.Code), responseMap["status_code"])
		assert.Equal(t, "loan not found", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("DeleteLoanHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/loans/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.DeleteLoanHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("DeleteLoanByID", "loan123").Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/loans/loan123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Loan deleted successfully", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})

	t.Run("DeleteLoanHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/loans/:id", controller.DeleteLoanHandler)

		req := httptest.NewRequest("DELETE", "/loans/loan123", nil)
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

	t.Run("DeleteLoanHandler - Error", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/loans/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.DeleteLoanHandler(c)
		})

		mockLoanUseCase.ExpectedCalls = nil
		mockLoanUseCase.Calls = nil

		mockLoanUseCase.On("DeleteLoanByID", "loan123").Return(errors.New("database error")).Once()

		req := httptest.NewRequest("DELETE", "/loans/loan123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Error", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusInternalServerError), responseMap["status_code"])
		assert.Equal(t, "database error", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockLoanUseCase.AssertExpectations(t)
	})
}
