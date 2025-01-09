package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/repositories"

	"github.com/gofiber/fiber/v2"
)

type AssUseCase interface {
	CreateAss(assancial entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error)
	GetAssByID(id string) (*entities.Asset, error)
	GetAssByUsername(username string) ([]entities.Asset, error)
	GetAssNextID() (string, error)
	UpdateAssByID(id string, asset entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error)
	DeleteAssByID(id string) error
	
}

type AssUseCaseImpl struct {
	assrepo 		repositories.AssRepository
}

func NewAssUseCase(assrepo repositories.AssRepository) *AssUseCaseImpl {
	return &AssUseCaseImpl{
		assrepo:  assrepo,
	}
}

func (u *AssUseCaseImpl) CreateAss(asset entities.Asset,ctx *fiber.Ctx) (*entities.Asset, error) {
	
	
	id, err := u.assrepo.GetAssNextID()
	if err != nil {
		return nil, err
	}
	if asset.TotalMoney <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "TotalMoney must be greater than zero",
		})
	}
	
	if asset.Name == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Name cannot be empty",
		})
	}
	
	if asset.MonthlyExpenses <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "MonthlyExpenses must be greater than zero",
		})
	}
	
	if asset.EndYear == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "EndYear cannot be empty",
		})
	}

	asset.ID = id

	var createdAss *entities.Asset
	createdAss, err = u.assrepo.CreateAss(&asset)
    if err != nil {
        return nil, err
    }
    
    return createdAss, nil
}


func (u *AssUseCaseImpl) GetAssByID(id string) (*entities.Asset, error) {
	return u.assrepo.GetAssByID(id)
}

func (u *AssUseCaseImpl) GetAssByUsername(username string) ([]entities.Asset, error) {
	return u.assrepo.GetAssByUsername(username)
}

func (u *AssUseCaseImpl) GetAssNextID() (string, error) {
	return u.assrepo.GetAssNextID()
}

func (u *AssUseCaseImpl) UpdateAssByID(id string, asset entities.Asset, ctx *fiber.Ctx) (*entities.Asset, error){
	if asset.TotalMoney <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "TotalMoney must be greater than zero",
		})
	}
	
	if asset.Name == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Asset Name cannot be empty",
		})
	}
	
	if asset.MonthlyExpenses <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "MonthlyExpenses must be greater than zero",
		})
	}
	
	if asset.EndYear == "" {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "EndYear cannot be empty",
		})
	}

	existingAss, err := u.assrepo.GetAssByID(id)
	if err != nil {
		return nil, err
	}

	existingAss.TotalMoney = asset.TotalMoney
	existingAss.Name = asset.Name
	existingAss.MonthlyExpenses = asset.MonthlyExpenses
	existingAss.EndYear = asset.EndYear

	updatedAss, err := u.assrepo.UpdateAssByID(existingAss)
    if err != nil {
        return nil, err
    }

	return updatedAss, nil
	
}

func (u *AssUseCaseImpl) DeleteAssByID(id string) error {
    existingAss, err := u.assrepo.GetAssByID(id)
    if err != nil {
        return err
    }

    if err := u.assrepo.DeleteAssByID(existingAss.ID); err != nil {
        return err
    }

    return nil
}