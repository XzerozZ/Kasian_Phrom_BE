package usecases_test

import (
	"errors"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateFav(t *testing.T) {
	mockRepo := new(mocks.MockFavRepository)
	useCase := usecases.NewFavUseCase(mockRepo)

	t.Run("Create Favorite - Success", func(t *testing.T) {
		fav := &entities.Favorite{
			UserID:         "user123",
			NursingHouseID: "house123",
		}

		mockRepo.On("CreateFav", fav).Return(nil).Once()

		err := useCase.CreateFav(fav)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Favorite - Error", func(t *testing.T) {
		fav := &entities.Favorite{
			UserID:         "user123",
			NursingHouseID: "house123",
		}

		expectedErr := errors.New("database error")
		mockRepo.On("CreateFav", fav).Return(expectedErr).Once()

		err := useCase.CreateFav(fav)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetFavByUserID(t *testing.T) {
	mockRepo := new(mocks.MockFavRepository)
	useCase := usecases.NewFavUseCase(mockRepo)

	t.Run("Get Favorite - Success", func(t *testing.T) {
		userID := "user123"
		expectedFavs := []entities.Favorite{
			{UserID: "user123", NursingHouseID: "house1"},
			{UserID: "user123", NursingHouseID: "house2"},
		}

		mockRepo.On("GetFavByUserID", userID).Return(expectedFavs, nil).Once()

		favs, err := useCase.GetFavByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedFavs, favs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get Favorite - Error", func(t *testing.T) {
		userID := "user123"
		expectedErr := errors.New("database error")
		var emptyFavs []entities.Favorite

		mockRepo.On("GetFavByUserID", userID).Return(emptyFavs, expectedErr).Once()

		favs, err := useCase.GetFavByUserID(userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, favs)
		mockRepo.AssertExpectations(t)
	})
}

func TestCheckFav(t *testing.T) {
	mockRepo := new(mocks.MockFavRepository)
	useCase := usecases.NewFavUseCase(mockRepo)

	t.Run("Check Favorite - Success", func(t *testing.T) {
		userID := "user123"
		nursingHouseID := "house123"

		mockRepo.On("CheckFav", userID, nursingHouseID).Return(nil).Once()

		err := useCase.CheckFav(userID, nursingHouseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Check Favorite - Not Found", func(t *testing.T) {
		userID := "user123"
		nursingHouseID := "house123"
		expectedErr := errors.New("favorite not found")

		mockRepo.On("CheckFav", userID, nursingHouseID).Return(expectedErr).Once()

		err := useCase.CheckFav(userID, nursingHouseID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteFavByID(t *testing.T) {
	mockRepo := new(mocks.MockFavRepository)
	useCase := usecases.NewFavUseCase(mockRepo)

	t.Run("Delete Favorite - Success", func(t *testing.T) {
		userID := "user123"
		nursingHouseID := "house123"

		mockRepo.On("DeleteFavByID", userID, nursingHouseID).Return(nil).Once()

		err := useCase.DeleteFavByID(userID, nursingHouseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete Favorite - Error", func(t *testing.T) {
		userID := "user123"
		nursingHouseID := "house123"
		expectedErr := errors.New("delete error")

		mockRepo.On("DeleteFavByID", userID, nursingHouseID).Return(expectedErr).Once()

		err := useCase.DeleteFavByID(userID, nursingHouseID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
