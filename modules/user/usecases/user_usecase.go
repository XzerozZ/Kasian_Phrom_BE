package usecases

import (
	"os"
	"time"
	"errors"
	"mime/multipart"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	retirementRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserUseCase interface {
	Register(user *entities.User, roleName string) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
	LoginAdmin(email, password string) (string, *entities.User, error)
	ResetPassword(userID, oldPassword, newPassword string) error
	GetUserByID(userID string) (*entities.User, error)
	GetSelectedHouse(userID string) (*entities.SelectedHouse, error)
	UpdateUserByID(id string, user entities.User, files *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	UpdateSelectedHouse(userID, nursingHouseID string)  (*entities.SelectedHouse, error)
	ForgotPassword(email string) error
	CalculateRetirement(userID string) (fiber.Map, error)
}

type UserUseCaseImpl struct {
	userrepo 		repositories.UserRepository
	retirementrepo 	retirementRepo.RetirementRepository
	jwtSecret		string
	supa			configs.Supabase
	mail			configs.Mail
}

func NewUserUseCase(userrepo repositories.UserRepository, retirementrepo retirementRepo.RetirementRepository, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userrepo: 		userrepo,
		retirementrepo: retirementrepo,
		jwtSecret:		jwt.Secret,
		supa:  			supa,
		mail:			mail,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string) (*entities.User, error) {
	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.ID = uuid.New().String()
	user.RoleID = role.ID
	user.Role = role
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := u.userrepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserUseCaseImpl) LoginAdmin(email, password string) (string, *entities.User, error) {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	if user.Role.RoleName != "Admin" {
		return "", nil, errors.New("access denied: only admins can login")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role.RoleName,
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	
	return tokenString, &user, nil
}


func (u *UserUseCaseImpl) Login(email, password string) (string, *entities.User, error) {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role.RoleName,
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	
	return tokenString, &user, nil
}

func (u *UserUseCaseImpl) ResetPassword(userID, oldPassword, newPassword string) error {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	if oldPassword == newPassword {
		return errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	_, err = u.userrepo.UpdateUserByID(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCaseImpl) GetUserByID(userID string) (*entities.User, error) {
	return u.userrepo.GetUserByID(userID)
}

func (u *UserUseCaseImpl) GetSelectedHouse(userID string) (*entities.SelectedHouse, error) {
	return u.userrepo.GetSelectedHouse(userID)
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user entities.User, file *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	if file != nil {
		fileName := uuid.New().String() + ".jpg"
		if err := ctx.SaveFile(file, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		existingUser.ImageLink = imageUrl
	}

	updatedUser, err := u.userrepo.UpdateUserByID(existingUser)
    if err != nil {
        return nil, err
    }

	return updatedUser, nil
}

func (u *UserUseCaseImpl) UpdateSelectedHouse(userID, nursingHouseID string) (*entities.SelectedHouse, error) {
	selectedHouse, err := u.userrepo.GetSelectedHouse(userID)
	if err != nil {
		return nil, err
	}

	selectedHouse.NursingHouseID = nursingHouseID
	updatedHouse, err := u.userrepo.UpdateSelectedHouse(selectedHouse)
    if err != nil {
        return nil, err
    }

	return updatedHouse, nil
}

func (u *UserUseCaseImpl) ForgotPassword(email string) error {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return errors.New("invalid email")
	}

	userID := user.ID
	otpCode, err := utils.GenerateRandomOTP(6, true)
    if err != nil {
        return err
    }

	expiresAt := time.Now().Add(5 * time.Minute)
	otp, err := u.userrepo.GetOTPByUserID(userID)
	if err == nil && otp != nil {
		if err := u.userrepo.DeleteOTP(userID); err != nil {
			return err
		}
	}

	newOTP := &entities.OTP{
		UserID:    userID,
		OTP:       otpCode,
		ExpiresAt: expiresAt,
	}

	if err := u.userrepo.CreateOTP(newOTP); err != nil {
		return err
	}

    if err := utils.SendMail("./assets/OTPMail.html", user, otpCode, u.mail); err != nil {
        return err
    }

	return nil
}

func (u *UserUseCaseImpl) CalculateRetirement(userID string) (fiber.Map, error) {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return fiber.Map{}, err
	}

	allCostAsset, err := utils.CalculateAllAssetsMonthlyExpenses(user)
	if err != nil {
		return fiber.Map{}, err
	}

	plan := user.RetirementPlan
	nursingHousePrice := 0.0
	if user.House.NursingHouse.ID != "" {
		nursingHousePrice, err = utils.CalculateNursingHouseMonthlyExpenses(user)
		if err != nil {
			return fiber.Map{}, err
		}
	}
	
	monthsUntilRetirement := (plan.RetirementAge * 12) - plan.AgeInMonths
	monthlyPlan := utils.MonthlyExpensesPlan{
		ExpectedMonthlyExpenses:      	plan.ExpectedMonthlyExpenses,
		AnnualExpenseIncrease: 			plan.AnnualExpenseIncrease,
		ExpectedInflation:    			plan.ExpectedInflation,
		Age:                  			plan.Age,
		RetirementAge:        			plan.RetirementAge,
		ExpectLifespan:       			plan.ExpectLifespan,
		MonthsUntilRetirement: 			monthsUntilRetirement,
		YearUntilLifeSpan:				plan.ExpectLifespan - plan.RetirementAge,
		AllCostAsset:           		allCostAsset,
		NursingHousePrice:    			nursingHousePrice,
	}

	requiredFunds, err := utils.CalculateMonthlySavings(monthlyPlan)
	if err != nil {
		return fiber.Map{}, err
	}

	requiredAllFunds, err := utils.CalculateRetirementFunds(monthlyPlan)
	if err != nil {
		return fiber.Map{}, err
	}

	assetSavings, err :=  utils.CalculateAllAssetSavings(user)
	if err != nil {
		return fiber.Map{}, err
	}
	
	yearUntilLifespan := plan.ExpectLifespan - plan.RetirementAge
	totalNursingHouseCost := float64(user.House.NursingHouse.Price * yearUntilLifespan) - user.House.CurrentMoney
	allRequiredFund := requiredAllFunds + allCostAsset + totalNursingHouseCost
	allSaving:= plan.CurrentSavings + assetSavings + user.House.CurrentMoney
	allMoney := allSaving + plan.CurrentTotalInvestment
	stillNeed := allRequiredFund - allMoney 
	response := fiber.Map{
		"allRequiredFund":				allRequiredFund,
		"stillneed":					stillNeed,
		"allretirementfund": 			requiredAllFunds,
        "monthly_expenses": 			requiredFunds,
		"all_money": 					float64(allMoney),
		"saving": 						float64(allSaving),
		"investment": 					float64(plan.CurrentTotalInvestment),
    }

	return response, nil
}