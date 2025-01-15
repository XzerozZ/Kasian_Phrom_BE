package controllers

import (
	"mime/multipart"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
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

func (c *UserController) LoginAdminHandler(ctx *fiber.Ctx) error {
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

	token, user, err := c.userusecase.LoginAdmin(req.Email, req.Password)
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

func (c *UserController) ResetPasswordHandler(ctx *fiber.Ctx) error {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Old password or new password is missing",
			"result":      nil,
		})
	}
	
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	err := c.userusecase.ResetPassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}


	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Password reset successfully",
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

func (c *UserController) GetUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.userusecase.GetUserByID(userID)
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
		"message":     	"Nursing house retrieved successfully",
		"result":      	data,
	})
}

func (c *UserController) GetSelectedHouseHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	selectedHouse, err := c.userusecase.GetSelectedHouse(userID)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	data, err := c.userusecase.GetUserByID(userID)
	monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpenses(data)
    if err != nil {
        return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
            "status":      fiber.ErrInternalServerError.Message,
            "status_code": fiber.ErrInternalServerError.Code,
            "message":     err.Error(),
            "result":      nil,
        })
    }

	response := fiber.Map{
        "selected": selectedHouse,
        "monthly_expenses": monthlyExpenses,
    }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Selected house retrieved successfully",
		"result":      response,
	})
}

func (c *UserController) UpdateUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var user entities.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
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

	
    var file *multipart.FileHeader
	files := form.File["images"]
	if len(files) > 0 {
		file = files[0]
	}

	updatedUser, err := c.userusecase.UpdateUserByID(userID, user, file, ctx)
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
		"message":     	"User retrieved successfully",
		"result":      	updatedUser,
	})
}

func (c *UserController) UpdateSelectedHouseHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	nursingHouseID := ctx.Params("nh_id")
	if nursingHouseID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request: Missing nursing house ID",
			"result":      nil,
		})
	}

	updatedHouse, err := c.userusecase.UpdateSelectedHouse(userID, nursingHouseID)
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
		"message":     	"House Updated to user successfully",
		"result": 		updatedHouse,
	})
}

func (c *UserController) GetRetirementPlanHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	requiredFunds, err := c.userusecase.CalculateRetirement(userID)
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
		"message":     	"This is user's retirement plan successfully",
		"result":		requiredFunds,
	})
}