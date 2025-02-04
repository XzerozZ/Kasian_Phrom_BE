package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AssetController struct {
	assetusecase usecases.AssetUseCase
}

func NewAssetController(assetusecase usecases.AssetUseCase) *AssetController {
	return &AssetController{assetusecase: assetusecase}
}

func (c *AssetController) CreateAssetHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var asset entities.Asset
	if err := ctx.BodyParser(&asset); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if asset.Name == "" || asset.Type == "" || asset.EndYear == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Name, Type or EndYear is missing.",
			"result":      nil,
		})
	}

	asset.UserID = userID
	createdAsset, err := c.assetusecase.CreateAsset(asset)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Asset created successfully",
		"result":      createdAsset,
	})
}

func (c *AssetController) GetAssetByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.assetusecase.GetAssetByID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	monthlyExpenses, err := utils.CalculateMonthlyExpenses(data)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	response := fiber.Map{
		"asset":            data,
		"monthly_expenses": monthlyExpenses,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Asset retrieved successfully",
		"result":      response,
	})
}

func (c *AssetController) GetAssetByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	assets, err := c.assetusecase.GetAssetByUserID(userID)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if len(assets) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "No Asset found for this user",
			"result":      nil,
		})
	}

	var assetsWithExpenses []fiber.Map
	for _, asset := range assets {
		monthlyExpenses, err := utils.CalculateMonthlyExpenses(&asset)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"status":      fiber.ErrInternalServerError.Message,
				"status_code": fiber.ErrInternalServerError.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}

		assetResponse := fiber.Map{
			"asset":            asset,
			"monthly_expenses": monthlyExpenses,
		}

		assetsWithExpenses = append(assetsWithExpenses, assetResponse)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Asset retrieved successfully",
		"result":      assetsWithExpenses,
	})
}

func (c *AssetController) UpdateAssetByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var asset entities.Asset
	asset.UserID = userID
	if err := ctx.BodyParser(&asset); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if asset.Name == "" || asset.Type == "" || asset.EndYear == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Name, Type or EndYear is empty.",
			"result":      nil,
		})
	}

	updatedAsset, err := c.assetusecase.UpdateAssetByID(id, asset)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Asset update successfully",
		"result":      updatedAsset,
	})
}

func (c *AssetController) DeleteAssetByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.assetusecase.DeleteAssetByID(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Asset deleted successfully",
		"result":      nil,
	})
}
