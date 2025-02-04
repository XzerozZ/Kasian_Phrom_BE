package usecases

import (
	"errors"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/history/repositories"
	retirementRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	userRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HistoryUseCase interface {
	CreateHistory(history entities.History) (*entities.History, error)
	GetHistoryByUserID(userID string) (fiber.Map, error)
	GetHistoryByMonth(userID string) (map[string]float64, error)
}

type HistoryUseCaseImpl struct {
	historyrepo    repositories.HistoryRepository
	userrepo       userRepo.UserRepository
	retirementrepo retirementRepo.RetirementRepository
	db             *gorm.DB
}

func NewHistoryUseCase(historyrepo repositories.HistoryRepository, userrepo userRepo.UserRepository, retirementrepo retirementRepo.RetirementRepository, db *gorm.DB) *HistoryUseCaseImpl {
	return &HistoryUseCaseImpl{
		historyrepo:    historyrepo,
		userrepo:       userrepo,
		retirementrepo: retirementrepo,
		db:             db,
	}
}

func (u *HistoryUseCaseImpl) CreateHistory(history entities.History) (*entities.History, error) {
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
			err := utils.DistributeSavingMoney(history.Money, user.Assets, user.House, user.RetirementPlan, u.db)
			if err != nil {
				return nil, err
			}
		} else if history.Type == "investment" {
			user.RetirementPlan.CurrentTotalInvestment += history.Money
			_, err := u.retirementrepo.UpdateRetirementPlan(&user.RetirementPlan)
			if err != nil {
				return nil, err
			}
		}

	case "withdraw":
		if history.Type == "saving_money" {
			err := utils.WithdrawSavingMoney(history.Money, user.Assets, user.House, user.RetirementPlan, u.db)
			if err != nil {
				return nil, err
			}
		} else if history.Type == "investment" {
			if user.RetirementPlan.CurrentTotalInvestment < history.Money {
				return nil, errors.New("insufficient investment funds")
			}

			user.RetirementPlan.CurrentTotalInvestment -= history.Money
			_, err := u.retirementrepo.UpdateRetirementPlan(&user.RetirementPlan)
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errors.New("invalid method type")
	}

	createdHistory, err := u.historyrepo.CreateHistory(&history)
	if err != nil {
		return nil, err
	}

	return createdHistory, nil
}

func (u *HistoryUseCaseImpl) GetHistoryByUserID(userID string) (fiber.Map, error) {
	data, err := u.historyrepo.GetHistoryByUserID(userID)
	if err != nil {
		return fiber.Map{}, err
	}

	currentMonthStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	histories, err := u.historyrepo.GetHistoryInRange(userID, currentMonthStart, currentMonthEnd)
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

func (u *HistoryUseCaseImpl) GetHistoryByMonth(userID string) (map[string]float64, error) {
	historyByMonth, err := u.historyrepo.GetUserHistoryByMonth(userID)
	if err != nil {
		return nil, err
	}

	return historyByMonth, nil
}
