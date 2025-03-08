package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	loanRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/repositories"
	notiRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/socket"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/repositories"
	"github.com/google/uuid"
)

type TransactionUseCase interface {
	CreateTransactionsForAllUsers() error
	MarkTransactiontoPaid(id, userID string) error
	GetTransactionByUserID(userID string) ([]map[string]interface{}, error)
}

type TransactionUseCaseImpl struct {
	transrepo repositories.TransRepository
	loanrepo  loanRepo.LoanRepository
	notirepo  notiRepo.NotiRepository
}

func NewTransactionUseCase(transrepo repositories.TransRepository, loanrepo loanRepo.LoanRepository, notirepo notiRepo.NotiRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		transrepo: transrepo,
		loanrepo:  loanrepo,
		notirepo:  notirepo,
	}
}

func (u *TransactionUseCaseImpl) CreateTransactionsForAllUsers() error {
	loans, err := u.loanrepo.GetAllLoansByStatus([]string{"In_Progress", "Paused"})
	if err != nil {
		return err
	}

	if len(loans) == 0 {
		return errors.New("no loans found for transaction creation")
	}

	loanIDs := make([]string, len(loans))
	for i, loan := range loans {
		loanIDs[i] = loan.ID
	}

	existingTransactions, err := u.transrepo.GetTransactionByLoanIDs(loanIDs)
	if err != nil {
		return err
	}

	for _, trans := range existingTransactions {
		if trans.Status == "ชำระแล้ว" || trans.Status == "หยุดพัก" {
			if err := u.transrepo.DeleteTransaction(trans.ID); err != nil {
				return err
			}
		} else if trans.Status == "ชำระ" {
			trans.Status = "ค้างชำระ"
			if err := u.transrepo.UpdateTransaction(&trans); err != nil {
				return err
			}
		}
	}

	for _, loan := range loans {
		existingTransactions, err := u.transrepo.CountTransactionsByLoanID(loan.ID)
		if err != nil {
			return err
		}

		if existingTransactions >= loan.RemainingMonths {
			continue
		}

		transactionStatus := "ชำระ"
		if loan.Status == "Paused" {
			transactionStatus = "หยุดพัก"
		} else if loan.Status == "In_Progress" {
			transactionStatus = "ชำระ"
		}

		transaction := &entities.Transaction{
			ID:        uuid.New().String(),
			Status:    transactionStatus,
			UserID:    loan.UserID,
			LoanID:    loan.ID,
			CreatedAt: time.Now(),
		}

		if err := u.transrepo.CreateTransaction(transaction); err != nil {
			return err
		}

	}

	return nil
}

func (u *TransactionUseCaseImpl) MarkTransactiontoPaid(id, userID string) error {
	transaction, err := u.transrepo.GetTransactionByID(id)
	if err != nil {
		return err
	}

	if transaction.Status == "หยุดพัก" {
		return errors.New("transaction is not in a payable state")
	}

	transaction.Status = "ชำระแล้ว"
	if err := u.transrepo.UpdateTransaction(transaction); err != nil {
		return err
	}

	loan, err := u.loanrepo.GetLoanByID(transaction.LoanID)
	if err != nil {
		return err
	}

	if loan.RemainingMonths > 0 {
		loan.RemainingMonths--
		if loan.RemainingMonths == 0 {
			loan.Status = "Completed"
			notification := &entities.Notification{
				ID:        uuid.New().String(),
				UserID:    userID,
				Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ : '%s' ได้เสร็จสิ้นแล้ว", loan.Name),
				CreatedAt: time.Now(),
			}

			_ = u.notirepo.CreateNotification(notification)
			socket.SendNotificationToUser(userID, fmt.Sprintf("Notification: %s", notification.Message))
		}

		if _, err := u.loanrepo.UpdateLoanByID(loan); err != nil {
			return err
		}
	}

	return nil
}

func (u *TransactionUseCaseImpl) GetTransactionByUserID(userID string) ([]map[string]interface{}, error) {
	return u.transrepo.GetTransactionByUserID(userID)
}
