package entities

import "time"

type News struct {
	ID				string		`json:"news_id" gorm:"primaryKey"`
	Title			string		`json:"title" gorm:"unique;not null"`
	Dialog 			[]Dialog 	`json:"role" gorm:"foreignKey:RoleID`
	PublishedAt   	time.Time	`json:"published_date"`
	UpdatedAt   	time.Time	`json:"updated_date"`
}

type CreateNewsRequest struct {
	Title   		string      `json:"title"`
	Dialogs 		[]CreateDialogRequest `json:"dialogs"`
}

type CreateDialogRequest struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
}