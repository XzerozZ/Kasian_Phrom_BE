package usecases

import (
	"errors"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/google/uuid"
)

type RetirementUseCase interface {
	CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, int, error)
	GetRetirementByID(id string) (*entities.RetirementPlan, error)
	GetRetirementByUserID(userID string) (*entities.RetirementPlan, error)
	UpdateRetirementByID(userID string, retirement entities.RetirementPlan) (*entities.RetirementPlan, error)
}

type RetirementUseCaseImpl struct {
	retirerepo repositories.RetirementRepository
}

func NewRetirementUseCase(retirerepo repositories.RetirementRepository) *RetirementUseCaseImpl {
	return &RetirementUseCaseImpl{retirerepo: retirerepo}
}

func (u *RetirementUseCaseImpl) CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, int, error) {
	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	age, err := utils.CalculateAge(retirement.BirthDate)
	if err != nil {
		return nil, 0, err
	}

	planAge, err := utils.CalculateRetirementPlanAge(retirement.BirthDate, time.Now())
	if err != nil {
		return nil, 0, err
	}

	if retirement.CurrentSavings < 0 {
		return nil, 0, errors.New("currentSavings must be greater than or equal to zero")
	}

	if retirement.MonthlyIncome < 0 {
		return nil, 0, errors.New("monthlyIncome must be greater than or equal to zero")
	}

	if retirement.MonthlyExpenses < 0 {
		return nil, 0, errors.New("monthlyExpenses must be greater than or equal to zero")
	}

	if retirement.CurrentSavingsReturns <= 0 {
		return nil, 0, errors.New("monthlyIncome must be greater than zero")
	}

	if retirement.CurrentTotalInvestment < 0 {
		return nil, 0, errors.New("currentTotalInvestment must be greater than zero or equal to zero")
	}

	if retirement.InvestmentReturn <= 0 {
		return nil, 0, errors.New("investmentReturn must be greater than zero")
	}

	if retirement.ExpectedInflation <= 0 {
		return nil, 0, errors.New("expectedInflation must be greater than zero")
	}

	if retirement.ExpectedMonthlyExpenses <= 0 {
		return nil, 0, errors.New("expectedMonthlyExpenses must be greater than zero")
	}

	if retirement.AnnualExpenseIncrease < 0 {
		return nil, 0, errors.New("annualExpenseIncrease must be greater than or equal to zero")
	}

	if retirement.AnnualSavingsReturn < 0 {
		return nil, 0, errors.New("annualSavingsReturn must be greater than or equal to zero")
	}

	if retirement.AnnualInvestmentReturn < 0 {
		return nil, 0, errors.New("annualInvestmentReturn must be greater than or equal to zero")
	}

	if age >= retirement.RetirementAge {
		return nil, 0, errors.New("age must be less than RetirementAge")
	}

	if retirement.RetirementAge >= retirement.ExpectLifespan {
		return nil, 0, errors.New("retirementAge must be less than ExpectLifespan")
	}

	requiredFunds, err := utils.CalculateRetirementFunds(&retirement, planAge)
	if err != nil {
		return nil, 0, err
	}

	monthlySavings, err := utils.CalculateMonthlySavings(&retirement, planAge, currentYear, currentMonth)
	if err != nil {
		return nil, 0, err
	}

	retirement.ID = uuid.New().String()
	retirement.Status = "In_Progress"
	retirement.LastCalculatedMonth = currentMonth
	retirement.LastRequiredFunds = requiredFunds
	retirement.LastMonthlyExpenses = monthlySavings
	createdRetire, err := u.retirerepo.CreateRetirement(&retirement)
	if err != nil {
		return nil, 0, err
	}

	return createdRetire, age, nil
}

func (u *RetirementUseCaseImpl) GetRetirementByID(id string) (*entities.RetirementPlan, error) {
	return u.retirerepo.GetRetirementByID(id)
}

func (u *RetirementUseCaseImpl) GetRetirementByUserID(userID string) (*entities.RetirementPlan, error) {
	return u.retirerepo.GetRetirementByUserID(userID)
}

func (u *RetirementUseCaseImpl) UpdateRetirementByID(userID string, retirement entities.RetirementPlan) (*entities.RetirementPlan, error) {
	existingRetirement, err := u.retirerepo.GetRetirementByUserID(userID)
	if err != nil {
		return nil, err
	}

	age, err := utils.CalculateRetirementPlanAge(retirement.BirthDate, time.Now())
	if err != nil {
		return nil, err
	}

	if retirement.MonthlyIncome < 0 {
		return nil, errors.New("monthlyIncome must be greater than or equal to zero")
	}

	if retirement.MonthlyExpenses < 0 {
		return nil, errors.New("monthlyExpenses must be greater than or equal to zero")
	}

	if retirement.CurrentSavingsReturns <= 0 {
		return nil, errors.New("monthlyIncome must be greater than zero")
	}

	if retirement.InvestmentReturn <= 0 {
		return nil, errors.New("investmentReturn must be greater than zero")
	}

	if retirement.ExpectedInflation <= 0 {
		return nil, errors.New("expectedInflation must be greater than zero")
	}

	if retirement.AnnualExpenseIncrease < 0 {
		return nil, errors.New("annualExpenseIncrease must be greater than or equal to zero")
	}

	if retirement.AnnualSavingsReturn < 0 {
		return nil, errors.New("annualSavingsReturn must be greater than or equal to zero")
	}

	if retirement.AnnualInvestmentReturn < 0 {
		return nil, errors.New("annualInvestmentReturn must be greater than or equal to zero")
	}

	if retirement.ExpectedMonthlyExpenses <= 0 {
		return nil, errors.New("expectedMonthlyExpenses must be greater than zero")
	}

	if age >= retirement.RetirementAge {
		return nil, errors.New("age must be less than RetirementAge")
	}

	if retirement.RetirementAge >= retirement.ExpectLifespan {
		return nil, errors.New("retirementAge must be less than ExpectLifespan")
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	needsRecalculation := false
	recalculateFunds := false
	if existingRetirement.ExpectLifespan != retirement.ExpectLifespan || existingRetirement.RetirementAge != retirement.RetirementAge || existingRetirement.ExpectedMonthlyExpenses != retirement.ExpectedMonthlyExpenses || existingRetirement.ExpectedInflation != retirement.ExpectedInflation {
		recalculateFunds = true
		needsRecalculation = true
	}

	if currentMonth != existingRetirement.LastCalculatedMonth {
		needsRecalculation = true
	}

	existingRetirement.BirthDate = retirement.BirthDate
	existingRetirement.ExpectLifespan = retirement.ExpectLifespan
	existingRetirement.RetirementAge = retirement.RetirementAge
	existingRetirement.PlanName = retirement.PlanName
	existingRetirement.MonthlyIncome = retirement.MonthlyIncome
	existingRetirement.ExpectedMonthlyExpenses = retirement.ExpectedMonthlyExpenses
	existingRetirement.MonthlyExpenses = retirement.MonthlyExpenses
	existingRetirement.CurrentSavingsReturns = retirement.CurrentSavingsReturns
	existingRetirement.InvestmentReturn = retirement.InvestmentReturn
	existingRetirement.ExpectedInflation = retirement.ExpectedInflation
	existingRetirement.AnnualExpenseIncrease = retirement.AnnualExpenseIncrease
	existingRetirement.AnnualSavingsReturn = retirement.AnnualSavingsReturn
	existingRetirement.AnnualInvestmentReturn = retirement.AnnualInvestmentReturn
	existingRetirement.CurrentTotalInvestment = retirement.CurrentTotalInvestment
	if needsRecalculation {
		if recalculateFunds {
			requiredFunds, err := utils.CalculateRetirementFunds(existingRetirement, age)
			if err != nil {
				return nil, err
			}
			existingRetirement.LastRequiredFunds = requiredFunds
		}

		monthlySavings, err := utils.CalculateMonthlySavings(existingRetirement, age, currentYear, currentMonth)
		if err != nil {
			return nil, err
		}

		existingRetirement.LastCalculatedMonth = currentMonth
		existingRetirement.LastMonthlyExpenses = monthlySavings
		currentTotalMoney := existingRetirement.CurrentSavings + existingRetirement.CurrentTotalInvestment
		if currentTotalMoney >= existingRetirement.LastRequiredFunds {
			existingRetirement.Status = "Completed"
			existingRetirement.LastMonthlyExpenses = 0
		} else {
			existingRetirement.Status = "In_Progress"
		}
	}

	return u.retirerepo.UpdateRetirementPlan(existingRetirement)
}
