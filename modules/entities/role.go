package entities

type Role struct {
	ID				int			`json:"r_id" gorm:"primaryKey"`
	RoleName		string		`json:"role"`
}