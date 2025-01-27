package entities

import "time"

type Loan struct {
	ID                 string  `json:"loan_id" gorm:"primaryKey"`
	Name               string  `json:"name" gorm:"not null"`
	Type               string  `json:"type" gorm:"not null"`
	MonthlyExpenses    float64 `json:"monthly_expenses" gorm:"not null"`
	InterestPercentage float64 `json:"interest_percentage" gorm:"not null"`
	RemainingMonths    int     `json:"remaining_months" gorm:"not null"`
	Installment        bool    `json:"installment" gorm:"not null"`
	Status             string  `json:"status" gorm:"not null"`
	UserID             string  `json:"-" gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
