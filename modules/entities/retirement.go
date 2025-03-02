package entities

import "time"

type RetirementPlan struct {
	ID                      string    `json:"financial_id" gorm:"primaryKey"`
	PlanName                string    `json:"planName" gorm:"not null"`
	BirthDate               string    `json:"birth_date" gorm:"not null"`
	RetirementAge           int       `json:"retirement_age" gorm:"not null"`
	ExpectLifespan          int       `json:"expect_lifespan" gorm:"not null"`
	CurrentSavings          float64   `json:"current_savings" gorm:"not null"`
	CurrentSavingsReturns   float64   `json:"current_savings_returns" gorm:"not null"`
	MonthlyIncome           float64   `json:"monthly_income" gorm:"not null"`
	MonthlyExpenses         float64   `json:"monthly_expenses" gorm:"not null"`
	CurrentTotalInvestment  float64   `json:"current_total_investment" gorm:"not null"`
	InvestmentReturn        float64   `json:"investment_return" gorm:"not null"`
	ExpectedMonthlyExpenses float64   `json:"expected_monthly_expenses"`
	ExpectedInflation       float64   `json:"expected_inflation" gorm:"not null"`
	AnnualExpenseIncrease   float64   `json:"annual_expense_increase" gorm:"not null"`
	AnnualSavingsReturn     float64   `json:"annual_savings_return" gorm:"not null"`
	AnnualInvestmentReturn  float64   `json:"annual_investment_return" gorm:"not null"`
	LastRequiredFunds       float64   `json:"last_required_funds" gorm:"default:0.0"`
	LastMonthlyExpenses     float64   `json:"last_monthly_expenses" gorm:"default:0.0"`
	LastCalculatedMonth     int       `json:"last_calculated_month" gorm:"default:0"`
	Status                  string    `json:"status" gorm:"not null"`
	UserID                  string    `json:"user_id" gorm:"unique;foreignKey:UserID"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
