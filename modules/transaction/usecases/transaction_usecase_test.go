package usecases_test

import (
	"errors"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionsForAllUsers(t *testing.T) {
	t.Run("Success - Create transactions for loans", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		loans := []entities.Loan{
			{
				ID:              "loan1",
				UserID:          "user1",
				Name:            "Loan 1",
				Status:          "In_Progress",
				RemainingMonths: 5,
			},
			{
				ID:              "loan2",
				UserID:          "user2",
				Name:            "Loan 2",
				Status:          "Paused",
				RemainingMonths: 3,
			},
		}

		existingTransactions := []entities.Transaction{
			{
				ID:     "trans1",
				Status: "ชำระ",
				UserID: "user1",
				LoanID: "loan1",
				Loan:   loans[0],
			},
		}

		loanRepo.On("GetAllLoansByStatus", []string{"In_Progress", "Paused"}).Return(loans, nil)
		transRepo.On("GetTransactionByLoanIDs", []string{"loan1", "loan2"}).Return(existingTransactions, nil)
		transRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		notiRepo.On("CreateNotification", mock.AnythingOfType("*entities.Notification")).Return(nil)

		transRepo.On("CountTransactionsByLoanID", "loan1").Return(1, nil)
		transRepo.On("CountTransactionsByLoanID", "loan2").Return(0, nil)
		transRepo.On("CreateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.CreateTransactionsForAllUsers()

		assert.NoError(t, err)
		transRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		notiRepo.AssertExpectations(t)
	})

	t.Run("Failed - No loans found", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		loanRepo.On("GetAllLoansByStatus", []string{"In_Progress", "Paused"}).Return([]entities.Loan{}, nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.CreateTransactionsForAllUsers()

		assert.Error(t, err)
		assert.Equal(t, "no loans found for transaction creation", err.Error())
		loanRepo.AssertExpectations(t)
	})

	t.Run("Failed - Error getting loans", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		loanRepo.On("GetAllLoansByStatus", []string{"In_Progress", "Paused"}).Return([]entities.Loan{}, errors.New("db error"))

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.CreateTransactionsForAllUsers()

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		loanRepo.AssertExpectations(t)
	})

	t.Run("Success - Delete paid and paused transactions", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		loans := []entities.Loan{
			{
				ID:              "loan1",
				UserID:          "user1",
				Name:            "Loan 1",
				Status:          "In_Progress",
				RemainingMonths: 5,
			},
		}

		existingTransactions := []entities.Transaction{
			{
				ID:     "trans1",
				Status: "ชำระแล้ว",
				UserID: "user1",
				LoanID: "loan1",
				Loan:   loans[0],
			},
			{
				ID:     "trans2",
				Status: "หยุดพัก",
				UserID: "user1",
				LoanID: "loan1",
				Loan:   loans[0],
			},
		}

		loanRepo.On("GetAllLoansByStatus", []string{"In_Progress", "Paused"}).Return(loans, nil)
		transRepo.On("GetTransactionByLoanIDs", []string{"loan1"}).Return(existingTransactions, nil)
		transRepo.On("DeleteTransaction", "trans1").Return(nil)
		transRepo.On("DeleteTransaction", "trans2").Return(nil)
		transRepo.On("CountTransactionsByLoanID", "loan1").Return(0, nil)
		transRepo.On("CreateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.CreateTransactionsForAllUsers()

		assert.NoError(t, err)
		transRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
	})
}

func TestMarkTransactiontoPaid(t *testing.T) {
	t.Run("Success - Mark transaction as paid", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		transaction := &entities.Transaction{
			ID:     "trans1",
			Status: "ชำระ",
			UserID: "user1",
			LoanID: "loan1",
		}

		loan := &entities.Loan{
			ID:              "loan1",
			UserID:          "user1",
			Name:            "Loan 1",
			Status:          "In_Progress",
			RemainingMonths: 2,
		}

		transRepo.On("GetTransactionByID", "trans1").Return(transaction, nil)
		transRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		loanRepo.On("GetLoanByID", "loan1").Return(loan, nil)
		loanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(loan, nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.MarkTransactiontoPaid("trans1", "user1")

		assert.NoError(t, err)
		assert.Equal(t, "ชำระแล้ว", transaction.Status)
		assert.Equal(t, 1, loan.RemainingMonths)
		transRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
	})

	t.Run("Success - Mark transaction as paid and complete loan", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		transaction := &entities.Transaction{
			ID:     "trans1",
			Status: "ชำระ",
			UserID: "user1",
			LoanID: "loan1",
		}

		loan := &entities.Loan{
			ID:              "loan1",
			UserID:          "user1",
			Name:            "Loan 1",
			Status:          "In_Progress",
			RemainingMonths: 1,
		}

		transRepo.On("GetTransactionByID", "trans1").Return(transaction, nil)
		transRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		loanRepo.On("GetLoanByID", "loan1").Return(loan, nil)
		loanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(loan, nil)
		notiRepo.On("CreateNotification", mock.AnythingOfType("*entities.Notification")).Return(nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.MarkTransactiontoPaid("trans1", "user1")

		assert.NoError(t, err)
		assert.Equal(t, "ชำระแล้ว", transaction.Status)
		assert.Equal(t, 0, loan.RemainingMonths)
		assert.Equal(t, "Completed", loan.Status)
		transRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		notiRepo.AssertExpectations(t)
	})

	t.Run("Failed - Transaction not found", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		var nilTransaction *entities.Transaction = nil
		transRepo.On("GetTransactionByID", "trans1").Return(nilTransaction, errors.New("transaction not found"))

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.MarkTransactiontoPaid("trans1", "user1")

		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())
		transRepo.AssertExpectations(t)
	})

	t.Run("Failed - Transaction is paused", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		transaction := &entities.Transaction{
			ID:     "trans1",
			Status: "หยุดพัก",
			UserID: "user1",
			LoanID: "loan1",
		}

		transRepo.On("GetTransactionByID", "trans1").Return(transaction, nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		err := useCase.MarkTransactiontoPaid("trans1", "user1")

		assert.Error(t, err)
		assert.Equal(t, "transaction is not in a payable state", err.Error())
		transRepo.AssertExpectations(t)
	})
}

func TestGetTransactionByUserID(t *testing.T) {
	t.Run("Success - Get transactions by user ID", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		expectedTransactions := []map[string]interface{}{
			{
				"id":     "trans1",
				"status": "ชำระ",
				"loan":   map[string]interface{}{"name": "Loan 1"},
			},
			{
				"id":     "trans2",
				"status": "ชำระแล้ว",
				"loan":   map[string]interface{}{"name": "Loan 2"},
			},
		}

		transRepo.On("GetTransactionByUserID", "user1").Return(expectedTransactions, nil)

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		transactions, err := useCase.GetTransactionByUserID("user1")

		assert.NoError(t, err)
		assert.Equal(t, expectedTransactions, transactions)
		transRepo.AssertExpectations(t)
	})

	t.Run("Failed - Error getting transactions", func(t *testing.T) {
		transRepo := new(mocks.MockTransRepository)
		loanRepo := new(mocks.MockLoanRepository)
		notiRepo := new(mocks.MockNotiRepository)

		transRepo.On("GetTransactionByUserID", "user1").Return([]map[string]interface{}(nil), errors.New("db error"))

		useCase := usecases.NewTransactionUseCase(transRepo, loanRepo, notiRepo)
		transactions, err := useCase.GetTransactionByUserID("user1")

		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.Equal(t, "db error", err.Error())
		transRepo.AssertExpectations(t)
	})
}
