package entities

import "time"

type Notification struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message" gorm:"not null"`
	Balance   float64   `json:"balance"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
