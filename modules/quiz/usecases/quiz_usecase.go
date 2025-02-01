package usecases

import (
	"errors"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
)

type QuizUseCase interface {
	CreateQuiz(userID string, weight []int) (*entities.Quiz, error)
	GetQuizByUserID(userID string) (*entities.Quiz, error)
}

type QuizUseCaseImpl struct {
	quizrepo repositories.QuizRepository
}

func NewQuizUseCase(quizrepo repositories.QuizRepository) *QuizUseCaseImpl {
	return &QuizUseCaseImpl{quizrepo: quizrepo}
}

func (u *QuizUseCaseImpl) CreateQuiz(userID string, weight []int) (*entities.Quiz, error) {
	quiz, err := u.quizrepo.GetQuizByUserID(userID)
	if err == nil && quiz != nil {
		if err := u.quizrepo.DeleteQuiz(userID); err != nil {
			return nil, err
		}
	}

	id, err := utils.CalculateRisk(weight)
	if err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, errors.New("invalid risk calculation: empty risk ID")
	}

	newRisk := &entities.Quiz{
		UserID: userID,
		RiskID: id,
	}

	createdQuiz, err := u.quizrepo.CreateQuiz(newRisk)
	if err != nil {
		return nil, err
	}

	return createdQuiz, nil
}

func (u *QuizUseCaseImpl) GetQuizByUserID(userID string) (*entities.Quiz, error) {
	return u.quizrepo.GetQuizByUserID(userID)
}
