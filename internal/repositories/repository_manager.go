package repositories

import (
	"gorm.io/gorm"
)

//Centralized repo

type GormRepositoryManager struct {
	db                 *gorm.DB
	borrowerRepository BorrowerRepository
	loanRepository     LoanRepository
	scheduleRepository ScheduleRepository
	paymentRepository  PaymentRepository
}

func NewGormRepositoryManager(db *gorm.DB) *GormRepositoryManager {
	return &GormRepositoryManager{
		db:                 db,
		borrowerRepository: NewGormBorrowerRepository(db),
		loanRepository:     NewGormLoanRepository(db),
		scheduleRepository: NewGormScheduleRepository(db),
		paymentRepository:  NewGormPaymentRepository(db),
	}
}

// Borrowers returns the borrower repository
func (r *GormRepositoryManager) Borrowers() BorrowerRepository {
	return r.borrowerRepository
}

// Loans returns the loan repository
func (r *GormRepositoryManager) Loans() LoanRepository {
	return r.loanRepository
}

// Schedules returns the schedule repository
func (r *GormRepositoryManager) Schedules() ScheduleRepository {
	return r.scheduleRepository
}

// Payments returns the payment repository
func (r *GormRepositoryManager) Payments() PaymentRepository {
	return r.paymentRepository
}

// WithTransaction runs a function within a database transaction
func (r *GormRepositoryManager) WithTransaction(fn func(repo RepositoryManager) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create a new repository manager with the transaction
		txRepo := &GormRepositoryManager{
			db:                 tx,
			borrowerRepository: NewGormBorrowerRepository(tx),
			loanRepository:     NewGormLoanRepository(tx),
			scheduleRepository: NewGormScheduleRepository(tx),
			paymentRepository:  NewGormPaymentRepository(tx),
		}

		// Run the provided function with the transaction-aware repository
		return fn(txRepo)
	})
}
