package usecases

import (
	"mime/multipart"
	"os"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NewsUseCase interface {
	CreateNews(news *entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error)
	GetAllNews() ([]entities.News, error)
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
	UpdateNewsByID(id string, news entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, shouldDeleteImageDesc bool, ctx *fiber.Ctx) (*entities.News, error)
	DeleteNewsByID(id string) error
}

type NewsUseCaseImpl struct {
	newsrepo repositories.NewsRepository
	config   configs.Supabase
}

func NewNewsUseCase(newsrepo repositories.NewsRepository, config configs.Supabase) *NewsUseCaseImpl {
	return &NewsUseCaseImpl{
		newsrepo: newsrepo,
		config:   config,
	}
}

func (u *NewsUseCaseImpl) CreateNews(news *entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, ctx *fiber.Ctx) (*entities.News, error) {
	id, err := u.newsrepo.GetNewsNextID()
	if err != nil {
		return nil, err
	}

	if imageTitleFile != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(imageTitleFile, "./uploads/"+fileName); err != nil {
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

		news.Image_Title = imageUrl
	}

	if imageDescFile != nil {
		fileName := uuid.New().String() + "_desc.jpg"
		if err := ctx.SaveFile(imageDescFile, "./uploads/"+fileName); err != nil {
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

		news.Image_Desc = imageUrl
	}

	for i, dialogReq := range news.Dialog {
		news.Dialog[i] = entities.Dialog{
			ID:     uuid.New().String(),
			Type:   dialogReq.Type,
			Desc:   dialogReq.Desc,
			Bold:   dialogReq.Bold,
			NewsID: id,
		}
	}

	news.ID = id
	createdNews, err := u.newsrepo.CreateNews(news)
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

func (u *NewsUseCaseImpl) UpdateNewsByID(id string, news entities.News, imageTitleFile *multipart.FileHeader, imageDescFile *multipart.FileHeader, shouldDeleteImageDesc bool, ctx *fiber.Ctx) (*entities.News, error) {
	existingNews, err := u.newsrepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}

	existingNews.Title = news.Title
	if imageTitleFile != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(imageTitleFile, "./uploads/"+fileName); err != nil {
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

		existingNews.Image_Title = imageUrl
	}

	if shouldDeleteImageDesc {
		existingNews.Image_Desc = ""
	} else if imageDescFile != nil {
		fileName := uuid.New().String() + "_desc.jpg"
		if err := ctx.SaveFile(imageDescFile, "./uploads/"+fileName); err != nil {
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

		existingNews.Image_Desc = imageUrl
	}

	for _, dialog := range existingNews.Dialog {
		if err := u.newsrepo.DeleteDialog(dialog.ID); err != nil {
			return nil, err
		}
	}

	for i, dialogReq := range news.Dialog {
		news.Dialog[i] = entities.Dialog{
			ID:     uuid.New().String(),
			Type:   dialogReq.Type,
			Desc:   dialogReq.Desc,
			Bold:   dialogReq.Bold,
			NewsID: existingNews.ID,
		}
	}

	existingNews.Dialog = news.Dialog
	updatedNews, err := u.newsrepo.UpdateNewsByID(existingNews)
	if err != nil {
		return nil, err
	}

	return updatedNews, nil
}

func (u *NewsUseCaseImpl) DeleteNewsByID(id string) error {
	existingNews, err := u.newsrepo.GetNewsByID(id)
	if err != nil {
		return err
	}

	for _, dialog := range existingNews.Dialog {
		if err := u.newsrepo.DeleteDialog(dialog.ID); err != nil {
			return err
		}
	}

	if err := u.newsrepo.DeleteNewsByID(id); err != nil {
		return err
	}

	return nil
}
