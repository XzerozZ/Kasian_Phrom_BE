package controllers_test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/controllers"
)

type MockNhUseCase struct {
	mock.Mock
}

func (m *MockNhUseCase) CreateNh(nh entities.NursingHouse, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	args := m.Called(nh, files, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) CreateNhMock(nh entities.NursingHouse, links []string, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	args := m.Called(nh, links, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetNhByID(id string) (*entities.NursingHouse, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetActiveNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetAllNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetInactiveNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetNhNextID() (string, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.(string), args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockNhUseCase) UpdateNhByID(id string, nh entities.NursingHouse, files []multipart.FileHeader, deleteImages []string, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	args := m.Called(id, nh, files, deleteImages, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) GetNhByIDForUser(id, userID string) (*entities.NursingHouse, error) {
	args := m.Called(id, userID)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) RecommendationCosine(userID string) ([]entities.NursingHouse, error) {
	args := m.Called(userID)
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNhUseCase) RecommendationLLM(userID string) ([]entities.NursingHouse, error) {
	args := m.Called(userID)
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNhHandlers(t *testing.T) {
	mockUseCase := new(MockNhUseCase)
	controller := controllers.NewNhController(mockUseCase)
	app := fiber.New()

	// Register all handlers
	app.Post("/nh", controller.CreateNhHandler)
	app.Get("/nh", controller.GetAllNhHandler)
	app.Get("/nh/active", controller.GetAllActiveNhHandler)
	app.Get("/nh/inactive", controller.GetAllInactiveNhHandler)
	app.Get("/nh/next-id", controller.GetNhNextIDHandler)
	app.Get("/nh/:id", controller.GetNhByIDHandler)
	app.Put("/nh/:id", controller.UpdateNhByIDHandler)
	app.Get("/nh/user/:id", controller.GetNhByIDForUserHandler)
	app.Get("/nh/recommend/cosine", controller.GetRecommendCosine)
	app.Get("/nh/recommend/llm", controller.GetRecommendLLM)

	t.Run("CreateNhHandler - Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Test Nursing Home")
		_ = writer.WriteField("address", "123 Test Address")
		_ = writer.WriteField("status", "active")

		part, _ := writer.CreateFormFile("images", "test.jpg")
		_, _ = part.Write([]byte("dummy image content"))
		writer.Close()

		expectedNh := &entities.NursingHouse{
			Name:    "Test Nursing Home",
			Address: "123 Test Address",
			Status:  "active",
		}

		mockUseCase.On("CreateNh", mock.Anything, mock.Anything, mock.Anything).Return(expectedNh, nil).Once()
		req := httptest.NewRequest("POST", "/nh", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateNhHandler - No Images", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Test Nursing Home")
		_ = writer.WriteField("address", "123 Test Address")
		_ = writer.WriteField("status", "active")
		writer.Close()

		req := httptest.NewRequest("POST", "/nh", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		responseBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(responseBody), "At least one image is required")
	})

	t.Run("CreateNhHandler - UseCase Error", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Test Nursing Home")
		_ = writer.WriteField("address", "123 Test Address")
		_ = writer.WriteField("status", "active")

		part, _ := writer.CreateFormFile("images", "test.jpg")
		_, _ = part.Write([]byte("dummy image content"))
		writer.Close()

		mockUseCase.On("CreateNh", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("usecase error")).Once()

		req := httptest.NewRequest("POST", "/nh", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAllNhHandler - Success", func(t *testing.T) {
		nhs := []entities.NursingHouse{
			{Name: "Home 1", Address: "Address 1"},
			{Name: "Home 2", Address: "Address 2"},
		}

		mockUseCase.On("GetAllNh").Return(nhs, nil).Once()

		req := httptest.NewRequest("GET", "/nh", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAllNhHandler - Error", func(t *testing.T) {
		mockUseCase.On("GetAllNh").Return(nil, errors.New("error getting nursing homes")).Once()

		req := httptest.NewRequest("GET", "/nh", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAllActiveNhHandler - Success", func(t *testing.T) {
		nhs := []entities.NursingHouse{
			{Name: "Home 1", Address: "Address 1", Status: "active"},
			{Name: "Home 2", Address: "Address 2", Status: "active"},
		}

		mockUseCase.On("GetActiveNh").Return(nhs, nil).Once()

		req := httptest.NewRequest("GET", "/nh/active", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNhByIDHandler - Success", func(t *testing.T) {
		id := "123"
		nh := &entities.NursingHouse{
			Name:    "Home 1",
			Address: "Address 1",
		}

		mockUseCase.On("GetNhByID", id).Return(nh, nil).Once()

		req := httptest.NewRequest("GET", "/nh/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNhByIDHandler - Not Found", func(t *testing.T) {
		id := "999"

		mockUseCase.On("GetNhByID", id).Return(nil, errors.New("not found")).Once()

		req := httptest.NewRequest("GET", "/nh/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNhByIDHandler - Success", func(t *testing.T) {
		id := "123"

		existingNh := &entities.NursingHouse{
			Name:    "Old Name",
			Address: "Old Address",
			Images:  []entities.Image{{ImageLink: "old_image.jpg"}},
		}
		mockUseCase.On("GetNhByID", id).Return(existingNh, nil).Once()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Updated Nursing Home")
		_ = writer.WriteField("address", "456 Updated Address")
		_ = writer.WriteField("status", "active")
		_ = writer.WriteField("delete_images", "old_image.jpg")

		part, _ := writer.CreateFormFile("images", "new_test.jpg")
		_, _ = part.Write([]byte("dummy image content"))

		writer.Close()

		updatedNh := &entities.NursingHouse{
			Name:    "Updated Nursing Home",
			Address: "456 Updated Address",
			Status:  "active",
		}

		mockUseCase.On("UpdateNhByID",
			id,
			mock.MatchedBy(func(nh entities.NursingHouse) bool {
				return nh.Name == "Updated Nursing Home" &&
					nh.Address == "456 Updated Address" &&
					nh.Status == "active"
			}),
			mock.AnythingOfType("[]multipart.FileHeader"),
			[]string{"old_image.jpg"},
			mock.AnythingOfType("*fiber.Ctx"),
		).Return(updatedNh, nil).Once()

		req := httptest.NewRequest("PUT", "/nh/"+id, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNhByIDHandler - Not Found", func(t *testing.T) {
		id := "999"

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Updated Nursing Home")
		_ = writer.WriteField("delete_images", "old_image.jpg")
		part, _ := writer.CreateFormFile("images", "new_test.jpg")
		_, _ = part.Write([]byte("dummy image content"))
		writer.Close()

		mockUseCase.On("GetNhByID", id).Return(nil, errors.New("not found")).Once()

		req := httptest.NewRequest("PUT", "/nh/"+id, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNhByIDHandler - No Remaining Images", func(t *testing.T) {
		id := "123"

		existingNh := &entities.NursingHouse{
			Name:    "Old Name",
			Address: "Old Address",
			Images:  []entities.Image{{ImageLink: "old_image.jpg"}},
		}
		mockUseCase.On("GetNhByID", id).Return(existingNh, nil).Once()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Updated Nursing Home")
		_ = writer.WriteField("delete_images", "old_image.jpg")
		writer.Close()

		req := httptest.NewRequest("PUT", "/nh/"+id, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNhByIDForUserHandler - Success", func(t *testing.T) {
		id := "123"
		userID := "user456"

		app := fiber.New()
		app.Get("/nh/user/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetNhByIDForUserHandler(c)
		})

		nh := &entities.NursingHouse{
			Name:    "Home 1",
			Address: "Address 1",
		}

		mockUseCase.On("GetNhByIDForUser", id, userID).Return(nh, nil).Once()

		req := httptest.NewRequest("GET", "/nh/user/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNhByIDForUserHandler - No User ID", func(t *testing.T) {
		id := "123"

		app := fiber.New()
		app.Get("/nh/user/:id", controller.GetNhByIDForUserHandler)

		req := httptest.NewRequest("GET", "/nh/user/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("GetRecommendCosine - Success", func(t *testing.T) {
		userID := "user456"

		app := fiber.New()
		app.Get("/recommend", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetRecommendCosine(c)
		})

		nhs := []entities.NursingHouse{
			{Name: "Home 1", Address: "Address 1"},
			{Name: "Home 2", Address: "Address 2"},
		}

		mockUseCase.On("RecommendationCosine", userID).Return(nhs, nil).Once()

		req := httptest.NewRequest("GET", "/recommend", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetRecommendLLM - Success", func(t *testing.T) {
		userID := "user456"

		app := fiber.New()
		app.Get("/recommend", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			return controller.GetRecommendLLM(c)
		})

		nhs := []entities.NursingHouse{
			{Name: "Home 1", Address: "Address 1"},
			{Name: "Home 2", Address: "Address 2"},
		}

		mockUseCase.On("RecommendationLLM", userID).Return(nhs, nil).Once()

		req := httptest.NewRequest("GET", "/recommend", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})
}
