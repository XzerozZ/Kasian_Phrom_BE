package entities

import "time"

type SelectedHouse struct {
	UserID              string       `json:"-" gorm:"primaryKey"`
	NursingHouseID      string       `json:"-"`
	CurrentMoney        float64      `json:"current_money" gorm:"default:0.0"`
	Status              string       `json:"status" gorm:"not null"`
	MonthlyExpenses     float64      `json:"monthly_expenses" gorm:"default:0.0"`
	LastCalculatedMonth int          `json:"last_calculated_month" gorm:"default:0"`
	NursingHouse        NursingHouse `gorm:"foreignKey:NursingHouseID"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
