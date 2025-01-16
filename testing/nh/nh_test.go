package controllers_test

import (
	"bytes"
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

func TestNhHandlers(t *testing.T) {
	mockUseCase := new(MockNhUseCase)
	controller := controllers.NewNhController(mockUseCase)
	app := fiber.New()
	app.Post("/nh", controller.CreateNhHandler)
	app.Put("/nh/:id", controller.UpdateNhByIDHandler)

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

		mockUseCase.On("CreateNh", mock.Anything, mock.Anything, mock.Anything).Return(expectedNh, nil)
		req := httptest.NewRequest("POST", "/nh", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("UpdateNhByIDHandler - Success", func(t *testing.T) {
		id := "123"

		existingNh := &entities.NursingHouse{
			Name:    "Old Name",
			Address: "Old Address",
			Images:  []entities.Image{{ImageLink: "old_image.jpg"}},
		}
		mockUseCase.On("GetNhByID", id).Return(existingNh, nil)

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
		).Return(updatedNh, nil)

		req := httptest.NewRequest("PUT", "/nh/"+id, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockUseCase.AssertExpectations(t)
	})
}
