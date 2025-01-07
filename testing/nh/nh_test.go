package controllers_test

import (
	"bytes"
	"encoding/json"
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

func (m *MockNhUseCase) GetAllNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
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

func (m *MockNhUseCase) GetInactiveNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.NursingHouse), args.Error(1)
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

func (m *MockNhUseCase) GetNhNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockNhUseCase) UpdateNhByID(id string, nh entities.NursingHouse, files []multipart.FileHeader, deleteImages []string, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	args := m.Called(id, nh, files, deleteImages, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.NursingHouse), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupApp(routeSetup func(*fiber.App)) *fiber.App {
	app := fiber.New()
	routeSetup(app)
	return app
}

func TestNhHandlers(t *testing.T) {
	mockUseCase := new(MockNhUseCase)
	controller := controllers.NewNhController(mockUseCase)

	app := setupApp(func(app *fiber.App) {
		app.Post("/nh", controller.CreateNhHandler)
		app.Get("/nh", controller.GetAllNhHandler)
		app.Put("/nh/:id", controller.UpdateNhByIDHandler) // Add the PUT route for update
	})

	t.Run("CreateNhHandler - Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// Mock Data
		nhData := map[string]string{
			"name":    "Test Nursing Home",
			"address": "123 Test Address",
			"status":  "active",
		}
		for key, value := range nhData {
			_ = writer.WriteField(key, value)
		}
		fileContents := []byte("dummy file content")
		part, _ := writer.CreateFormFile("images", "test.jpg")
		_, _ = part.Write(fileContents)
		writer.Close()

		// Mock Response
		expectedNh := &entities.NursingHouse{
			Name:    "Test Nursing Home",
			Address: "123 Test Address",
			Status:  "active",
		}
		mockUseCase.On("CreateNh", mock.Anything, mock.Anything, mock.Anything).Return(expectedNh, nil)

		req := httptest.NewRequest("POST", "/nh", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNhByIDHandler - Success", func(t *testing.T) {
		// Mock input
		id := "123"
		nursingHouse := entities.NursingHouse{
			Name:    "Updated Nursing Home",
			Address: "456 Updated Address",
			Status:  "active",
		}
		files := []multipart.FileHeader{
			{Filename: "new_image.jpg"},
		}
		deleteImages := []string{"old_image.jpg"}

		expectedNh := &entities.NursingHouse{
			Name:    "Updated Nursing Home",
			Address: "456 Updated Address",
			Status:  "active",
		}

		mockUseCase.On("GetNhByID", id).Return(&entities.NursingHouse{
			Images: []string{"old_image.jpg"},
		}, nil)

		mockUseCase.On("UpdateNhByID", id, nursingHouse, files, deleteImages, mock.Anything).Return(expectedNh, nil)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", nursingHouse.Name)
		_ = writer.WriteField("address", nursingHouse.Address)
		_ = writer.WriteField("status", nursingHouse.Status)
		part, _ := writer.CreateFormFile("images", "new_image.jpg")
		_, _ = part.Write([]byte("new image content"))
		_ = writer.WriteField("delete_images", "old_image.jpg")
		writer.Close()

		req := httptest.NewRequest("PUT", "/nh/123", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var actualNh entities.NursingHouse
		err = json.NewDecoder(resp.Body).Decode(&actualNh)
		assert.NoError(t, err)
		assert.Equal(t, expectedNh, &actualNh)
		mockUseCase.AssertExpectations(t)
	})
}