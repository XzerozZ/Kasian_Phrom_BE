package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/usecases"

	"github.com/gofiber/fiber/v2"
)

type RetirementController struct {
	retirementusecase usecases.RetirementUseCase
}

func NewRetirementController(retirementusecase usecases.RetirementUseCase) *RetirementController {
	return &RetirementController{retirementusecase: retirementusecase}
}

func (c *RetirementController) CreateRetirementHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var retirement entities.RetirementPlan
	if err := ctx.BodyParser(&retirement); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if retirement.PlanName == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "PlanName is missing.",
			"result":      nil,
		})
	}

	retirement.UserID = userID
	createdRetirement, err := c.retirementusecase.CreateRetirement(retirement)
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
		"result":      createdRetirement,
	})
}
