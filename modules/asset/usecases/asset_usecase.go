package usecases

import (
	"time"
	"errors"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
)

type AssetUseCase interface {
	CreateAsset(assancial entities.Asset) (*entities.Asset, error)
	GetAssetByID(id string) (*entities.Asset, error)
	GetAssetByUserID(userID string) ([]entities.Asset, error)
	UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error)
	DeleteAssetByID(id string) error
	CalculateMonthlyExpenses(asset *entities.Asset) (float64, error)
}

type AssetUseCaseImpl struct {
	assetrepo 		repositories.AssetRepository
}

func NewAssetUseCase(assetrepo repositories.AssetRepository) *AssetUseCaseImpl {
	return &AssetUseCaseImpl{assetrepo:  assetrepo}
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
	return u.assetrepo.GetAssetByUserID(userID)
}

func (u *AssetUseCaseImpl) UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error){
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
	
    currentYear := time.Now().Year()
    if endYear < currentYear {
        return nil, errors.New("end year must be greater than or equal to current year")
    }

	existingAsset.TotalCost = asset.TotalCost
	existingAsset.Name = asset.Name
	existingAsset.Type = asset.Type
	existingAsset.EndYear = asset.EndYear

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

func (u *AssetUseCaseImpl) CalculateMonthlyExpenses(asset *entities.Asset) (float64, error) {
    endYear, err := strconv.Atoi(asset.EndYear)
    if err != nil {
        return 0, err
    }

    currentYear := time.Now().Year()
    if endYear < currentYear {
        return 0, errors.New("end year must be greater than or equal to current year")
    }

    remainingMonths := (endYear - currentYear) * 12 + (12 - int(time.Now().Month()) + 1)
    remainingCost := asset.TotalCost - asset.CurrentMoney
    if remainingCost < 0 {
        return 0, errors.New("current money cannot exceed total cost")
    }

    return remainingCost / float64(remainingMonths), nil
}