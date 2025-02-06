package entities

import "time"

type User struct {
	ID             string         `json:"u_id" gorm:"primaryKey" `
	Firstname      string         `json:"fname"`
	Lastname       string         `json:"lname"`
	Username       string         `json:"uname" gorm:"not null"`
	Email          string         `json:"email" gorm:"unique;not null"`
	Password       string         `json:"-"`
	Provider       string         `json:"provider" gorm:"not null"`
	ImageLink      string         `json:"image_link"`
	RoleID         int            `json:"-" gorm:"not null"`
	Role           Role           `json:"role" gorm:"foreignKey:RoleID"`
	Favorites      []Favorite     `json:"favorites,omitempty" gorm:"foreignKey:UserID"`
	Assets         []Asset        `json:"assets,omitempty" gorm:"foreignKey:UserID"`
	Loans          []Loan         `json:"loans,omitempty" gorm:"foreignKey:UserID"`
	House          SelectedHouse  `json:"house" gorm:"foreignKey:UserID"`
	RetirementPlan RetirementPlan `json:"retirement" gorm:"foreignKey:UserID"`
	Quiz           Quiz           `json:"risk" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
