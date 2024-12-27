package entities

type Dialog struct {
	ID				int			`json:"d_id" gorm:"primaryKey;autoIncrement"`
	Type			string		`json:"type" gorm:"not null"`
	Desc			string		`json:"desc" gorm:"not null"`
	NewsID			string		`json:"news_id" gorm:"not null"`
}