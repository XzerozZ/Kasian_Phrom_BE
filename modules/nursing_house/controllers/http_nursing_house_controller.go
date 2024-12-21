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
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Nursing house created successfully",
		"result":      nursingHouse,
	})
  }