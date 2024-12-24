package controllers

import (
	"strconv"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"

	"github.com/gofiber/fiber/v2"
)

type NhController struct {
	nhusecase usecases.NhUseCase
}

func NewNhController(nhusecase usecases.NhUseCase) *NhController {
	return &NhController{nhusecase: nhusecase}
}

func (c *NhController) CreateNhHandler(ctx *fiber.Ctx) error {
	var nursingHouse entities.NursingHouse
	if err := ctx.BodyParser(&nursingHouse); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	data, err := c.nhusecase.CreateNh(nursingHouse)
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
		"message":     	"Nursing house created successfully",
		"result":      	data,
	})
}

func (c *NhController) GetAllNhHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetAllNh()
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing houses retrieved successfully",
		"result":      	data,
	})
}

func (c *NhController) GetAllActiveNhHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetActiveNh()
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
		"message":     	"Nursing house retrieved successfully",
		"result":      	data,
	})
}

func (c *NhController) GetNhByIDHandler(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	data, err := c.nhusecase.GetNhByID(id)
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
		"message":     	"Nursing house retrieved successfully",
		"result":      	data,
	})
}

func (c *NhController) UpdateNhByIDHandler(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	var nursingHouse entities.NursingHouse
	if err := ctx.BodyParser(&nursingHouse); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	updatedNh, err := c.nhusecase.UpdateNhByID(int(id), nursingHouse)
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
		"message":     	"Nursing house retrieved successfully",
		"result":      	updatedNh,
	})
}

func (c *NhController) DeleteNhByIDHandler(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	err = c.nhusecase.DeleteNhByID(id)
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
		"message":     	"Nursing house Deleted successfully",
	})
}