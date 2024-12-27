package entities

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        		uuid.UUID 	`json:"u_id" gorm:"type:uuid;primaryKey" `
	Firstname      	string    	`json:"fname"`
	Lastname		string		`json:"lname"`
	Username		string		`json:"uname" gorm:"not null"`
	Email     		string    	`json:"email" gorm:"unique;not null"`
	Password  		string    	`json:"-" gorm:"not null"`
	RoleID			int			`json:"r_id" gorm:"not null"`
	Role			Role		`json:"role" gorm:"foreignKey:RoleID`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
}