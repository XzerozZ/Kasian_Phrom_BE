package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormQuizRepository struct {
	db *gorm.DB
}

func NewGormQuizRepository(db *gorm.DB) *GormQuizRepository {
	return &GormQuizRepository{db: db}
}

type QuizRepository interface {
	CreateQuiz(quiz *entities.Quiz) (*entities.Quiz, error)
	GetQuizByUserID(userID string) (*entities.Quiz, error)
	DeleteQuiz(userID string) error
}

func (r *GormQuizRepository) CreateQuiz(quiz *entities.Quiz) (*entities.Quiz, error) {
	if err := r.db.Create(quiz).Error; err != nil {
		return nil, err
	}

	return r.GetQuizByUserID(quiz.UserID)
}

func (r *GormQuizRepository) GetQuizByUserID(userID string) (*entities.Quiz, error) {
	var quiz entities.Quiz
	if err := r.db.Preload("Risk").Where("user_id = ?", userID).First(&quiz).Error; err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (r *GormQuizRepository) DeleteQuiz(userID string) error {
	if err := r.db.Delete(&entities.Quiz{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}
