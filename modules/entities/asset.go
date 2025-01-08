package entities

type Asset struct {
	ID              string  `json:"asset_id" gorm:"primaryKey"`
	Name            string  `json:"name" gorm:"not null"`
	TotalMoney      float64 `json:"total_money" gorm:"not null"`
	MonthlyExpenses float64 `json:"monthly_expenses" gorm:"not null"`
	EndYear         string  `json:"end_year" gorm:"not null"`
	UserID          string  `json:"user_id" gorm:"not null"` // One-to-many
}
