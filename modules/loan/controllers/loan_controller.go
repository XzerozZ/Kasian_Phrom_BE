package controller

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type LoanController struct {
	loanusecase usecases.LoanUseCase
}

func NewLoanController(loanusecase usecases.LoanUseCase) *LoanController {
	return &LoanController{loanusecase: loanusecase}
}

func (c *LoanController) CreateLoanHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var loan entities.Loan
	if err := ctx.BodyParser(&loan); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if loan.Name == "" || loan.Type == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Name or Type is missing.",
			"result":      nil,
		})
	}

	loan.UserID = userID
	createdLoan, err := c.loanusecase.CreateLoan(loan)
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
		"message":     "Loan created successfully",
		"result":      createdLoan,
	})
}

func (c *LoanController) GetLoanByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.loanusecase.GetLoanByID(id)
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

func (c *LoanController) GetLoanByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	loans, err := c.loanusecase.GetLoanByUserID(userID)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if len(loans) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "No Asset found for this user",
			"result":      nil,
		})
	}

	totalLoans, err := utils.CalculateAllLoan(loans)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	response := fiber.Map{
		"loan":   totalLoans,
		"number": len(loans),
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Loans retrieved successfully",
		"result":      response,
	})
}

func (c *LoanController) DeleteLoanHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.loanusecase.DeleteLoanByID(id); err != nil {
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
		"message":     "Loan deleted successfully",
		"result":      nil,
	})
}
