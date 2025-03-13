package repositories

import (
	"loan-billing-system/internal/models"
	"time"

	"github.com/google/uuid"
)

//Centralized interface for easier use

// BorrowerRepository defines the interface for borrower data access
type BorrowerRepository interface {
	GetByID(id uuid.UUID) (*models.Borrower, error)
	GetAll() ([]models.Borrower, error)
	GetDelinquent() ([]models.Borrower, error)
	Create(name, contactInfo string) (*models.Borrower, error)
	Update(borrower *models.Borrower) error
	UpdateDelinquencyStatus(id uuid.UUID, isDelinquent bool) error
}

// LoanRepository defines the interface for loan data access
type LoanRepository interface {
	GetByID(id uuid.UUID) (*models.Loan, error)
	GetByBorrowerID(borrowerID uuid.UUID) ([]models.Loan, error)
	GetAllActive() ([]models.Loan, error)
	Create(loan *models.Loan) error
	Update(loan *models.Loan) error
	UpdateStatus(id uuid.UUID, status string) error
	UpdateBalance(id uuid.UUID, balance int64) error
	UpdateLastPaymentDate(id uuid.UUID, date time.Time) error
	GetPotentialDelinquent() ([]models.Loan, error)
}

// ScheduleRepository defines the interface for schedule data access
type ScheduleRepository interface {
	GetByID(id uuid.UUID) (*models.Schedule, error)
	GetByLoanID(loanID uuid.UUID) ([]models.Schedule, error)
	GetUnpaidByLoanID(loanID uuid.UUID) ([]models.Schedule, error)
	Create(schedule *models.Schedule) error
	CreateBatch(schedules []models.Schedule) error
	UpdatePaidStatus(id uuid.UUID, paid bool) error
	CountUnpaidByLoanID(loanID uuid.UUID) (int64, error)
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	GetByID(id uuid.UUID) (*models.Payment, error)
	GetByLoanID(loanID uuid.UUID) ([]models.Payment, error)
	Create(payment *models.Payment) error
}

// RepositoryManager provides access to all repositories
type RepositoryManager interface {
	Borrowers() BorrowerRepository
	Loans() LoanRepository
	Schedules() ScheduleRepository
	Payments() PaymentRepository
	WithTransaction(fn func(repo RepositoryManager) error) error
}
