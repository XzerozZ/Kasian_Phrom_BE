package utils

import (
	"math"
	"errors"
)

type MonthlyExpensesPlan struct {
	MonthlyExpenses     	float64
	AnnualExpenseIncrease 	float64
	ExpectedInflation   	float64
	Age                  	int
	RetirementAge        	int
	ExpectLifespan       	int
	YearsUntilRetirement 	int
	YearUntilLifeSpan		int
	AllCostAsset            float64
	NursingHousePrice   	float64
}

func CalculateRetirementFunds(plan MonthlyExpensesPlan) (float64, error) {
	if plan.MonthlyExpenses < 0 || plan.AnnualExpenseIncrease < 0 || plan.ExpectedInflation < 0 {
		return 0, errors.New("MonthlyExpenses, AnnualExpenseIncrease and ExpectedInflation must be greater than zero")
	}

	totalRequiredFunds := 0.0
	for i := 0; i < plan.YearsUntilRetirement; i++ {
		monthlyExpenses := plan.MonthlyExpenses * 12
		for j := 0; j < plan.ExpectLifespan-plan.Age-1; j++ {
			increaseFactor := math.Pow(1+plan.AnnualExpenseIncrease+plan.ExpectedInflation, float64(j))
			totalRequiredFunds += monthlyExpenses * increaseFactor
		}
	}

	return totalRequiredFunds, nil
}

func CalculateMonthlySavings(plan MonthlyExpensesPlan) (float64, error) {
	requiredFunds, err := CalculateRetirementFunds(plan)
	if err != nil {
		return 0, err
	}
	
	totalNursingHouseCost := plan.NursingHousePrice * float64(plan.YearUntilLifeSpan)
	monthlySavings := (requiredFunds + totalNursingHouseCost) / float64(plan.YearsUntilRetirement*12) + plan.AllCostAsset
	return monthlySavings, nil
}
