package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/usecases"
	"github.com/gofiber/fiber/v2"
)

type NotiController struct {
	notiusecase usecases.NotiUsecase
}

func NewNotiController(notiusecase usecases.NotiUsecase) *NotiController {
	return &NotiController{notiusecase: notiusecase}
}

func (c *NotiController) GetNotificationsByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	notis, err := c.notiusecase.GetNotificationsByUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "No Notification found for this user",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Notification retrieved successfully",
		"result":      notis,
	})
}

func (c *NotiController) MarkAsReadHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.notiusecase.MarkNotificationsAsRead(userID); err != nil {
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
		"message":     "Read notification successfully",
	})
}
