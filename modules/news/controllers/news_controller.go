package controllers

import (
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/usecases"

	"github.com/gofiber/fiber/v2"
)

type NewsController struct {
	newsusecase usecases.NewsUseCase
}

func NewNewsController(newsusecase usecases.NewsUseCase) *NewsController {
	return &NewsController{newsusecase: newsusecase}
}

func (c *NewsController) CreateNewsHandler(ctx *fiber.Ctx) error {
	req := new(entities.CreateNewsRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	if err := c.newsusecase.CreateNews(req); err != nil {
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
		"result":      	req,
	})
}