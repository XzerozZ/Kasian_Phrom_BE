package controllers

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/usecases"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userusecase usecases.UserUseCase
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
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Username == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Username is missing",
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	if req.RoleName == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "RoleName is missing",
			"result":      nil,
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
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "user created successfully",
		"result":      data,
	})
}

func (c *UserController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	token, user, err := c.userusecase.Login(req.Email, req.Password)
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
		"message":     "Login successful",
		"result": fiber.Map{
			"token": token,
			"u_id":  user.ID,
			"uname": user.Username,
			"role":  user.Role.RoleName,
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
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	token, user, err := c.userusecase.LoginAdmin(req.Email, req.Password)
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
		"message":     "Login successful",
		"result": fiber.Map{
			"token": token,
			"u_id":  user.ID,
			"uname": user.Username,
			"role":  user.Role.RoleName,
		},
	})
}

func (c *UserController) LoginWithGoogleHandler(ctx *fiber.Ctx) error {
	var req struct {
		Firstname string `json:"fname"`
		Lastname  string `json:"lname"`
		Username  string `json:"uname"`
		Email     string `json:"email"`
		ImageLink string `json:"image_link"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Firstname == "" || req.Lastname == "" || req.Username == "" || req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Firstname, Lastname, Username, Email or ImageLink is missing",
			"result":      nil,
		})
	}

	if req.ImageLink == "" {
		req.ImageLink = "https://mvfxlcnhrtduomirjeir.supabase.co/storage/v1/object/public/photos/seProfile/UserProfileDefault.jpg"
	}

	user := &entities.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Username,
		Email:     req.Email,
		ImageLink: req.ImageLink,
	}

	token, user, err := c.userusecase.LoginWithGoogle(user)
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
		"message":     "Login successful",
		"result": fiber.Map{
			"token": token,
			"u_id":  user.ID,
			"uname": user.Username,
			"role":  user.Role.RoleName,
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

func (c *UserController) ForgotPasswordHandler(ctx *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.ForgotPassword(req.Email)
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
		"message":     "Sent OTP successfully",
	})
}

func (c *UserController) VerifyOTPHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp"`
	}

	var req OTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.OTP == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "OTP is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.VerifyOTP(req.Email, req.OTP)
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
		"message":     "OTP is correct",
	})
}

func (c *UserController) ChangedPasswordHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email       string `json:"email" validate:"required,email"`
		NewPassword string `json:"newPassword"`
	}

	var req OTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "NewPassword is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.ChangedPassword(req.Email, req.NewPassword)
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
		"message":     "changed password successfully",
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
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "User Info retrieved successfully",
		"result":      data,
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Selected house retrieved successfully",
		"result":      selectedHouse,
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
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
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

	var file *multipart.FileHeader
	files := form.File["images"]
	if len(files) > 0 {
		file = files[0]
	}

	updatedUser, err := c.userusecase.UpdateUserByID(userID, user, file, ctx)
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
		"message":     "User retrieved successfully",
		"result":      updatedUser,
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

	var transfers []entities.TransferRequest
	if nursingHouseID == "00001" {
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Invalid form data",
				"result":      nil,
			})
		}

		types, names, amounts := form.Value["type"], form.Value["name"], form.Value["amount"]
		if len(types) == 0 || len(names) == 0 || len(amounts) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Missing required fields: 'type', 'name', or 'amount'",
				"result":      nil,
			})
		}

		if len(types) != len(names) || len(types) != len(amounts) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Mismatch in count of 'type', 'name', and 'amount'",
				"result":      nil,
			})
		}

		transfers = make([]entities.TransferRequest, 0, len(types))
		for i := 0; i < len(types); i++ {
			amount, err := strconv.ParseFloat(amounts[i], 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusBadRequest,
					"message":     fmt.Sprintf("Invalid amount format at index %d: %s", i, amounts[i]),
					"result":      nil,
				})
			}

			if amount < 0 {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusBadRequest,
					"message":     fmt.Sprintf("Amount cannot be negative at index %d: %f", i, amount),
					"result":      nil,
				})
			}

			transfers = append(transfers, entities.TransferRequest{
				Type:   types[i],
				Name:   names[i],
				Amount: amount,
			})
		}
	}

	updatedHouse, err := c.userusecase.UpdateSelectedHouse(userID, nursingHouseID, transfers)
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
		"message":     "House Updated to user successfully",
		"result":      updatedHouse,
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
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "This is user's retirement plan successfully",
		"result":      requiredFunds,
	})
}

func (c *UserController) CreateHistoryHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	var history entities.History
	if err := ctx.BodyParser(&history); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if history.Method == "" || history.Type == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Method or Type is missing.",
			"result":      nil,
		})
	}

	history.UserID = userID
	createdHistory, err := c.userusecase.CreateHistory(history)
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
		"message":     "History created successfully",
		"result":      createdHistory,
	})
}

func (c *UserController) GetHistoryByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	history, err := c.userusecase.GetHistoryByUserID(userID)
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
		"message":     "Retirement retrieved successfully",
		"result":      history,
	})
}

func (c *UserController) GetSummaryHistoryByUserIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	history, err := c.userusecase.GetHistoryByMonth(userID)
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
		"message":     "Retirement retrieved successfully",
		"result":      history,
	})
}
