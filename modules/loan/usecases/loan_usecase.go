package usecases

import (
	"errors"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/repositories"
	transRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/repositories"
	"github.com/google/uuid"
)

type LoanUseCase interface {
	CreateLoan(loan entities.Loan) (*entities.Loan, error)
	GetLoanByID(id string) (*entities.Loan, error)
	GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error)
	UpdateLoanStatusByID(id string, loan entities.Loan) (*entities.Loan, error)
	DeleteLoanByID(id string) error
}

type LoanUseCaseImpl struct {
	loanrepo  repositories.LoanRepository
	transrepo transRepo.TransRepository
}

func NewLoanUseCase(loanrepo repositories.LoanRepository, transrepo transRepo.TransRepository) *LoanUseCaseImpl {
	return &LoanUseCaseImpl{
		loanrepo:  loanrepo,
		transrepo: transrepo,
	}
}

func (u *LoanUseCaseImpl) CreateLoan(loan entities.Loan) (*entities.Loan, error) {
	if loan.MonthlyExpenses <= 0 {
		return nil, errors.New("monthly expense must be greater than zero")
	}

	if loan.RemainingMonths <= 0 {
		return nil, errors.New("remaining months must be greater than zero")
	}

	loan.ID = uuid.New().String()
	if loan.Installment {
		loan.Status = "In_Progress"
	} else {
		loan.Status = "Paused"
	}

	createdLoan, err := u.loanrepo.CreateLoan(&loan)
	if err != nil {
		return nil, err
	}

	status := "ชำระ"
	if loan.Status == "Paused" {
		status = "หยุดพัก"
	}

	transaction := &entities.Transaction{
		ID:        uuid.New().String(),
		Status:    status,
		UserID:    loan.UserID,
		LoanID:    loan.ID,
		CreatedAt: time.Now(),
	}

	if err := u.transrepo.CreateTransaction(transaction); err != nil {
		u.loanrepo.DeleteLoanByID(loan.ID)
		return nil, err
	}
	return createdLoan, nil
}

func (u *LoanUseCaseImpl) GetLoanByID(id string) (*entities.Loan, error) {
	return u.loanrepo.GetLoanByID(id)
}

func (u *LoanUseCaseImpl) GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error) {
	return u.loanrepo.GetLoanByUserID(userID)
}

func (u *LoanUseCaseImpl) UpdateLoanStatusByID(id string, loan entities.Loan) (*entities.Loan, error) {
	existingLoan, err := u.loanrepo.GetLoanByID(id)
	if err != nil {
		return nil, err
	}

	existingLoan.Name = loan.Name
	installmentChangedToFalse := existingLoan.Installment && !loan.Installment
	installmentChangedToTrue := !existingLoan.Installment && loan.Installment
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

	latestTransaction, err := u.transrepo.GetLatestTransactionByLoanID(id)
	if err == nil {
		if installmentChangedToFalse && latestTransaction.Status == "ชำระ" {
			latestTransaction.Status = "หยุดพัก"
			if err := u.transrepo.UpdateTransaction(latestTransaction); err != nil {
				return nil, err
			}
		}

		if installmentChangedToTrue && latestTransaction.Status == "หยุดพัก" {
			latestTransaction.Status = "ชำระ"
			if err := u.transrepo.UpdateTransaction(latestTransaction); err != nil {
				return nil, err
			}
		}
	}

	updatedLoan, err := u.loanrepo.UpdateLoanByID(existingLoan)
	if err != nil {
		return nil, err
	}

	return updatedLoan, nil

}

func (u *LoanUseCaseImpl) DeleteLoanByID(id string) error {
	if err := u.transrepo.DeleteTransactionsByLoanID(id); err != nil {
		return err
	}

	if err := u.loanrepo.DeleteLoanByID(id); err != nil {
		return err
	}

	return nil
}
