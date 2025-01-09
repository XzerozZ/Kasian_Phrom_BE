package entities

type Financial struct {
	ID                      string  `json:"financial_id" gorm:"primaryKey"`
	// PlanName                string  `json:"planName" gorm:"not null"`
	BirthDate               string  `json:"birth_date" gorm:"not null"`
	Age                     uint8   `json:"age" gorm:"not null"`
	RetirementAge           uint8   `json:"retirement_age" gorm:"not null"`
	ExpLifespan             uint8   `json:"exp_lifespan" gorm:"not null"`
	CurrentSavingsReturns   float64 `json:"current_savings_returns" gorm:"not null"`
	MonthlyIncome           float64 `json:"monthly_income" gorm:"not null"`
	MonthlyExpenses         float64 `json:"monthly_expenses" gorm:"not null"`
	CurrentTotalInvestment  float64 `json:"current_total_investment"`
	InvestmentReturn        float64 `json:"investment_return"`
	ExpectedMonthlyExpenses float64 `json:"expected_monthly_expenses" gorm:"not null"`
	ExpectedInflation       float64 `json:"expected_inflation" gorm:"not null"`
	AnnualExpenseIncrease   float64 `json:"annual_expense_increase" gorm:"not null"`
	AnnualSavingsReturn     float64 `json:"annual_savings_return" gorm:"not null"`
	AnnualInvestmentReturn  float64 `json:"annual_investment_return" gorm:"not null"`
	UserID                  string  `json:"user_id" gorm:"foreignKey:UserID"`
}
