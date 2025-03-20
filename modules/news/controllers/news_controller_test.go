package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/usecases/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMultipartFormWithFile(fieldName, fileName string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("title", "Test News Title")
	writer.WriteField("type", "text")
	writer.WriteField("desc", "Test description")
	writer.WriteField("bold", "false")

	if fileName != "" {
		tmpFile, err := os.CreateTemp("", "test-*.jpg")
		if err != nil {
			return nil, "", err
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte("test image content")); err != nil {
			return nil, "", err
		}

		part, err := writer.CreateFormFile(fieldName, filepath.Base(tmpFile.Name()))
		if err != nil {
			return nil, "", err
		}

		if _, err := tmpFile.Seek(0, 0); err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, tmpFile); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func TestNewsHandlers(t *testing.T) {
	mockUseCase := new(mocks.MockNewsUseCase)
	controller := controllers.NewNewsController(mockUseCase)

	t.Run("CreateNewsHandler - Success", func(t *testing.T) {
		expectedNews := &entities.News{
			ID:    "1",
			Title: "Test News Title",
			Dialog: []entities.Dialog{
				{Type: "text", Desc: "Test description", Bold: false},
			},
		}

		mockUseCase.On("CreateNews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedNews, nil).Once()

		body, contentType, err := createMultipartFormWithFile("image_title", "test.jpg")
		assert.NoError(t, err)

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", contentType)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News created successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateNewsHandler - Missing Image Title", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("title", "Test News Title")
		writer.WriteField("type", "text")
		writer.WriteField("desc", "Test description")
		writer.WriteField("bold", "false")

		writer.Close()

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "Invalid image_title file")
	})

	t.Run("CreateNewsHandler - Missing Title", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("image_title", "test-image.jpg")
		if err != nil {
			t.Fatal(err)
		}

		dummyImageData := []byte("fake image content")
		fileWriter.Write(dummyImageData)
		writer.WriteField("type", "text")
		writer.WriteField("desc", "Test description")
		writer.WriteField("bold", "false")
		writer.Close()

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "title cannot be empty")
	})

	t.Run("CreateNewsHandler - Mismatch Dialog Fields", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("image_title", "test-image.jpg")
		if err != nil {
			t.Fatal(err)
		}

		dummyImageData := []byte("fake image content")
		fileWriter.Write(dummyImageData)
		writer.WriteField("title", "Test News Title")
		writer.WriteField("type", "text")
		writer.WriteField("desc", "Test description")
		writer.Close()

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "Mismatch in count of 'type', 'desc', and 'bold'")
	})

	t.Run("CreateNewsHandler - Invalid Bold Value", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("image_title", "test-image.jpg")
		if err != nil {
			t.Fatal(err)
		}

		dummyImageData := []byte("fake image content")
		fileWriter.Write(dummyImageData)
		writer.WriteField("title", "Test News Title")
		writer.WriteField("type", "text")
		writer.WriteField("desc", "Test description")
		writer.WriteField("bold", "not-a-boolean")
		writer.Close()

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "Invalid value for 'bold', must be true or false")
	})

	t.Run("CreateNewsHandler - UseCase Error", func(t *testing.T) {
		mockUseCase.On("CreateNews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*entities.News)(nil), errors.New("database error")).Once()

		body, contentType, err := createMultipartFormWithFile("image_title", "test.jpg")
		assert.NoError(t, err)

		app := fiber.New()
		app.Post("/api/news", controller.CreateNewsHandler)

		req := httptest.NewRequest("POST", "/api/news", body)
		req.Header.Set("Content-Type", contentType)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response["status"])
		assert.Contains(t, response["message"], "database error")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAllNewsHandler - Success", func(t *testing.T) {
		expectedNews := []entities.News{
			{
				ID:    "1",
				Title: "Test News 1",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "Test description 1", Bold: false},
				},
			},
			{
				ID:    "2",
				Title: "Test News 2",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "Test description 2", Bold: true},
				},
			},
		}

		mockUseCase.On("GetAllNews").Return(expectedNews, nil).Once()

		app := fiber.New()
		app.Get("/api/news", controller.GetAllNewsHandler)

		req := httptest.NewRequest("GET", "/api/news", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetAllNewsHandler - UseCase Error", func(t *testing.T) {
		mockUseCase.On("GetAllNews").Return(([]entities.News)(nil), errors.New("database error")).Once()
		app := fiber.New()
		app.Get("/api/news", controller.GetAllNewsHandler)

		req := httptest.NewRequest("GET", "/api/news", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response["status"])
		assert.Contains(t, response["message"], "database error")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNewsByIDHandler - Success", func(t *testing.T) {
		newsID := "1"
		expectedNews := &entities.News{
			ID:    newsID,
			Title: "Test News",
			Dialog: []entities.Dialog{
				{Type: "text", Desc: "Test description", Bold: false},
			},
		}

		mockUseCase.On("GetNewsByID", newsID).Return(expectedNews, nil).Once()

		app := fiber.New()
		app.Get("/api/news/:id", controller.GetNewsByIDHandler)

		req := httptest.NewRequest("GET", "/api/news/"+newsID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNewsByIDHandler - Not Found", func(t *testing.T) {
		newsID := "999"

		mockUseCase.On("GetNewsByID", newsID).Return((*entities.News)(nil), errors.New("news not found")).Once()

		app := fiber.New()
		app.Get("/api/news/:id", controller.GetNewsByIDHandler)

		req := httptest.NewRequest("GET", "/api/news/"+newsID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Not Found", response["status"])
		assert.Contains(t, response["message"], "news not found")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNewsNextIDHandler - Success", func(t *testing.T) {
		expectedNextID := "00005"
		mockUseCase.On("GetNewsNextID").Return(expectedNextID, nil).Once()

		app := fiber.New()
		app.Get("/api/news/next-id", controller.GetNewsNextIDHandler)

		req := httptest.NewRequest("GET", "/api/news/next-id", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("GetNewsNextIDHandler - UseCase Error", func(t *testing.T) {
		mockUseCase.On("GetNewsNextID").Return("", errors.New("database error")).Once()

		app := fiber.New()
		app.Get("/api/news/next-id", controller.GetNewsNextIDHandler)

		req := httptest.NewRequest("GET", "/api/news/next-id", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Not Found", response["status"])
		assert.Contains(t, response["message"], "database error")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNewsByIDHandler - Success", func(t *testing.T) {
		newsID := "1"
		expectedNews := &entities.News{
			ID:    newsID,
			Title: "Updated News Title",
			Dialog: []entities.Dialog{
				{Type: "text", Desc: "Updated description", Bold: false},
			},
		}

		mockUseCase.On("UpdateNewsByID", newsID, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedNews, nil).Once()

		body, contentType, err := createMultipartFormWithFile("image_title", "test.jpg")
		assert.NoError(t, err)

		app := fiber.New()
		app.Put("/api/news/:id", controller.UpdateNewsByIDHandler)

		req := httptest.NewRequest("PUT", "/api/news/"+newsID, body)
		req.Header.Set("Content-Type", contentType)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News retrieved successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("UpdateNewsByIDHandler - Missing Title", func(t *testing.T) {
		newsID := "1"

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("type", "text")
		writer.WriteField("desc", "Updated description")
		writer.WriteField("bold", "false")
		writer.Close()

		app := fiber.New()
		app.Put("/api/news/:id", controller.UpdateNewsByIDHandler)

		req := httptest.NewRequest("PUT", "/api/news/"+newsID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response["status"])
		assert.Contains(t, response["message"], "title cannot be empty")
	})

	t.Run("UpdateNewsByIDHandler - UseCase Error", func(t *testing.T) {
		newsID := "1"

		mockUseCase.On("UpdateNewsByID", newsID, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*entities.News)(nil), errors.New("news not found")).Once()

		body, contentType, err := createMultipartFormWithFile("image_title", "test.jpg")
		assert.NoError(t, err)

		app := fiber.New()
		app.Put("/api/news/:id", controller.UpdateNewsByIDHandler)

		req := httptest.NewRequest("PUT", "/api/news/"+newsID, body)
		req.Header.Set("Content-Type", contentType)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Not Found", response["status"])
		assert.Contains(t, response["message"], "news not found")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteNewsByIDHandler - Success", func(t *testing.T) {
		newsID := "1"

		mockUseCase.On("DeleteNewsByID", newsID).Return(nil).Once()

		app := fiber.New()
		app.Delete("/api/news/:id", controller.DeleteNewsByIDHandler)

		req := httptest.NewRequest("DELETE", "/api/news/"+newsID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Success", response["status"])
		assert.Equal(t, "News deleted successfully", response["message"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("DeleteNewsByIDHandler - UseCase Error", func(t *testing.T) {
		newsID := "999"

		mockUseCase.On("DeleteNewsByID", newsID).Return(errors.New("news not found")).Once()

		app := fiber.New()
		app.Delete("/api/news/:id", controller.DeleteNewsByIDHandler)

		req := httptest.NewRequest("DELETE", "/api/news/"+newsID, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response["status"])
		assert.Contains(t, response["message"], "news not found")

		mockUseCase.AssertExpectations(t)
	})
}
