package entities

type Quiz struct {
	UserID string `json:"-" gorm:"primaryKey"`
	RiskID int    `json:"-" gorm:"primaryKey"`
	Risk   Risk   `gorm:"foreignKey:RiskID;references:ID"`
}
