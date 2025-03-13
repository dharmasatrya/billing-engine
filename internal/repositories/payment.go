package repositories

import (
	"loan-billing-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormPaymentRepository struct {
	db *gorm.DB
}

func NewGormPaymentRepository(db *gorm.DB) *GormPaymentRepository {
	return &GormPaymentRepository{db: db}
}

// GetByID retrieves a payment by ID
func (r *GormPaymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByLoanID retrieves payments by loan ID
func (r *GormPaymentRepository) GetByLoanID(loanID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("loan_id = ?", loanID).Order("payment_date").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// Create creates a new payment
func (r *GormPaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}
