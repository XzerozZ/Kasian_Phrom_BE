package controllers

import (

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/usecases"

	"github.com/gofiber/fiber/v2"
	"fmt"
)

type RetController struct {
	retusecase usecases.RetUseCase
}

func NewRetController(retusecase usecases.RetUseCase) *RetController {
	return &RetController{retusecase: retusecase}
}

func (c *RetController) CreateRetHandler(ctx *fiber.Ctx) error {
	fmt.Println("เข้ามาใน controller")
	type RequestBody struct {
		UserID string `json:"userID"`
	}
	var reqBody RequestBody
	
	if err := ctx.BodyParser(&reqBody); err != nil {
		fmt.Println("อยู่ใน BodyParser")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	// ตรวจสอบค่า userID
	if reqBody.UserID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":       fiber.ErrBadRequest.Message,
			"status_code":  fiber.ErrBadRequest.Code,
			"message":      "userID is required",
			"result":       nil,
		})
	}

	fmt.Println("ผ่าน BodyParser")
	data, err := c.retusecase.CreateRet(reqBody.UserID, ctx)
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
		"message":     	"Retirement Plan created successfully",
		"result":      	data,
	})
}

func (c *RetController) GetRetByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.retusecase.GetRetByID(id)
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
		"message":     	"Retirement plan retrieved successfully",
		"result":      	data,
	})
}

func (c *RetController) GetRetNextIDHandler(ctx *fiber.Ctx) error {
	data, err := c.retusecase.GetRetNextID()
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
		"message":     	"Retirement plan retrieved successfully",
		"result":      	data,
	})
}