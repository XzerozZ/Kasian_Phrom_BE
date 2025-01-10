package entities

import "time"

type User struct {
	ID        		string	 		`json:"u_id" gorm:"primaryKey" `
	Firstname      	string    		`json:"fname"`
	Lastname		string			`json:"lname"`
	Username		string			`json:"uname" gorm:"not null"`
	Email     		string    		`json:"email" gorm:"unique;not null"`
	Password  		string    		`json:"-" gorm:"not null"`
	ImageLink 		string 			`json:"image_link"`
	RoleID			int				`json:"-" gorm:"not null"`
	Role			Role			`json:"role" gorm:"foreignKey:RoleID`
	Favorites   	[]Favorite 		`json:"favorites" gorm:"foreignKey:UserID"`
	Assets			[]Asset			`json:"assets" gorm:"foreignKey:UserID"`
	HouseID			*string			`json:"-"`
	House			NursingHouse 	`json:"house" gorm:"foreignKey:HouseID"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
}