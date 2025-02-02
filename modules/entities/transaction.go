package entities

import "time"

type Transaction struct {
	ID              string    `json:"trans_id" gorm:"primaryKey"`
	Status          string    `json:"status" gorm:"not null"`
	Type            string    `json:"type" gorm:"not null"`
	MonthlyExpenses float64   `json:"monthly_expenses" gorm:"not null"`
	UserID          string    `json:"-" gorm:"not null"`
	LoanID          string    `json:"-" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at"`
}
