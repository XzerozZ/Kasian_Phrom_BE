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

func (u *AssetUseCaseImpl) CreateNotification(userID, assetName string) error {
	notification := &entities.Notification{
		ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), assetName),
		UserID:    userID,
		Message:   fmt.Sprintf("สินทรัพย์ '%s' ถูกหยุดพักชั่วคราวเนื่องจากหมดเวลา", assetName),
		CreatedAt: time.Now(),
	}
	return u.notirepo.CreateNotification(notification)
}

func (u *AssetUseCaseImpl) UpdateAssetStatus(asset *entities.Asset, currentYear int) error {
	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return err
	}

	if endYear <= currentYear {
		asset.Status = "Paused"
		asset.LastCalculatedMonth = 0
		asset.MonthlyExpenses = 0
		return u.CreateNotification(asset.UserID, asset.Name)
	}

	if asset.Status == "Paused" {
		asset.LastCalculatedMonth = 0
		asset.MonthlyExpenses = 0
		return nil
	}

	if asset.TotalCost <= asset.CurrentMoney {
		asset.Status = "Completed"
	} else {
		asset.Status = "In_Progress"
	}
	return nil
}

func (u *AssetUseCaseImpl) CreateAsset(asset entities.Asset) (*entities.Asset, error) {
	if asset.TotalCost <= 0 {
		return nil, errors.New("totalcost must be greater than zero")
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return nil, err
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	if endYear <= currentYear {
		return nil, errors.New("end year must be greater than or equal to current year")
	}

	id, err := u.assetrepo.GetAssetNextID()
	if err != nil {
		return nil, err
	}

	asset.ID = id
	asset.MonthlyExpenses = utils.CalculateMonthlyExpenses(&asset, currentYear, currentMonth)
	asset.LastCalculatedMonth = currentMonth
	return u.assetrepo.CreateAsset(&asset)
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
		if err := u.UpdateAssetStatus(&assets[i], currentYear); err != nil {
			return nil, err
		}

		if assets[i].LastCalculatedMonth != currentMonth {
			assets[i].MonthlyExpenses = utils.CalculateMonthlyExpenses(&assets[i], currentYear, currentMonth)
			assets[i].LastCalculatedMonth = currentMonth
			if _, err = u.assetrepo.UpdateAssetByID(&assets[i]); err != nil {
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

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	totalCostChanged := existingAsset.TotalCost != asset.TotalCost
	existingAsset.TotalCost = asset.TotalCost
	existingAsset.Name = asset.Name
	existingAsset.Type = asset.Type
	existingAsset.EndYear = asset.EndYear
	existingAsset.Status = asset.Status
	existingAsset.LastCalculatedMonth = currentMonth
	if asset.Status == "Paused" {
		existingAsset.LastCalculatedMonth = 0
		existingAsset.MonthlyExpenses = 0
	} else {
		if err := u.UpdateAssetStatus(existingAsset, currentYear); err != nil {
			return nil, err
		}

		if totalCostChanged || existingAsset.LastCalculatedMonth != currentMonth {
			existingAsset.LastCalculatedMonth = currentMonth
			existingAsset.MonthlyExpenses = utils.CalculateMonthlyExpenses(existingAsset, currentYear, currentMonth)
		}
	}

	return u.assetrepo.UpdateAssetByID(existingAsset)
}

func (u *AssetUseCaseImpl) DeleteAssetByID(id string) error {
	asset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return err
	}

	return u.assetrepo.DeleteAssetByID(asset.ID)
}
