package entities

import "time"

type News struct {
	ID          string    `json:"news_id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"unique;not null"`
	Image_Title string    `json:"image_title" gorm:"not null"`
	Image_Desc  string    `json:"image_desc" gorm:"not null"`
	Dialog      []Dialog  `json:"dialog" gorm:"foreignKey:NewsID"`
	PublishedAt time.Time `json:"published_date"`
	UpdatedAt   time.Time `json:"updated_date"`
}
