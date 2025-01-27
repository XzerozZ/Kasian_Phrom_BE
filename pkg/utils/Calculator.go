package utils

import (
	"errors"
	"math"
	"strconv"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type MonthlyExpensesPlan struct {
	ExpectedMonthlyExpenses float64
	AnnualExpenseIncrease   float64
	ExpectedInflation       float64
	Age                     int
	RetirementAge           int
	ExpectLifespan          int
	YearsUntilRetirement    int
	AllCostAsset            float64
	NursingHousePrice       float64
}

func CalculateRetirementFunds(plan MonthlyExpensesPlan) (float64, error) {
	if plan.ExpectedMonthlyExpenses <= 0 || plan.AnnualExpenseIncrease < 0 || plan.ExpectedInflation < 0 {
		return 0, errors.New("expected monthly expenses, annual expense increase, and inflation must be greater than zero")
	}

	yearsUntilRetirement := plan.RetirementAge - plan.Age
	yearsInRetirement := plan.ExpectLifespan - plan.RetirementAge
	if yearsUntilRetirement <= 0 {
		return 0, errors.New("retirement age must be greater than current age")
	}

	if yearsInRetirement <= 0 {
		return 0, errors.New("expected lifespan must be greater than retirement age")
	}

	totalRequiredFunds := 0.0
	annualExpenses := plan.ExpectedMonthlyExpenses * 12
	firstYearFactor := math.Pow(1+(plan.ExpectedInflation/100), float64(1+yearsUntilRetirement))
	totalRequiredFunds += annualExpenses * firstYearFactor
	for year := 2; year <= yearsInRetirement; year++ {
		remainingYears := yearsUntilRetirement + year
		compoundingFactor := math.Pow(1+(plan.ExpectedInflation/100), float64(remainingYears))
		totalRequiredFunds += annualExpenses * compoundingFactor
	}

	return totalRequiredFunds, nil
}

func CalculateMonthlySavings(plan MonthlyExpensesPlan) (float64, error) {
	requiredFunds, err := CalculateRetirementFunds(plan)
	if err != nil {
		return 0, err
	}

	monthsUntilRetirement := plan.YearsUntilRetirement * 12
	if monthsUntilRetirement <= 0 {
		return 0, errors.New("years until retirement must be greater than zero")
	}

	monthlySavings := requiredFunds/float64(monthsUntilRetirement) + plan.AllCostAsset + plan.NursingHousePrice
	return monthlySavings, nil
}

func CalculateMonthlyExpenses(asset *entities.Asset) (float64, error) {
	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return 0, err
	}

	currentYear := asset.UpdatedAt.Year()
	if endYear < currentYear {
		return 0, errors.New("end year must be greater than or equal to current year")
	}

	remainingMonths := (endYear-currentYear-1)*12 + (12 - int(asset.UpdatedAt.Month()) + 1)
	remainingCost := asset.TotalCost - asset.CurrentMoney
	if remainingCost < 0 {
		return 0, errors.New("current money cannot exceed total cost")
	}

	return remainingCost / float64(remainingMonths), nil
}

func CalculateAllAssetsMonthlyExpenses(user *entities.User) (float64, error) {
	var total float64
	for _, asset := range user.Assets {
		monthlyExpense, err := CalculateMonthlyExpenses(&asset)
		if err != nil {
			return 0, err
		}

		total += monthlyExpense
	}

	return total, nil
}

func CalculateNursingHouseMonthlyExpenses(user *entities.User) (float64, error) {
	monthsUntilRetirement := user.RetirementPlan.RetirementAge * 12
	yearUntilLifespan := user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge
	totalNursingHouseCost := user.House.NursingHouse.Price * yearUntilLifespan
	return float64(totalNursingHouseCost) / float64(monthsUntilRetirement), nil
}

func CalculateAllAssetSavings(user *entities.User) (float64, error) {
	var total float64
	for _, asset := range user.Assets {
		total += asset.CurrentMoney
	}

	return total, nil
}

func CalculateAllLoan(loans []entities.Loan) (float64, error) {
	var total float64
	for _, loan := range loans {
		totalLoan := loan.MonthlyExpenses * float64(loan.RemainingMonths)
		total += totalLoan
	}

	return total, nil
}
