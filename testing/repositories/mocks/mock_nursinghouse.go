package mocks

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockNhRepository struct {
	mock.Mock
}

func (m *MockNhRepository) CreateNh(nursingHouse *entities.NursingHouse, images []entities.Image) (*entities.NursingHouse, error) {
	args := m.Called(nursingHouse, images)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetAllNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	return args.Get(0).([]entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetActiveNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	return args.Get(0).([]entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetInactiveNh() ([]entities.NursingHouse, error) {
	args := m.Called()
	return args.Get(0).([]entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetNhByID(id string) (*entities.NursingHouse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetNhNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockNhRepository) UpdateNhByID(nursingHouse *entities.NursingHouse) (*entities.NursingHouse, error) {
	args := m.Called(nursingHouse)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) RemoveImages(nursingHouseID string, imageID *string) error {
	args := m.Called(nursingHouseID, imageID)
	return args.Error(0)
}

func (m *MockNhRepository) AddImages(nursingHouseID string, images []entities.Image) (*entities.NursingHouse, error) {
	args := m.Called(nursingHouseID, images)
	return args.Get(0).(*entities.NursingHouse), args.Error(1)
}

func (m *MockNhRepository) GetNhHistory(userID string) (*entities.NursingHouseHistory, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.NursingHouseHistory), args.Error(1)
}

func (m *MockNhRepository) CreateNhHistory(nhHistory *entities.NursingHouseHistory) error {
	args := m.Called(nhHistory)
	return args.Error(0)
}

func (m *MockNhRepository) UpdateNhHistory(nhHistory *entities.NursingHouseHistory) error {
	args := m.Called(nhHistory)
	return args.Error(0)
}

func (m *MockNhRepository) GetNhByName(name string) (entities.NursingHouse, error) {
	args := m.Called(name)
	return args.Get(0).(entities.NursingHouse), args.Error(1)
}
