package mocks

import (
	"mime/multipart"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type MockNewsUseCase struct {
	mock.Mock
}

func (m *MockNewsUseCase) CreateNews(news *entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error) {
	args := m.Called(news, imageTitleFile, imageDescFile, ctx)
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsUseCase) GetAllNews() ([]entities.News, error) {
	args := m.Called()
	return args.Get(0).([]entities.News), args.Error(1)
}

func (m *MockNewsUseCase) GetNewsByID(id string) (*entities.News, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsUseCase) GetNewsNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockNewsUseCase) UpdateNewsByID(id string, news entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, shouldDeleteImageDesc bool, ctx *fiber.Ctx) (*entities.News, error) {
	args := m.Called(id, news, imageTitleFile, imageDescFile, shouldDeleteImageDesc, ctx)
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsUseCase) DeleteNewsByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
