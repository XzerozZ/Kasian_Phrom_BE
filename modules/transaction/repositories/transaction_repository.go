package repositories

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

type GormTransRepository struct {
	db *gorm.DB
}

func NewGormTransRepository(db *gorm.DB) *GormTransRepository {
	return &GormTransRepository{db: db}
}

type TransRepository interface {
	CreateTransaction(transaction *entities.Transaction) error
	GetTransactionByID(id string) (*entities.Transaction, error)
	GetTransactionByUserID(userID string) ([]map[string]interface{}, error)
	GetTransactionByLoanIDs(loanIDs []string) ([]entities.Transaction, error)
	GetLatestTransactionByLoanID(loanID string) (*entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) error
	DeleteTransaction(id string) error
	DeleteTransactionsByLoanID(loanID string) error
	CountTransactionsByLoanID(loanID string) (int, error)
}

func (r *GormTransRepository) CreateTransaction(transaction *entities.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *GormTransRepository) GetTransactionByID(id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.Preload("Loan").Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *GormTransRepository) GetTransactionByUserID(userID string) ([]map[string]interface{}, error) {
	var transactions []entities.Transaction
	if err := r.db.Preload("Loan").Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, trans := range transactions {
		totalAmount := float64(trans.Loan.RemainingMonths) * trans.Loan.MonthlyExpenses
		transactionData := map[string]interface{}{
			"transaction_id": trans.ID,
			"status":         trans.Status,
			"created_at":     trans.CreatedAt,
			"loan": map[string]interface{}{
				"loan_id":             trans.Loan.ID,
				"name":                trans.Loan.Name,
				"type":                trans.Loan.Type,
				"monthly_expenses":    trans.Loan.MonthlyExpenses,
				"interest_percentage": trans.Loan.InterestPercentage,
				"remaining_months":    trans.Loan.RemainingMonths,
				"installment":         trans.Loan.Installment,
				"status":              trans.Loan.Status,
			},
			"total_amount": totalAmount,
		}

		result = append(result, transactionData)
	}

	return result, nil
}

func (r *GormTransRepository) GetTransactionByLoanIDs(loanIDs []string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := r.db.Preload("Loan").Where("loan_id IN ?", loanIDs).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *GormTransRepository) GetLatestTransactionByLoanID(loanID string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.Preload("Loan").Where("loan_id = ?", loanID).Order("created_at DESC").First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *GormTransRepository) UpdateTransaction(transaction *entities.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *GormTransRepository) DeleteTransaction(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Transaction{}).Error
}

func (r *GormTransRepository) CountTransactionsByLoanID(loanID string) (int, error) {
	var count int64
	if err := r.db.Model(&entities.Transaction{}).Where("loan_id = ?", loanID).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *GormTransRepository) DeleteTransactionsByLoanID(loanID string) error {
	return r.db.Where("loan_id = ?", loanID).Delete(&entities.Transaction{}).Error
}
