package usecases

import (
	"os"
	"errors"
	"strconv"
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
	CalculateRetirement(userID string) (float64, error)
}

type UserUseCaseImpl struct {
	userrepo 		repositories.UserRepository
	retirementrepo 	retirementRepo.RetirementRepository
	jwtSecret		string
	supa			configs.Supabase
}

func NewUserUseCase(userrepo repositories.UserRepository, retirementrepo retirementRepo.RetirementRepository, jwt configs.JWT, supa configs.Supabase) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userrepo: 		userrepo,
		retirementrepo: retirementrepo,
		jwtSecret:		jwt.Secret,
		supa:  			supa,
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

	if nursingHouseID == "00001" {
        retirementPlan, err := u.retirementrepo.GetRetirementByID(userID)
        if err != nil {
            return nil, err
        }

        retirementPlan.CurrentSavings += selectedHouse.CurrentMoney
		selectedHouse.CurrentMoney = 0
        _, err = u.retirementrepo.UpdateRetirementPlan(retirementPlan)
        if err != nil {
            return nil, err
        }

		updatedHouse, err = u.userrepo.UpdateSelectedHouse(selectedHouse)
        if err != nil {
            return nil, err
        }

		return updatedHouse, nil
    }

	return updatedHouse, nil
}

func CalculateAllAssets(user *entities.User) (float64, error) {
	total := 0.0
	for _, asset := range user.Assets {
		endYear, err := strconv.Atoi(asset.EndYear)
		if err != nil {
			return 0, err
		}

		remainingMonths := (endYear - asset.UpdatedAt.Year() -1 ) * 12 + (12 - int(asset.UpdatedAt.Month()) + 1)
		if remainingMonths <= 0 { 
			return 0, err
		}
		remainingCost := (asset.TotalCost - asset.CurrentMoney) / float64(remainingMonths)
		total += remainingCost
	}

	return total, nil
}

func (u *UserUseCaseImpl) CalculateRetirement(userID string) (float64, error) {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	allCostAsset, err := CalculateAllAssets(user)
	if err != nil {
		return 0, err
	}

	plan := user.RetirementPlan
	nursingHousePrice := 0.0
	if user.House.NursingHouse.ID != "" {
		nursingHousePrice = float64(user.House.NursingHouse.Price)
	}

	monthlyPlan := utils.MonthlyExpensesPlan{
		MonthlyExpenses:      	plan.MonthlyExpenses,
		AnnualExpenseIncrease: 	plan.AnnualExpenseIncrease,
		ExpectedInflation:    	plan.ExpectedInflation,
		Age:                  	plan.Age,
		RetirementAge:        	plan.RetirementAge,
		ExpectLifespan:       	plan.ExpectLifespan,
		YearsUntilRetirement: 	plan.RetirementAge - plan.Age,
		YearUntilLifeSpan:		plan.ExpectLifespan - plan.RetirementAge,
		AllCostAsset:           allCostAsset,
		NursingHousePrice:    	nursingHousePrice,
	}

	requiredFunds, err := utils.CalculateMonthlySavings(monthlyPlan)
	if err != nil {
		return 0, err
	}

	return requiredFunds, nil
}