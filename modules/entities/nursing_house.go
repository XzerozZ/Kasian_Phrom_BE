package entities

import "time"

type NursingHouse struct {
	ID				uint		`json:"nh_id" gorm:"primaryKey;autoIncrement"`
	Name			string		`json:"name" gorm:"type:varchar(255);unique;not null"`
	Province		string		`json:"province" gorm:"type:varchar(255);not null"`
	Address			string		`json:"address" gorm:"not null"`
	Price			uint		`json:"price" gorm:"not null"`
	Google_map		string		`json:"map" gorm:"type:varchar(50);not null"`
	Phone_number	string		`json:"phone_number" gorm:"not null"`
	Web_site		string		`json:"site"`
	Time			string		`json:"Date" gorm:"not null"`
	Status    		string    	`jsoon:"status" gorm:"type:varchar(50);default:'Active'"`
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}