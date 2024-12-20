package entities

import (
	"gorm.io/gorm"
)

type NhUsecase interface {
	CreateNh(req *Nursing_House) (*Nursing_House, error)
}

type NhRepository interface {
	CreateNh(req *Nursing_House) (*Nursing_House, error)
}

type Nursing_House struct {
	gorm.Model
	Nh_ID			uint	`json:"nh_id" gorm:"primaryKey;autoIncrement"`
	Name			string	`json:"name" gorm:"unique;not null"`
	Province		string	`json:"province" gorm:"not null"`
	Address			string	`json:"address" gorm:"not null"`
	Price			uint	`json:"price" gorm:"not null"`
	Google_map		string	`json:"map" gorm:"not null"`
	Phone_number	string	`json:"phone_number" gorm:"not null"`
	Web_site		string	`json:"site"`
	Time			string	`json:"Date'`
}