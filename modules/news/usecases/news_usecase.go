package usecases

import (
	"mime/multipart"
	"os"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/repositories"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

type NewsUseCase interface {
	CreateNews(news *entities.News, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error)
	GetAllNews() ([]entities.News, error)
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
}

type NewsUseCaseImpl struct {
	newsrepo 	repositories.NewsRepository
	config		configs.Supabase
}

func NewNewsUseCase(newsrepo repositories.NewsRepository, config configs.Supabase) *NewsUseCaseImpl {
	return &NewsUseCaseImpl{
		newsrepo:  newsrepo,
		config:    config,
	}
}

func (u *NewsUseCaseImpl) CreateNews(news *entities.News, files []multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error) {
	id, err := u.newsrepo.GetNewsNextID()
	if err != nil {
		return nil, err
	}

	for i, dialogReq := range news.Dialog {
		news.Dialog[i] = entities.Dialog{
			ID:		uuid.New().String(),
			Type:   dialogReq.Type,
			Desc:   dialogReq.Desc,
			NewsID: id,
		}
	}

	if len(files) == 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "at least one image is required",
        })
	}
	
	news.ID = id
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

	var createdNews *entities.News
	createdNews, err = u.newsrepo.CreateNews(news, images)
    if err != nil {
        return nil, err
    }
    
    return createdNews, nil
}

func (u *NewsUseCaseImpl) GetAllNews() ([]entities.News, error) {
	return u.newsrepo.GetAllNews()
}

func (u *NewsUseCaseImpl) GetNewsByID(id string) (*entities.News, error) {
	return u.newsrepo.GetNewsByID(id)
}

func (u *NewsUseCaseImpl) GetNewsNextID() (string, error) {
	return u.newsrepo.GetNewsNextID()
}