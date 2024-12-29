package usecases

import (
	"mime/multipart"
	"os"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

type NhUseCase interface {
	CreateNh(nursingHouse entities.NursingHouse, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetActiveNh() ([]entities.NursingHouse, error)
	GetInactiveNh() ([]entities.NursingHouse, error)
	GetNhByID(id string) (*entities.NursingHouse, error)
	GetNhNextID() (string, error)
	UpdateNhByID(id string, nursingHouse entities.NursingHouse, files []multipart.FileHeader, imagesToDelete []string, ctx *fiber.Ctx) (*entities.NursingHouse, error)
}

type NhUseCaseImpl struct {
	nhrepo 		repositories.NhRepository
	config		configs.Supabase
}

func NewNhUseCase(nhrepo repositories.NhRepository, config configs.Supabase) *NhUseCaseImpl {
	return &NhUseCaseImpl{
		nhrepo:  nhrepo,
		config:  config,
	}
}

func (u *NhUseCaseImpl) CreateNh(nursingHouse entities.NursingHouse, files []multipart.FileHeader,ctx *fiber.Ctx) (*entities.NursingHouse, error) {
	id, err := u.nhrepo.GetNhNextID()
	if err != nil {
		return nil, err
	}

	if nursingHouse.Price <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "price must be greater than zero",
        })
    }

	if len(files) == 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "at least one image is required",
        })
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
            ID:  		uuid.New().String(),
            ImageLink:  imageUrl,
        })
    }

	var createdNh *entities.NursingHouse
	createdNh, err = u.nhrepo.CreateNh(&nursingHouse, images)
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
	if nursingHouse.Price <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "price must be greater than zero",
        })
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