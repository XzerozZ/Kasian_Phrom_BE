package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormNotiRepository struct {
	db *gorm.DB
}

func NewGormNotiRepository(db *gorm.DB) *GormNotiRepository {
	return &GormNotiRepository{db: db}
}

type NotiRepository interface {
	CreateNotification(notification *entities.Notification) error
	GetNotificationsByUserID(userID string) ([]entities.Notification, error)
}

func (r *GormNotiRepository) CreateNotification(notification *entities.Notification) error {
	return r.db.Create(notification).Error
}

func (r *GormNotiRepository) GetNotificationsByUserID(userID string) ([]entities.Notification, error) {
	var notifications []entities.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *GormNotiRepository) MarkNotificationAsRead(notificationID string) error {
	return r.db.Model(&entities.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}
