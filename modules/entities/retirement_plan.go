package entities

type RetirementPlan struct {
	ID                     string  `json:"retirement_plan_id" gorm:"primaryKey"`
	UserID                 string  `json:"user_id" gorm:"not null"`
	TotalRetirementExpense float64 `json:"total_retirement_expense" gorm:"not null"`//เงินที่ต้องใช้เกษียณทั้งหมด
	RequiredMonthlySavings float64 `json:"required_monthly_savings" gorm:"not null"`//เงินที่ต้องออมต่อเดือน
	TotalAssetCost         float64 `json:"total_asset_cost" gorm:"not null"` //มูลค่าของทรัพย์สินที่ต้องการครอบครองหลังเกษียณ
}
