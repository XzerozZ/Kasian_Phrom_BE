package usecases

import (
	"time"
	"errors"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
)

type RetirementUseCase interface {
	CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, error)
	GetRetirementByID(id string) (*entities.RetirementPlan, error)
}

type RetirementUseCaseImpl struct {
	retirerepo 		repositories.RetirementRepository
}

func NewRetirementUseCase(retirerepo repositories.RetirementRepository) *RetirementUseCaseImpl {
	return &RetirementUseCaseImpl{retirerepo:  retirerepo}
}

func (u *RetirementUseCaseImpl) CreateRetirement(retirement entities.RetirementPlan) (*entities.RetirementPlan, error) {
	id, err := u.retirerepo.GetRetirementNextID()
	if err != nil {
		return nil, err
	}
	
	layout := "02-01-2006"
	birthDate, err := time.Parse(layout, retirement.BirthDate)
	if err != nil {
		return nil, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	now := time.Now()
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	if now.Day() < birthDate.Day() {
		months--
	}
	
	if months < 0 {
		years--
		months += 12
	}

	ageInMonths := (years * 12) + months
	retirement.AgeInMonths = ageInMonths
	retirement.Age = years
	if retirement.CurrentSavings < 0 {
		return nil, errors.New("CurrentSavings must be greater than or equal to zero")
	}

	if retirement.MonthlyIncome < 0 {
		return nil, errors.New("MonthlyIncome must be greater than or equal to zero")
	}

	if retirement.MonthlyExpenses < 0 {
		return nil, errors.New("MonthlyExpenses must be greater than or equal to zero")
	}

	if retirement.CurrentSavingsReturns <= 0 {
		return nil, errors.New("MonthlyIncome must be greater than zero")
	}

	if retirement.CurrentTotalInvestment < 0 {
		return nil, errors.New("CurrentTotalInvestment must be greater than zero or equal to zero")
	}

	if retirement.InvestmentReturn <= 0 {
		return nil, errors.New("InvestmentReturn must be greater than zero")
	}

	if retirement.ExpectedInflation <= 0 {
		return nil, errors.New("ExpectedInflation must be greater than zero")
	}

	if retirement.AnnualExpenseIncrease < 0 {
		return nil, errors.New("AnnualExpenseIncrease must be greater than or equal to zero")
	}

	if retirement.AnnualSavingsReturn < 0 {
		return nil, errors.New("AnnualSavingsReturn must be greater than or equal to zero")
	}

	if retirement.AnnualInvestmentReturn < 0 {
		return nil, errors.New("AnnualInvestmentReturn must be greater than or equal to zero")
	}

	if years >= retirement.RetirementAge {
		return nil, errors.New("Age must be less than RetirementAge")
	}

	if retirement.RetirementAge >= retirement.ExpectLifespan {
		return nil, errors.New("RetirementAge must be less than ExpectLifespan")
	}

	retirement.ID = id
	createdRetire, err := u.retirerepo.CreateRetirement(&retirement)
    if err != nil { 
        return nil, err
    }
    
    return createdRetire, nil
}

func (u *RetirementUseCaseImpl) GetRetirementByID(id string) (*entities.RetirementPlan, error) {
	return u.retirerepo.GetRetirementByID(id)
}

