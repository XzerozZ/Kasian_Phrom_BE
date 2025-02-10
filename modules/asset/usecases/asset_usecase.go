package usecases

import (
	"errors"
	"strconv"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type AssetUseCase interface {
	CreateAsset(asset entities.Asset) (*entities.Asset, error)
	GetAssetByID(id string) (*entities.Asset, error)
	GetAssetByUserID(userID string) ([]entities.Asset, error)
	UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error)
	DeleteAssetByID(id string) error
}

type AssetUseCaseImpl struct {
	assetrepo repositories.AssetRepository
}

func NewAssetUseCase(assetrepo repositories.AssetRepository) *AssetUseCaseImpl {
	return &AssetUseCaseImpl{assetrepo: assetrepo}
}

func (u *AssetUseCaseImpl) CreateAsset(asset entities.Asset) (*entities.Asset, error) {
	id, err := u.assetrepo.GetAssetNextID()
	if err != nil {
		return nil, err
	}

	if asset.TotalCost <= 0 {
		return nil, errors.New("totalcost must be greater than zero")
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return nil, err
	}

	currentYear := time.Now().Year()
	if endYear < currentYear {
		return nil, errors.New("end year must be greater than or equal to current year")
	}

	asset.ID = id
	createdAsset, err := u.assetrepo.CreateAsset(&asset)
	if err != nil {
		return nil, err
	}

	return createdAsset, nil
}

func (u *AssetUseCaseImpl) GetAssetByID(id string) (*entities.Asset, error) {
	return u.assetrepo.GetAssetByID(id)
}

func (u *AssetUseCaseImpl) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	assets, err := u.assetrepo.GetAssetByUserID(userID)
	if err != nil {
		return nil, err
	}

	currentYear := time.Now().Year()
	for i := range assets {
		endYear, err := strconv.Atoi(assets[i].EndYear)
		if err != nil {
			return nil, err
		}

		if assets[i].Status != "completed" && endYear <= currentYear {
			assets[i].Status = "paused"
			// message := fmt.Sprintf("สินทรัพย์ '%s' ถูกหยุดพักชั่วคราวเนื่องจากหมดเวลา", assets[i].Name)
			// notification := &entities.Notification{
			// 	ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), assets[i].ID),
			// 	UserID:    userID,
			// 	Message:   message,
			// 	CreatedAt: time.Now(),
			// }
			// u.notificationrepo.CreateNotification(notification)

			_, err := u.assetrepo.UpdateAssetByID(&assets[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return assets, nil
}

func (u *AssetUseCaseImpl) UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error) {
	existingAsset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return nil, err
	}

	if asset.TotalCost <= 0 {
		return nil, errors.New("totalcost must be greater than zero")
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return nil, err
	}

	existingAsset.TotalCost = asset.TotalCost
	existingAsset.Name = asset.Name
	existingAsset.Type = asset.Type
	existingAsset.EndYear = asset.EndYear
	existingAsset.Status = asset.Status
	currentYear := time.Now().Year()
	if endYear <= currentYear {
		existingAsset.Status = "Paused"
	}

	if existingAsset.Status != "Paused" {
		if existingAsset.TotalCost <= existingAsset.CurrentMoney {
			existingAsset.Status = "Completed"
		} else {
			existingAsset.Status = "In_Progress"
		}
	}

	updatedAsset, err := u.assetrepo.UpdateAssetByID(existingAsset)
	if err != nil {
		return nil, err
	}

	return updatedAsset, nil
}

func (u *AssetUseCaseImpl) DeleteAssetByID(id string) error {
	existingAsset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return err
	}

	if err := u.assetrepo.DeleteAssetByID(existingAsset.ID); err != nil {
		return err
	}

	return nil
}
