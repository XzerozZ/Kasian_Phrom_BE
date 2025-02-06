package usecases

import (
	"errors"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock NewsRepository
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

// Test CreateNews
func TestCreateNews(t *testing.T) {
	testCases := []struct {
		name            string
		prepareMockRepo func(*MockNewsRepository)
		news            *entities.News
		expectedError   bool
	}{
		{
			name: "Successful News Creation",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsNextID").Return("NEWS001", nil)
				m.On("CreateNews", mock.Anything).Return(&entities.News{ID: "NEWS001"}, nil)
			},
			news: &entities.News{
				Title: "Test News",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "Test Dialog"},
				},
			},
			expectedError: false,
		},
		{
			name: "Failed to Get Next ID",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsNextID").Return("", errors.New("id generation failed"))
			},
			news:          &entities.News{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockNewsRepository)
			mockConfig := configs.Supabase{}
			tc.prepareMockRepo(mockRepo)
			useCase := usecases.NewNewsUseCase(mockRepo, mockConfig)
			result, err := useCase.CreateNews(tc.news, nil, nil, &fiber.Ctx{})

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Test GetAllNews
func TestGetAllNews(t *testing.T) {
	testCases := []struct {
		name            string
		prepareMockRepo func(*MockNewsRepository)
		expectedError   bool
	}{
		{
			name: "Successful Retrieval",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetAllNews").Return([]entities.News{
					{ID: "NEWS001", Title: "News 1"},
					{ID: "NEWS002", Title: "News 2"},
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "Retrieval Failure",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetAllNews").Return([]entities.News{}, errors.New("retrieval failed"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockNewsRepository)
			mockConfig := configs.Supabase{}
			tc.prepareMockRepo(mockRepo)
			useCase := usecases.NewNewsUseCase(mockRepo, mockConfig)
			news, err := useCase.GetAllNews()
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, news)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Test GetNewsByID
func TestGetNewsByID(t *testing.T) {
	testCases := []struct {
		name            string
		newsID          string
		prepareMockRepo func(*MockNewsRepository)
		expectedError   bool
	}{
		{
			name:   "Successful Retrieval",
			newsID: "NEWS001",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsByID", "NEWS001").Return(&entities.News{
					ID:    "NEWS001",
					Title: "Test News",
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Retrieval Failure",
			newsID: "NONEXISTENT",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsByID", "NONEXISTENT").Return(nil, errors.New("news not found"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockNewsRepository)
			mockConfig := configs.Supabase{}
			tc.prepareMockRepo(mockRepo)
			useCase := usecases.NewNewsUseCase(mockRepo, mockConfig)
			news, err := useCase.GetNewsByID(tc.newsID)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, news)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, news)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Test UpdateNewsByID
func TestUpdateNewsByID(t *testing.T) {
	testCases := []struct {
		name            string
		newsID          string
		prepareMockRepo func(*MockNewsRepository)
		updateNews      entities.News
		expectedError   bool
	}{
		{
			name:   "Successful Update",
			newsID: "NEWS001",
			prepareMockRepo: func(m *MockNewsRepository) {
				existingNews := &entities.News{
					ID:    "NEWS001",
					Title: "Old Title",
					Dialog: []entities.Dialog{
						{ID: "DIALOG001"},
					},
				}
				m.On("GetNewsByID", "NEWS001").Return(existingNews, nil)
				m.On("DeleteDialog", "DIALOG001").Return(nil)
				m.On("UpdateNewsByID", mock.Anything).Return(existingNews, nil)
			},
			updateNews: entities.News{
				Title: "New Title",
				Dialog: []entities.Dialog{
					{Type: "text", Desc: "New Dialog"},
				},
			},
			expectedError: false,
		},
		{
			name:   "News Not Found",
			newsID: "NONEXISTENT",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsByID", "NONEXISTENT").Return(nil, errors.New("news not found"))
			},
			updateNews:    entities.News{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockNewsRepository)
			mockConfig := configs.Supabase{}
			tc.prepareMockRepo(mockRepo)
			useCase := usecases.NewNewsUseCase(mockRepo, mockConfig)
			result, err := useCase.UpdateNewsByID(
				tc.newsID,
				tc.updateNews,
				nil,
				nil,
				false,
				&fiber.Ctx{},
			)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Test DeleteNewsByID
func TestDeleteNewsByID(t *testing.T) {
	testCases := []struct {
		name            string
		newsID          string
		prepareMockRepo func(*MockNewsRepository)
		expectedError   bool
	}{
		{
			name:   "Successful Deletion",
			newsID: "NEWS001",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsByID", "NEWS001").Return(&entities.News{
					ID: "NEWS001",
					Dialog: []entities.Dialog{
						{ID: "DIALOG001"},
						{ID: "DIALOG002"},
					},
				}, nil)
				m.On("DeleteDialog", "DIALOG001").Return(nil)
				m.On("DeleteDialog", "DIALOG002").Return(nil)
				m.On("DeleteNewsByID", "NEWS001").Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "News Not Found",
			newsID: "NONEXISTENT",
			prepareMockRepo: func(m *MockNewsRepository) {
				m.On("GetNewsByID", "NONEXISTENT").Return(nil, errors.New("news not found"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockNewsRepository)
			mockConfig := configs.Supabase{}
			tc.prepareMockRepo(mockRepo)
			useCase := usecases.NewNewsUseCase(mockRepo, mockConfig)

			err := useCase.DeleteNewsByID(tc.newsID)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
