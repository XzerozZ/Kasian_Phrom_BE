package mocks

import (
	"mime/multipart"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
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
