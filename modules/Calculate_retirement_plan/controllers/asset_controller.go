package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type AssController struct {
	assusecase usecases.AssUseCase
}

func NewAssController(assusecase usecases.AssUseCase) *AssController {
	return &AssController{assusecase: assusecase}
}

func (c *AssController) CreateAssHandler(ctx *fiber.Ctx) error {

	var asset entities.Asset

	if err := ctx.BodyParser(&asset); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	data, err := c.assusecase.CreateAss(asset, ctx)
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
		"result":      data,
	})
}

func (c *AssController) GetAssByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.assusecase.GetAssByID(id)
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
		"message":     "Asset retrieved successfully",
		"result":      data,
	})
}

func (c *AssController) GetAssNextIDHandler(ctx *fiber.Ctx) error {
	data, err := c.assusecase.GetAssNextID()
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
		"message":     "Asset get next id successfully",
		"result":      data,
	})
}


func (c *AssController) UpdateAssByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var asset entities.Asset

	if err := ctx.BodyParser(&asset); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	updatedAss, err := c.assusecase.UpdateAssByID(id, asset, ctx)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Asset update successfully",
		"result":      	updatedAss,
	})
}


func (c *AssController) DeleteAssByIDHandler(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    err := c.assusecase.DeleteAssByID(id)
    if err != nil {
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

func (c *AssController) GetAssByUsernameHandler(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	data, err := c.assusecase.GetAssByUsername(username)
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
		"message":     "Asset retrieved successfully",
		"result":      data,
	})
}
