package entities

type Favorite struct {
	UserID         string       `json:"u_id" gorm:"primaryKey"`
	NursingHouseID string       `json:"nh_id" gorm:"primaryKey"`
	User           User         `json:"-" gorm:"foreignKey:UserID;references:ID"`
	NursingHouse   NursingHouse `gorm:"foreignKey:NursingHouseID;references:ID"`
}
