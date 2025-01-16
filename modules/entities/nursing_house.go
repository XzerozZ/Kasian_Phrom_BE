package entities

import "time"

type NursingHouse struct {
	ID           string  `json:"nh_id" gorm:"primaryKey"`
	Name         string  `json:"name"`
	Province     string  `json:"province"`
	Address      string  `json:"address"`
	Price        int     `json:"price" gorm:"not null"`
	Google_map   string  `json:"map"`
	Phone_number string  `json:"phone_number"`
	Web_site     string  `json:"site"`
	Time         string  `json:"Date"`
	Status       string  `jsoon:"status" gorm:"type:varchar(50);default:'Active'"`
	Images       []Image `json:"images" gorm:"many2many:nh_images;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
