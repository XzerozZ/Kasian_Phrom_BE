package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRetirementHandlers(t *testing.T) {
	mockUseCase := new(mocks.MockRetirementUseCase)
	controller := controllers.NewRetirementController(mockUseCase)
	app := fiber.New()
	app.Post("/retirement", controller.CreateRetirementHandler)
	app.Get("/retirement", controller.GetRetirementByUserIDHandler)
	app.Put("/retirement", controller.UpdateRetirementHandler)

	t.Run("CreateRetirementHandler - Success", func(t *testing.T) {
		retirement := entities.RetirementPlan{
			PlanName:       "Test Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
		}

		jsonData, _ := json.Marshal(retirement)
		expectedRetirement := &entities.RetirementPlan{
			PlanName:       "Test Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
			UserID:         "user123",
		}

		expectedAge := 33
		mockUseCase.On("CreateRetirement", mock.MatchedBy(func(r entities.RetirementPlan) bool {
			return r.PlanName == retirement.PlanName &&
				r.BirthDate == retirement.BirthDate &&
				r.RetirementAge == retirement.RetirementAge &&
				r.ExpectLifespan == retirement.ExpectLifespan
		})).Return(expectedRetirement, expectedAge, nil).Once()

		req := httptest.NewRequest("POST", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Post("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateRetirementHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset created successfully", response["message"])
		result, ok := response["result"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotNil(t, result["plan"])
		assert.Equal(t, float64(expectedAge), result["age"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateRetirementHandler - No User ID", func(t *testing.T) {
		retirement := entities.RetirementPlan{
			PlanName:       "Test Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
		}

		jsonData, _ := json.Marshal(retirement)
		req := httptest.NewRequest("POST", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Error", response["status"])
		assert.Contains(t, response["message"], "Unauthorized")
	})

	t.Run("CreateRetirementHandler - Missing PlanName", func(t *testing.T) {
		retirement := entities.RetirementPlan{
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
		}

		jsonData, _ := json.Marshal(retirement)
		req := httptest.NewRequest("POST", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Post("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateRetirementHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "PlanName is missing")
	})

	t.Run("CreateRetirementHandler - UseCase Error", func(t *testing.T) {
		retirement := entities.RetirementPlan{
			PlanName:       "Test Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
		}

		jsonData, _ := json.Marshal(retirement)
		mockUseCase.On("CreateRetirement", mock.Anything).Return(nil, 0, errors.New("usecase error")).Once()
		req := httptest.NewRequest("POST", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Post("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateRetirementHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response["status"])
		assert.Contains(t, response["message"], "usecase error")
		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetRetirementByUserIDHandler - Success", func(t *testing.T) {
		userID := "user123"
		expectedRetirement := &entities.RetirementPlan{
			PlanName:       "Test Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  65,
			ExpectLifespan: 85,
			UserID:         userID,
		}

		mockUseCase.On("GetRetirementByUserID", userID).Return(expectedRetirement, nil).Once()
		req := httptest.NewRequest("GET", "/retirement", nil)
		app := fiber.New()
		app.Get("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetRetirementByUserIDHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Retirement retrieved successfully", response["message"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetRetirementByUserIDHandler - No User ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/retirement", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("GetRetirementByUserIDHandler - Not Found", func(t *testing.T) {
		userID := "user123"
		mockUseCase.On("GetRetirementByUserID", userID).Return(nil, errors.New("not found")).Once()
		req := httptest.NewRequest("GET", "/retirement", nil)
		app := fiber.New()
		app.Get("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetRetirementByUserIDHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateRetirementHandler - Success", func(t *testing.T) {
		userID := "user123"
		retirement := entities.RetirementPlan{
			PlanName:       "Updated Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  67,
			ExpectLifespan: 90,
		}

		jsonData, _ := json.Marshal(retirement)
		expectedRetirement := &entities.RetirementPlan{
			PlanName:       "Updated Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  67,
			ExpectLifespan: 90,
			UserID:         userID,
		}

		mockUseCase.On("UpdateRetirementByID", userID, mock.MatchedBy(func(r entities.RetirementPlan) bool {
			return r.PlanName == retirement.PlanName &&
				r.BirthDate == retirement.BirthDate &&
				r.RetirementAge == retirement.RetirementAge &&
				r.ExpectLifespan == retirement.ExpectLifespan
		})).Return(expectedRetirement, nil).Once()

		req := httptest.NewRequest("PUT", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Put("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.UpdateRetirementHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset created successfully", response["message"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateRetirementHandler - No User ID", func(t *testing.T) {
		retirement := entities.RetirementPlan{
			PlanName:       "Updated Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  67,
			ExpectLifespan: 90,
		}

		jsonData, _ := json.Marshal(retirement)
		req := httptest.NewRequest("PUT", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UpdateRetirementHandler - UseCase Error", func(t *testing.T) {
		userID := "user123"
		retirement := entities.RetirementPlan{
			PlanName:       "Updated Plan",
			BirthDate:      "01-01-1990",
			RetirementAge:  67,
			ExpectLifespan: 90,
		}

		jsonData, _ := json.Marshal(retirement)
		mockUseCase.On("UpdateRetirementByID", userID, mock.Anything).Return(nil, errors.New("usecase error")).Once()
		req := httptest.NewRequest("PUT", "/retirement", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Put("/retirement", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.UpdateRetirementHandler(c)
		})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
