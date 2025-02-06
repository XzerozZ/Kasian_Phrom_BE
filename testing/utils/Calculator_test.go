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
	})

	t.Run("birthday not yet occurred this year", func(t *testing.T) {
		now := time.Now()
		futureBirthday := now.AddDate(-30, 0, 1)
		birthDateStr := futureBirthday.Format("02-01-2006")

		age, err := utils.CalculateAge(birthDateStr)
		assert.NoError(t, err)
		assert.Equal(t, 29, age)
	})

	t.Run("edge case - today's birthday", func(t *testing.T) {
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
		assert.Equal(t, 0, age)
	})

	t.Run("birthday not yet occurred this year", func(t *testing.T) {
		birthDate := "15-12-1990"
		planCreationDate := time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC)

		age, err := utils.CalculateRetirementPlanAge(birthDate, planCreationDate)
		assert.NoError(t, err)
		assert.Equal(t, 33, age)
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

	t.Run("invalid retirement age", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			Age:           65,
			RetirementAge: 60,
		}
		funds, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), funds)
	})

	t.Run("zero monthly expenses", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			ExpectedMonthlyExpenses: 0,
			AnnualExpenseIncrease:   2,
			ExpectedInflation:       3,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
		}
		funds, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), funds)
	})

	t.Run("negative inflation rate", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			ExpectedMonthlyExpenses: 3000,
			AnnualExpenseIncrease:   2,
			ExpectedInflation:       -1,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          80,
		}
		funds, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), funds)
	})

	t.Run("lifespan less than retirement age", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			ExpectedMonthlyExpenses: 3000,
			AnnualExpenseIncrease:   2,
			ExpectedInflation:       3,
			Age:                     30,
			RetirementAge:           60,
			ExpectLifespan:          50,
		}
		funds, err := utils.CalculateRetirementFunds(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), funds)
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
			YearsUntilRetirement: 0,
		}
		savings, err := utils.CalculateMonthlySavings(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), savings)
	})

	t.Run("invalid retirement funds calculation", func(t *testing.T) {
		plan := utils.MonthlyExpensesPlan{
			Age:                  65,
			RetirementAge:        60,
			YearsUntilRetirement: 10,
		}
		savings, err := utils.CalculateMonthlySavings(plan)
		assert.Error(t, err)
		assert.Equal(t, float64(0), savings)
	})
}

func TestCalculateMonthlyExpenses(t *testing.T) {
	t.Run("valid asset calculation", func(t *testing.T) {
		now := time.Now()
		asset := entities.Asset{
			TotalCost:    60000,
			CurrentMoney: 10000,
			EndYear:      strconv.Itoa(now.Year() + 5),
			UpdatedAt:    now,
		}
		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.NoError(t, err)
		assert.Greater(t, expenses, float64(0))
	})

	t.Run("invalid end year", func(t *testing.T) {
		asset := entities.Asset{
			TotalCost:    60000,
			CurrentMoney: 10000,
			EndYear:      "invalid",
			UpdatedAt:    time.Now(),
		}
		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.Error(t, err)
		assert.Equal(t, float64(0), expenses)
	})

	t.Run("current money exceeds total cost", func(t *testing.T) {
		asset := entities.Asset{
			TotalCost:    60000,
			CurrentMoney: 70000,
			EndYear:      strconv.Itoa(time.Now().Year() + 5),
			UpdatedAt:    time.Now(),
		}
		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.Error(t, err)
		assert.Equal(t, float64(0), expenses)
	})

	t.Run("end year less than current year", func(t *testing.T) {
		asset := entities.Asset{
			TotalCost:    60000,
			CurrentMoney: 10000,
			EndYear:      strconv.Itoa(time.Now().Year() - 1),
			UpdatedAt:    time.Now(),
		}
		expenses, err := utils.CalculateMonthlyExpenses(&asset)
		assert.Error(t, err)
		assert.Equal(t, float64(0), expenses)
	})
}

func TestCalculateAllAssetsMonthlyExpenses(t *testing.T) {
	t.Run("valid user with multiple assets", func(t *testing.T) {
		now := time.Now()
		user := &entities.User{
			Assets: []entities.Asset{
				{
					TotalCost:    60000,
					CurrentMoney: 10000,
					EndYear:      strconv.Itoa(now.Year() + 5),
					UpdatedAt:    now,
				},
				{
					TotalCost:    30000,
					CurrentMoney: 5000,
					EndYear:      strconv.Itoa(now.Year() + 3),
					UpdatedAt:    now,
				},
			},
		}
		totalExpenses, err := utils.CalculateAllAssetsMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Greater(t, totalExpenses, float64(0))
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

func TestCalculateNursingHouseMonthlyExpenses(t *testing.T) {
	t.Run("valid nursing house calculation", func(t *testing.T) {
		user := &entities.User{
			RetirementPlan: entities.RetirementPlan{
				RetirementAge:  60,
				ExpectLifespan: 80,
			},
			House: entities.SelectedHouse{
				NursingHouse: entities.NursingHouse{
					Price: 5000,
				},
			},
		}
		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Greater(t, monthlyExpenses, float64(0))
	})

	t.Run("zero lifespan", func(t *testing.T) {
		user := &entities.User{
			RetirementPlan: entities.RetirementPlan{
				RetirementAge:  60,
				ExpectLifespan: 80,
			},
		}
		monthlyExpenses, err := utils.CalculateNursingHouseMonthlyExpenses(user)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), monthlyExpenses)
	})
}

func TestCalculateAllAssetSavings(t *testing.T) {
	t.Run("valid user with multiple assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{
				{CurrentMoney: 10000},
				{CurrentMoney: 5000},
			},
		}
		totalSavings, err := utils.CalculateAllAssetSavings(user)
		assert.NoError(t, err)
		assert.Equal(t, float64(15000), totalSavings)
	})

	t.Run("user with no assets", func(t *testing.T) {
		user := &entities.User{
			Assets: []entities.Asset{},
		}
		totalSavings, err := utils.CalculateAllAssetSavings(user)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), totalSavings)
	})
}
