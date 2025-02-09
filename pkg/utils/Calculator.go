package utils

import (
	"errors"
	"math"
	"strconv"
	"time"

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

func CalculateRetirementPlanAge(birthDateStr string, planCreationDate time.Time) (int, error) {
	layout := "02-01-2006"
	birthDate, err := time.Parse(layout, birthDateStr)
	if err != nil {
		return 0, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	years := planCreationDate.Year() - birthDate.Year()
	if planCreationDate.YearDay() < birthDate.YearDay() {
		years--
	}
	return years, nil
}

func CalculateAge(birthDateStr string) (int, error) {
	layout := "02-01-2006"
	birthDate, err := time.Parse(layout, birthDateStr)
	if err != nil {
		return 0, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	now := time.Now()
	years := now.Year() - birthDate.Year()

	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	return years, nil
}

func CalculateRetirementFunds(plan MonthlyExpensesPlan) (float64, error) {
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
	for year := 1; year <= yearsInRetirement; year++ {
		remainingYears := yearsUntilRetirement + year
		compoundingFactor := math.Pow(1+(plan.ExpectedInflation/100), float64(remainingYears))
		totalRequiredFunds += annualExpenses * compoundingFactor
	}

	totalRequiredFunds = math.Round(totalRequiredFunds)
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

	funds := requiredFunds / float64(monthsUntilRetirement)
	funds = math.Round(funds)
	monthlySavings := funds + plan.AllCostAsset + plan.NursingHousePrice
	return monthlySavings, nil
}

func CalculateMonthlyExpenses(asset *entities.Asset) (float64, error) {
	if asset.Status == "Completed" {
		return 0, nil
	}

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

	monthlyExpenses := remainingCost / float64(remainingMonths)
	monthlyExpenses = math.Round(monthlyExpenses)
	return monthlyExpenses, nil
}

func CalculateAllAssetsMonthlyExpenses(user *entities.User) (float64, error) {
	var total float64
	for _, asset := range user.Assets {
		if asset.Status == "In_Progress" {
			monthlyExpense, err := CalculateMonthlyExpenses(&asset)
			if err != nil {
				return 0, err
			}

			total += monthlyExpense
		}
	}

	total = math.Round(total)
	return total, nil
}

func CalculateNursingHouseMonthlyExpenses(user *entities.User) (float64, error) {
	if user.House.NursingHouseID != "00001" {
		return 0, nil
	}

	monthsUntilRetirement := user.RetirementPlan.RetirementAge * 12
	yearUntilLifespan := user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge
	totalNursingHouseCost := user.House.NursingHouse.Price * yearUntilLifespan
	cost := float64(totalNursingHouseCost) / float64(monthsUntilRetirement)
	cost = math.Round(cost)
	return cost, nil
}

func CalculateAllAssetSavings(user *entities.User, method string) float64 {
	var total float64
	if method == "All" {
		for _, asset := range user.Assets {
			total += asset.CurrentMoney
		}
	} else if method == "Plan" {
		for _, asset := range user.Assets {
			if asset.Status == "Completed" {
				total += asset.TotalCost
			} else {
				total += asset.CurrentMoney
			}
		}
	}

	total = math.Round(total)
	return total
}

func DistributeSavingMoney(amount float64, count int) []float64 {
	portion := amount / float64(count)
	amounts := make([]float64, count)
	for i := range amounts {
		amounts[i] = portion
	}

	return amounts
}
