package controllers

import (
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/usecases"

	"github.com/gofiber/fiber/v2"
)

type FinController struct {
	finusecase usecases.FinUseCase
}

func NewFinController(finusecase usecases.FinUseCase) *FinController {
	return &FinController{finusecase: finusecase}
}

func (c *FinController) CreateFinHandler(ctx *fiber.Ctx) error {
	var financial entities.Financial

	
	if err := ctx.BodyParser(&financial); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	data, err := c.finusecase.CreateFin(financial, ctx)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":		"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"financial created successfully",
		"result":      	data,
	})
}


func (c *FinController) GetFinByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.finusecase.GetFinByID(id)
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
		"message":     	"get financial successfully",
		"result":      	data,
	})
}

func (c *FinController) GetFinByUserIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.finusecase.GetFinByUserID(id)
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
		"message":     	"get financial successfully",
		"result":      	data,
	})
}

func (c *FinController) GetFinNextIDHandler(ctx *fiber.Ctx) error {
	data, err := c.finusecase.GetFinNextID()
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
		"message":     	"get fin next id successfully",
		"result":      	data,
	})
}