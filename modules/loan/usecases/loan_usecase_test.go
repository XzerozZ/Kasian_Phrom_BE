package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/testing/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	t.Run("success with installment true", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		expectedLoan := loan
		expectedLoan.ID = mock.Anything
		expectedLoan.Status = "In_Progress"

		mockLoanRepo.On("CreateLoan", mock.AnythingOfType("*entities.Loan")).Return(&expectedLoan, nil)
		mockTransRepo.On("CreateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "In_Progress", result.Status)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("success with installment false", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     false,
		}

		expectedLoan := loan
		expectedLoan.ID = mock.Anything
		expectedLoan.Status = "Paused"

		mockLoanRepo.On("CreateLoan", mock.AnythingOfType("*entities.Loan")).Return(&expectedLoan, nil)
		mockTransRepo.On("CreateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Paused", result.Status)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("fail with zero monthly expenses", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 0,
			RemainingMonths: 12,
			Installment:     true,
		}

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "monthly expense must be greater than zero", err.Error())
		mockLoanRepo.AssertNotCalled(t, "CreateLoan")
		mockTransRepo.AssertNotCalled(t, "CreateTransaction")
	})

	t.Run("fail with zero remaining months", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 0,
			Installment:     true,
		}

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "remaining months must be greater than zero", err.Error())
		mockLoanRepo.AssertNotCalled(t, "CreateLoan")
		mockTransRepo.AssertNotCalled(t, "CreateTransaction")
	})

	t.Run("fail if loan repository returns error", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
		}

		expectedError := errors.New("database error")
		mockLoanRepo.On("CreateLoan", mock.AnythingOfType("*entities.Loan")).Return(&loan, expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertNotCalled(t, "CreateTransaction")
	})

	t.Run("fail if transaction repository returns error", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loan := entities.Loan{
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
		}

		expectedLoan := loan
		expectedLoan.ID = "test-id"
		expectedLoan.Status = "In_Progress"

		expectedError := errors.New("transaction error")

		mockLoanRepo.On("CreateLoan", mock.AnythingOfType("*entities.Loan")).Return(&expectedLoan, nil)
		mockTransRepo.On("CreateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(expectedError)
		mockLoanRepo.On("DeleteLoanByID", mock.AnythingOfType("string")).Return(nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.CreateLoan(loan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})
}

func TestGetLoanByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		expectedLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		mockLoanRepo.On("GetLoanByID", loanID).Return(expectedLoan, nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.GetLoanByID(loanID)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoan, result)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestGetLoanByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		userID := "user-123"
		expectedLoans := []entities.Loan{
			{
				ID:              "loan-1",
				UserID:          userID,
				Name:            "Test Loan 1",
				MonthlyExpenses: 10000,
				RemainingMonths: 12,
				Installment:     true,
				Status:          "In_Progress",
			},
			{
				ID:              "loan-2",
				UserID:          userID,
				Name:            "Test Loan 2",
				MonthlyExpenses: 5000,
				RemainingMonths: 6,
				Installment:     false,
				Status:          "Paused",
			},
		}

		expectedMeta := map[string]interface{}{
			"total": 2,
		}

		mockLoanRepo.On("GetLoanByUserID", userID).Return(expectedLoans, expectedMeta, nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, meta, err := useCase.GetLoanByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoans, result)
		assert.Equal(t, expectedMeta, meta)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("fail if repository returns error", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		userID := "user-123"
		expectedError := errors.New("database error")

		var emptyLoans []entities.Loan
		emptyMeta := make(map[string]interface{})

		mockLoanRepo.On("GetLoanByUserID", userID).Return(emptyLoans, emptyMeta, expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, meta, err := useCase.GetLoanByUserID(userID)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Empty(t, meta)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestUpdateLoanStatusByID(t *testing.T) {
	t.Run("success with installment change to true", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		existingLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     false,
			Status:          "Paused",
		}

		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: true,
		}

		expectedLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Updated Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		latestTransaction := &entities.Transaction{
			ID:        "trans-123",
			LoanID:    loanID,
			Status:    "หยุดพัก",
			CreatedAt: time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", loanID).Return(existingLoan, nil)
		mockTransRepo.On("GetLatestTransactionByLoanID", loanID).Return(latestTransaction, nil)
		mockTransRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		mockLoanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(expectedLoan, nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoan, result)
		assert.Equal(t, "ชำระ", latestTransaction.Status)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("success with zero monthly expenses", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		existingLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 0,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: true,
		}

		expectedLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Updated Loan",
			MonthlyExpenses: 0,
			RemainingMonths: 12,
			Installment:     false,
			Status:          "Completed",
		}

		mockLoanRepo.On("GetLoanByID", loanID).Return(existingLoan, nil)

		var nilTransaction *entities.Transaction
		mockTransRepo.On("GetLatestTransactionByLoanID", loanID).Return(nilTransaction, errors.New("transaction not found"))

		mockLoanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(expectedLoan, nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoan, result)
		assert.Equal(t, "Completed", result.Status)
		assert.False(t, result.Installment)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("fail if loan not found", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "non-existent-loan"
		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: true,
		}

		var nilLoan *entities.Loan
		mockLoanRepo.On("GetLoanByID", loanID).Return(nilLoan, errors.New("loan not found"))

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "loan not found", err.Error())
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("fail if transaction update fails", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		existingLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     false,
			Status:          "Paused",
		}

		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: true,
		}

		latestTransaction := &entities.Transaction{
			ID:        "trans-123",
			LoanID:    loanID,
			Status:    "หยุดพัก",
			CreatedAt: time.Now(),
		}

		expectedError := errors.New("transaction update failed")

		mockLoanRepo.On("GetLoanByID", loanID).Return(existingLoan, nil)
		mockTransRepo.On("GetLatestTransactionByLoanID", loanID).Return(latestTransaction, nil)
		mockTransRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
		mockLoanRepo.AssertNotCalled(t, "UpdateLoanByID")
	})

	t.Run("fail if loan update fails", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		existingLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: false,
		}

		latestTransaction := &entities.Transaction{
			ID:        "trans-123",
			LoanID:    loanID,
			Status:    "ชำระ",
			CreatedAt: time.Now(),
		}

		expectedError := errors.New("loan update failed")

		mockLoanRepo.On("GetLoanByID", loanID).Return(existingLoan, nil)
		mockTransRepo.On("GetLatestTransactionByLoanID", loanID).Return(latestTransaction, nil)
		mockTransRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		var nilLoan *entities.Loan
		mockLoanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(nilLoan, expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("success with installment change to false", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		existingLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Test Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     true,
			Status:          "In_Progress",
		}

		updateLoan := entities.Loan{
			Name:        "Updated Loan",
			Installment: false,
		}

		expectedLoan := &entities.Loan{
			ID:              loanID,
			UserID:          "user-123",
			Name:            "Updated Loan",
			MonthlyExpenses: 10000,
			RemainingMonths: 12,
			Installment:     false,
			Status:          "Paused",
		}

		latestTransaction := &entities.Transaction{
			ID:        "trans-123",
			LoanID:    loanID,
			Status:    "ชำระ",
			CreatedAt: time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", loanID).Return(existingLoan, nil)
		mockTransRepo.On("GetLatestTransactionByLoanID", loanID).Return(latestTransaction, nil)
		mockTransRepo.On("UpdateTransaction", mock.AnythingOfType("*entities.Transaction")).Return(nil)
		mockLoanRepo.On("UpdateLoanByID", mock.AnythingOfType("*entities.Loan")).Return(expectedLoan, nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		result, err := useCase.UpdateLoanStatusByID(loanID, updateLoan)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoan, result)
		assert.Equal(t, "หยุดพัก", latestTransaction.Status)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})
}

func TestDeleteLoanByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"

		mockTransRepo.On("DeleteTransactionsByLoanID", loanID).Return(nil)
		mockLoanRepo.On("DeleteLoanByID", loanID).Return(nil)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		err := useCase.DeleteLoanByID(loanID)

		assert.NoError(t, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})

	t.Run("fail if transaction deletion fails", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		expectedError := errors.New("transaction deletion failed")

		mockTransRepo.On("DeleteTransactionsByLoanID", loanID).Return(expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		err := useCase.DeleteLoanByID(loanID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockTransRepo.AssertExpectations(t)
		mockLoanRepo.AssertNotCalled(t, "DeleteLoanByID")
	})

	t.Run("fail if loan deletion fails", func(t *testing.T) {
		mockLoanRepo := new(mocks.MockLoanRepository)
		mockTransRepo := new(mocks.MockTransRepository)

		loanID := "loan-123"
		expectedError := errors.New("loan deletion failed")

		mockTransRepo.On("DeleteTransactionsByLoanID", loanID).Return(nil)
		mockLoanRepo.On("DeleteLoanByID", loanID).Return(expectedError)

		useCase := usecases.NewLoanUseCase(mockLoanRepo, mockTransRepo)
		err := useCase.DeleteLoanByID(loanID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockLoanRepo.AssertExpectations(t)
		mockTransRepo.AssertExpectations(t)
	})
}
