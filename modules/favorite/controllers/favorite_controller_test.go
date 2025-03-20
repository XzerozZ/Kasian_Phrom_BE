package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFavHandlers(t *testing.T) {
	mockUseCase := new(mocks.MockFavUseCase)
	controller := controllers.NewFavController(mockUseCase)

	t.Run("CreateFavHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Post("/favorites", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CreateFavHandler(c)
		})

		reqBody := &entities.Favorite{
			NursingHouseID: "00001",
		}

		reqBodyJSON, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		mockUseCase.On("CreateFav", mock.MatchedBy(func(a *entities.Favorite) bool {
			return a != nil && a.NursingHouseID == "00001" && a.UserID == "user_123"
		})).Return(nil).Once()

		req := httptest.NewRequest("POST", "/favorites", bytes.NewReader(reqBodyJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Favorite successfully", responseMap["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateFavHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Post("/favorites", controller.CreateFavHandler)

		fav := entities.Favorite{
			NursingHouseID: "nursing_house_123",
		}

		reqBody, _ := json.Marshal(fav)
		req := httptest.NewRequest("POST", "/favorites", bytes.NewReader(reqBody))
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

	t.Run("CreateFavHandler - Missing NursingHouseID", func(t *testing.T) {
		app := fiber.New()
		app.Post("/favorites", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CreateFavHandler(c)
		})

		fav := entities.Favorite{}

		reqBody, _ := json.Marshal(fav)
		req := httptest.NewRequest("POST", "/favorites", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.ErrBadRequest.Code, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, fiber.ErrBadRequest.Message, responseMap["status"])
		assert.Equal(t, float64(fiber.ErrBadRequest.Code), responseMap["status_code"])
		assert.Equal(t, "NursingHouseID is missing", responseMap["message"])
		assert.Nil(t, responseMap["result"])
	})

	t.Run("CreateFavHandler - InternalError", func(t *testing.T) {
		app := fiber.New()
		app.Post("/favorites", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CreateFavHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		fav := entities.Favorite{
			NursingHouseID: "nursing_house_123",
		}

		mockUseCase.On("CreateFav", mock.MatchedBy(func(f *entities.Favorite) bool {
			return f.UserID == "user_123" && f.NursingHouseID == "nursing_house_123"
		})).Return(errors.New("database error")).Once()

		reqBody, _ := json.Marshal(fav)
		req := httptest.NewRequest("POST", "/favorites", bytes.NewReader(reqBody))
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

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetFavByUserIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.GetFavByUserIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		favs := []entities.Favorite{
			{UserID: "user_123", NursingHouseID: "nursing_house_123"},
			{UserID: "user_123", NursingHouseID: "nursing_house_456"},
		}

		mockUseCase.On("GetFavByUserID", "user_123").Return(favs, nil).Once()
		req := httptest.NewRequest("GET", "/favorites", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Favorites retrieved successfully", responseMap["message"])
		assert.NotNil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetFavByUserIDHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites", controller.GetFavByUserIDHandler)
		req := httptest.NewRequest("GET", "/favorites", nil)

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

	t.Run("GetFavByUserIDHandler - NotFound", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.GetFavByUserIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("GetFavByUserID", "user_123").Return([]entities.Favorite{}, errors.New("not found")).Once()
		req := httptest.NewRequest("GET", "/favorites", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "No favorites found for this user", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CheckFavHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CheckFavHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("CheckFav", "user_123", "nursing_house_123").Return(nil).Once()
		req := httptest.NewRequest("GET", "/favorites/nursing_house_123", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Favorited Nursing House", responseMap["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CheckFavHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites/:nh_id", controller.CheckFavHandler)
		req := httptest.NewRequest("GET", "/favorites/nursing_house_123", nil)

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

	t.Run("CheckFavHandler - NotFavorited", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CheckFavHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("CheckFav", "user_123", "nursing_house_123").Return(errors.New("not favorited nursing house")).Once()
		req := httptest.NewRequest("GET", "/favorites/nursing_house_123", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "Not Favorited Nursing House", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CheckFavHandler - InternalError", func(t *testing.T) {
		app := fiber.New()
		app.Get("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.CheckFavHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("CheckFav", "user_123", "nursing_house_123").Return(errors.New("database error")).Once()
		req := httptest.NewRequest("GET", "/favorites/nursing_house_123", nil)

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

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteFavByIDHandler - Success", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.DeleteFavByIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("DeleteFavByID", "user_123", "nursing_house_123").Return(nil).Once()
		req := httptest.NewRequest("DELETE", "/favorites/nursing_house_123", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Success", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusOK), responseMap["status_code"])
		assert.Equal(t, "Favorite deleted successfully", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteFavByIDHandler - Unauthorized", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/favorites/:nh_id", controller.DeleteFavByIDHandler)
		req := httptest.NewRequest("DELETE", "/favorites/nursing_house_123", nil)

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

	t.Run("DeleteFavByIDHandler - NotFound", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.DeleteFavByIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("DeleteFavByID", "user_123", "nursing_house_123").Return(errors.New("record not found")).Once()
		req := httptest.NewRequest("DELETE", "/favorites/nursing_house_123", nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var responseMap map[string]interface{}
		responseBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(responseBody, &responseMap)

		assert.Equal(t, "Not Found", responseMap["status"])
		assert.Equal(t, float64(fiber.StatusNotFound), responseMap["status_code"])
		assert.Equal(t, "Favorite not found", responseMap["message"])
		assert.Nil(t, responseMap["result"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteFavByIDHandler - InternalError", func(t *testing.T) {
		app := fiber.New()
		app.Delete("/favorites/:nh_id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user_123")
			return controller.DeleteFavByIDHandler(c)
		})

		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		mockUseCase.On("DeleteFavByID", "user_123", "nursing_house_123").Return(errors.New("database error")).Once()
		req := httptest.NewRequest("DELETE", "/favorites/nursing_house_123", nil)

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

		mockUseCase.AssertExpectations(t)
	})
}
