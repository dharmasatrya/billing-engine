package services

import (
	"errors"
	"loan-billing-system/internal/models"
	"loan-billing-system/internal/repositories"
	"time"

	"github.com/google/uuid"
)

// LoanService handles loan business logic
type LoanService struct {
	repos repositories.RepositoryManager
}

// NewLoanService creates a new loan service
func NewLoanService(repos repositories.RepositoryManager) *LoanService {
	return &LoanService{
		repos: repos,
	}
}

// GetLoan retrieves a loan by ID
func (s *LoanService) GetLoan(id uuid.UUID) (*models.Loan, error) {
	return s.repos.Loans().GetByID(id)
}

// CreateLoan creates a new loan and generates payment schedules
func (s *LoanService) CreateLoan(borrowerID uuid.UUID, amount int64, interestRate float64, termWeeks uint) (*models.Loan, error) {
	// Check if borrower exists
	_, err := s.repos.Borrowers().GetByID(borrowerID)
	if err != nil {
		return nil, errors.New("borrower not found")
	}

	// Calculate total with interest
	totalDue := calculateTotalDue(amount, interestRate, termWeeks)

	// Create new loan
	loan := models.Loan{
		BorrowerID:     borrowerID,
		Amount:         amount,
		InterestRate:   interestRate,
		TermWeeks:      termWeeks,
		StartDate:      time.Now(),
		Status:         "active",
		CurrentBalance: totalDue,
	}

	// Save the loan
	if err := s.repos.Loans().Create(&loan); err != nil {
		return nil, err
	}

	// Generate weekly payment schedule
	if err := s.generateLoanSchedule(&loan); err != nil {
		return nil, err
	}

	return &loan, nil
}

// Helper function to calculate total due with interest
func calculateTotalDue(principal int64, interestRate float64, termWeeks uint) int64 {
	interestFactor := interestRate / 100 * float64(termWeeks) / 52
	return int64(float64(principal) * (1 + interestFactor))
}

// generateLoanSchedule creates payment schedules for a loan
func (s *LoanService) generateLoanSchedule(loan *models.Loan) error {
	weeklyPayment := loan.CalculateWeeklyPayment()

	// Create 50 weekly schedules
	var schedules []models.Schedule
	for week := uint(1); week <= loan.TermWeeks; week++ {
		dueDate := loan.StartDate.AddDate(0, 0, int(week)*7)
		schedule := models.Schedule{
			LoanID:     loan.ID,
			WeekNumber: week,
			DueDate:    dueDate,
			Amount:     weeklyPayment,
			Paid:       false,
		}
		schedules = append(schedules, schedule)
	}

	// Save all schedules to database
	return s.repos.Schedules().CreateBatch(schedules)
}

// GetOutstanding returns the current outstanding balance for a loan
func (s *LoanService) GetOutstanding(loanID uuid.UUID) (int64, error) {

	loan, err := s.repos.Loans().GetByID(loanID)
	if err != nil {
		return 0, err
	}

	return loan.CurrentBalance, nil
}

// IsDelinquent checks if a loan is delinquent (2+ missed payments)
func (s *LoanService) IsDelinquent(loanID uuid.UUID) (bool, error) {
	// Get the loan with schedules
	loan, err := s.repos.Loans().GetByID(loanID)
	if err != nil {
		return false, err
	}

	// Get all schedules for this loan
	schedules, err := s.repos.Schedules().GetByLoanID(loanID)
	if err != nil {
		return false, err
	}

	// Get current date
	currentDate := time.Now()

	// Count consecutive unpaid schedules that are past due
	var consecutiveMissed int
	maxConsecutiveMissed := 0

	for i := 0; i < len(schedules); i++ {
		if !schedules[i].Paid && schedules[i].DueDate.Before(currentDate) {
			consecutiveMissed++
			if consecutiveMissed > maxConsecutiveMissed {
				maxConsecutiveMissed = consecutiveMissed
			}
		} else {
			consecutiveMissed = 0
		}
	}

	isDelinquent := maxConsecutiveMissed >= 2

	// Update borrower's delinquent status if needed
	if isDelinquent {
		if err := s.repos.Borrowers().UpdateDelinquencyStatus(loan.BorrowerID, true); err != nil {
			return isDelinquent, err
		}
	}

	return isDelinquent, nil
}

// MakePayment records a payment for a loan
func (s *LoanService) MakePayment(loanID uuid.UUID, amount int64) error {
	return s.repos.WithTransaction(func(repo repositories.RepositoryManager) error {
		// Get the loan
		loan, err := repo.Loans().GetByID(loanID)
		if err != nil {
			return err
		}

		// Find the earliest unpaid schedule
		unpaidSchedules, err := repo.Schedules().GetUnpaidByLoanID(loanID)
		if err != nil || len(unpaidSchedules) == 0 {
			return errors.New("no unpaid schedules found")
		}
		schedule := unpaidSchedules[0]

		// Check if payment amount matches the schedule amount
		if amount != schedule.Amount {
			return errors.New("payment amount must match the scheduled amount")
		}

		// Record the payment
		paymentDate := time.Now()
		payment := models.Payment{
			LoanID:      loanID,
			ScheduleID:  schedule.ID,
			Amount:      amount,
			PaymentDate: paymentDate,
		}
		if err := repo.Payments().Create(&payment); err != nil {
			return err
		}

		// Update the schedule as paid
		if err := repo.Schedules().UpdatePaidStatus(schedule.ID, true); err != nil {
			return err
		}

		// Update the current balance
		newBalance := loan.CurrentBalance - amount
		if err := repo.Loans().UpdateBalance(loanID, newBalance); err != nil {
			return err
		}

		// Update the last payment date
		if err := repo.Loans().UpdateLastPaymentDate(loanID, paymentDate); err != nil {
			return err
		}

		// Check if all schedules are paid
		unpaidCount, err := repo.Schedules().CountUnpaidByLoanID(loanID)
		if err != nil {
			return err
		}

		// If all schedules are paid, update loan status to closed
		if unpaidCount == 0 {
			if err := repo.Loans().UpdateStatus(loanID, "closed"); err != nil {
				return err
			}
		}

		// Re-check delinquency status
		schedules, err := repo.Schedules().GetByLoanID(loanID)
		if err != nil {
			return err
		}

		// Count consecutive unpaid schedules that are past due
		currentDate := time.Now()
		var consecutiveMissed int
		maxConsecutiveMissed := 0

		for i := 0; i < len(schedules); i++ {
			if !schedules[i].Paid && schedules[i].DueDate.Before(currentDate) {
				consecutiveMissed++
				if consecutiveMissed > maxConsecutiveMissed {
					maxConsecutiveMissed = consecutiveMissed
				}
			} else {
				consecutiveMissed = 0
			}
		}

		isDelinquent := maxConsecutiveMissed >= 2

		// Update borrower's delinquent status
		return repo.Borrowers().UpdateDelinquencyStatus(loan.BorrowerID, isDelinquent)
	})
}

// GetPotentialDelinquentLoans returns loans that haven't been paid recently
func (s *LoanService) GetPotentialDelinquentLoans() ([]models.Loan, error) {
	return s.repos.Loans().GetPotentialDelinquent()
}
