package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/repositories"
)

type NotiUsecase interface {
}

type NotiUseCaseImpl struct {
	notirepo repositories.NotiRepository
}

func NewNotiUseCase(notirepo repositories.NotiRepository) *NotiUseCaseImpl {
	return &NotiUseCaseImpl{notirepo: notirepo}
}

func (u *NotiUseCaseImpl) GetNotificationsByUserID(userID string) ([]entities.Notification, error) {
	return u.notirepo.GetNotificationsByUserID(userID)
}
