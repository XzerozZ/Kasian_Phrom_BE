package entities

type Favorite struct {
	UserID          string 			`json:"-" gorm:"primaryKey"`
	NursingHouseID  string 			`json:"-" gorm:"primaryKey"`
	User            User   			`json:"-" gorm:"foreignKey:UserID;references:ID"`
	NursingHouse    NursingHouse 	`gorm:"foreignKey:NursingHouseID;references:ID"`
}