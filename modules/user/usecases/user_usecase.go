package usecases

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	assetRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	notiRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/repositories"
	nhRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	retirementRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(user *entities.User, roleName string) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
	LoginAdmin(email, password string) (string, *entities.User, error)
	LoginWithGoogle(user *entities.User) (string, *entities.User, error)
	ResetPassword(userID, oldPassword, newPassword string) error
	GetUserByID(userID string) (*entities.User, error)
	UpdateUserByID(id string, user entities.User, files *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	ForgotPassword(email string) error
	VerifyOTP(email, otpCode string) error
	ChangedPassword(email, newPassword string) error
	CalculateRetirement(userID string) (fiber.Map, error)

	GetSelectedHouse(userID string) (*entities.SelectedHouse, error)
	UpdateSelectedHouse(userID, nursingHouseID string) (*entities.SelectedHouse, error)

	CreateHistory(history entities.History) (*entities.History, error)
	GetHistoryByUserID(userID string) (fiber.Map, error)
	GetHistoryByMonth(userID string) (map[string]float64, error)
}

type UserUseCaseImpl struct {
	userrepo       repositories.UserRepository
	retirementrepo retirementRepo.RetirementRepository
	assetrepo      assetRepo.AssetRepository
	notirepo       notiRepo.NotiRepository
	nhrepo         nhRepo.NhRepository
	jwtSecret      string
	supa           configs.Supabase
	mail           configs.Mail
}

func NewUserUseCase(userrepo repositories.UserRepository, retirementrepo retirementRepo.RetirementRepository, assetrepo assetRepo.AssetRepository, notirepo notiRepo.NotiRepository, nhrepo nhRepo.NhRepository, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userrepo:       userrepo,
		retirementrepo: retirementrepo,
		assetrepo:      assetrepo,
		notirepo:       notirepo,
		nhrepo:         nhrepo,
		jwtSecret:      jwt.Secret,
		supa:           supa,
		mail:           mail,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string) (*entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	user.Email = normalizedEmail
	if _, err := u.userrepo.FindUserByEmail(user.Email); err == nil {
		return nil, errors.New("this email already have account")
	}

	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.ID = uuid.New().String()
	user.RoleID = role.ID
	user.Role = role
	user.Provider = "Credentials"
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
	normalizedEmail, err := utils.NormalizeEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email format")
	}

	email = normalizedEmail
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
	normalizedEmail, err := utils.NormalizeEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email format")
	}

	email = normalizedEmail
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email")
	}

	if user.Provider != "Credentials" {
		return "", nil, errors.New("this email is already registered with another authentication method")
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

func (u *UserUseCaseImpl) LoginWithGoogle(user *entities.User) (string, *entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return "", nil, errors.New("invalid email format")
	}

	user.Email = normalizedEmail
	account, err := u.userrepo.FindUserByEmail(user.Email)
	if err == nil {
		if account.Provider != "Google" {
			return "", nil, errors.New("this email is already registered with another authentication method")
		}
	} else {
		role, err := u.userrepo.GetRoleByName("User")
		if err != nil {
			return "", nil, errors.New("role not found")
		}

		user.ID = uuid.New().String()
		user.RoleID = role.ID
		user.Role = role
		user.Provider = "Google"

		if _, err := u.userrepo.CreateUser(user); err != nil {
			return "", nil, err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": account.ID,
		"role":    account.Role.RoleName,
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &account, nil
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
	house, err := u.userrepo.GetSelectedHouse(userID)
	if err != nil {
		return nil, err
	}

	currentMonth := int(time.Now().Month())
	if house.NursingHouseID == "00001" {
		house.MonthlyExpenses = 0
		house.LastCalculatedMonth = 0
		_, err = u.userrepo.UpdateSelectedHouse(house)
		if err != nil {
			return nil, err
		}
	}

	if house.LastCalculatedMonth != currentMonth {
		user, err := u.userrepo.GetUserByID(userID)
		if err != nil {
			return nil, err
		}

		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpense(user)
		if err != nil {
			return nil, err
		}

		requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * house.NursingHouse.Price
		if float64(requiredMoney) > user.House.CurrentMoney {
			house.Status = "In_Progress"
			house.MonthlyExpenses = monthlyExpenses
			house.LastCalculatedMonth = currentMonth
		} else {
			house.Status = "Completed"
			house.MonthlyExpenses = 0
			house.LastCalculatedMonth = 0
		}

		_, err = u.userrepo.UpdateSelectedHouse(house)
		if err != nil {
			return nil, err
		}
	}

	return house, nil
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user entities.User, file *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	existingUser.Email = normalizedEmail
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Username = user.Username
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

	currentMonth := int(time.Now().Month())
	selectedHouse.NursingHouseID = nursingHouseID
	if nursingHouseID == "00001" {
		selectedHouse.Status = "Completed"
		selectedHouse.LastCalculatedMonth = 0
		selectedHouse.MonthlyExpenses = 0
	} else if nursingHouseID != selectedHouse.NursingHouseID {
		user, err := u.userrepo.GetUserByID(userID)
		if err != nil {
			return nil, err
		}

		nursingHouse, err := u.nhrepo.GetNhByID(nursingHouseID)
		if err != nil {
			return nil, err
		}

		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpense(user)
		if err != nil {
			return nil, err
		}

		requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * nursingHouse.Price
		if float64(requiredMoney) > user.House.CurrentMoney {
			selectedHouse.Status = "In_Progress"
			selectedHouse.MonthlyExpenses = monthlyExpenses
			selectedHouse.LastCalculatedMonth = currentMonth
		} else {
			selectedHouse.Status = "Completed"
			selectedHouse.MonthlyExpenses = 0
			selectedHouse.LastCalculatedMonth = 0
		}
	} else {
		if selectedHouse.LastCalculatedMonth != currentMonth {
			user, err := u.userrepo.GetUserByID(userID)
			if err != nil {
				return nil, err
			}

			monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpense(user)
			if err != nil {
				return nil, err
			}

			requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
			if float64(requiredMoney) > user.House.CurrentMoney {
				selectedHouse.Status = "In_Progress"
				selectedHouse.MonthlyExpenses = monthlyExpenses
				selectedHouse.LastCalculatedMonth = currentMonth
			} else {
				selectedHouse.Status = "Completed"
				selectedHouse.MonthlyExpenses = 0
				selectedHouse.LastCalculatedMonth = 0
			}
		}
	}

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

func (u *UserUseCaseImpl) VerifyOTP(email, otpCode string) error {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return err
	}

	otp, err := u.userrepo.GetOTPByUserID(user.ID)
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP is expired")
	}

	if otp.OTP != otpCode {
		return errors.New("OTP is incorrect")
	}

	if err := u.userrepo.DeleteOTP(user.ID); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCaseImpl) ChangedPassword(email, newPassword string) error {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		return errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	_, err = u.userrepo.UpdateUserByID(&user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCaseImpl) CalculateRetirement(userID string) (fiber.Map, error) {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return fiber.Map{}, err
	}

	plan := user.RetirementPlan
	expectedMonthlyExpenses := plan.ExpectedMonthlyExpenses
	age, err := utils.CalculateRetirementPlanAge(plan.BirthDate, plan.CreatedAt)
	if err != nil {
		return nil, err
	}

	allCostAsset, err := utils.CalculateAllAssetsMonthlyExpenses(user)
	if err != nil {
		return fiber.Map{}, err
	}

	nursingHousePrice := 0.0
	cost := user.House.CurrentMoney
	if user.House.NursingHouse.ID != "" && user.House.Status != "Completed" {
		nursingHousePrice = float64(user.House.NursingHouse.Price)
	} else if user.House.Status == "Completed" {
		cost = float64((plan.ExpectLifespan - plan.RetirementAge) * 12 * user.House.NursingHouse.Price)
	}

	monthlyPlan := utils.MonthlyExpensesPlan{
		ExpectedMonthlyExpenses: expectedMonthlyExpenses,
		AnnualExpenseIncrease:   plan.AnnualExpenseIncrease,
		ExpectedInflation:       plan.ExpectedInflation,
		Age:                     age,
		RetirementAge:           plan.RetirementAge,
		ExpectLifespan:          plan.ExpectLifespan,
		AllCostAsset:            allCostAsset,
		NursingHousePrice:       nursingHousePrice,
	}

	requiredFunds := 0.0
	if plan.Status != "Completed" {
		requiredFunds, err = utils.CalculateMonthlySavings(monthlyPlan)
		if err != nil {
			return fiber.Map{}, err
		}
	}

	requiredAllFunds, err := utils.CalculateRetirementFunds(monthlyPlan)
	if err != nil {
		return fiber.Map{}, err
	}

	assetSavingsforAll := utils.CalculateAllAssetSavings(user, "All")
	assetSavingsforPlan := utils.CalculateAllAssetSavings(user, "Plan")
	currentMonthStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	deposits, err := u.userrepo.GetUserDepositsInRange(userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return fiber.Map{}, err
	}

	totalDeposits := 0.0
	for _, history := range deposits {
		totalDeposits += history.Money
	}

	moneyForPlan := plan.CurrentSavings + plan.CurrentTotalInvestment
	if moneyForPlan >= requiredAllFunds {
		moneyForPlan = requiredAllFunds
	}

	adjustedMonthlyExpenses := requiredFunds - totalDeposits // เดือนนี้เหลือเท่าไร
	totalNursingHouseCost := float64((plan.ExpectLifespan-plan.RetirementAge)*12*user.House.NursingHouse.Price) - user.House.CurrentMoney
	allRequiredFund := requiredAllFunds + allCostAsset + totalNursingHouseCost
	allSavingforPlan := moneyForPlan + plan.CurrentSavings + assetSavingsforPlan + cost
	allSavingforAll := plan.CurrentSavings + assetSavingsforAll + user.House.CurrentMoney // เงินออม
	allMoney := allSavingforAll + plan.CurrentTotalInvestment                             // เงินสุทธิ
	stillNeed := allRequiredFund - allSavingforPlan                                       // เงินที่ต้องเก็บอีก
	response := fiber.Map{
		"plan_name":                plan.PlanName,
		"allRequiredFund":          math.Round(allRequiredFund),
		"stillneed":                math.Round(stillNeed),
		"allretirementfund":        math.Round(requiredAllFunds),
		"monthly_expenses":         math.Round(adjustedMonthlyExpenses),
		"plan_saving":              math.Round(plan.CurrentSavings),
		"all_money":                math.Round(allMoney),
		"saving":                   math.Round(allSavingforAll),
		"investment":               math.Round(plan.CurrentTotalInvestment),
		"totalHouseCost":           math.Round(totalNursingHouseCost),
		"totalAssetCost":           math.Round(allCostAsset),
		"annual_savings_return":    plan.AnnualSavingsReturn,
		"annual_investment_return": plan.AnnualInvestmentReturn,
	}

	return response, nil
}

func (u *UserUseCaseImpl) CreateHistory(history entities.History) (*entities.History, error) {
	user, err := u.userrepo.GetUserByID(history.UserID)
	if err != nil {
		return nil, err
	}

	if history.Money <= 0 {
		return nil, errors.New("money must be greater than zero")
	}

	retirementData, err := u.CalculateRetirement(history.UserID)
	if err != nil {
		return nil, err
	}

	history.ID = uuid.New().String()
	history.TrackDate = time.Now()
	switch history.Method {
	case "deposit":
		if history.Type == "saving_money" {
			switch history.Category {
			case "spread":
				var validAssets []entities.Asset
				for _, asset := range user.Assets {
					if asset.Status == "In_Progress" {
						validAssets = append(validAssets, asset)
					}
				}

				var validHouse *entities.SelectedHouse
				if user.House.NursingHouseID != "00001" && user.House.Status != "Completed" {
					validHouse = &user.House
				}

				count := len(validAssets)
				if validHouse != nil {
					count++
				}

				count++
				amounts := history.Money / float64(count)
				for i := range validAssets {
					validAssets[i].CurrentMoney += amounts
					if validAssets[i].CurrentMoney >= validAssets[i].TotalCost {
						validAssets[i].Status = "Completed"
					}

					notification := &entities.Notification{
						ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), validAssets[i].ID),
						UserID:    user.ID,
						Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ : '%s' ได้เสร็จสิ้นแล้ว", validAssets[i].Name),
						CreatedAt: time.Now(),
					}

					_ = u.notirepo.CreateNotification(notification)
					if _, err := u.assetrepo.UpdateAssetByID(&validAssets[i]); err != nil {
						return nil, err
					}
				}

				if validHouse != nil {
					user.House.CurrentMoney += amounts
					requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
					if user.House.CurrentMoney >= float64(requiredMoney) {
						user.House.Status = "Completed"
						user.House.MonthlyExpenses = 0
						user.House.LastCalculatedMonth = 0
						notification := &entities.Notification{
							ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), validHouse.NursingHouse.Name),
							UserID:    user.ID,
							Message:   fmt.Sprintf("สุดยอดมาก บ้านพัก : '%s' ได้เสร็จสิ้นแล้ว", validHouse.NursingHouse.Name),
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
					}

				}

				user.RetirementPlan.CurrentSavings += amounts
				allRequiredFund := retirementData["allretirementfund"].(float64)
				allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
				if allMoney >= allRequiredFund {
					user.RetirementPlan.Status = "Completed"
				}

			case "retirementplan":
				user.RetirementPlan.CurrentSavings += history.Money
				allRequiredFund := retirementData["allretirementfund"].(float64)
				allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
				if allMoney >= allRequiredFund {
					user.RetirementPlan.Status = "Completed"
				}

			case "house":
				if user.House.NursingHouseID != "00001" || user.House.Status != "Completed" {
					user.House.CurrentMoney += history.Money
					requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
					if user.House.CurrentMoney >= float64(requiredMoney) {
						user.House.Status = "Completed"
						user.House.MonthlyExpenses = 0
						user.House.LastCalculatedMonth = 0
					}
				} else {
					return nil, errors.New("cannot update completed nursing house")
				}

			case "asset":
				asset, err := u.assetrepo.FindAssetByNameandUserID(history.Name, history.UserID)
				if err != nil {
					return nil, err
				}

				if asset.Status == "In_Progress" {
					asset.CurrentMoney += history.Money
					if asset.CurrentMoney >= asset.TotalCost {
						asset.Status = "Completed"
						asset.MonthlyExpenses = 0
						asset.LastCalculatedMonth = 0
					}

					_, err = u.assetrepo.UpdateAssetByID(asset)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New("cannot update completed or paused asset")
				}

			default:
				return nil, errors.New("invalid category for saving_money")
			}
		} else if history.Type == "investment" {
			user.RetirementPlan.CurrentTotalInvestment += history.Money
			allRequiredFund := retirementData["allretirementfund"].(float64)
			allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
			if allMoney >= allRequiredFund {
				user.RetirementPlan.Status = "Completed"
			}
		}

	case "withdraw":
		if history.Type == "saving_money" {
			switch history.Category {
			case "retirementplan":
				if user.RetirementPlan.CurrentSavings < history.Money {
					return nil, errors.New("insufficient funds in retirement savings")
				}

				user.RetirementPlan.CurrentSavings -= history.Money
				allRequiredFund := retirementData["allretirementfund"].(float64)
				allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
				if allMoney >= allRequiredFund {
					user.RetirementPlan.Status = "Completed"
				} else {
					user.RetirementPlan.Status = "In_Progress"
				}
			case "house":
				if user.House.NursingHouseID != "00001" || user.House.Status != "Completed" {
					if user.House.CurrentMoney < history.Money {
						return nil, errors.New("insufficient funds for house savings")
					}

					user.House.CurrentMoney -= history.Money
				} else {
					return nil, errors.New("cannot update completed nursing house")
				}

			case "asset":
				asset, err := u.assetrepo.FindAssetByNameandUserID(history.Name, history.UserID)
				if err != nil {
					return nil, err
				}

				if asset.Status != "Completed" {
					if asset.CurrentMoney < history.Money {
						return nil, errors.New("insufficient funds for asset savings")
					}

					asset.CurrentMoney -= history.Money
					_, err = u.assetrepo.UpdateAssetByID(asset)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New("cannot update completed asset")
				}

			default:
				return nil, errors.New("invalid category for saving_money")
			}
		} else if history.Type == "investment" {
			if user.RetirementPlan.CurrentTotalInvestment < history.Money {
				return nil, errors.New("insufficient investment funds")
			}

			user.RetirementPlan.CurrentTotalInvestment -= history.Money
			allRequiredFund := retirementData["allretirementfund"].(float64)
			allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
			if allMoney >= allRequiredFund {
				user.RetirementPlan.Status = "Completed"
			} else {
				user.RetirementPlan.Status = "In_Progress"
			}
		}

	default:
		return nil, errors.New("invalid method type")
	}

	_, err = u.userrepo.UpdateSelectedHouse(&user.House)
	if err != nil {
		return nil, err
	}

	_, err = u.retirementrepo.UpdateRetirementPlan(&user.RetirementPlan)
	if err != nil {
		return nil, err
	}

	createdHistory, err := u.userrepo.CreateHistory(&history)
	if err != nil {
		return nil, err
	}

	return createdHistory, nil
}

func (u *UserUseCaseImpl) GetHistoryByUserID(userID string) (fiber.Map, error) {
	data, err := u.userrepo.GetHistoryByUserID(userID)
	if err != nil {
		return fiber.Map{}, err
	}

	currentMonthStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	histories, err := u.userrepo.GetHistoryInRange(userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return fiber.Map{}, err
	}

	total := 0.0
	for _, history := range histories {
		if history.Method == "deposit" {
			total += history.Money
		} else if history.Method == "withdraw" {
			total -= history.Money
		}
	}

	response := fiber.Map{
		"data":  data,
		"total": total,
	}

	return response, nil
}

func (u *UserUseCaseImpl) GetHistoryByMonth(userID string) (map[string]float64, error) {
	historyByMonth, err := u.userrepo.GetUserHistoryByMonth(userID)
	if err != nil {
		return nil, err
	}

	return historyByMonth, nil
}
