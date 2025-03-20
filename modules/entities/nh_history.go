package entities

type NursingHouseHistory struct {
	UserID         string       `json:"user_id" gorm:"primaryKey"`
	NursingHouseID string       `json:"nh_id" gorm:"not null"`
	NursingHouse   NursingHouse `gorm:"foreignKey:NursingHouseID"`
}
