package controllers

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/usecases"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userusecase	usecases.UserUseCase
}

func NewUserController(userusecase usecases.UserUseCase) *UserController {
	return &UserController{userusecase: userusecase}
}

func (c *UserController) RegisterHandler(ctx *fiber.Ctx) error {
	var req struct {
		Username string `json:"uname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleName string `json:"role"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	data, err := c.userusecase.Register(user, req.RoleName)
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
		"message":     	"Nursing house created successfully",
		"result":      	data,
	})
}

func (c *UserController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	token, user, err := c.userusecase.Login(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message": 		"Login successful",
		"result":     	fiber.Map{
			"token":       token,
			"u_id":        user.ID,
			"uname":	   user.Username,
			"role":        user.Role.RoleName,
		},
	})
}

func (c *UserController) LogoutHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Logout successful",
		"result":      nil,
	})
}