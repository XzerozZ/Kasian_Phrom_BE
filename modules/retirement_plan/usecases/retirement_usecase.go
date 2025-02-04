package usecases

import (
	"errors"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/google/uuid"
)

type RetirementUseCase interface {
	CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, int, error)
	GetRetirementByID(id string) (*entities.RetirementPlan, error)
	GetRetirementByUserID(userID string) (*entities.RetirementPlan, error)
	UpdateRetirementByID(id string, retirement entities.RetirementPlan) (*entities.RetirementPlan, error)
}

type RetirementUseCaseImpl struct {
	retirerepo repositories.RetirementRepository
}

func NewRetirementUseCase(retirerepo repositories.RetirementRepository) *RetirementUseCaseImpl {
	return &RetirementUseCaseImpl{retirerepo: retirerepo}
}

func (u *RetirementUseCaseImpl) CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, int, error) {
	age, err := utils.CalculateAge(retirement.BirthDate)
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

	retirement.ID = uuid.New().String()
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

func (u *RetirementUseCaseImpl) UpdateRetirementByID(id string, retirement entities.RetirementPlan) (*entities.RetirementPlan, error) {
	existingRetirement, err := u.retirerepo.GetRetirementByUserID(id)
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

	existingRetirement.PlanName = retirement.PlanName
	existingRetirement.CurrentSavings = retirement.ExpectedMonthlyExpenses
	existingRetirement.MonthlyExpenses = retirement.MonthlyExpenses
	existingRetirement.CurrentSavingsReturns = retirement.CurrentSavingsReturns
	existingRetirement.InvestmentReturn = retirement.InvestmentReturn
	existingRetirement.ExpectedInflation = retirement.ExpectedInflation
	existingRetirement.AnnualExpenseIncrease = retirement.AnnualExpenseIncrease
	existingRetirement.AnnualSavingsReturn = retirement.AnnualSavingsReturn
	existingRetirement.AnnualInvestmentReturn = retirement.AnnualInvestmentReturn
	updatedRetirement, err := u.retirerepo.UpdateRetirementPlan(existingRetirement)
	if err != nil {
		return nil, err
	}

	return updatedRetirement, nil
}
