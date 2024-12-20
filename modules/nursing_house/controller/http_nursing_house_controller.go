package controller

import (
	"github.com/gofiber/fiber/v2"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type HttpNhController struct {
	NhUse entities.NhUsecase
}

func NewHttpNhController(r fiber.Router,useCase entities.NhUsecase) {
	controllers := &HttpNhController{
		NhUse: NhUse,
	}
	r.Post("/Nursing_house", controllers.CreateNh)
}

func (h *HttpNhController) CreateNh(c *fiber.Ctx) error {
	req := new(entities.Nursing_House)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	res, err := h.NhUse.CreateNh(req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "",
		"result":      res,
	})
  }