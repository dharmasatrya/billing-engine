package repositories

import (
	"loan-billing-system/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormLoanRepository struct {
	db *gorm.DB
}

func NewGormLoanRepository(db *gorm.DB) *GormLoanRepository {
	return &GormLoanRepository{db: db}
}

// GetByID retrieves a loan by ID
func (r *GormLoanRepository) GetByID(id uuid.UUID) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.Preload("Borrower").Preload("Schedules").First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// GetByBorrowerID retrieves loans by borrower ID
func (r *GormLoanRepository) GetByBorrowerID(borrowerID uuid.UUID) ([]models.Loan, error) {
	var loans []models.Loan
	if err := r.db.Where("borrower_id = ?", borrowerID).Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

// GetAllActive retrieves all active loans
func (r *GormLoanRepository) GetAllActive() ([]models.Loan, error) {
	var loans []models.Loan
	if err := r.db.Where("status = ?", "active").Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

// Create creates a new loan
func (r *GormLoanRepository) Create(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

// Update updates a loan
func (r *GormLoanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

// UpdateStatus updates a loan's status
func (r *GormLoanRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Loan{}).Where("id = ?", id).
		Update("status", status).Error
}

// UpdateBalance updates the current balance of a loan
func (r *GormLoanRepository) UpdateBalance(id uuid.UUID, balance int64) error {
	return r.db.Model(&models.Loan{}).Where("id = ?", id).
		Update("current_balance", balance).Error
}

// UpdateLastPaymentDate updates the last payment date of a loan
func (r *GormLoanRepository) UpdateLastPaymentDate(id uuid.UUID, date time.Time) error {
	return r.db.Model(&models.Loan{}).Where("id = ?", id).
		Update("last_payment_date", date).Error
}

// GetPotentialDelinquent retrieves active loans that haven't received a payment recently
func (r *GormLoanRepository) GetPotentialDelinquent() ([]models.Loan, error) {
	var loans []models.Loan

	// Current date
	now := time.Now()

	// Setting cutoff date to 2 weeks ago
	// This means we only check loans that haven't had a payment in the last 2 weeks
	cutoffDate := now.AddDate(0, 0, -14)

	// Get loans that are active and either have no payments or
	// haven't had a payment since the cutoff date
	err := r.db.Where("status = ? AND (last_payment_date IS NULL OR last_payment_date < ?)",
		"active", cutoffDate).Find(&loans).Error

	return loans, err
}
