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
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/socket"
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
	UpdateSelectedHouse(userID, nursingHouseID string, transfers []entities.TransferRequest) (*entities.SelectedHouse, error)

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

const (
	defaultHouseID   = "00001"
	dateLayout       = "02-01-2006"
	statusCompleted  = "Completed"
	statusInProgress = "In_Progress"
)

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
	if house.NursingHouseID == defaultHouseID {
		house.MonthlyExpenses = 0
		house.LastCalculatedMonth = 0
		return u.userrepo.UpdateSelectedHouse(house)
	}

	if house.LastCalculatedMonth != currentMonth {
		user, err := u.userrepo.GetUserByID(userID)
		if err != nil {
			return nil, err
		}

		currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpense(user, float64(house.NursingHouse.Price), int(currentYear), currentMonth)
		if err != nil {
			return nil, err
		}

		house.MonthlyExpenses = monthlyExpenses
		house.LastCalculatedMonth = currentMonth
		return u.userrepo.UpdateSelectedHouse(house)
	}

	return house, nil
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user entities.User, file *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

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

func (u *UserUseCaseImpl) UpdateSelectedHouse(userID, nursingHouseID string, transfers []entities.TransferRequest) (*entities.SelectedHouse, error) {
	selectedHouse, err := u.userrepo.GetSelectedHouse(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if nursingHouseID == defaultHouseID {
		totalTransfer := 0.0
		for _, transfer := range transfers {
			totalTransfer += transfer.Amount
		}

		if totalTransfer > selectedHouse.CurrentMoney {
			return nil, errors.New("transfer amount exceeds House's current money")
		}

		for _, transfer := range transfers {
			switch transfer.Type {
			case "asset":
				selectedItem, err := u.assetrepo.FindAssetByNameandUserID(transfer.Name, userID)
				if err != nil {
					return nil, err
				}

				if selectedItem.Status == "In_Progress" {
					selectedItem.CurrentMoney += transfer.Amount
					his := entities.History{
						ID:           uuid.New().String(),
						Method:       "deposit",
						Type:         "saving_money",
						Category:     "asset",
						Name:         selectedItem.Name,
						Money:        transfer.Amount,
						TransferFrom: selectedHouse.NursingHouse.Name,
						UserID:       userID,
						TrackDate:    time.Now(),
					}

					_, err = u.userrepo.CreateHistory(&his)
					if err != nil {
						return nil, err
					}

					if selectedItem.CurrentMoney >= selectedItem.TotalCost {
						selectedItem.Status = "Completed"
						selectedItem.MonthlyExpenses = 0
						selectedItem.LastCalculatedMonth = 0
						notification := &entities.Notification{
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ %s ได้เสร็จสิ้นแล้ว", selectedItem.Name),
							Type:      "asset",
							ObjectID:  selectedItem.ID,
							Balance:   selectedItem.CurrentMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(userID, *notification)
					}

					_, err = u.assetrepo.UpdateAssetByID(selectedItem)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New("cannot update completed or paused asset")
				}

			case "retirementplan":
				retirement, err := u.retirementrepo.GetRetirementByUserID(userID)
				if err != nil {
					return nil, err
				}

				retirement.CurrentSavings += transfer.Amount
				his := entities.History{
					ID:           uuid.New().String(),
					Method:       "deposit",
					Type:         "saving_money",
					Category:     "retirementplan",
					Name:         retirement.PlanName,
					Money:        transfer.Amount,
					TransferFrom: selectedHouse.NursingHouse.Name,
					UserID:       userID,
					TrackDate:    time.Now(),
				}

				_, err = u.userrepo.CreateHistory(&his)
				if err != nil {
					return nil, err
				}

				allMoney := retirement.CurrentSavings + retirement.CurrentTotalInvestment
				if allMoney >= retirement.LastRequiredFunds {
					retirement.Status = "Completed"
					retirement.LastMonthlyExpenses = 0
					retirement.LastMonthlyExpenses = 0
					notification := &entities.Notification{
						ID:        uuid.New().String(),
						UserID:    user.ID,
						Message:   fmt.Sprintf("สุดยอดมาก แผนเกษียณ %s ของคุณได้ถึงเป้าแล้ว", user.RetirementPlan.PlanName),
						Type:      "retirementplan",
						ObjectID:  user.House.NursingHouseID,
						Balance:   allMoney,
						CreatedAt: time.Now(),
					}

					_ = u.notirepo.CreateNotification(notification)
					socket.SendNotificationToUser(userID, *notification)
				}
				_, err = u.retirementrepo.UpdateRetirementPlan(retirement)
				if err != nil {
					return nil, err
				}
			default:
				continue
			}
		}

		his := entities.History{
			ID:        uuid.New().String(),
			Method:    "withdraw",
			Type:      "saving_money",
			Category:  "asset",
			Name:      selectedHouse.NursingHouse.Name,
			Money:     selectedHouse.CurrentMoney - totalTransfer,
			UserID:    userID,
			TrackDate: time.Now(),
		}

		_, err = u.userrepo.CreateHistory(&his)
		if err != nil {
			return nil, err
		}

		selectedHouse.NursingHouseID = nursingHouseID
		selectedHouse.Status = statusCompleted
		selectedHouse.LastCalculatedMonth = 0
		selectedHouse.CurrentMoney = 0
		selectedHouse.MonthlyExpenses = 0
	}

	if nursingHouseID != selectedHouse.NursingHouseID || selectedHouse.LastCalculatedMonth != int(time.Now().Month()) {
		nursingHouse, err := u.nhrepo.GetNhByID(nursingHouseID)
		if err != nil {
			return nil, err
		}

		currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
		selectedHouse.Status = "In_Progress"
		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpense(user, float64(nursingHouse.Price), int(currentYear), currentMonth)
		if err != nil {
			return nil, err
		}

		selectedHouse.MonthlyExpenses = monthlyExpenses
		selectedHouse.NursingHouseID = nursingHouseID
		selectedHouse.LastCalculatedMonth = currentMonth
		requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * nursingHouse.Price
		if float64(requiredMoney) < user.House.CurrentMoney {
			selectedHouse.Status = statusCompleted
		}
	}

	return u.userrepo.UpdateSelectedHouse(selectedHouse)
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
	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	age, err := utils.CalculateRetirementPlanAge(plan.BirthDate, plan.CreatedAt)
	if err != nil {
		return nil, err
	}

	var allAssetsExpense, allTotalCost float64
	assetSavingsforAll := utils.CalculateAllAssetSavings(user, "All")
	assetSavingsforPlan := utils.CalculateAllAssetSavings(user, "Plan")
	for _, asset := range user.Assets {
		allTotalCost += asset.TotalCost
		if asset.LastCalculatedMonth == currentMonth {
			allAssetsExpense += asset.MonthlyExpenses
		} else if asset.LastCalculatedMonth != currentMonth {
			asset.MonthlyExpenses = utils.CalculateMonthlyExpenses(&asset, currentYear, currentMonth)
			asset.LastCalculatedMonth = currentMonth
			allAssetsExpense += asset.MonthlyExpenses
			_, err := u.assetrepo.UpdateAssetByID(&asset)
			if err != nil {
				return nil, err
			}
		}
	}

	var nursingHousePrice float64
	cost := user.House.CurrentMoney
	if user.House.Status == "Completed" {
		cost = float64((plan.ExpectLifespan - plan.RetirementAge) * 12 * user.House.NursingHouse.Price)
	}

	if user.House.LastCalculatedMonth == currentMonth {
		nursingHousePrice = user.House.MonthlyExpenses
	} else if user.House.LastCalculatedMonth != currentMonth {
		nursingHousePrice, err = utils.CalculateNursingHouseMonthlyExpense(user, float64(user.House.NursingHouse.Price), currentYear, currentMonth)
		if err != nil {
			return nil, err
		}

		user.House.MonthlyExpenses = nursingHousePrice
		user.House.LastCalculatedMonth = currentMonth
		_, err = u.userrepo.UpdateSelectedHouse(&user.House)
		if err != nil {
			return nil, err
		}
	}

	var planExpense float64
	if plan.LastCalculatedMonth == currentMonth {
		planExpense = plan.LastMonthlyExpenses
	} else if plan.LastCalculatedMonth != currentMonth {
		planExpense, err := utils.CalculateMonthlySavings(&plan, age, currentYear, currentMonth)
		if err != nil {
			return nil, err
		}

		plan.MonthlyExpenses = planExpense
		plan.LastCalculatedMonth = currentMonth
		_, err = u.retirementrepo.UpdateRetirementPlan(&plan)
		if err != nil {
			return nil, err
		}
	}

	moneyForPlan := plan.CurrentSavings + plan.CurrentTotalInvestment
	if moneyForPlan >= plan.LastRequiredFunds {
		moneyForPlan = plan.LastRequiredFunds
	}

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

	totalNursingHouseCost := float64((plan.ExpectLifespan - plan.RetirementAge) * 12 * user.House.NursingHouse.Price)
	allRequiredFund := plan.LastRequiredFunds + totalNursingHouseCost + allTotalCost
	adjustedMonthlyExpenses := (planExpense + nursingHousePrice + allAssetsExpense) - totalDeposits
	savingforPlan := moneyForPlan + assetSavingsforPlan + cost
	savingforAll := assetSavingsforAll + user.House.CurrentMoney + plan.CurrentSavings
	allMoney := savingforAll + plan.CurrentTotalInvestment
	stillNeed := allRequiredFund - savingforPlan
	response := fiber.Map{
		"plan_name":                plan.PlanName,
		"allRequiredFund":          math.Round(allRequiredFund),             //จำนวนเงินที่ต้องการทั้งหมด
		"stillneed":                math.Round(stillNeed),                   //ขาดอีก
		"allretirementfund":        math.Round(plan.LastRequiredFunds),      //เงินเกษียณที่ต้องการทั้งหมด
		"monthly_expenses":         math.Round(adjustedMonthlyExpenses),     //เงินที่ต้องผ่อนเดือนนี้ทั้งหมด - เงินที่เก็บเดือนนี้ทั้งหมด
		"plan_saving":              math.Round(plan.CurrentSavings),         //เงินออมของแผน
		"all_money":                math.Round(allMoney),                    //เงินสุทธิ
		"saving":                   math.Round(savingforAll),                //เงินออมทั้งหมด
		"investment":               math.Round(plan.CurrentTotalInvestment), //เงินลงทุน
		"all_assets_expense":       math.Round(allAssetsExpense),            //ราคาบ้านพักต่อเดือน
		"nursingHouse_expense":     math.Round(nursingHousePrice),           //ราคาของทรัพย์สินที่ต้องผ่อนต่อเดือนทั้งหมด
		"plan_expense":             math.Round(planExpense),
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

				var validPlan *entities.RetirementPlan
				if user.RetirementPlan.Status != "Completed" {
					validPlan = &user.RetirementPlan
				}

				if validPlan != nil {
					count++
				}

				amounts := history.Money / float64(count)
				for i := range validAssets {
					validAssets[i].CurrentMoney += amounts
					if validAssets[i].CurrentMoney >= validAssets[i].TotalCost {
						validAssets[i].Status = "Completed"
						validAssets[i].MonthlyExpenses = 0
						validAssets[i].LastCalculatedMonth = 0
						notification := &entities.Notification{
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ %s ได้เสร็จสิ้นแล้ว", validAssets[i].Name),
							Type:      "asset",
							ObjectID:  validAssets[i].ID,
							Balance:   validAssets[i].CurrentMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(user.ID, *notification)
					}

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
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   "สุดยอดมาก บ้านพักคนชรา ได้เสร็จสิ้นแล้ว",
							Type:      "house",
							ObjectID:  user.House.NursingHouseID,
							Balance:   validHouse.CurrentMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(user.ID, *notification)
					}

				}

				if validPlan != nil {
					user.RetirementPlan.CurrentSavings += amounts
					allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
					if allMoney >= user.RetirementPlan.LastRequiredFunds {
						user.RetirementPlan.Status = "Completed"
						user.RetirementPlan.LastMonthlyExpenses = 0
						user.RetirementPlan.LastMonthlyExpenses = 0
						notification := &entities.Notification{
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   fmt.Sprintf("สุดยอดมาก แผนเกษียณ %s ของคุณได้ถึงเป้าแล้ว", validHouse.NursingHouse.Name),
							Type:      "retirementplan",
							ObjectID:  user.RetirementPlan.ID,
							Balance:   allMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(user.ID, *notification)
					}
				}
			case "retirementplan":
				user.RetirementPlan.CurrentSavings += history.Money
				allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
				if allMoney >= user.RetirementPlan.LastRequiredFunds {
					user.RetirementPlan.Status = "Completed"
					user.RetirementPlan.LastMonthlyExpenses = 0
					user.RetirementPlan.LastMonthlyExpenses = 0
					notification := &entities.Notification{
						ID:        uuid.New().String(),
						UserID:    user.ID,
						Message:   fmt.Sprintf("สุดยอดมาก แผนเกษียณ %s ของคุณได้ถึงเป้าแล้ว", user.RetirementPlan.PlanName),
						Type:      "retirementplan",
						ObjectID:  user.RetirementPlan.ID,
						Balance:   allMoney,
						CreatedAt: time.Now(),
					}

					_ = u.notirepo.CreateNotification(notification)
					socket.SendNotificationToUser(user.ID, *notification)
				}

			case "house":
				if user.House.NursingHouseID != "00001" || user.House.Status != "Completed" {
					user.House.CurrentMoney += history.Money
					requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
					if user.House.CurrentMoney >= float64(requiredMoney) {
						user.House.Status = "Completed"
						user.House.MonthlyExpenses = 0
						user.House.LastCalculatedMonth = 0
						notification := &entities.Notification{
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   "สุดยอดมาก บ้านพักคนชรา ได้เสร็จสิ้นแล้ว",
							Type:      "house",
							ObjectID:  user.House.NursingHouseID,
							Balance:   user.House.CurrentMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(user.ID, *notification)
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
						notification := &entities.Notification{
							ID:        uuid.New().String(),
							UserID:    user.ID,
							Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ %s ได้เสร็จสิ้นแล้ว", asset.Name),
							Type:      "asset",
							ObjectID:  asset.ID,
							Balance:   asset.CurrentMoney,
							CreatedAt: time.Now(),
						}

						_ = u.notirepo.CreateNotification(notification)
						socket.SendNotificationToUser(user.ID, *notification)
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
			allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
			if allMoney >= user.RetirementPlan.LastRequiredFunds {
				user.RetirementPlan.Status = "Completed"
				user.RetirementPlan.LastMonthlyExpenses = 0
				user.RetirementPlan.LastMonthlyExpenses = 0
				notification := &entities.Notification{
					ID:        uuid.New().String(),
					UserID:    user.ID,
					Message:   fmt.Sprintf("สุดยอดมาก แผนเกษียณ %s ของคุณได้ถึงเป้าแล้ว", user.RetirementPlan.PlanName),
					Type:      "retirementplan",
					ObjectID:  user.RetirementPlan.ID,
					Balance:   allMoney,
					CreatedAt: time.Now(),
				}

				_ = u.notirepo.CreateNotification(notification)
				socket.SendNotificationToUser(user.ID, *notification)
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
				allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
				if allMoney >= user.RetirementPlan.LastRequiredFunds {
					user.RetirementPlan.Status = "Completed"
					user.RetirementPlan.LastMonthlyExpenses = 0
					user.RetirementPlan.LastMonthlyExpenses = 0
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
			allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
			if allMoney >= user.RetirementPlan.LastRequiredFunds {
				user.RetirementPlan.Status = "Completed"
				user.RetirementPlan.LastMonthlyExpenses = 0
				user.RetirementPlan.LastMonthlyExpenses = 0
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
