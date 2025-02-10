package utils_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {
	t.Run("valid birth date", func(t *testing.T) {
		birthDate := time.Now().AddDate(-30, 0, 0).Format("02-01-2006")
		age, err := utils.CalculateAge(birthDate)
		assert.NoError(t, err)
		assert.Equal(t, 30, age)
	})

	t.Run("invalid date format", func(t *testing.T) {
		birthDate := "2000-01-02"
		age, err := utils.CalculateAge(birthDate)
		assert.Error(t, err)
		assert.Equal(t, 0, age)
		assert.Contains(t, err.Error(), "invalid BirthDate format")
	})

	t.Run("birthday not yet occurred this year", func(t *testing.T) {
		now := time.Now()
		futureBirthday := now.AddDate(-30, 1, 0) // One month in the future
		birthDateStr := futureBirthday.Format("02-01-2006")
		age, err := utils.CalculateAge(birthDateStr)
		assert.NoError(t, err)
		assert.Equal(t, 29, age)
	})

	t.Run("edge case - born today", func(t *testing.T) {
		now := time.Now()
		birthDateStr := now.Format("02-01-2006")
		age, err := utils.CalculateAge(birthDateStr)
		assert.NoError(t, err)
		assert.Equal(t, 0, age)
	})
}

func TestCalculateRetirementPlanAge(t *testing.T) {
	t.Run("valid birth date", func(t *testing.T) {
		birthDate := "15-05-1990"
		planCreationDate := time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC)
		age, err := utils.CalculateRetirementPlanAge(birthDate, planCreationDate)
		assert.NoError(t, err)
		assert.Equal(t, 33, age)
	})

	t.Run("invalid date format", func(t *testing.T) {
		birthDate := "1990-05-15"
		planCreationDate := time.Now()
		age, err := utils.CalculateRetirementPlanAge(birthDate, planCreationDate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid BirthDate format")
		assert.Equal(t, 0, age)
	})

	t.Run("future birth date", func(t *testing.T) {
		now := time.Now()
		birthDate := now.AddDate(1, 0, 0).Format("02-01-2006")
		age, err := utils.CalculateRetirementPlanAge(birthDate, now)
		assert.NoError(t, err)
		assert.Equal(t, -1, age)
	})
}

func TestCalculateRetirementFunds(t *testing.T) {
	t.Run("valid retirement plan", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			ExpectedMonthlyExpenses: 3000,
			AnnualExpenseIncrease:   2,
			ExpectedInflation:       3,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
		}

		funds, err := utils.CalculateRetirementFunds(plan)
		assert.NoError(t, err)
		assert.Greater(t, funds, float64(0))
	})

	t.Run("retirement age equal to current age", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			Age:            60,
			RetirementAge:  60,
			ExpectLifespan: 80,
		}

		_, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "retirement age must be greater than current age")
	})

	t.Run("expected lifespan equal to retirement age", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			Age:            30,
			RetirementAge:  60,
			ExpectLifespan: 60,
		}

		_, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected lifespan must be greater than retirement age")
	})
}

func TestCalculateMonthlySavings(t *testing.T) {
	t.Run("valid savings plan", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			ExpectedMonthlyExpenses: 3000,
			AnnualExpenseIncrease:   2,
			ExpectedInflation:       3,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
			YearsUntilRetirement:    30,
			AllCostAsset:            50000,
			NursingHousePrice:       10000,
		}

		savings, err := utils.CalculateMonthlySavings(plan)
		assert.NoError(t, err)
		assert.Greater(t, savings, float64(0))
	})

	t.Run("zero years until retirement", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			YearsUntilRetirement:    0,
			ExpectedMonthlyExpenses: 3000,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
		}

		_, err := utils.CalculateMonthlySavings(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "years until retirement must be greater than zero")
	})

	t.Run("negative years until retirement", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			YearsUntilRetirement:    -5,
			ExpectedMonthlyExpenses: 3000,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
		}

		_, err := utils.CalculateMonthlySavings(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "years until retirement must be greater than zero")
	})
}

func TestCalculateMonthlyExpenses(t *testing.T) {
	t.Run("completed asset", func(t *testing.T) {
		asset := entities.Asset{
			Status:       "Completed",
			TotalCost:    60000,
			CurrentMoney: 60000,
		}

		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), expenses)
	})

	t.Run("valid in-progress asset", func(t *testing.T) {
		now := time.Now()
		asset := entities.Asset{
			Status:       "In_Progress",
			TotalCost:    60000,
			CurrentMoney: 10000,
			EndYear:      strconv.Itoa(now.Year() + 5),
			UpdatedAt:    now,
		}

		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.NoError(t, err)
		assert.Greater(t, expenses, float64(0))
	})

	t.Run("current money equals total cost", func(t *testing.T) {
		now := time.Now()
		asset := entities.Asset{
			Status:       "In_Progress",
			TotalCost:    60000,
			CurrentMoney: 60000,
			EndYear:      strconv.Itoa(now.Year() + 5),
			UpdatedAt:    now,
		}

		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), expenses)
	})
}

func TestCalculateAllAssetsMonthlyExpenses(t *testing.T) {
	now := time.Now()
	t.Run("mixed status assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{
				{
					Status:       "In_Progress",
					TotalCost:    60000,
					CurrentMoney: 10000,
					EndYear:      strconv.Itoa(now.Year() + 5),
					UpdatedAt:    now,
				},
				{
					Status:       "Completed",
					TotalCost:    30000,
					CurrentMoney: 30000,
					EndYear:      strconv.Itoa(now.Year() + 3),
					UpdatedAt:    now,
				},
			},
		}

		totalExpenses, err := utils.CalculateAllAssetsMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Greater(t, totalExpenses, float64(0))
	})

	t.Run("all completed assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{
				{
					Status:       "Completed",
					TotalCost:    60000,
					CurrentMoney: 60000,
				},
				{
					Status:       "Completed",
					TotalCost:    30000,
					CurrentMoney: 30000,
				},
			},
		}

		totalExpenses, err := utils.CalculateAllAssetsMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), totalExpenses)
	})

	t.Run("user with no assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{},
		}

		totalExpenses, err := utils.CalculateAllAssetsMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), totalExpenses)
	})
}

func TestCalculateAllAssetSavings(t *testing.T) {
	t.Run("all method with mixed status", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{
				{
					Status:       "In_Progress",
					CurrentMoney: 10000,
					TotalCost:    20000,
				},
				{
					Status:       "Completed",
					CurrentMoney: 30000,
					TotalCost:    30000,
				},
			},
		}

		total := utils.CalculateAllAssetSavings(user, "All")
		assert.Equal(t, float64(40000), total)
	})

	t.Run("plan method with mixed status", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{
				{
					Status:       "In_Progress",
					CurrentMoney: 10000,
					TotalCost:    20000,
				},
				{
					Status:       "Completed",
					CurrentMoney: 30000,
					TotalCost:    30000,
				},
			},
		}

		total := utils.CalculateAllAssetSavings(user, "Plan")
		assert.Equal(t, float64(40000), total)
	})

	t.Run("empty assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{},
		}

		total := utils.CalculateAllAssetSavings(user, "All")
		assert.Equal(t, float64(0), total)
	})
}
