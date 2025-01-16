package entities

import "time"

type OTP struct {
	UserID      string  		`json:"-" gorm:"primaryKey"`
	OTP       	string    		`json:"otp" gorm:"not null;unique"`
	ExpiresAt 	time.Time 		`json:"expires_at" gorm:"not null"`
	User   		User   			`json:"-" gorm:"foreignKey:UserID;references:ID"`
}