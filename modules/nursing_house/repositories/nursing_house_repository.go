package repositories

import (
	"fmt"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
)

type GormNhRepository struct {
	db *gorm.DB
}

func NewGormNhRepository(db *gorm.DB) *GormNhRepository {
	return &GormNhRepository{db: db}
}

type NhRepository interface {
	CreateNh(nursingHouse *entities.NursingHouse, images []entities.Image) (*entities.NursingHouse, error)
	GetAllNh() ([]entities.NursingHouse, error)
	GetActiveNh() ([]entities.NursingHouse, error)
	GetInactiveNh() ([]entities.NursingHouse, error)
	GetNhByID(id string) (*entities.NursingHouse, error)
	GetNhNextID() (string, error)
	UpdateNhByID(nursingHouse *entities.NursingHouse) (*entities.NursingHouse, error)
	AddImages(id string, images []entities.Image) (*entities.NursingHouse, error)
	RemoveImages(id string, imageID *string) error
}

func (r *GormNhRepository) CreateNh(nursingHouse *entities.NursingHouse, images []entities.Image) (*entities.NursingHouse, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(nursingHouse).Error; err != nil {
			return err
		}

		for _, image := range images {
			if err := tx.Create(&image).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(nursingHouse).Association("Images").Append(images); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetNhByID(nursingHouse.ID)
}

func (r *GormNhRepository) GetAllNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Preload("Images").Where("id != ?", "00001").Find(&nursingHouses).Error; err != nil {
		return nil, err
	}

	return nursingHouses, nil
}

func (r *GormNhRepository) GetActiveNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Preload("Images").Where("status = ? AND id != ?", "Active", "00001").Find(&nursingHouses).Error; err != nil {
		return nil, err
	}

	return nursingHouses, nil
}

func (r *GormNhRepository) GetInactiveNh() ([]entities.NursingHouse, error) {
	var nursingHouses []entities.NursingHouse
	if err := r.db.Preload("Images").Where("status = ? AND id != ?", "Inactive", "00001").Find(&nursingHouses).Error; err != nil {
		return nil, err
	}

	return nursingHouses, nil
}

func (r *GormNhRepository) GetNhByID(id string) (*entities.NursingHouse, error) {
	var nursingHouse entities.NursingHouse
	if err := r.db.Preload("Images").First(&nursingHouse, id).Error; err != nil {
		return nil, err
	}

	return &nursingHouse, nil
}

func (r *GormNhRepository) GetNhNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.NursingHouse{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
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

func (r *GormNhRepository) UpdateNhByID(nursingHouse *entities.NursingHouse) (*entities.NursingHouse, error) {
	if err := r.db.Save(&nursingHouse).Error; err != nil {
		return nil, err
	}

	return r.GetNhByID(nursingHouse.ID)
}

func (r *GormNhRepository) AddImages(id string, images []entities.Image) (*entities.NursingHouse, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var nursingHouse entities.NursingHouse
		if err := tx.First(&nursingHouse, id).Error; err != nil {
			return err
		}

		for i := range images {
			if err := tx.Create(&images[i]).Error; err != nil {
				return err
			}
		}

		return tx.Model(&nursingHouse).Association("Images").Append(images)
	})

	if err != nil {
		return nil, err
	}

	return r.GetNhByID(id)
}

func (r *GormNhRepository) RemoveImages(id string, imageID *string) error {
	var imagesToDelete []entities.Image
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var nursingHouse entities.NursingHouse
		if err := tx.Preload("Images").Where("id = ?", id).First(&nursingHouse).Error; err != nil {
			return err
		}

		for _, img := range nursingHouse.Images {
			if img.ID == *imageID {
				imagesToDelete = append(imagesToDelete, img)
				break
			}
		}

		var imageIDs []string
		for _, img := range imagesToDelete {
			imageIDs = append(imageIDs, img.ID)
		}

		if err := tx.Model(&nursingHouse).Association("Images").Delete(imagesToDelete); err != nil {
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

	return nil
}
