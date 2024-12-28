package entities

import "time"

type News struct {
	ID				string		`json:"news_id" gorm:"primaryKey"`
	Title			string		`json:"title" gorm:"unique;not null"`
	Dialog 			[]Dialog 	`json:"role" gorm:"foreignKey:NewsID"`
	Images 			[]Image 	`json:"images" gorm:"many2many:news_images;" `
	PublishedAt   	time.Time	`json:"published_date"`
	UpdatedAt   	time.Time	`json:"updated_date"`
}