package usecases

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	notiRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
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
	notirepo  notiRepositories.NotiRepository
}

func NewAssetUseCase(assetrepo repositories.AssetRepository, notirepo notiRepositories.NotiRepository) *AssetUseCaseImpl {
	return &AssetUseCaseImpl{
		assetrepo: assetrepo,
		notirepo:  notirepo,
	}
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

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	if endYear < currentYear {
		return nil, errors.New("end year must be greater than or equal to current year")
	}

	monthlyExpenses, err := utils.CalculateMonthlyExpenses(&asset)
	if err != nil {
		return nil, err
	}

	asset.ID = id
	asset.MonthlyExpenses = monthlyExpenses
	asset.LastCalculatedMonth = currentMonth
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

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	for i := range assets {
		endYear, err := strconv.Atoi(assets[i].EndYear)
		if err != nil {
			return nil, err
		}

		if assets[i].Status != "completed" && endYear <= currentYear {
			assets[i].Status = "paused"
			message := fmt.Sprintf("สินทรัพย์ '%s' ถูกหยุดพักชั่วคราวเนื่องจากหมดเวลา", assets[i].Name)
			notification := &entities.Notification{
				ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), assets[i].ID),
				UserID:    userID,
				Message:   message,
				CreatedAt: time.Now(),
			}

			_ = u.notirepo.CreateNotification(notification)
			_, err := u.assetrepo.UpdateAssetByID(&assets[i])
			if err != nil {
				return nil, err
			}
		}

		if assets[i].LastCalculatedMonth != currentMonth {
			newMonthlyExpenses, err := utils.CalculateMonthlyExpenses(&assets[i])
			if err == nil {
				assets[i].MonthlyExpenses = newMonthlyExpenses
				assets[i].LastCalculatedMonth = currentMonth
				_, err = u.assetrepo.UpdateAssetByID(&assets[i])
				if err != nil {
					return nil, err
				}
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

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	if existingAsset.LastCalculatedMonth != currentMonth || existingAsset.TotalCost != asset.TotalCost {
		monthlyExpenses, err := utils.CalculateMonthlyExpenses(existingAsset)
		if err != nil {
			return nil, err
		}
		existingAsset.MonthlyExpenses = monthlyExpenses
		existingAsset.LastCalculatedMonth = currentMonth
	}

	existingAsset.TotalCost = asset.TotalCost
	existingAsset.Name = asset.Name
	existingAsset.Type = asset.Type
	existingAsset.EndYear = asset.EndYear
	existingAsset.Status = asset.Status
	existingAsset.LastCalculatedMonth = currentMonth
	if endYear <= currentYear {
		existingAsset.Status = "Paused"
	} else {
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
