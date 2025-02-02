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
	GetTransactionByUserID(userID string) ([]entities.Transaction, error)
	GetTransactionByLoanIDs(loanIDs []string) ([]entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) error
	DeleteTransaction(id string) error
}

func (r *GormTransRepository) CreateTransaction(transaction *entities.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *GormTransRepository) GetTransactionByID(id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *GormTransRepository) GetTransactionByUserID(userID string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := r.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *GormTransRepository) GetTransactionByLoanIDs(loanIDs []string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := r.db.Where("loan_id IN ?", loanIDs).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *GormTransRepository) UpdateTransaction(transaction *entities.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *GormTransRepository) DeleteTransaction(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Transaction{}).Error
}
