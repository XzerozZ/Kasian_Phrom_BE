package entities

type Asset struct {
	ID              	string  	`json:"asset_id" gorm:"primaryKey"`
	Name            	string  	`json:"name" gorm:"not null"`
	Type				string		`json:"type" gorm:"not null"`
	TotalCost	     	float64	 	`json:"total_cost" gorm:"not null"`
	CurrentMoney		float64		`json:"current_money" gorm:'default:0.0'`
	Status				string		`json:"status" gorm:"default:'In_Progress'"`
	EndYear         	string  	`json:"end_year" gorm:"not null"`
	UserID          	string  	`json:"-" gorm:"not null"`
}
