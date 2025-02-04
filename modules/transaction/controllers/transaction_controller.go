package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/usecases"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	transusecase usecases.TransactionUseCase
}

func NewTransactionController(transusecase usecases.TransactionUseCase) *TransactionController {
	return &TransactionController{transusecase: transusecase}
}

func (c *TransactionController) CreateTransactionsForAllUsersHandler(ctx *fiber.Ctx) error {
	if err := c.transusecase.CreateTransactionsForAllUsers(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Create Loans' Transaction successfully",
	})
}

func (c *TransactionController) MarkTransactiontoPaidHandler(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.transusecase.MarkTransactiontoPaid(transactionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Update Transaction successfully",
	})
}

func (c *TransactionController) GetTransactionByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	transactions, err := c.transusecase.GetTransactionByUserID(userID)
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
		"message":     "User's Transactions retrieved successfully",
		"result":      transactions,
	})
}
