package entities

import "time"

type History struct {
	ID        string    `json:"history_id" gorm:"primaryKey"`
	Method    string    `json:"method" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Category  string    `json:"category" gorm:"not null"`
	Money     float64   `json:"money" gorm:"not null"`
	UserID    string    `json:"-" gorm:"not null"`
	TrackDate time.Time `json:"track_at"`
}
