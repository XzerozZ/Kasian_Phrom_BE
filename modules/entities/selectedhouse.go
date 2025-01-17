package entities

type SelectedHouse struct {
	UserID         string       `json:"-" gorm:"primaryKey"`
	NursingHouseID string       `json:"-"`
	CurrentMoney   float64      `json:"current_money" gorm:"default:0.0"`
	NursingHouse   NursingHouse `gorm:"foreignKey:NursingHouseID"`
}
