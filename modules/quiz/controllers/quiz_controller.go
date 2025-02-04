package controllers

import (
	"fmt"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/usecases"
	"github.com/gofiber/fiber/v2"
)

type QuizController struct {
	quizusecase usecases.QuizUseCase
}

func NewQuizController(quizusecase usecases.QuizUseCase) *QuizController {
	return &QuizController{quizusecase: quizusecase}
}

func (c *QuizController) CreateQuizHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Failed to parse form data",
			"result":      nil,
		})
	}

	weightValues, exists := form.Value["weight"]
	if !exists || len(weightValues) != 10 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Must answer 12 quiz",
			"result":      nil,
		})
	}

	weights := make([]int, 10)
	for i, weightStr := range weightValues {
		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     fmt.Sprintf("Invalid weight value at position %d", i+1),
				"result":      nil,
			})
		}
		weights[i] = weight
	}

	quiz, err := c.quizusecase.CreateQuiz(userID, weights)
	if err != nil {
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
		"message":     "Quiz created successfully",
		"result":      quiz,
	})
}

func (c *QuizController) GetQuizByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.quizusecase.GetQuizByUserID(userID)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     "This user has not answered quiz yet.",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Quiz retrieved successfully",
		"result":      data,
	})
}
