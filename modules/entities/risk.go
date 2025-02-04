package entities

type Risk struct {
	ID       int    `json:"risk_id" gorm:"primaryKey"`
	RiskName string `json:"risk"`
}
