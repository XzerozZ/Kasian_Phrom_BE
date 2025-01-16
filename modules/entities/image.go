package entities

type Image struct {
	ID        string `json:"image_id" gorm:"primaryKey"`
	ImageLink string `json:"image_link" gorm:"not null"`
}
