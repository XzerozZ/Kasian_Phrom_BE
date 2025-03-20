package entities

import "time"

type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Balance   float64   `json:"balance"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	ObjectID  string    `json:"object_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
