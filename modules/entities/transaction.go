package entities

import "time"

type Transaction struct {
	ID        string    `json:"trans_id" gorm:"primaryKey"`
	Status    string    `json:"status" gorm:"not null"`
	UserID    string    `json:"-" gorm:"not null"`
	LoanID    string    `json:"-" gorm:"not null"`
	Loan      Loan      `gorm:"foreignKey:LoanID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
}
