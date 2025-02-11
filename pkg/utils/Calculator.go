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

	yearsUntilRetirement := plan.RetirementAge - plan.Age
	monthsUntilRetirement := yearsUntilRetirement * 12
	if monthsUntilRetirement <= 0 {
		return 0, errors.New("years until retirement must be greater than zero")
	}

	funds := requiredFunds / float64(monthsUntilRetirement)
	funds = math.Round(funds)
	monthlySavings := funds + plan.AllCostAsset + plan.NursingHousePrice
	return monthlySavings, nil
}

func CalculateMonthlyExpenses(asset *entities.Asset) float64 {
	if asset.Status == "Completed" {
		return 0
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return 0
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	remainingMonths := (endYear-currentYear-1)*12 + (12 - currentMonth + 1)
	remainingCost := asset.TotalCost - asset.CurrentMoney
	monthlyExpenses := remainingCost / float64(remainingMonths)
	monthlyExpenses = math.Round(monthlyExpenses)
	return monthlyExpenses
}

func CalculateAllAssetsMonthlyExpenses(user *entities.User) (float64, error) {
	var total float64
	for _, asset := range user.Assets {
		if asset.Status == "In_Progress" {
			total += asset.MonthlyExpenses
		}
	}

	total = math.Round(total)
	return total, nil
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

func CalculateNursingHouseMonthlyExpense(user *entities.User) (float64, error) {
	if user.House.Status == "Completed" || user.House.NursingHouseID == "00001" {
		return 0, nil
	}

	layout := "02-01-2006"
	birthDate, err := time.Parse(layout, user.RetirementPlan.BirthDate)
	if err != nil {
		return 0, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	retirementYear := birthDate.Year() + user.RetirementPlan.RetirementAge
	retirementMonth := birthDate.Month()
	remainingMonths := (retirementYear-currentYear)*12 + (int(retirementMonth) - currentMonth)
	if remainingMonths <= 0 {
		return 0, nil
	}

	totalCost := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
	remainingCost := float64(totalCost) - user.House.CurrentMoney
	monthlyExpenses := remainingCost / float64(remainingMonths)
	monthlyExpenses = math.Round(monthlyExpenses)
	return monthlyExpenses, nil
}
