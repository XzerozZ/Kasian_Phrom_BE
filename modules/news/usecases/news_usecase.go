package usecases

import (
	"time"
	"errors"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/repositories"
)

type NewsUseCase interface {
	CreateNews(news *entities.CreateNewsRequest) error
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
}

type NewsUseCaseImpl struct {
	newsrepo 	repositories.NewsRepository
}

func NewNewsUseCase(newsrepo repositories.NewsRepository) *NewsUseCaseImpl {
	return &NewsUseCaseImpl{
		newsrepo:  newsrepo,
	}
}

func (u *NewsUseCaseImpl) CreateNews(req *entities.CreateNewsRequest) error {
	id, err := u.newsrepo.GetNewsNextID()
	if err != nil {
		return err
	}

	news := &entities.News{
		ID:          id,
		Title:       req.Title,	
		PublishedAt: time.Now(),
		UpdatedAt:   time.Now(),
		Dialog:      make([]entities.Dialog, len(req.Dialogs)),
	}

	if len(req.Dialogs) == 0 {
		return errors.New("dialogs cannot be empty")
	}

	for i, dialogReq := range req.Dialogs {
		news.Dialog[i] = entities.Dialog{
			Type:   dialogReq.Type,
			Desc:   dialogReq.Desc,
			NewsID: id,
		}
	}

	return u.newsrepo.CreateNews(news)
}

func (u *NewsUseCaseImpl) GetNewsByID(id string) (*entities.News, error) {
	return u.newsrepo.GetNewsByID(id)
}

func (u *NewsUseCaseImpl) GetNewsNextID() (string, error) {
	return u.newsrepo.GetNewsNextID()
}