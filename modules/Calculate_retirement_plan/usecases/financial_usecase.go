package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type FinUseCase interface {
	CreateFin(financial entities.Financial, ctx *fiber.Ctx) (*entities.Financial, error)
	GetFinByID(id string) (*entities.Financial, error)
	GetFinByUserID(id string) (*entities.Financial, error)
	GetFinNextID() (string, error)
}

type FinUseCaseImpl struct {
	finrepo repositories.FinRepository
}

func NewFinUseCase(finrepo repositories.FinRepository) *FinUseCaseImpl {
	return &FinUseCaseImpl{
		finrepo: finrepo,
	}
}

func (u *FinUseCaseImpl) CreateFin(financial entities.Financial, ctx *fiber.Ctx) (*entities.Financial, error) {
	//เปลี่ยนเป็น nextID
	// financial.ID = uuid.New().String()
	id, err := u.finrepo.GetFinNextID()
	if err != nil {
		return nil, err
	}


	if financial.Age <= 0 {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Age must be greater than zero",
		})
	}

	financial.ID = id

	var createdFin *entities.Financial
	createdFin, err = u.finrepo.CreateFin(&financial)
	if err != nil {
		return nil, err
	}

	return createdFin, nil
}

func (u *FinUseCaseImpl) GetFinByID(financial_id string) (*entities.Financial, error) {
	return u.finrepo.GetFinByID(financial_id)
}

func (u *FinUseCaseImpl) GetFinByUserID(user_id string) (*entities.Financial, error) {
	return u.finrepo.GetFinByUserID(user_id)
}

func (u *FinUseCaseImpl) GetFinNextID() (string, error) {
	return u.finrepo.GetFinNextID()
}
