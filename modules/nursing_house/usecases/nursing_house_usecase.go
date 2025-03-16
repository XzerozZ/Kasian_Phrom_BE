package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"golang.org/x/exp/rand"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NhUseCase interface {
	CreateNh(nursingHouse entities.NursingHouse, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetActiveNh() ([]entities.NursingHouse, error)
	GetInactiveNh() ([]entities.NursingHouse, error)
	GetNhByID(id string) (*entities.NursingHouse, error)
	GetNhNextID() (string, error)
	UpdateNhByID(id string, nursingHouse entities.NursingHouse, files []multipart.FileHeader, imagesToDelete []string, ctx *fiber.Ctx) (*entities.NursingHouse, error)

	GetNhByIDForUser(id, userID string) (*entities.NursingHouse, error)
	RecommendationCosine(userID string) ([]entities.NursingHouse, error)
	RecommendationLLM(userID string) ([]entities.NursingHouse, error)

	CreateNhMock(nursingHouse entities.NursingHouse, links []string, ctx *fiber.Ctx) (*entities.NursingHouse, error)
}

type NhUseCaseImpl struct {
	nhrepo repositories.NhRepository
	config configs.Supabase
	recom  configs.Recommend
}

func NewNhUseCase(nhrepo repositories.NhRepository, config configs.Supabase, recom configs.Recommend) *NhUseCaseImpl {
	return &NhUseCaseImpl{
		nhrepo: nhrepo,
		config: config,
		recom:  recom,
	}
}

func (u *NhUseCaseImpl) CreateNh(nursingHouse entities.NursingHouse, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	id, err := u.nhrepo.GetNhNextID()
	if err != nil {
		return nil, err
	}

	if nursingHouse.Price < 0 {
		return nil, errors.New("price must be greater than zero")
	}

	if len(files) == 0 {
		return nil, errors.New("at least one image is required")
	}

	nursingHouse.ID = id
	var images []entities.Image
	for _, file := range files {
		fileName := uuid.New().String() + ".jpg"
		if err := ctx.SaveFile(&file, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.config)
		if err != nil {
			return nil, err
		}

		err = os.Remove("./uploads/" + fileName)
		if err != nil {
			return nil, err
		}

		images = append(images, entities.Image{
			ID:        uuid.New().String(),
			ImageLink: imageUrl,
		})
	}

	createdNh, err := u.nhrepo.CreateNh(&nursingHouse, images)
	if err != nil {
		return nil, err
	}

	return createdNh, nil
}

func (u *NhUseCaseImpl) GetAllNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetAllNh()
}

func (u *NhUseCaseImpl) GetActiveNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetActiveNh()
}

func (u *NhUseCaseImpl) GetInactiveNh() ([]entities.NursingHouse, error) {
	return u.nhrepo.GetInactiveNh()
}

func (u *NhUseCaseImpl) GetNhByID(id string) (*entities.NursingHouse, error) {
	return u.nhrepo.GetNhByID(id)
}

func (u *NhUseCaseImpl) GetNhNextID() (string, error) {
	return u.nhrepo.GetNhNextID()
}

func (u *NhUseCaseImpl) UpdateNhByID(id string, nursingHouse entities.NursingHouse, files []multipart.FileHeader, imagesToDelete []string, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	if nursingHouse.Price < 0 {
		return nil, errors.New("price must be greater than zero")
	}

	existingNh, err := u.nhrepo.GetNhByID(id)
	if err != nil {
		return nil, err
	}

	existingNh.Name = nursingHouse.Name
	existingNh.Province = nursingHouse.Province
	existingNh.Address = nursingHouse.Address
	existingNh.Price = nursingHouse.Price
	existingNh.Google_map = nursingHouse.Google_map
	existingNh.Phone_number = nursingHouse.Phone_number
	existingNh.Web_site = nursingHouse.Web_site
	existingNh.Time = nursingHouse.Time
	existingNh.Status = nursingHouse.Status

	if len(imagesToDelete) > 0 {
		for _, imageID := range imagesToDelete {
			if err := u.nhrepo.RemoveImages(id, &imageID); err != nil {
				return nil, err
			}
		}
	}

	var newImages []entities.Image
	if len(files) > 0 {
		for _, file := range files {
			fileName := uuid.New().String() + ".jpg"
			if err := ctx.SaveFile(&file, "./uploads/"+fileName); err != nil {
				return nil, err
			}

			imageUrl, err := utils.UploadImage(fileName, "", u.config)
			if err != nil {
				os.Remove("./uploads/" + fileName)
				return nil, err
			}

			if err := os.Remove("./uploads/" + fileName); err != nil {
				return nil, err
			}

			newImages = append(newImages, entities.Image{
				ID:        uuid.New().String(),
				ImageLink: imageUrl,
			})
		}
	}

	existingNh.Images = newImages
	if len(newImages) > 0 {
		existingNh.Images = append(existingNh.Images, newImages...)
		_, err = u.nhrepo.AddImages(id, newImages)
		if err != nil {
			return nil, err
		}
	}

	updatedNh, err := u.nhrepo.UpdateNhByID(existingNh)
	if err != nil {
		return nil, err
	}

	return updatedNh, nil
}

func (u *NhUseCaseImpl) GetNhByIDForUser(id, userID string) (*entities.NursingHouse, error) {
	nhHistory, err := u.nhrepo.GetNhHistory(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newNhHistory := &entities.NursingHouseHistory{
				UserID:         userID,
				NursingHouseID: id,
			}

			if err := u.nhrepo.CreateNhHistory(newNhHistory); err != nil {
				return nil, err
			}

			return u.nhrepo.GetNhByID(id)
		}

		return nil, err
	}

	if nhHistory.NursingHouseID != id {
		nhHistory.NursingHouseID = id
		if err := u.nhrepo.UpdateNhHistory(nhHistory); err != nil {
			return nil, err
		}
	}

	return u.nhrepo.GetNhByID(id)
}

func (u *NhUseCaseImpl) RecommendationCosine(userID string) ([]entities.NursingHouse, error) {
	nhHistory, err := u.nhrepo.GetNhHistory(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			recommend, err := u.nhrepo.GetAllNh()
			if err != nil {
				return nil, err
			}

			limit := 5
			if len(recommend) < limit {
				limit = len(recommend)
			}

			rand.Shuffle(len(recommend), func(i, j int) {
				recommend[i], recommend[j] = recommend[j], recommend[i]
			})

			return recommend[:limit], nil
		}

		return nil, err
	}

	nhNameEncoded := url.QueryEscape(nhHistory.NursingHouse.Name)
	url := fmt.Sprintf("%s/cosine?nh_name=%s", u.recom, nhNameEncoded)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get recommendation: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	resultList, ok := result["result"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	var nursingHomes []entities.NursingHouse
	for _, name := range resultList {
		nhName, ok := name.(string)
		if !ok {
			continue
		}

		nursingHome, err := u.nhrepo.GetNhByName(nhName)
		if err != nil {
			continue
		}

		nursingHomes = append(nursingHomes, nursingHome)
	}

	return nursingHomes, nil
}

func (u *NhUseCaseImpl) RecommendationLLM(userID string) ([]entities.NursingHouse, error) {
	nhHistory, err := u.nhrepo.GetNhHistory(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			recommend, err := u.nhrepo.GetAllNh()
			if err != nil {
				return nil, err
			}

			limit := 5
			if len(recommend) < limit {
				limit = len(recommend)
			}

			rand.Shuffle(len(recommend), func(i, j int) {
				recommend[i], recommend[j] = recommend[j], recommend[i]
			})

			return recommend[:limit], nil
		}

		return nil, err
	}

	nhNameEncoded := url.QueryEscape(nhHistory.NursingHouse.Name)
	url := fmt.Sprintf("%s/llm?nh_name=%s", u.recom, nhNameEncoded)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get recommendation: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	resultList, ok := result["result"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	var nursingHomes []entities.NursingHouse
	for _, name := range resultList {
		nhName, ok := name.(string)
		if !ok {
			continue
		}

		nursingHome, err := u.nhrepo.GetNhByName(nhName)
		if err != nil {
			continue
		}

		nursingHomes = append(nursingHomes, nursingHome)
	}

	return nursingHomes, nil
}

func (u *NhUseCaseImpl) CreateNhMock(nursingHouse entities.NursingHouse, links []string, ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	id, err := u.nhrepo.GetNhNextID()
	if err != nil {
		return nil, err
	}

	if nursingHouse.Price < 0 {
		return nil, errors.New("price must be greater than zero")
	}

	nursingHouse.ID = id
	var images []entities.Image
	for _, links := range links {
		images = append(images, entities.Image{
			ID:        uuid.New().String(),
			ImageLink: links,
		})
	}

	createdNh, err := u.nhrepo.CreateNh(&nursingHouse, images)
	if err != nil {
		return nil, err
	}

	return createdNh, nil
}
