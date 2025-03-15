package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssetHandlers(t *testing.T) {
	mockUseCase := new(mocks.MockAssetUseCase)
	controller := controllers.NewAssetController(mockUseCase)

	t.Run("CreateAssetHandler - Success", func(t *testing.T) {
		asset := entities.Asset{
			Name:      "House",
			Type:      "Car",
			TotalCost: 10000,
			EndYear:   "2026",
		}

		jsonData, _ := json.Marshal(asset)
		expectedAsset := &entities.Asset{
			ID:        "ASSET001",
			Name:      "House",
			Type:      "Car",
			TotalCost: 10000,
			EndYear:   "2026",
			UserID:    "user123",
		}

		mockUseCase.On("CreateAsset", mock.MatchedBy(func(a entities.Asset) bool {
			return a.Name == asset.Name &&
				a.TotalCost == asset.TotalCost &&
				a.EndYear == asset.EndYear &&
				a.UserID == "user123"
		})).Return(expectedAsset, nil).Once()

		app := fiber.New()
		app.Post("/assets", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateAssetHandler(c)
		})

		req := httptest.NewRequest("POST", "/assets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset created successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateAssetHandler - No User ID", func(t *testing.T) {
		asset := entities.Asset{
			Name:      "House",
			Type:      "Car",
			TotalCost: 10000,
			EndYear:   "2026",
		}

		jsonData, _ := json.Marshal(asset)

		app := fiber.New()
		app.Post("/assets", controller.CreateAssetHandler)

		req := httptest.NewRequest("POST", "/assets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Error", response["status"])
		assert.Contains(t, response["message"], "Unauthorized")
	})

	t.Run("CreateAssetHandler - Missing Name", func(t *testing.T) {
		asset := entities.Asset{
			TotalCost: 10000,
			Type:      "Car",
			EndYear:   "2026",
		}

		jsonData, _ := json.Marshal(asset)

		app := fiber.New()
		app.Post("/assets", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateAssetHandler(c)
		})

		req := httptest.NewRequest("POST", "/assets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "Name, Type or EndYear is missing.")
	})

	t.Run("CreateAssetHandler - UseCase Error", func(t *testing.T) {
		asset := entities.Asset{
			Name:      "House",
			Type:      "Car",
			TotalCost: 10000,
			EndYear:   "2026",
		}

		jsonData, _ := json.Marshal(asset)

		mockUseCase.On("CreateAsset", mock.Anything).Return(nil, errors.New("usecase error")).Once()

		app := fiber.New()
		app.Post("/assets", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.CreateAssetHandler(c)
		})

		req := httptest.NewRequest("POST", "/assets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response["status"])
		assert.Contains(t, response["message"], "usecase error")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAssetByUserIDHandler - Success", func(t *testing.T) {
		userID := "user123"
		expectedAssets := []entities.Asset{
			{
				ID:        "ASSET001",
				Name:      "House",
				Type:      "Car",
				TotalCost: 10000,
				EndYear:   "2026",
				UserID:    userID,
			},
			{
				ID:        "ASSET002",
				Name:      "Car",
				TotalCost: 5000,
				EndYear:   "2024",
				UserID:    userID,
			},
		}

		mockUseCase.On("GetAssetByUserID", userID).Return(expectedAssets, nil).Once()

		app := fiber.New()
		app.Get("/assets", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetAssetByUserIDHandler(c)
		})

		req := httptest.NewRequest("GET", "/assets", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAssetByUserIDHandler - No User ID", func(t *testing.T) {
		app := fiber.New()
		app.Get("/assets", controller.GetAssetByUserIDHandler)

		req := httptest.NewRequest("GET", "/assets", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("GetAssetByUserIDHandler - Not Found", func(t *testing.T) {
		userID := "user123"

		mockUseCase.On("GetAssetByUserID", userID).Return(nil, errors.New("not found")).Once()

		app := fiber.New()
		app.Get("/assets", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetAssetByUserIDHandler(c)
		})

		req := httptest.NewRequest("GET", "/assets", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAssetByIDHandler - Success", func(t *testing.T) {
		assetID := "ASSET001"
		expectedAsset := &entities.Asset{
			ID:        assetID,
			Name:      "House",
			Type:      "Car",
			TotalCost: 10000,
			EndYear:   "2026",
			UserID:    "user123",
		}

		mockUseCase.On("GetAssetByID", assetID).Return(expectedAsset, nil).Once()

		app := fiber.New()
		app.Get("/assets/:id", controller.GetAssetByIDHandler)

		req := httptest.NewRequest("GET", "/assets/"+assetID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAssetByIDHandler - Not Found", func(t *testing.T) {
		assetID := "NONEXISTENT"

		mockUseCase.On("GetAssetByID", assetID).Return(nil, errors.New("not found")).Once()

		app := fiber.New()
		app.Get("/assets/:id", controller.GetAssetByIDHandler)

		req := httptest.NewRequest("GET", "/assets/"+assetID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateAssetHandler - Success", func(t *testing.T) {
		assetID := "ASSET001"
		asset := entities.Asset{
			Name:      "Updated House",
			Type:      "Car",
			TotalCost: 12000,
			EndYear:   "2027",
		}

		jsonData, _ := json.Marshal(asset)
		expectedAsset := &entities.Asset{
			ID:        assetID,
			Name:      "Updated House",
			Type:      "Car",
			TotalCost: 12000,
			EndYear:   "2027",
			UserID:    "user123",
		}

		mockUseCase.On("UpdateAssetByID", assetID, mock.MatchedBy(func(a entities.Asset) bool {
			return a.Name == asset.Name &&
				a.Type == asset.Type &&
				a.TotalCost == asset.TotalCost &&
				a.EndYear == asset.EndYear
		})).Return(expectedAsset, nil).Once()

		app := fiber.New()
		app.Put("/assets/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.UpdateAssetByIDHandler(c)
		})

		req := httptest.NewRequest("PUT", "/assets/"+assetID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset update successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateAssetHandler - No User ID", func(t *testing.T) {
		assetID := "ASSET001"
		asset := entities.Asset{
			Name:      "Updated House",
			Type:      "Car",
			TotalCost: 12000,
			EndYear:   "2027",
		}

		jsonData, _ := json.Marshal(asset)

		app := fiber.New()
		app.Put("/assets/:id", controller.UpdateAssetByIDHandler)

		req := httptest.NewRequest("PUT", "/assets/"+assetID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UpdateAssetHandler - UseCase Error", func(t *testing.T) {
		assetID := "ASSET001"
		asset := entities.Asset{
			Name:      "Updated House",
			Type:      "Car",
			TotalCost: 12000,
			EndYear:   "2027",
		}

		jsonData, _ := json.Marshal(asset)

		mockUseCase.On("UpdateAssetByID", assetID, mock.Anything).Return(nil, errors.New("usecase error")).Once()

		app := fiber.New()
		app.Put("/assets/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user123")
			return controller.UpdateAssetByIDHandler(c)
		})

		req := httptest.NewRequest("PUT", "/assets/"+assetID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteAssetHandler - Success", func(t *testing.T) {
		assetID := "ASSET001"
		userID := "user123"

		mockUseCase.On("DeleteAssetByID", assetID, userID, mock.Anything).Return(nil).Once()

		app := fiber.New()
		app.Delete("/assets/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.DeleteAssetByIDHandler(c)
		})

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("type", "Stock")
		_ = writer.WriteField("name", "Apple Inc")
		_ = writer.WriteField("amount", "1500.50")
		writer.Close()

		req := httptest.NewRequest("DELETE", "/assets/"+assetID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "Asset deleted successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteAssetHandler - With Transfer Data", func(t *testing.T) {
		assetID := "ASSET001"
		userID := "user123"

		mockUseCase.On("DeleteAssetByID", assetID, userID, mock.MatchedBy(func(t []entities.TransferRequest) bool {
			return len(t) == 1 && t[0].Type == "asset" && t[0].Name == "Investment" && t[0].Amount == 500
		})).Return(nil).Once()

		app := fiber.New()
		app.Delete("/assets/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.DeleteAssetByIDHandler(c)
		})

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("type", "asset")
		_ = writer.WriteField("name", "Investment")
		_ = writer.WriteField("amount", "500")
		writer.Close()

		req := httptest.NewRequest("DELETE", "/assets/"+assetID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteAssetHandler - No User ID", func(t *testing.T) {
		assetID := "ASSET001"

		app := fiber.New()
		app.Delete("/assets/:id", controller.DeleteAssetByIDHandler)

		req := httptest.NewRequest("DELETE", "/assets/"+assetID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("DeleteAssetHandler - UseCase Error", func(t *testing.T) {
		assetID := "ASSET001"
		userID := "user123"

		mockUseCase.On("DeleteAssetByID", assetID, userID, mock.Anything).Return(errors.New("usecase error")).Once()

		app := fiber.New()
		app.Delete("/assets/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.DeleteAssetByIDHandler(c)
		})

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("type", "asset")
		_ = writer.WriteField("name", "Investment")
		_ = writer.WriteField("amount", "500")
		writer.Close()

		req := httptest.NewRequest("DELETE", "/assets/"+assetID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})
}
