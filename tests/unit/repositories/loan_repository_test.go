package scheduler_test

import (
	"loan-billing-system/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// LoanServiceInterface defines the interface for the loan service
type LoanServiceInterface interface {
	GetLoan(id uuid.UUID) (*models.Loan, error)
	CreateLoan(borrowerID uuid.UUID, amount int64, interestRate float64, termWeeks uint) (*models.Loan, error)
	GetOutstanding(loanID uuid.UUID) (int64, error)
	IsDelinquent(loanID uuid.UUID) (bool, error)
	MakePayment(loanID uuid.UUID, amount int64) error
}

// MockLoanService mocks the loan service for scheduler testing
type MockLoanService struct {
	mock.Mock
}

func (m *MockLoanService) GetLoan(id uuid.UUID) (*models.Loan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanService) CreateLoan(borrowerID uuid.UUID, amount int64, interestRate float64, termWeeks uint) (*models.Loan, error) {
	args := m.Called(borrowerID, amount, interestRate, termWeeks)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanService) GetOutstanding(loanID uuid.UUID) (int64, error) {
	args := m.Called(loanID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLoanService) IsDelinquent(loanID uuid.UUID) (bool, error) {
	args := m.Called(loanID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLoanService) MakePayment(loanID uuid.UUID, amount int64) error {
	args := m.Called(loanID, amount)
	return args.Error(0)
}

// TestScheduler is a simplified scheduler for testing
type TestScheduler struct {
	DB          *gorm.DB
	LoanService LoanServiceInterface
}

// NewTestScheduler creates a new test scheduler
func NewTestScheduler(db *gorm.DB, loanService LoanServiceInterface) *TestScheduler {
	return &TestScheduler{
		DB:          db,
		LoanService: loanService,
	}
}

// RunNow simulates running the scheduler's delinquency check
func (s *TestScheduler) RunNow() {
	var loans []models.Loan
	s.DB.Where("status = ?", "active").Find(&loans)

	for _, loan := range loans {
		s.LoanService.IsDelinquent(loan.ID)
	}
}

// SchedulerTestSuite defines the test suite for scheduler
type SchedulerTestSuite struct {
	suite.Suite
	DB          *gorm.DB
	LoanService *MockLoanService
	Scheduler   *TestScheduler
}

// SetupSuite prepares the test suite before any tests run
func (s *SchedulerTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing with specific configs for UUID compatibility
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal(err)
	}

	// For SQLite tests, we'll modify the models slightly to work with SQLite's limitations
	// by creating tables manually instead of using AutoMigrate with UUID default values

	// Create tables manually for SQLite compatibility
	db.Exec(`CREATE TABLE borrowers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		contact_info TEXT,
		is_delinquent BOOLEAN DEFAULT false,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`)

	db.Exec(`CREATE TABLE loans (
		id TEXT PRIMARY KEY,
		borrower_id TEXT NOT NULL,
		amount INTEGER NOT NULL,
		interest_rate REAL NOT NULL,
		term_weeks INTEGER NOT NULL,
		start_date DATETIME NOT NULL,
		status TEXT NOT NULL DEFAULT 'active',
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY (borrower_id) REFERENCES borrowers(id)
	)`)

	db.Exec(`CREATE TABLE schedules (
		id TEXT PRIMARY KEY,
		loan_id TEXT NOT NULL,
		week_number INTEGER NOT NULL,
		due_date DATETIME NOT NULL,
		amount INTEGER NOT NULL,
		paid BOOLEAN DEFAULT false,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY (loan_id) REFERENCES loans(id)
	)`)

	db.Exec(`CREATE TABLE payments (
		id TEXT PRIMARY KEY,
		loan_id TEXT NOT NULL,
		schedule_id TEXT NOT NULL,
		amount INTEGER NOT NULL,
		payment_date DATETIME NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY (loan_id) REFERENCES loans(id),
		FOREIGN KEY (schedule_id) REFERENCES schedules(id)
	)`)

	s.DB = db
	s.LoanService = new(MockLoanService)
	s.Scheduler = NewTestScheduler(db, s.LoanService)
}

// TearDownTest cleans up after each test
func (s *SchedulerTestSuite) TearDownTest() {
	// Clean up the database after each test
	s.DB.Exec("DELETE FROM payments")
	s.DB.Exec("DELETE FROM schedules")
	s.DB.Exec("DELETE FROM loans")
	s.DB.Exec("DELETE FROM borrowers")
}

// TestCheckDelinquency tests the delinquency check functionality
func (s *SchedulerTestSuite) TestCheckDelinquency() {
	// Create test borrowers with manual UUID assignment
	borrower1ID := uuid.New()
	borrower2ID := uuid.New()

	s.DB.Exec("INSERT INTO borrowers (id, name, contact_info, is_delinquent, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		borrower1ID.String(),
		"Borrower 1",
		"borrower1@example.com",
		false,
		time.Now(),
		time.Now(),
	)

	s.DB.Exec("INSERT INTO borrowers (id, name, contact_info, is_delinquent, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		borrower2ID.String(),
		"Borrower 2",
		"borrower2@example.com",
		false,
		time.Now(),
		time.Now(),
	)

	// Create test loans
	loan1ID := uuid.New()
	loan2ID := uuid.New()

	s.DB.Exec("INSERT INTO loans (id, borrower_id, amount, interest_rate, term_weeks, start_date, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		loan1ID.String(),
		borrower1ID.String(),
		5000000,
		10.0,
		50,
		time.Now().AddDate(0, 0, -10*7), // 10 weeks ago
		"active",
		time.Now(),
		time.Now(),
	)

	s.DB.Exec("INSERT INTO loans (id, borrower_id, amount, interest_rate, term_weeks, start_date, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		loan2ID.String(),
		borrower2ID.String(),
		3000000,
		10.0,
		50,
		time.Now().AddDate(0, 0, -5*7), // 5 weeks ago
		"active",
		time.Now(),
		time.Now(),
	)

	// Setup expectations
	s.LoanService.On("IsDelinquent", loan1ID).Return(true, nil)
	s.LoanService.On("IsDelinquent", loan2ID).Return(false, nil)

	// Run the delinquency check
	s.Scheduler.RunNow()

	// Verify that the service was called with the correct loans
	s.LoanService.AssertExpectations(s.T())
}

func TestSchedulerSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
