package entities

type Dialog struct {
	ID				string		`json:"d_id" gorm:"primaryKey"`
	Type			string		`json:"type" gorm:"not null"`
	Desc			string		`json:"desc" gorm:"not null"`
	Bold			bool		`json:"bold" gorm:"not null"`
	NewsID			string		`json:"news_id" gorm:"not null"`
}