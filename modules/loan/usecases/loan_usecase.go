package usecases

import (
	"errors"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/repositories"
)

type LoanUseCase interface {
	CreateLoan(loan entities.Loan) (*entities.Loan, error)
	GetLoanByID(id string) (*entities.Loan, error)
	GetLoanByUserID(userID string) ([]entities.Loan, error)
	UpdateLoanStatusByID(id string, loan entities.Loan) (*entities.Loan, error)
	DeleteLoanByID(id string) error
}

type LoanUseCaseImpl struct {
	loanrepo repositories.LoanRepository
}

func NewLoanUseCase(loanrepo repositories.LoanRepository) *LoanUseCaseImpl {
	return &LoanUseCaseImpl{loanrepo: loanrepo}
}

func (u *LoanUseCaseImpl) CreateLoan(loan entities.Loan) (*entities.Loan, error) {
	id, err := u.loanrepo.GetLoanNextID()
	if err != nil {
		return nil, err
	}

	if loan.MonthlyExpenses <= 0 {
		return nil, errors.New("monthly expense must be greater than zero")
	}

	if loan.InterestPercentage <= 0 {
		return nil, errors.New("interest percentage must be greater than zero")
	}

	if loan.RemainingMonths <= 0 {
		return nil, errors.New("remaining months must be greater than zero")
	}

	loan.ID = id
	loan.Status = "In_Progress"
	createdLoan, err := u.loanrepo.CreateLoan(&loan)
	if err != nil {
		return nil, err
	}

	return createdLoan, nil
}

func (u *LoanUseCaseImpl) GetLoanByID(id string) (*entities.Loan, error) {
	return u.loanrepo.GetLoanByID(id)
}

func (u *LoanUseCaseImpl) GetLoanByUserID(userID string) ([]entities.Loan, error) {
	return u.loanrepo.GetLoanByUserID(userID)
}

func (u *LoanUseCaseImpl) UpdateLoanStatusByID(id string, loan entities.Loan) (*entities.Loan, error) {
	existingLoan, err := u.loanrepo.GetLoanByID(id)
	if err != nil {
		return nil, err
	}

	existingLoan.Installment = loan.Installment
	if existingLoan.MonthlyExpenses > 0 {
		if loan.Installment {
			existingLoan.Status = "In_Progress"
		} else {
			existingLoan.Status = "Paused"
		}
	} else {
		existingLoan.Status = "Completed"
		existingLoan.Installment = false
	}

	updatedLoan, err := u.loanrepo.UpdateLoanByID(existingLoan)
	if err != nil {
		return nil, err
	}

	return updatedLoan, nil

}

func (u *LoanUseCaseImpl) DeleteLoanByID(id string) error {
	return u.loanrepo.DeleteLoanByID(id)
}
