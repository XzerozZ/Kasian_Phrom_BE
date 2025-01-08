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
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/controllers"
)

type MockNewsUseCase struct {
	mock.Mock
}

func (m *MockNewsUseCase) CreateNews(news *entities.News, imageTitle, imageDesc *multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error) {
	args := m.Called(news, imageTitle, imageDesc, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.News), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNewsUseCase) GetNewsByID(id string) (*entities.News, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entities.News), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNewsUseCase) GetAllNews() ([]entities.News, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]entities.News), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNewsUseCase) GetNewsNextID() (string, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.(string), args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockNewsUseCase) UpdateNewsByID(id string, news entities.News, imageTitle, imageDesc *multipart.FileHeader, shouldDeleteImageDesc bool, ctx *fiber.Ctx) (*entities.News, error) {
	args := m.Called(id, news, imageTitle, imageDesc, shouldDeleteImageDesc, ctx)
	if result := args.Get(0); result != nil {
		return result.(*entities.News), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNewsUseCase) DeleteNewsByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNewsHandlers(t *testing.T) {
	mockUseCase := new(MockNewsUseCase)
	controller := controllers.NewNewsController(mockUseCase)

	app := fiber.New()
	app.Post("/news", controller.CreateNewsHandler)
	app.Get("/news", controller.GetAllNewsHandler)
	app.Get("/news/next-id", controller.GetNewsNextIDHandler)  // Moved before :id route
	app.Get("/news/:id", controller.GetNewsByIDHandler)
	app.Put("/news/:id", controller.UpdateNewsByIDHandler)
	app.Delete("/news/:id", controller.DeleteNewsByIDHandler)

	t.Run("CreateNewsHandler - Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("title", "Test News")
		_ = writer.WriteField("type", "text")
		_ = writer.WriteField("desc", "Test Description")
		_ = writer.WriteField("bold", "true")

		imageTitlePart, _ := writer.CreateFormFile("image_title", "title.jpg")
		_, _ = imageTitlePart.Write([]byte("dummy title image"))
		
		imageDescPart, _ := writer.CreateFormFile("image_desc", "desc.jpg")
		_, _ = imageDescPart.Write([]byte("dummy desc image"))
		
		writer.Close()

		expectedNews := &entities.News{
			Title: "Test News",
			Dialog: []entities.Dialog{
				{
					Type: "text",
					Desc: "Test Description",
					Bold: true,
				},
			},
		}

		mockUseCase.On("CreateNews", 
			mock.MatchedBy(func(n *entities.News) bool {
				return n.Title == "Test News" && len(n.Dialog) == 1
			}),
			mock.AnythingOfType("*multipart.FileHeader"),
			mock.AnythingOfType("*multipart.FileHeader"),
			mock.AnythingOfType("*fiber.Ctx"),
		).Return(expectedNews, nil)

		req := httptest.NewRequest("POST", "/news", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("GetAllNewsHandler - Success", func(t *testing.T) {
		expectedNews := []entities.News{
			{
				Title: "News 1",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "Description 1", Bold: true},
				},
			},
			{
				Title: "News 2",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "Description 2", Bold: false},
				},
			},
		}

		mockUseCase.On("GetAllNews").Return(expectedNews, nil)

		req := httptest.NewRequest("GET", "/news", nil)
		resp, err := app.Test(req, -1)
		
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("UpdateNewsByIDHandler - Success", func(t *testing.T) {
		id := "123"
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("title", "Updated News")
		_ = writer.WriteField("type", "text")
		_ = writer.WriteField("desc", "Updated Description")
		_ = writer.WriteField("bold", "true")
		_ = writer.WriteField("image_desc", "del_img")

		imageTitlePart, _ := writer.CreateFormFile("image_title", "new_title.jpg")
		_, _ = imageTitlePart.Write([]byte("new dummy title image"))

		writer.Close()

		expectedNews := &entities.News{
			Title: "Updated News",
			Dialog: []entities.Dialog{
				{
					Type: "text",
					Desc: "Updated Description",
					Bold: true,
				},
			},
		}

		mockUseCase.On("UpdateNewsByID",
			id,
			mock.MatchedBy(func(n entities.News) bool {
				return n.Title == "Updated News" && len(n.Dialog) == 1
			}),
			mock.AnythingOfType("*multipart.FileHeader"),
			mock.AnythingOfType("*multipart.FileHeader"),
			true,
			mock.AnythingOfType("*fiber.Ctx"),
		).Return(expectedNews, nil)

		req := httptest.NewRequest("PUT", "/news/"+id, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("DeleteNewsByIDHandler - Success", func(t *testing.T) {
		id := "123"
		mockUseCase.On("DeleteNewsByID", id).Return(nil)

		req := httptest.NewRequest("DELETE", "/news/"+id, nil)
		resp, err := app.Test(req, -1)
		
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("GetNewsByIDHandler - Success", func(t *testing.T) {
		id := "123"
		expectedNews := &entities.News{
			Title: "Test News",
			Dialog: []entities.Dialog{
				{Type: "text", Desc: "Description", Bold: true},
			},
		}

		mockUseCase.On("GetNewsByID", id).Return(expectedNews, nil)

		req := httptest.NewRequest("GET", "/news/"+id, nil)
		resp, err := app.Test(req, -1)
		
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("GetNewsNextIDHandler - Success", func(t *testing.T) {
		expectedID := "NEWS-001"
		mockUseCase.On("GetNewsNextID").Return(expectedID, nil)

		req := httptest.NewRequest("GET", "/news/next-id", nil)
		resp, err := app.Test(req, -1)
		
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}