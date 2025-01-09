package usecases

import (
	"time"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"

	"github.com/gofiber/fiber/v2"
)

type AssetUseCase interface {
	CreateAsset(assancial entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error)
	GetAssetByID(id string) (*entities.Asset, error)
	GetAssetByUserID(userID string) ([]entities.Asset, error)
	UpdateAssetByID(id string, asset entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error)
	DeleteAssetByID(id string) error
	CalculateMonthlyExpenses(asset *entities.Asset, ctx *fiber.Ctx) (float64, error)
}

type AssetUseCaseImpl struct {
	assetrepo 		repositories.AssetRepository
}

func NewAssetUseCase(assetrepo repositories.AssetRepository) *AssetUseCaseImpl {
	return &AssetUseCaseImpl{assetrepo:  assetrepo}
}

func (u *AssetUseCaseImpl) CreateAsset(asset entities.Asset,ctx *fiber.Ctx) (*entities.Asset, error) {
	id, err := u.assetrepo.GetAssetNextID()
	if err != nil {
		return nil, err
	}
	
	if asset.TotalCost <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "TotalMoney must be greater than zero",
		})
	}
	
	if asset.Name == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Name cannot be empty",
		})
	}
	
	if asset.Type == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Type cannot be empty",
		})
	}
	
	if asset.EndYear == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "EndYear cannot be empty",
		})
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

func (u *AssetUseCaseImpl) UpdateAssetByID(id string, asset entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error){
	if asset.TotalCost <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "TotalMoney must be greater than zero",
		})
	}
	
	if asset.Name == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Name cannot be empty",
		})
	}
	
	if asset.Type == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Type cannot be empty",
		})
	}

	if asset.EndYear == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "EndYear cannot be empty",
		})
	}

	existingAsset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return nil, err
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

func (u *AssetUseCaseImpl) CalculateMonthlyExpenses(asset *entities.Asset, ctx *fiber.Ctx) (float64, error) {
    endYear, err := strconv.Atoi(asset.EndYear)
    if err != nil {
        return 0, err
    }

    currentYear := time.Now().Year()
    if endYear < currentYear {
        return 0, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "end year must be greater than or equal to current year",
		})
    }

    remainingMonths := (endYear - currentYear) * 12 + (12 - int(time.Now().Month()) + 1)
    remainingCost := asset.TotalCost - asset.CurrentMoney
    if remainingCost < 0 {
        return 0, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "current money cannot exceed total cost",
		})
    }

    return remainingCost / float64(remainingMonths), nil
}