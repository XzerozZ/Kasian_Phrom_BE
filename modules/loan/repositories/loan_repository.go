package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormLoanRepository struct {
	db *gorm.DB
}

func NewGormLoanRepository(db *gorm.DB) *GormLoanRepository {
	return &GormLoanRepository{db: db}
}

type LoanRepository interface {
	CreateLoan(loan *entities.Loan) (*entities.Loan, error)
	GetLoanByID(id string) (*entities.Loan, error)
	GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error)
	GetAllLoansByStatus(statuses []string) ([]entities.Loan, error)
	UpdateLoanByID(loan *entities.Loan) (*entities.Loan, error)
	DeleteLoanByID(id string) error
}

func (r *GormLoanRepository) CreateLoan(loan *entities.Loan) (*entities.Loan, error) {
	if err := r.db.Create(&loan).Error; err != nil {
		return nil, err
	}

	return r.GetLoanByID(loan.ID)
}

func (r *GormLoanRepository) GetLoanByID(id string) (*entities.Loan, error) {
	var loan entities.Loan
	if err := r.db.First(&loan, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *GormLoanRepository) GetLoanByUserID(userID string) ([]entities.Loan, map[string]interface{}, error) {
	var loans []entities.Loan
	if err := r.db.Where("user_id = ? AND (status = ? OR status = ?)", userID, "In_Progress", "Paused").Find(&loans).Error; err != nil {
		return nil, nil, err
	}

	var totalLoan int
	var totalLoanAmount float64
	var totalTransactionAmount float64

	for _, loan := range loans {
		loanTotalAmount := float64(loan.RemainingMonths) * loan.MonthlyExpenses
		totalLoanAmount += loanTotalAmount
		totalLoan++

		var transactions []entities.Transaction
		if err := r.db.Where("loan_id = ? AND (status = ? OR status = ?)", loan.ID, "ชำระ", "ค้างชำระ").Find(&transactions).Error; err != nil {
			return nil, nil, err
		}

		for range transactions {
			totalTransactionAmount += loan.MonthlyExpenses
		}
	}

	loanSummary := map[string]interface{}{
		"total_loan":               totalLoan,
		"total_amount":             totalLoanAmount,
		"total_transaction_amount": totalTransactionAmount,
	}

	return loans, loanSummary, nil
}

func (r *GormLoanRepository) GetAllLoansByStatus(statuses []string) ([]entities.Loan, error) {
	var loans []entities.Loan
	if err := r.db.Where("status IN ?", statuses).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *GormLoanRepository) UpdateLoanByID(loan *entities.Loan) (*entities.Loan, error) {
	if err := r.db.Save(&loan).Error; err != nil {
		return nil, err
	}

	return r.GetLoanByID(loan.ID)
}

func (r *GormLoanRepository) DeleteLoanByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Loan{}).Error
}
