package repositories

import (
	"loan-billing-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormBorrowerRepository implements BorrowerRepository using GORM
type GormBorrowerRepository struct {
	db *gorm.DB
}

// NewGormBorrowerRepository creates a new GORM-based borrower repository
func NewGormBorrowerRepository(db *gorm.DB) *GormBorrowerRepository {
	return &GormBorrowerRepository{db: db}
}

// GetByID retrieves a borrower by ID
func (r *GormBorrowerRepository) GetByID(id uuid.UUID) (*models.Borrower, error) {
	var borrower models.Borrower
	if err := r.db.Preload("Loans").First(&borrower, id).Error; err != nil {
		return nil, err
	}
	return &borrower, nil
}

// GetAll retrieves all borrowers
func (r *GormBorrowerRepository) GetAll() ([]models.Borrower, error) {
	var borrowers []models.Borrower
	if err := r.db.Find(&borrowers).Error; err != nil {
		return nil, err
	}
	return borrowers, nil
}

// GetDelinquent retrieves all delinquent borrowers
func (r *GormBorrowerRepository) GetDelinquent() ([]models.Borrower, error) {
	var borrowers []models.Borrower
	if err := r.db.Where("is_delinquent = ?", true).Find(&borrowers).Error; err != nil {
		return nil, err
	}
	return borrowers, nil
}

// Create creates a new borrower
func (r *GormBorrowerRepository) Create(name, contactInfo string) (*models.Borrower, error) {
	borrower := models.Borrower{
		Name:         name,
		ContactInfo:  contactInfo,
		IsDelinquent: false,
	}
	if err := r.db.Create(&borrower).Error; err != nil {
		return nil, err
	}
	return &borrower, nil
}

// Update updates a borrower
func (r *GormBorrowerRepository) Update(borrower *models.Borrower) error {
	return r.db.Save(borrower).Error
}

// UpdateDelinquencyStatus updates a borrower's delinquency status
func (r *GormBorrowerRepository) UpdateDelinquencyStatus(id uuid.UUID, isDelinquent bool) error {
	return r.db.Model(&models.Borrower{}).Where("id = ?", id).
		Update("is_delinquent", isDelinquent).Error
}
