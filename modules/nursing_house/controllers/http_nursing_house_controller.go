package controllers

import (
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"

	"github.com/gofiber/fiber/v2"
)

type NhController struct {
	nhusecase usecases.NhUsecase
}

func NewNhController(nhusecase usecases.NhUsecase) *NhController {
	return &NhController{nhusecase: nhusecase}
}

func (c *NhController) CreateNhHandler(ctx *fiber.Ctx) error {
	nursingHouse := new(entities.NursingHouse)
	if err := ctx.BodyParser(nursingHouse); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	err := c.nhusecase.CreateNh(nursingHouse)
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
		"message":     "Nursing house created successfully",
		"result":      nursingHouse,
	})
}

func (c *NhController) GetAllNhHandler(ctx *fiber.Ctx) error {
	nhList, err := c.NhUsecase.GetAllNh()
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Nursing houses retrieved successfully",
		"result":      nhList,
	})
}

func (c *NhController) GetNhByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	nh, err := c.NhUsecase.GetNhByID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     "Nursing house retrieved successfully",
		"result":      nh,
	})
}