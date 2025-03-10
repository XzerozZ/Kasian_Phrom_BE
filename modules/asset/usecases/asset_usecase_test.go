package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) GetAssetNextID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) CreateAsset(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAssetByID(id string) (*entities.Asset, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) UpdateAssetByID(asset *entities.Asset) (*entities.Asset, error) {
	args := m.Called(asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

func (m *MockAssetRepository) DeleteAssetByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAssetRepository) FindAssetByNameandUserID(name, userID string) (*entities.Asset, error) {
	args := m.Called(name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Asset), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateOTP(otp *entities.OTP) error {
	panic("unimplemented")
}

func (m *MockUserRepository) CreateSelectedHouse(selectedHouse *entities.SelectedHouse) error {
	panic("unimplemented")
}

func (m *MockUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) DeleteOTP(userID string) error {
	panic("unimplemented")
}

func (m *MockUserRepository) FindUserByEmail(email string) (entities.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetHistoryByUserID(userID string) ([]entities.History, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetHistoryInRange(userID string, startDate time.Time, endDate time.Time) ([]entities.History, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetRoleByName(name string) (entities.Role, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetUserByID(id string) (*entities.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetUserDepositsInRange(userID string, startDate time.Time, endDate time.Time) ([]entities.History, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetUserHistoryByMonth(userID string) (map[string]float64, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) UpdateSelectedHouse(selectedHouse *entities.SelectedHouse) (*entities.SelectedHouse, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) UpdateUserByID(user *entities.User) (*entities.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) CreateHistory(history *entities.History) (*entities.History, error) {
	args := m.Called(history)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.History), args.Error(1)
}

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) AddImages(id string, images []entities.Image) (*entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) CreateNh(nursingHouse *entities.NursingHouse, images []entities.Image) (*entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) GetActiveNh() ([]entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) GetAllNh() ([]entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) GetInactiveNh() ([]entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) GetNhByID(id string) (*entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) GetNhNextID() (string, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) RemoveImages(id string, imageID *string) error {
	panic("unimplemented")
}

func (m *MockNotificationRepository) UpdateNhByID(nursingHouse *entities.NursingHouse) (*entities.NursingHouse, error) {
	panic("unimplemented")
}

func (m *MockNotificationRepository) CreateNotification(notification *entities.Notification) (*entities.Notification, error) {
	args := m.Called(notification)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Notification), args.Error(1)
}

func TestCreateAsset(t *testing.T) {
	testCases := []struct {
		name            string
		prepareMockRepo func(*MockAssetRepository, *MockUserRepository, *MockNotificationRepository)
		asset           *entities.Asset
		expectedError   bool
	}{
		{
			name: "Successful Asset Creation",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				m.On("GetAssetNextID").Return("ASSET001", nil)
				m.On("CreateAsset", mock.MatchedBy(func(asset *entities.Asset) bool {
					return asset != nil
				})).Return(&entities.Asset{ID: "ASSET001"}, nil)
			},
			asset:         &entities.Asset{Name: "House", TotalCost: 10000, EndYear: "2026"},
			expectedError: false,
		},
		{
			name: "Failed to Get Next ID",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				m.On("GetAssetNextID").Return("", errors.New("id generation failed"))
			},
			asset:         &entities.Asset{Name: "House", TotalCost: 10000},
			expectedError: true,
		},
		{
			name: "Failed to Create Asset",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				m.On("GetAssetNextID").Return("ASSET001", nil)
				m.On("CreateAsset", mock.MatchedBy(func(asset *entities.Asset) bool {
					return asset != nil
				})).Return(nil, errors.New("database error"))
			},
			asset:         &entities.Asset{Name: "House", TotalCost: 10000},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockAssetRepository)
			mockUserRepo := new(MockUserRepository)
			mockNotifRepo := new(MockNotificationRepository)
			tc.prepareMockRepo(mockRepo, mockUserRepo, mockNotifRepo)
			useCase := usecases.NewAssetUseCase(mockRepo, mockUserRepo, mockNotifRepo, nil, nil)
			result, err := useCase.CreateAsset(*tc.asset)
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

func TestUpdateAssetByID(t *testing.T) {
	testCases := []struct {
		name            string
		assetID         string
		prepareMockRepo func(*MockAssetRepository, *MockUserRepository, *MockNotificationRepository)
		updateAsset     entities.Asset
		expectedError   bool
	}{
		{
			name:    "Successful Update",
			assetID: "ASSET001",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				existingAsset := &entities.Asset{ID: "ASSET001", Name: "Old House", TotalCost: 8000, EndYear: "2025"}
				updatedAsset := &entities.Asset{ID: "ASSET001", Name: "New House", TotalCost: 10000, EndYear: "2025"}
				m.On("GetAssetByID", "ASSET001").Return(existingAsset, nil)
				m.On("UpdateAssetByID", mock.Anything).Return(updatedAsset, nil)
				n.On("CreateNotification", mock.Anything).Return(&entities.Notification{}, nil)
			},

			updateAsset:   entities.Asset{ID: "ASSET001", Name: "New House", TotalCost: 10000, EndYear: "2025"},
			expectedError: false,
		},
		{
			name:    "Asset Not Found",
			assetID: "NONEXISTENT",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				m.On("GetAssetByID", "NONEXISTENT").Return(nil, errors.New("asset not found"))
			},

			updateAsset:   entities.Asset{},
			expectedError: true,
		},
		{
			name:    "Update Failed",
			assetID: "ASSET001",
			prepareMockRepo: func(m *MockAssetRepository, u *MockUserRepository, n *MockNotificationRepository) {
				existingAsset := &entities.Asset{ID: "ASSET001", Name: "Old House", TotalCost: 8000, EndYear: "2025"}
				m.On("GetAssetByID", "ASSET001").Return(existingAsset, nil)
				m.On("UpdateAssetByID", mock.Anything).Return(nil, errors.New("update failed"))
				n.On("CreateNotification", mock.Anything).Return(&entities.Notification{}, nil)
			},

			updateAsset:   entities.Asset{Name: "New House", TotalCost: 10000, EndYear: "2025"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockAssetRepository)
			mockUserRepo := new(MockUserRepository)
			mockNotifRepo := new(MockNotificationRepository)
			tc.prepareMockRepo(mockRepo, mockUserRepo, mockNotifRepo)
			useCase := usecases.NewAssetUseCase(mockRepo, mockUserRepo, mockNotifRepo, nil, nil)
			result, err := useCase.UpdateAssetByID(tc.assetID, tc.updateAsset)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.assetID, result.ID)
				assert.Equal(t, tc.updateAsset.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
			mockNotifRepo.AssertExpectations(t)
		})
	}
}
