package usecases_test

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

func TestGetAllNh(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouses := []entities.NursingHouse{
		{ID: "NH001", Name: "Test Home 1", Price: 1000},
		{ID: "NH002", Name: "Test Home 2", Price: 2000},
	}

	mockRepo.On("GetAllNh").Return(mockNursingHouses, nil)

	nursingHouses, err := useCase.GetAllNh()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(nursingHouses))
	assert.Equal(t, "NH001", nursingHouses[0].ID)
	assert.Equal(t, "NH002", nursingHouses[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestGetActiveNh(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouses := []entities.NursingHouse{
		{ID: "NH001", Name: "Test Home 1", Status: "Active"},
		{ID: "NH002", Name: "Test Home 2", Status: "Active"},
	}

	mockRepo.On("GetActiveNh").Return(mockNursingHouses, nil)

	nursingHouses, err := useCase.GetActiveNh()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(nursingHouses))
	assert.Equal(t, true, nursingHouses[0].Status == "Active")
	assert.Equal(t, true, nursingHouses[1].Status == "Active")

	mockRepo.AssertExpectations(t)
}

func TestGetInactiveNh(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouses := []entities.NursingHouse{
		{ID: "NH001", Name: "Test Home 1", Status: "Inactive"},
		{ID: "NH002", Name: "Test Home 2", Status: "Inactive"},
	}

	mockRepo.On("GetInactiveNh").Return(mockNursingHouses, nil)

	nursingHouses, err := useCase.GetInactiveNh()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(nursingHouses))
	assert.Equal(t, true, nursingHouses[0].Status == "Inactive")
	assert.Equal(t, true, nursingHouses[1].Status == "Inactive")

	mockRepo.AssertExpectations(t)
}

func TestGetNhByID(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouse := &entities.NursingHouse{
		ID:    "NH001",
		Name:  "Test Home",
		Price: 1000,
	}

	mockRepo.On("GetNhByID", "NH001").Return(mockNursingHouse, nil)

	nursingHouse, err := useCase.GetNhByID("NH001")

	assert.NoError(t, err)
	assert.Equal(t, "NH001", nursingHouse.ID)
	assert.Equal(t, "Test Home", nursingHouse.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetNhByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockRepo.On("GetNhByID", "NH999").Return(nil, errors.New("nursing house not found"))

	nursingHouse, err := useCase.GetNhByID("NH999")

	assert.Error(t, err)
	assert.Nil(t, nursingHouse)
	assert.Equal(t, "nursing house not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetNhNextID(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockRepo.On("GetNhNextID").Return("NH005", nil)

	id, err := useCase.GetNhNextID()

	assert.NoError(t, err)
	assert.Equal(t, "NH005", id)

	mockRepo.AssertExpectations(t)
}

func TestCreateNh_InvalidPrice(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nursingHouse := entities.NursingHouse{
		Name:  "Test Home",
		Price: -1000,
	}

	mockRepo.On("GetNhNextID").Return("NH005", nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	var files []multipart.FileHeader
	files = append(files, multipart.FileHeader{
		Filename: "test.jpg",
		Size:     1024,
	})

	result, err := useCase.CreateNh(nursingHouse, files, ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "price must be greater than zero", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestCreateNh_NoImages(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nursingHouse := entities.NursingHouse{
		Name:  "Test Home",
		Price: 1000,
	}

	mockRepo.On("GetNhNextID").Return("NH005", nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	var files []multipart.FileHeader

	result, err := useCase.CreateNh(nursingHouse, files, ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "at least one image is required", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetNhByIDForUser_NewHistory(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouse := &entities.NursingHouse{
		ID:    "NH001",
		Name:  "Test Home",
		Price: 1000,
	}

	userID := "user123"
	nhID := "NH001"

	mockRepo.On("GetNhHistory", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateNhHistory", mock.AnythingOfType("*entities.NursingHouseHistory")).Return(nil)
	mockRepo.On("GetNhByID", nhID).Return(mockNursingHouse, nil)

	result, err := useCase.GetNhByIDForUser(nhID, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockNursingHouse, result)

	mockRepo.AssertExpectations(t)
}

func TestGetNhByIDForUser_ExistingHistory(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouse := &entities.NursingHouse{
		ID:    "NH001",
		Name:  "Test Home",
		Price: 1000,
	}

	userID := "user123"
	nhID := "NH001"
	differentNhID := "NH002"

	existingHistory := &entities.NursingHouseHistory{
		UserID:         userID,
		NursingHouseID: differentNhID,
		NursingHouse: entities.NursingHouse{
			ID:   differentNhID,
			Name: "Different Home",
		},
	}

	mockRepo.On("GetNhHistory", userID).Return(existingHistory, nil)
	mockRepo.On("UpdateNhHistory", mock.AnythingOfType("*entities.NursingHouseHistory")).Return(nil)
	mockRepo.On("GetNhByID", nhID).Return(mockNursingHouse, nil)

	result, err := useCase.GetNhByIDForUser(nhID, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockNursingHouse, result)

	assert.Equal(t, nhID, existingHistory.NursingHouseID)

	mockRepo.AssertExpectations(t)
}

func TestGetNhByIDForUser_SameHistory(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	mockNursingHouse := &entities.NursingHouse{
		ID:    "NH001",
		Name:  "Test Home",
		Price: 1000,
	}

	userID := "user123"
	nhID := "NH001"

	existingHistory := &entities.NursingHouseHistory{
		UserID:         userID,
		NursingHouseID: nhID,
		NursingHouse: entities.NursingHouse{
			ID:   nhID,
			Name: "Test Home",
		},
	}

	mockRepo.On("GetNhHistory", userID).Return(existingHistory, nil)
	mockRepo.On("GetNhByID", nhID).Return(mockNursingHouse, nil)

	result, err := useCase.GetNhByIDForUser(nhID, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockNursingHouse, result)

	mockRepo.AssertExpectations(t)
}

func TestRecommendationCosine_NoHistory(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	userID := "user123"

	mockNursingHouses := []entities.NursingHouse{
		{ID: "NH001", Name: "Test Home 1", Price: 1000},
		{ID: "NH002", Name: "Test Home 2", Price: 2000},
		{ID: "NH003", Name: "Test Home 3", Price: 3000},
		{ID: "NH004", Name: "Test Home 4", Price: 4000},
		{ID: "NH005", Name: "Test Home 5", Price: 5000},
		{ID: "NH006", Name: "Test Home 6", Price: 6000},
	}

	mockRepo.On("GetNhHistory", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("GetAllNh").Return(mockNursingHouses, nil)

	results, err := useCase.RecommendationCosine(userID)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(results))

	mockRepo.AssertExpectations(t)
}

func TestRecommendationLLM_NoHistory(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	userID := "user123"

	mockNursingHouses := []entities.NursingHouse{
		{ID: "NH001", Name: "Test Home 1", Price: 1000},
		{ID: "NH002", Name: "Test Home 2", Price: 2000},
		{ID: "NH003", Name: "Test Home 3", Price: 3000},
		{ID: "NH004", Name: "Test Home 4", Price: 4000},
		{ID: "NH005", Name: "Test Home 5", Price: 5000},
		{ID: "NH006", Name: "Test Home 6", Price: 6000},
	}

	mockRepo.On("GetNhHistory", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("GetAllNh").Return(mockNursingHouses, nil)

	results, err := useCase.RecommendationLLM(userID)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(results))

	mockRepo.AssertExpectations(t)
}

func TestUpdateNhByID(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nhID := "NH001"

	existingNursingHouse := &entities.NursingHouse{
		ID:           nhID,
		Name:         "Original Name",
		Province:     "Original Province",
		Address:      "Original Address",
		Price:        1000,
		Google_map:   "Original Google Map Link",
		Phone_number: "Original Phone",
		Web_site:     "Original Website",
		Time:         "Original Time",
		Status:       "Active",
		Images: []entities.Image{
			{ID: "img1", ImageLink: "link1"},
			{ID: "img2", ImageLink: "link2"},
		},
	}

	updatedNursingHouse := entities.NursingHouse{
		Name:         "Updated Name",
		Province:     "Updated Province",
		Address:      "Updated Address",
		Price:        2000,
		Google_map:   "Updated Google Map Link",
		Phone_number: "Updated Phone",
		Web_site:     "Updated Website",
		Time:         "Updated Time",
		Status:       "Inactive",
	}

	expectedResult := &entities.NursingHouse{
		ID:           nhID,
		Name:         updatedNursingHouse.Name,
		Province:     updatedNursingHouse.Province,
		Address:      updatedNursingHouse.Address,
		Price:        updatedNursingHouse.Price,
		Google_map:   updatedNursingHouse.Google_map,
		Phone_number: updatedNursingHouse.Phone_number,
		Web_site:     updatedNursingHouse.Web_site,
		Time:         updatedNursingHouse.Time,
		Status:       updatedNursingHouse.Status,
	}

	imagesToDelete := []string{"img1"}

	mockRepo.On("GetNhByID", nhID).Return(existingNursingHouse, nil)
	mockRepo.On("RemoveImages", nhID, &imagesToDelete[0]).Return(nil)
	mockRepo.On("UpdateNhByID", mock.AnythingOfType("*entities.NursingHouse")).Return(expectedResult, nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	var files []multipart.FileHeader

	result, err := useCase.UpdateNhByID(nhID, updatedNursingHouse, files, imagesToDelete, ctx)

	assert.NoError(t, err)
	assert.Equal(t, nhID, result.ID)
	assert.Equal(t, updatedNursingHouse.Name, result.Name)
	assert.Equal(t, updatedNursingHouse.Price, result.Price)
	assert.Equal(t, updatedNursingHouse.Status, result.Status)

	mockRepo.AssertExpectations(t)
}

func TestUpdateNhByID_InvalidPrice(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nhID := "NH001"

	updatedNursingHouse := entities.NursingHouse{
		Name:  "Updated Name",
		Price: -1000,
	}

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	var files []multipart.FileHeader
	var imagesToDelete []string

	result, err := useCase.UpdateNhByID(nhID, updatedNursingHouse, files, imagesToDelete, ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "price must be greater than zero", err.Error())

	mockRepo.AssertNotCalled(t, "GetNhByID")
}

func TestUpdateNhByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nhID := "NH999"

	updatedNursingHouse := entities.NursingHouse{
		Name:  "Updated Name",
		Price: 1000,
	}

	mockRepo.On("GetNhByID", nhID).Return(nil, errors.New("nursing house not found"))

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	var files []multipart.FileHeader
	var imagesToDelete []string

	result, err := useCase.UpdateNhByID(nhID, updatedNursingHouse, files, imagesToDelete, ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "nursing house not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestCreateNhMock(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nursingHouse := entities.NursingHouse{
		Name:  "Test Home",
		Price: 1000,
	}

	links := []string{"link1", "link2", "link3"}
	mockRepo.On("GetNhNextID").Return("NH005", nil)

	mockResult := &entities.NursingHouse{
		ID:    "NH005",
		Name:  "Test Home",
		Price: 1000,
		Images: []entities.Image{
			{ID: "mock-id-1", ImageLink: "link1"},
			{ID: "mock-id-2", ImageLink: "link2"},
			{ID: "mock-id-3", ImageLink: "link3"},
		},
	}

	mockRepo.On("CreateNh", mock.AnythingOfType("*entities.NursingHouse"), mock.AnythingOfType("[]entities.Image")).Return(mockResult, nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	result, err := useCase.CreateNhMock(nursingHouse, links, ctx)

	assert.NoError(t, err)
	assert.Equal(t, "NH005", result.ID)
	assert.Equal(t, "Test Home", result.Name)
	assert.Equal(t, 1000, result.Price)
	assert.Equal(t, 3, len(result.Images))

	mockRepo.AssertExpectations(t)
}

func TestCreateNhMock_InvalidPrice(t *testing.T) {
	mockRepo := new(mocks.MockNhRepository)
	useCase := usecases.NewNhUseCase(mockRepo, configs.Supabase{}, configs.Recommend{})

	nursingHouse := entities.NursingHouse{
		Name:  "Test Home",
		Price: -1000,
	}

	links := []string{"link1", "link2", "link3"}

	mockRepo.On("GetNhNextID").Return("NH005", nil)

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	result, err := useCase.CreateNhMock(nursingHouse, links, ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "price must be greater than zero", err.Error())

	mockRepo.AssertNotCalled(t, "CreateNh")
	mockRepo.AssertExpectations(t)
}
