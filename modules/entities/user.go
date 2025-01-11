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
	Favorites   	[]Favorite 		`json:"-" gorm:"foreignKey:UserID"`
	Assets			[]Asset			`json:"-" gorm:"foreignKey:UserID"`
	NursingHouseID	*string			`json:"-"`
	House			NursingHouse 	`json:"house" gorm:"foreignKey:NursingHouseID;references:ID"`
	RetirementPlan	RetirementPlan	`json:"-" gorm:"foreignKey:UserID"`		
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
}