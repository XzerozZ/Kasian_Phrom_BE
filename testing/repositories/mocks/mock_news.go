package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockNewsRepository struct {
	mock.Mock
}

func (m *MockNewsRepository) GetNewsNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockNewsRepository) CreateNews(news *entities.News) (*entities.News, error) {
	args := m.Called(news)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsRepository) GetAllNews() ([]entities.News, error) {
	args := m.Called()
	return args.Get(0).([]entities.News), args.Error(1)
}

func (m *MockNewsRepository) GetNewsByID(id string) (*entities.News, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsRepository) UpdateNewsByID(news *entities.News) (*entities.News, error) {
	args := m.Called(news)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.News), args.Error(1)
}

func (m *MockNewsRepository) DeleteNewsByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNewsRepository) DeleteDialog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
