package repositories

import (
	"fmt"
	"strconv"

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
	GetLoanByUserID(userID string) ([]entities.Loan, error)
	GetLoanNextID() (string, error)
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

func (r *GormLoanRepository) GetLoanNextID() (string, error) {
	var maxID string
	if err := r.db.Model(&entities.Loan{}).Select("COALESCE(MAX(CAST(id AS INT)), 0)").Scan(&maxID).Error; err != nil {
		return "", err
	}

	maxIDInt := 0
	if maxID != "" {
		maxIDInt, _ = strconv.Atoi(maxID)
	}

	nextID := maxIDInt + 1
	formattedID := fmt.Sprintf("%05d", nextID)
	return formattedID, nil
}

func (r *GormLoanRepository) GetLoanByUserID(userID string) ([]entities.Loan, error) {
	var loans []entities.Loan
	if err := r.db.Where("user_id = ?", userID).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *GormLoanRepository) DeleteLoanByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Loan{}).Error
}
