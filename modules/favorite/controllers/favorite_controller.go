package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/usecases"

	"github.com/gofiber/fiber/v2"
)

type FavController struct {
	favusecase usecases.FavUseCase
}

func NewFavController(favusecase usecases.FavUseCase) *FavController {
	return &FavController{favusecase: favusecase}
}

func (c *FavController) CreateFavHandler(ctx *fiber.Ctx) error {
	var fav entities.Favorite
	if err := ctx.BodyParser(&fav); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid input data",
			"result":      nil,
		})
	}

	if err := c.favusecase.CreateFav(&fav); err != nil {
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
		"message":     "Favorite successfully",
	})
}

func (c *FavController) GetFavByUserIDHandler(ctx *fiber.Ctx) error {
	userID := ctx.Params("user_id")
	user, err := c.favusecase.GetFavByUserID(userID) 
	if err != nil {
    	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        	"status":      "Internal Server Error",
        	"status_code": fiber.StatusInternalServerError,
        	"message":     err.Error(),
        	"result":      nil,
    	})
	}

	favs := user.Favorites
	if len(favs) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "No favorites found for the given user",
			"result":      nil,
		})
	}

	var result []fiber.Map
	for _, fav := range favs {
		result = append(result, fiber.Map{
			"user_id":       fav.UserID,
			"nh_id":         fav.NursingHouseID,
			"NursingHouse":  fav.NursingHouse,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Favorites retrieved successfully",
		"result":      result,
	})
}

func (c *FavController) CheckFavHandler(ctx *fiber.Ctx) error {
	userID := ctx.Params("user_id")
	nursingHouseID := ctx.Params("nh_id")
	if err := c.favusecase.CheckFav(userID, nursingHouseID); err != nil {
		if err.Error() == "not favorited nursing house" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":      "Not Found",
				"status_code": fiber.StatusNotFound,
				"message":     "Not Favorited Nursing House",
				"result":      nil,
			})
		}
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
		"message":     "Favorited Nursing House",
	})
}

func (c *FavController) DeleteFavByIDHandler(ctx *fiber.Ctx) error {
	userID := ctx.Params("user_id")
	nursingHouseID := ctx.Params("nh_id")

	err := c.favusecase.DeleteFavByID(userID, nursingHouseID)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":      "Not Found",
				"status_code": fiber.StatusNotFound,
				"message":     "Favorite not found",
				"result":      nil,
			})
		}
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
		"message":     "Favorite deleted successfully",
		"result":      nil,
	})
}
