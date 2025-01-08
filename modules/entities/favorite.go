package entities

type Favorite struct {
	UserID			string		`json:"user_id" gorm:"primaryKey"`
	NursingHouseID	string		`json:"nh_id" gorm:"primaryKey"`
}