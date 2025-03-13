package services

import (
	"loan-billing-system/internal/models"
	"loan-billing-system/internal/repositories"

	"github.com/google/uuid"
)

// BorrowerService handles borrower business logic
type BorrowerService struct {
	repos repositories.RepositoryManager
}

// NewBorrowerService creates a new borrower service
func NewBorrowerService(repos repositories.RepositoryManager) *BorrowerService {
	return &BorrowerService{repos: repos}
}

// GetBorrower retrieves a borrower by ID
func (s *BorrowerService) GetBorrower(id uuid.UUID) (*models.Borrower, error) {
	return s.repos.Borrowers().GetByID(id)
}

// ListBorrowers retrieves all borrowers
func (s *BorrowerService) ListBorrowers() ([]models.Borrower, error) {
	return s.repos.Borrowers().GetAll()
}

// CreateBorrower creates a new borrower
func (s *BorrowerService) CreateBorrower(name, contactInfo string) (*models.Borrower, error) {
	return s.repos.Borrowers().Create(name, contactInfo)
}

// GetDelinquentBorrowers retrieves all delinquent borrowers
func (s *BorrowerService) GetDelinquentBorrowers() ([]models.Borrower, error) {
	return s.repos.Borrowers().GetDelinquent()
}
