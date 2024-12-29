package repositories

import (
	"fmt"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormNewsRepository struct {
  	db *gorm.DB
}

func NewGormNewsRepository(db *gorm.DB) *GormNewsRepository {
	return &GormNewsRepository{db: db}
}

type NewsRepository interface {
	CreateNews(news *entities.News, images []entities.Image) (*entities.News, error)
	GetAllNews() ([]entities.News, error)
	GetNewsByID(id string) (*entities.News, error)
	GetNewsNextID() (string, error)
	UpdateNewsByID(news *entities.News) (*entities.News, error)
	AddImages(id string, images []entities.Image) (*entities.News, error)
    RemoveImages(id string, imageID *string) error
	DeleteDialog(id string) error
	DeleteNewsByID(id string) error
}

func (r *GormNewsRepository) CreateNews(news *entities.News, images []entities.Image) (*entities.News, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(news).Error; err != nil {
			return err
		}

		for _, image := range images {
			if err := tx.Create(&image).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(news).Association("Images").Append(images); err != nil {
			return err
		}

		return nil
	})
	
	if err != nil {
		return nil, err
	}

	return r.GetNewsByID(news.ID)
}

func(r *GormNewsRepository) GetAllNews() ([]entities.News, error){
	var news []entities.News
	if err := r.db.Preload("Images").Find(&news).Error; err != nil {
		return nil, err
	}

	return news , nil
}

func (r *GormNewsRepository) GetNewsByID(id string) (*entities.News, error) {
	var news entities.News
	if err := r.db.Preload("Dialog").Preload("Images").First(&news, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &news , nil
}

func (r *GormNewsRepository) GetNewsNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.News{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
		return "", err
	}
	
	maxIDInt := 0
	if maxID != "" {
		maxIDInt, _ = strconv.Atoi(maxID)
	}
	nextID := maxIDInt + 1
	formattedID := fmt.Sprintf("%05d", nextID)
	return formattedID, nil
}

func (r *GormNewsRepository) UpdateNewsByID(news *entities.News) (*entities.News, error) {
	if err := r.db.Save(&news).Error; err != nil {
		return nil, err
	}

	return r.GetNewsByID(news.ID)
}

func (r *GormNewsRepository) AddImages(id string, images []entities.Image) (*entities.News, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
        var news entities.News
        if err := tx.First(&news, id).Error; err != nil {
            return err
        }
        
        for i := range images {
            if err := tx.Create(&images[i]).Error; err != nil {
                return err
            }
        }
        
        return tx.Model(&news).Association("Images").Append(images)
    })

    if err != nil {
        return nil, err
    }

    return r.GetNewsByID(id)
}

func (r *GormNewsRepository) RemoveImages(id string, imageID *string) error {
    var imagesToDelete []entities.Image
	err := r.db.Transaction(func(tx *gorm.DB) error {
        var news entities.News
		if err := tx.Preload("Images").Where("id = ?", id).First(&news).Error; err != nil {
            return err
        }

        for _, img := range news.Images {
			if img.ID == *imageID {
				imagesToDelete = append(imagesToDelete, img)
				break
			}
		}

		var imageIDs []string
        for _, img := range imagesToDelete {
            imageIDs = append(imageIDs, img.ID)
        }

        if err := tx.Model(&news).Association("Images").Delete(imagesToDelete); err != nil {
            return err
        }

		if err := tx.Where("id IN ?", imageIDs).Delete(&entities.Image{}).Error; err != nil {
            return err
        }

        return nil
    })

	if err != nil {
        return err
    }

	return  nil
}

func (r *GormNewsRepository) DeleteDialog(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Dialog{}).Error
}

func (r *GormNewsRepository) DeleteNewsByID(id string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Table("news_images").Where("news_id = ?", id).Delete(nil).Error; err != nil {
            return err
        }

        if err := tx.Where("news_id = ?", id).Delete(&entities.Dialog{}).Error; err != nil {
            return err
        }

        if err := tx.Where("id = ?", id).Delete(&entities.News{}).Error; err != nil {
            return err
        }

        return nil
    })
}
