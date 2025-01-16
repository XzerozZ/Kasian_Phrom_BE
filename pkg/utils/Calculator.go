package utils

import (
	"math"
	"errors"
	"strconv"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

type MonthlyExpensesPlan struct {
	ExpectedMonthlyExpenses float64
	AnnualExpenseIncrease 	float64
	ExpectedInflation   	float64
	Age                  	int
	RetirementAge        	int
	ExpectLifespan       	int
	MonthsUntilRetirement 	int
	YearUntilLifeSpan		int
	AllCostAsset            float64
	NursingHousePrice   	float64
}

func CalculateRetirementFunds(plan MonthlyExpensesPlan) (float64, error) {
	if plan.ExpectedMonthlyExpenses < 0 || plan.AnnualExpenseIncrease < 0 || plan.ExpectedInflation < 0 {
		return 0, errors.New("ExpectedMonthlyExpenses, AnnualExpenseIncrease and ExpectedInflation must be greater than zero")
	}

	totalRequiredFunds := 0.0
	remainingMonthsInFirstYear := plan.MonthsUntilRetirement % 12
	if remainingMonthsInFirstYear == 0 {
		remainingMonthsInFirstYear = 12
	}

	totalRequiredFunds += plan.ExpectedMonthlyExpenses * float64(remainingMonthsInFirstYear)
	for year := 1; year <= (plan.MonthsUntilRetirement / 12) + 1; year++ {
		annualExpenses := plan.ExpectedMonthlyExpenses * 12
		increaseFactor := math.Pow(1 + (plan.AnnualExpenseIncrease/100) + (plan.ExpectedInflation/100), float64(year - 1))
		totalRequiredFunds += annualExpenses * increaseFactor
	}

	return totalRequiredFunds, nil
}

func CalculateMonthlySavings(plan MonthlyExpensesPlan) (float64, error) {
	requiredFunds, err := CalculateRetirementFunds(plan)
	if err != nil {
		return 0, err
	}
	
	monthlySavings := requiredFunds / float64(plan.MonthsUntilRetirement) + plan.AllCostAsset + plan.NursingHousePrice
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

    remainingMonths := (endYear - currentYear - 1) * 12 + (12 - int(asset.UpdatedAt.Month()) + 1)
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
	monthsUntilRetirement := (user.RetirementPlan.RetirementAge * 12) - user.RetirementPlan.AgeInMonths
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