package usecases_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		user := &entities.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		role := &entities.Role{
			ID:       1,
			RoleName: "User",
		}

		userRepo.On("FindUserByEmail", "test@example.com").Return((*entities.User)(nil), errors.New("user not found"))
		userRepo.On("GetRoleByName", "User").Return(*role, nil)
		userRepo.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(user, nil)

		result, err := useCase.Register(user, "User")

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		assert.Equal(t, role.ID, user.RoleID)
		assert.Equal(t, "Credentials", user.Provider)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		user := &entities.User{
			Email:    "existing@example.com",
			Password: "password123",
		}

		userRepo.On("FindUserByEmail", "existing@example.com").Return(user, nil)

		result, err := useCase.Register(user, "User")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "this email already have account", err.Error())

		userRepo.AssertExpectations(t)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		user := &entities.User{
			Email:    "new@example.com",
			Password: "password123",
		}

		userRepo.On("FindUserByEmail", "new@example.com").Return((*entities.User)(nil), errors.New("user not found"))
		userRepo.On("GetRoleByName", "InvalidRole").Return(entities.Role{}, errors.New("role not found"))

		result, err := useCase.Register(user, "InvalidRole")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "role not found", err.Error())

		userRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-id-1",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Provider: "Credentials",
			Role: entities.Role{
				RoleName: "User",
			},
		}

		userRepo.On("FindUserByEmail", "test@example.com").Return(user, nil)

		token, result, err := useCase.Login("test@example.com", password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, user, result)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "user-id-1", claims["user_id"])
		assert.Equal(t, "User", claims["role"])

		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		userRepo.On("FindUserByEmail", "nonexistent@example.com").Return((*entities.User)(nil), errors.New("user not found"))

		token, result, err := useCase.Login("nonexistent@example.com", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, result)
		assert.Equal(t, "invalid email", err.Error())

		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-id-1",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Provider: "Credentials",
		}

		userRepo.On("FindUserByEmail", "test@example.com").Return(user, nil)

		token, result, err := useCase.Login("test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, result)
		assert.Equal(t, "invalid password", err.Error())

		userRepo.AssertExpectations(t)
	})
}

func TestLoginAdmin(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "admin-id-1",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Role: entities.Role{
				RoleName: "Admin",
			},
		}

		userRepo.On("FindUserByEmail", "admin@example.com").Return(user, nil)

		token, result, err := useCase.LoginAdmin("admin@example.com", password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, user, result)

		userRepo.AssertExpectations(t)
	})

	t.Run("Not Admin Role", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-id-1",
			Email:    "user@example.com",
			Password: string(hashedPassword),
			Role: entities.Role{
				RoleName: "User",
			},
		}

		userRepo.On("FindUserByEmail", "user@example.com").Return(user, nil)

		token, result, err := useCase.LoginAdmin("user@example.com", password)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, result)
		assert.Equal(t, "access denied: only admins can login", err.Error())

		userRepo.AssertExpectations(t)
	})
}

func TestLoginWithGoogle(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Existing User", func(t *testing.T) {
		existingUser := &entities.User{
			ID:       "google-user-id",
			Email:    "google@example.com",
			Provider: "Google",
			Role: entities.Role{
				RoleName: "User",
			},
		}

		googleUser := &entities.User{
			Email: "google@example.com",
		}

		userRepo.On("FindUserByEmail", "google@example.com").Return(existingUser, nil)

		token, user, err := useCase.LoginWithGoogle(googleUser)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, existingUser, user)

		userRepo.AssertExpectations(t)
	})

	t.Run("Wrong Provider", func(t *testing.T) {
		existingUser := &entities.User{
			ID:       "user-id",
			Email:    "existing@example.com",
			Provider: "Credentials",
		}

		googleUser := &entities.User{
			Email: "existing@example.com",
		}

		userRepo.On("FindUserByEmail", "existing@example.com").Return(existingUser, nil)

		token, user, err := useCase.LoginWithGoogle(googleUser)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Equal(t, "this email is already registered with another authentication method", err.Error())

		userRepo.AssertExpectations(t)
	})
}

func TestResetPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		userID := "user-123"
		oldPassword := "oldpassword123"
		newPassword := "newpassword456"

		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       userID,
			Password: string(hashedOldPassword),
		}

		userRepo.On("GetUserByID", userID).Return(user, nil)
		userRepo.On("UpdateUserByID", mock.MatchedBy(func(u *entities.User) bool {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(newPassword))
			return u.ID == userID && err == nil
		})).Return(user, nil)

		err := useCase.ResetPassword(userID, oldPassword, newPassword)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := "nonexistent-user"
		oldPassword := "oldpassword123"
		newPassword := "newpassword456"

		userRepo.On("GetUserByID", userID).Return((*entities.User)(nil), errors.New("user not found"))

		err := useCase.ResetPassword(userID, oldPassword, newPassword)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Old Password", func(t *testing.T) {
		userID := "user-123"
		oldPassword := "oldpassword123"
		wrongOldPassword := "wrongoldpassword"
		newPassword := "newpassword456"

		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       userID,
			Password: string(hashedOldPassword),
		}

		userRepo.On("GetUserByID", userID).Return(user, nil)

		err := useCase.ResetPassword(userID, wrongOldPassword, newPassword)

		assert.Error(t, err)
		assert.Equal(t, "invalid old password", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Same Old and New Password", func(t *testing.T) {
		userID := "user-123"
		password := "samepassword123"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       userID,
			Password: string(hashedPassword),
		}

		userRepo.On("GetUserByID", userID).Return(user, nil)

		err := useCase.ResetPassword(userID, password, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid old password", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Update Error", func(t *testing.T) {
		userID := "user-123"
		oldPassword := "oldpassword123"
		newPassword := "newpassword456"

		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       userID,
			Password: string(hashedOldPassword),
		}

		userRepo.On("GetUserByID", userID).Return(user, nil)
		userRepo.On("UpdateUserByID", mock.AnythingOfType("*entities.User")).Return((*entities.User)(nil), errors.New("invalid old password"))

		err := useCase.ResetPassword(userID, oldPassword, newPassword)

		assert.Error(t, err)
		assert.Equal(t, "invalid old password", err.Error())
		userRepo.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		expectedUser := &entities.User{
			ID:        "user-123",
			Email:     "user@example.com",
			Firstname: "Test User",
		}

		userRepo.On("GetUserByID", "user-123").Return(expectedUser, nil)

		user, err := useCase.GetUserByID("user-123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		userRepo.AssertExpectations(t)
	})
}

func TestGetSelectedHouse(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	defaultHouseID := "default-house-id"

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Default House", func(t *testing.T) {
		inputHouse := &entities.SelectedHouse{
			UserID:              "user-123",
			NursingHouseID:      defaultHouseID,
			MonthlyExpenses:     1000.0,
			LastCalculatedMonth: 2,
		}

		expectedUpdatedHouse := &entities.SelectedHouse{
			UserID:              "user-123",
			NursingHouseID:      defaultHouseID,
			MonthlyExpenses:     0,
			LastCalculatedMonth: 3,
		}

		user := &entities.User{
			ID: "user-123",
			RetirementPlan: entities.RetirementPlan{
				BirthDate: "01-01-1990",
			},
		}

		userRepo.On("GetUserByID", "user-123").Return(user, nil)
		userRepo.On("GetSelectedHouse", "user-123").Return(inputHouse, nil)
		userRepo.On("UpdateSelectedHouse", mock.MatchedBy(func(h *entities.SelectedHouse) bool {
			return h.UserID == "user-123" &&
				h.NursingHouseID == defaultHouseID &&
				h.MonthlyExpenses == 0 &&
				h.LastCalculatedMonth == 3
		})).Return(expectedUpdatedHouse, nil)

		house, err := useCase.GetSelectedHouse("user-123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUpdatedHouse, house)

		userRepo.AssertExpectations(t)
	})
}

func TestUpdateUserByID(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	app := fiber.New()

	t.Run("Successful Update Without File", func(t *testing.T) {
		existingUser := &entities.User{
			ID:        "user-123",
			Username:  "oldusername",
			Firstname: "Old",
			Lastname:  "User",
			ImageLink: "old-image.jpg",
		}

		updateUser := entities.User{
			Username:  "newusername",
			Firstname: "New",
			Lastname:  "User",
		}

		expectedUser := &entities.User{
			ID:        "user-123",
			Username:  "newusername",
			Firstname: "New",
			Lastname:  "User",
			ImageLink: "old-image.jpg",
		}

		userRepo.On("GetUserByID", "user-123").Return(existingUser, nil).Once()
		userRepo.On("UpdateUserByID", mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "user-123" &&
				u.Username == "newusername" &&
				u.Firstname == "New" &&
				u.Lastname == "User" &&
				u.ImageLink == "old-image.jpg"
		})).Return(expectedUser, nil).Once()

		fiberCtx := app.AcquireCtx(&fasthttp.RequestCtx{})
		fiberCtx.Request().SetRequestURI("/api/user/user-123")
		fiberCtx.Request().Header.SetMethod("PUT")
		defer app.ReleaseCtx(fiberCtx)

		updatedUser, err := useCase.UpdateUserByID("user-123", updateUser, nil, fiberCtx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, updatedUser)
		userRepo.AssertExpectations(t)
	})

	t.Run("Successful Update With File", func(t *testing.T) {
		existingUser := &entities.User{
			ID:        "user-123",
			Username:  "oldusername",
			Firstname: "Old",
			Lastname:  "User",
			ImageLink: "old-image.jpg",
		}

		updateUser := entities.User{
			Username:  "newusername",
			Firstname: "New",
			Lastname:  "User",
		}

		expectedUser := &entities.User{
			ID:        "user-123",
			Username:  "newusername",
			Firstname: "New",
			Lastname:  "User",
			ImageLink: "new-image-url.jpg",
		}

		userRepo.On("GetUserByID", "user-123").Return(existingUser, nil).Once()
		userRepo.On("UpdateUserByID", mock.AnythingOfType("*entities.User")).Return(expectedUser, nil).Once()

		ctx := new(fiber.Ctx)

		updatedUser, err := useCase.UpdateUserByID("user-123", updateUser, nil, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, updatedUser)

		userRepo.AssertExpectations(t)
	})

	t.Run("Error Getting User", func(t *testing.T) {
		expectedError := errors.New("user not found")

		userRepo.On("GetUserByID", "user-123").Return((*entities.User)(nil), expectedError).Once()

		updatedUser, err := useCase.UpdateUserByID("user-123", entities.User{}, nil, nil)

		assert.EqualError(t, err, expectedError.Error())
		assert.Nil(t, updatedUser)

		userRepo.AssertExpectations(t)
	})

	t.Run("Error Updating User", func(t *testing.T) {
		existingUser := &entities.User{
			ID:        "user-123",
			Username:  "oldusername",
			Firstname: "Old",
			Lastname:  "User",
			ImageLink: "old-image.jpg",
		}

		updateUser := entities.User{
			Username:  "newusername",
			Firstname: "New",
			Lastname:  "User",
		}

		expectedError := errors.New("update error")

		userRepo.On("GetUserByID", "user-123").Return(existingUser, nil).Once()
		userRepo.On("UpdateUserByID", mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "user-123" &&
				u.Username == "newusername" &&
				u.Firstname == "New" &&
				u.Lastname == "User" &&
				u.ImageLink == "old-image.jpg"
		})).Return((*entities.User)(nil), expectedError).Once()

		fiberCtx := app.AcquireCtx(&fasthttp.RequestCtx{})
		fiberCtx.Request().SetRequestURI("/api/user/user-123")
		fiberCtx.Request().Header.SetMethod("PUT")
		defer app.ReleaseCtx(fiberCtx)

		updatedUser, err := useCase.UpdateUserByID("user-123", updateUser, nil, fiberCtx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, updatedUser)
		userRepo.AssertExpectations(t)
	})
}

func TestForgotPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	type otpGenerator func(length int, onlyDigits bool) (string, error)

	createUserUseCase := func(
		generateOTP otpGenerator,
	) *usecases.UserUseCaseImpl {
		uc := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)
		ucValue := reflect.ValueOf(uc).Elem()

		if generateOTP != nil {
			field := ucValue.FieldByName("GenerateOTP")
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(generateOTP))
			}
		}

		return uc
	}

	t.Run("Invalid Email", func(t *testing.T) {
		email := "nonexistent@example.com"

		userRepo.On("FindUserByEmail", email).Return((*entities.User)(nil), errors.New("user not found")).Once()

		useCase := createUserUseCase(nil)

		err := useCase.ForgotPassword(email)

		assert.Error(t, err)
		assert.Equal(t, "invalid email", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Error Creating OTP", func(t *testing.T) {
		email := "test@example.com"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		otpCode := "123456"
		expectedError := errors.New("database error")

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return((*entities.OTP)(nil), errors.New("not found")).Once()
		userRepo.On("CreateOTP", mock.MatchedBy(func(otp *entities.OTP) bool {
			return otp.UserID == "user-123" && len(otp.OTP) == 6
		})).Return(expectedError).Once()

		mockGenerateOTP := func(length int, onlyDigits bool) (string, error) {
			return otpCode, nil
		}

		useCase := createUserUseCase(mockGenerateOTP)

		err := useCase.ForgotPassword(email)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

}

func TestVerifyOTP(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Successful OTP Verification", func(t *testing.T) {
		email := "test@example.com"
		otpCode := "123456"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		otp := &entities.OTP{
			UserID:    "user-123",
			OTP:       otpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return(otp, nil).Once()
		userRepo.On("DeleteOTP", "user-123").Return(nil).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		email := "nonexistent@example.com"
		otpCode := "123456"
		expectedError := errors.New("user not found")

		userRepo.On("FindUserByEmail", email).Return((*entities.User)(nil), expectedError).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("OTP Not Found", func(t *testing.T) {
		email := "test@example.com"
		otpCode := "123456"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		expectedError := errors.New("otp not found")

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return((*entities.OTP)(nil), expectedError).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("Expired OTP", func(t *testing.T) {
		email := "test@example.com"
		otpCode := "123456"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		otp := &entities.OTP{
			UserID:    "user-123",
			OTP:       otpCode,
			ExpiresAt: time.Now().Add(-5 * time.Minute),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return(otp, nil).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.Error(t, err)
		assert.Equal(t, "OTP is expired", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Incorrect OTP", func(t *testing.T) {
		email := "test@example.com"
		otpCode := "123456"
		wrongOtpCode := "654321"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		otp := &entities.OTP{
			UserID:    "user-123",
			OTP:       wrongOtpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return(otp, nil).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.Error(t, err)
		assert.Equal(t, "OTP is incorrect", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Error Deleting OTP", func(t *testing.T) {
		email := "test@example.com"
		otpCode := "123456"
		user := &entities.User{
			ID:    "user-123",
			Email: email,
		}
		otp := &entities.OTP{
			UserID:    "user-123",
			OTP:       otpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		expectedError := errors.New("database error")

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("GetOTPByUserID", "user-123").Return(otp, nil).Once()
		userRepo.On("DeleteOTP", "user-123").Return(expectedError).Once()

		err := useCase.VerifyOTP(email, otpCode)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})
}

func TestChangedPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Successful Password Change", func(t *testing.T) {
		email := "test@example.com"
		oldPassword := "oldPassword123"
		newPassword := "newPassword456"

		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-123",
			Email:    email,
			Password: string(hashedOldPassword),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("UpdateUserByID", mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "user-123" &&
				u.Email == email &&
				u.Password != string(hashedOldPassword)
		})).Return(user, nil).Once()

		err := useCase.ChangedPassword(email, newPassword)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		email := "nonexistent@example.com"
		newPassword := "newPassword456"
		expectedError := errors.New("user not found")

		userRepo.On("FindUserByEmail", email).Return((*entities.User)(nil), expectedError).Once()

		err := useCase.ChangedPassword(email, newPassword)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("Same Password Error", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-123",
			Email:    email,
			Password: string(hashedPassword),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()

		err := useCase.ChangedPassword(email, password)

		assert.Error(t, err)
		assert.Equal(t, "new password cannot be the same as the old password", err.Error())
		userRepo.AssertExpectations(t)
	})

	t.Run("Update Error", func(t *testing.T) {
		email := "test@example.com"
		oldPassword := "oldPassword123"
		newPassword := "newPassword456"
		expectedError := errors.New("database error")

		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		user := &entities.User{
			ID:       "user-123",
			Email:    email,
			Password: string(hashedOldPassword),
		}

		userRepo.On("FindUserByEmail", email).Return(user, nil).Once()
		userRepo.On("UpdateUserByID", mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "user-123" && u.Email == email
		})).Return((*entities.User)(nil), expectedError).Once()

		err := useCase.ChangedPassword(email, newPassword)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})
}

func TestUpdateSelectedHouse(t *testing.T) {
	t.Run("should return error when GetSelectedHouse fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		retirementRepo := new(mocks.MockRetirementRepository)
		assetRepo := new(mocks.MockAssetRepository)
		notiRepo := new(mocks.MockNotiRepository)
		nhRepo := new(mocks.MockNhRepository)

		jwtConfig := configs.JWT{Secret: "test-secret"}
		supaConfig := configs.Supabase{}
		mailConfig := configs.Mail{}

		userID := "user-123"
		nursingHouseID := "house-123"
		expectedError := errors.New("selected house not found")

		userRepo.On("GetSelectedHouse", userID).Return((*entities.SelectedHouse)(nil), expectedError)

		useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

		result, err := useCase.UpdateSelectedHouse(userID, nursingHouseID, []entities.TransferRequest{})

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetUserByID fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		retirementRepo := new(mocks.MockRetirementRepository)
		assetRepo := new(mocks.MockAssetRepository)
		notiRepo := new(mocks.MockNotiRepository)
		nhRepo := new(mocks.MockNhRepository)

		jwtConfig := configs.JWT{Secret: "test-secret"}
		supaConfig := configs.Supabase{}
		mailConfig := configs.Mail{}

		userID := "user-123"
		nursingHouseID := "house-123"
		selectedHouse := &entities.SelectedHouse{
			UserID:         userID,
			NursingHouseID: "house-456",
			CurrentMoney:   1000,
			Status:         "In_Progress",
		}
		expectedError := errors.New("user not found")

		userRepo.On("GetSelectedHouse", userID).Return(selectedHouse, nil)
		userRepo.On("GetUserByID", userID).Return((*entities.User)(nil), expectedError)

		useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

		result, err := useCase.UpdateSelectedHouse(userID, nursingHouseID, []entities.TransferRequest{})

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("should return error when FindAssetByNameandUserID fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		retirementRepo := new(mocks.MockRetirementRepository)
		assetRepo := new(mocks.MockAssetRepository)
		notiRepo := new(mocks.MockNotiRepository)
		nhRepo := new(mocks.MockNhRepository)

		jwtConfig := configs.JWT{Secret: "test-secret"}
		supaConfig := configs.Supabase{}
		mailConfig := configs.Mail{}

		userID := "user-123"
		nursingHouseID := "default"
		selectedHouse := &entities.SelectedHouse{
			UserID:         userID,
			NursingHouseID: "house-456",
			CurrentMoney:   1000,
			Status:         "In_Progress",
			NursingHouse:   entities.NursingHouse{Name: "Premium House"},
		}
		user := &entities.User{
			ID:        userID,
			Firstname: "Test User",
			RetirementPlan: entities.RetirementPlan{
				BirthDate:      "02-01-2006",
				ExpectLifespan: 80,
				RetirementAge:  60,
			},
		}

		transfers := []entities.TransferRequest{
			{Type: "asset", Name: "Car", Amount: 500},
		}

		expectedError := errors.New("asset not found")

		userRepo.On("GetSelectedHouse", userID).Return(selectedHouse, nil)
		userRepo.On("GetUserByID", userID).Return(user, nil)
		assetRepo.On("FindAssetByNameandUserID", "Car", userID).Return(nil, expectedError)
		nhRepo.On("GetNhByID", nursingHouseID).Return(&entities.NursingHouse{ID: nursingHouseID}, nil)
		userRepo.On("UpdateSelectedHouse", mock.Anything).Return((*entities.SelectedHouse)(nil), expectedError).Times(0)

		useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

		result, err := useCase.UpdateSelectedHouse(userID, nursingHouseID, transfers)

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetNhByID fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		retirementRepo := new(mocks.MockRetirementRepository)
		assetRepo := new(mocks.MockAssetRepository)
		notiRepo := new(mocks.MockNotiRepository)
		nhRepo := new(mocks.MockNhRepository)

		jwtConfig := configs.JWT{Secret: "test-secret"}
		supaConfig := configs.Supabase{}
		mailConfig := configs.Mail{}

		useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

		userID := "user-123"
		nursingHouseID := "new-house-123"
		selectedHouse := &entities.SelectedHouse{
			UserID:         userID,
			NursingHouseID: "house-456",
			CurrentMoney:   1000,
			Status:         "In_Progress",
		}
		user := &entities.User{
			ID:        userID,
			Firstname: "Test User",
			RetirementPlan: entities.RetirementPlan{
				BirthDate:      "02-01-2006",
				ExpectLifespan: 80,
				RetirementAge:  60,
			},
		}

		expectedError := errors.New("nursing house not found")

		userRepo.On("GetSelectedHouse", userID).Return(selectedHouse, nil)
		userRepo.On("GetUserByID", userID).Return(user, nil)
		nhRepo.On("GetNhByID", nursingHouseID).Return(nil, expectedError)

		result, err := useCase.UpdateSelectedHouse(userID, nursingHouseID, []entities.TransferRequest{})

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
		nhRepo.AssertExpectations(t)
	})

	t.Run("should return error when UpdateSelectedHouse fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		retirementRepo := new(mocks.MockRetirementRepository)
		assetRepo := new(mocks.MockAssetRepository)
		notiRepo := new(mocks.MockNotiRepository)
		nhRepo := new(mocks.MockNhRepository)

		jwtConfig := configs.JWT{Secret: "test-secret"}
		supaConfig := configs.Supabase{}
		mailConfig := configs.Mail{}

		useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

		userID := "user-123"
		nursingHouseID := "new-house-123"
		selectedHouse := &entities.SelectedHouse{
			UserID:              userID,
			NursingHouseID:      "house-456",
			CurrentMoney:        1000,
			Status:              "In_Progress",
			LastCalculatedMonth: 0,
		}
		user := &entities.User{
			ID:        userID,
			Firstname: "Test User",
			RetirementPlan: entities.RetirementPlan{
				BirthDate:      "02-01-2006",
				ExpectLifespan: 80,
				RetirementAge:  60,
			},
			House: entities.SelectedHouse{
				CurrentMoney: 1000,
			},
		}
		nursingHouse := &entities.NursingHouse{
			ID:    nursingHouseID,
			Name:  "New Premium House",
			Price: 1500,
		}

		expectedError := errors.New("failed to update selected house")

		userRepo.On("GetSelectedHouse", userID).Return(selectedHouse, nil)
		userRepo.On("GetUserByID", userID).Return(user, nil)
		nhRepo.On("GetNhByID", nursingHouseID).Return(nursingHouse, nil)
		userRepo.On("UpdateSelectedHouse", mock.AnythingOfType("*entities.SelectedHouse")).Return((*entities.SelectedHouse)(nil), expectedError)

		result, err := useCase.UpdateSelectedHouse(userID, nursingHouseID, []entities.TransferRequest{})

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		userRepo.AssertExpectations(t)
		nhRepo.AssertExpectations(t)
	})
}

func TestCalculateRetirement(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("GetUserByID Error", func(t *testing.T) {
		expectedError := errors.New("user not found")

		userRepo.On("GetUserByID", "invalid-user-id").Return((*entities.User)(nil), expectedError)

		result, err := useCase.CalculateRetirement("invalid-user-id")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, fiber.Map{}, result)

		userRepo.AssertExpectations(t)
	})

	t.Run("UpdateRetirementPlan Error", func(t *testing.T) {
		currentMonth := int(time.Now().Month())
		lastMonth := currentMonth - 1
		if lastMonth == 0 {
			lastMonth = 12
		}

		user := &entities.User{
			ID: "user-123",
			RetirementPlan: entities.RetirementPlan{
				BirthDate:           "01-02-2006",
				CreatedAt:           time.Now().AddDate(0, -1, 0),
				LastCalculatedMonth: lastMonth,
				RetirementAge:       65,
				ExpectLifespan:      85,
			},
			Assets: []entities.Asset{
				{
					ID:                  "asset-123",
					LastCalculatedMonth: currentMonth,
					TotalCost:           50000,
					MonthlyExpenses:     500,
				},
			},
			House: entities.SelectedHouse{
				Status:              "Owned",
				CurrentMoney:        100000,
				LastCalculatedMonth: currentMonth,
				MonthlyExpenses:     2000,
				NursingHouse: entities.NursingHouse{
					Price: 5000,
				},
			},
		}

		expectedError := errors.New("database error when getting deposits")

		userRepo.On("GetUserByID", "user-123").Return(user, nil)
		retirementRepo.On("UpdateRetirementPlan", mock.AnythingOfType("*entities.RetirementPlan")).Return((*entities.RetirementPlan)(nil), expectedError)

		result, err := useCase.CalculateRetirement("user-123")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)

		userRepo.AssertExpectations(t)
		retirementRepo.AssertExpectations(t)
	})
}

func TestCreateHistory(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Success", func(t *testing.T) {
		currentMonth := int(time.Now().Month())
		lastMonth := currentMonth - 1
		if lastMonth == 0 {
			lastMonth = 12
		}

		user := &entities.User{
			ID: "user-123",
			RetirementPlan: entities.RetirementPlan{
				BirthDate:           "01-02-2006",
				CreatedAt:           time.Now().AddDate(0, -1, 0),
				LastCalculatedMonth: lastMonth,
				RetirementAge:       65,
				ExpectLifespan:      85,
			},
			Assets: []entities.Asset{
				{
					ID:                  "asset-123",
					LastCalculatedMonth: currentMonth,
					TotalCost:           50000,
					MonthlyExpenses:     500,
				},
			},
			House: entities.SelectedHouse{
				Status:              "Owned",
				CurrentMoney:        100000,
				LastCalculatedMonth: currentMonth,
				MonthlyExpenses:     2000,
				NursingHouse: entities.NursingHouse{
					Price: 5000,
				},
			},
		}

		userRepo.On("GetUserByID", "user-123").Return(user, nil)
		retirementRepo.On("UpdateRetirementPlan", mock.AnythingOfType("*entities.RetirementPlan")).Return(&user.RetirementPlan, nil)

		userRepo.On("GetUserDepositsInRange", mock.Anything, mock.Anything, mock.Anything).Return([]entities.History{}, nil)

		result, err := useCase.CalculateRetirement("user-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)

		userRepo.AssertExpectations(t)
		retirementRepo.AssertExpectations(t)
	})

	t.Run("Negative Case - User Not Found", func(t *testing.T) {
		userRepo.On("GetUserByID", "invalid-user").Return((*entities.User)(nil), errors.New("user not found"))

		result, err := useCase.CalculateRetirement("invalid-user")

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.NotNil(t, result)

		userRepo.AssertExpectations(t)
	})

	t.Run("Negative Case - Money Must Be Greater Than Zero", func(t *testing.T) {
		history := entities.History{
			UserID:   "user-123",
			Method:   "deposit",
			Type:     "saving_money",
			Category: "retirementplan",
			Money:    0,
			Name:     "",
		}

		mockUser := &entities.User{
			ID: "user-123",
			RetirementPlan: entities.RetirementPlan{
				ID:     "plan-123",
				Status: "In_Progress",
			},
		}

		userRepo.On("GetUserByID", "user-123").Return(mockUser, nil)

		result, err := useCase.CreateHistory(history)

		assert.Error(t, err)
		assert.Equal(t, "money must be greater than zero", err.Error())
		assert.Nil(t, result)

		userRepo.AssertExpectations(t)
	})
}

func TestGetHistoryByUserID(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	retirementRepo := new(mocks.MockRetirementRepository)
	assetRepo := new(mocks.MockAssetRepository)
	notiRepo := new(mocks.MockNotiRepository)
	nhRepo := new(mocks.MockNhRepository)

	jwtConfig := configs.JWT{Secret: "test-secret"}
	supaConfig := configs.Supabase{}
	mailConfig := configs.Mail{}

	useCase := usecases.NewUserUseCase(userRepo, retirementRepo, assetRepo, notiRepo, nhRepo, jwtConfig, supaConfig, mailConfig)

	t.Run("Positive Case - Retrieve History Successfully", func(t *testing.T) {
		mockHistories := []entities.History{
			{
				ID:        "history1",
				UserID:    "user123",
				Method:    "deposit",
				Type:      "saving_money",
				Category:  "retirementplan",
				Money:     1000,
				TrackDate: time.Now(),
			},
			{
				ID:        "history2",
				UserID:    "user123",
				Method:    "withdraw",
				Type:      "saving_money",
				Category:  "retirementplan",
				Money:     500,
				TrackDate: time.Now(),
			},
		}

		userRepo.On("GetHistoryByUserID", "user123").Return(mockHistories, nil)
		userRepo.On("GetHistoryInRange", "user123", mock.Anything, mock.Anything).Return(mockHistories, nil)

		result, err := useCase.GetHistoryByUserID("user123")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockHistories, result["data"])
		assert.Equal(t, 500.0, result["total"])

		userRepo.AssertCalled(t, "GetHistoryByUserID", "user123")
		userRepo.AssertCalled(t, "GetHistoryInRange", "user123", mock.Anything, mock.Anything)
	})

	t.Run("Negative Case - User Not Found", func(t *testing.T) {
		expectedError := errors.New("user not found")
		userRepo.On("GetHistoryByUserID", "invalid_user").Return([]entities.History{}, expectedError)

		result, err := useCase.GetHistoryByUserID("invalid_user")

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, fiber.Map{}, result)

		userRepo.AssertCalled(t, "GetHistoryByUserID", "invalid_user")
		userRepo.AssertNotCalled(t, "GetHistoryInRange")
	})
}
