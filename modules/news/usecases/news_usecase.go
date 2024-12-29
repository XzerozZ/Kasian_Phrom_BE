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
	UpdateNewsByID(id string, news entities.News, files []multipart.FileHeader, imagesToDelete []string, ctx *fiber.Ctx) (*entities.News, error)
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

func (u *NewsUseCaseImpl) UpdateNewsByID(id string, news entities.News, files []multipart.FileHeader, imagesToDelete []string, ctx *fiber.Ctx) (*entities.News, error) {
	existingNews, err := u.newsrepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}
	
	existingNews.Title = news.Title
	
	for _, dialog := range existingNews.Dialog {
		if err := u.newsrepo.DeleteDialog(dialog.ID); err != nil {
			return nil, err
		}
	}

	if len(imagesToDelete) > 0 {
        for _, imageID := range imagesToDelete {
            if err := u.newsrepo.RemoveImages(id, &imageID); err != nil {
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
	
	for i, dialogReq := range news.Dialog {
		news.Dialog[i] = entities.Dialog{
			ID:		uuid.New().String(),
			Type:   dialogReq.Type,
			Desc:   dialogReq.Desc,
			NewsID: existingNews.ID,
		}
	}

	existingNews.Dialog = news.Dialog
	existingNews.Images = newImages
	if len(newImages) > 0 {
		existingNews.Images = append(existingNews.Images, newImages...)
		_, err = u.newsrepo.AddImages(id, newImages)
		if err != nil {
			return nil, err
		}
	}

	updatedNews, err := u.newsrepo.UpdateNewsByID(existingNews)
    if err != nil {
        return nil, err
    }

	return updatedNews, nil
}