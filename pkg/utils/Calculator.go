package utils

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
)

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

func CalculateRetirementFunds(plan *entities.RetirementPlan, age int) (float64, error) {
	yearsUntilRetirement := plan.RetirementAge - age
	yearsInRetirement := plan.ExpectLifespan - plan.RetirementAge
	if yearsUntilRetirement <= 0 {
		return 0, errors.New("retirement age must be greater than current age")
	}

	if yearsInRetirement <= 0 {
		return 0, errors.New("expected lifespan must be greater than retirement age")
	}

	var totalRequiredFunds float64
	annualExpenses := plan.ExpectedMonthlyExpenses * 12
	for year := 1; year <= yearsInRetirement; year++ {
		remainingYears := yearsUntilRetirement + year
		compoundingFactor := math.Pow(1+(plan.ExpectedInflation/100), float64(remainingYears))
		totalRequiredFunds += annualExpenses * compoundingFactor
	}

	return math.Round(totalRequiredFunds), nil
}

func CalculateMonthlySavings(plan *entities.RetirementPlan, age, currentYear, currentMonth int) (float64, error) {
	requiredFunds, err := CalculateRetirementFunds(plan, age)
	if err != nil {
		return 0, err
	}

	birthDate, err := time.Parse("02-01-2006", plan.BirthDate)
	if err != nil {
		return 0, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	retirementYear := birthDate.Year() + plan.RetirementAge
	retirementMonth := int(birthDate.Month())
	remainingMonths := (retirementYear-currentYear)*12 + (retirementMonth - currentMonth + 1)
	if remainingMonths <= 0 {
		return 0, errors.New("months until retirement must be greater than zero")
	}

	remainingMoney := requiredFunds - (plan.CurrentSavings + plan.CurrentTotalInvestment)
	return math.Round(remainingMoney / float64(remainingMonths)), nil
}

func CalculateMonthlyExpenses(asset *entities.Asset, currentYear, currentMonth int) float64 {
	if asset.Status != "In_Progress" {
		return 0
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return 0
	}

	remainingMonths := (endYear-currentYear)*12 - (currentMonth - 1)
	if currentYear >= endYear {
		return 0
	}

	if remainingMonths <= 0 {
		return 0
	}

	remainingCost := asset.TotalCost - asset.CurrentMoney
	return math.Round(remainingCost / float64(remainingMonths))
}

func CalculateAllAssetsMonthlyExpenses(user *entities.User) (float64, error) {
	var total float64
	currentMonth := int(time.Now().Month())
	for _, asset := range user.Assets {
		if asset.Status == "In_Progress" && asset.LastCalculatedMonth == currentMonth {
			total += asset.MonthlyExpenses
		} else if asset.Status == "In_Progress" && asset.LastCalculatedMonth != currentMonth {

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

func CalculateNursingHouseMonthlyExpense(user *entities.User, currentYear, currentMonth int) (float64, error) {
	if user.House.Status == "Completed" || user.House.NursingHouseID == "00001" {
		return 0, nil
	}

	birthDate, err := time.Parse("02-01-2006", user.RetirementPlan.BirthDate)
	if err != nil {
		return 0, errors.New("invalid BirthDate format, expected DD-MM-YYYY")
	}

	retirementYear := birthDate.Year() + user.RetirementPlan.RetirementAge
	retirementMonth := int(birthDate.Month())
	remainingMonths := (retirementYear-currentYear)*12 + (retirementMonth - currentMonth + 1)
	if remainingMonths <= 0 {
		return 0, nil
	}

	totalCost := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
	remainingCost := float64(totalCost) - user.House.CurrentMoney
	return math.Round(remainingCost / float64(remainingMonths)), nil
}
