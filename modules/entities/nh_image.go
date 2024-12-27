package entities

type NHImage struct {
	NHID 		string 	`json:"nh_id" gorm:"primaryKey"`
	ImageID     string 	`json:"image_id" gorm:"primaryKey"`
}