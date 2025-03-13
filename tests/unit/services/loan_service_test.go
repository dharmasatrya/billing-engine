package services_test

import (
	"errors"
	"loan-billing-system/internal/models"
	"loan-billing-system/internal/repositories"
	"loan-billing-system/internal/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock repositories
type MockRepoManager struct {
	mock.Mock
	borrowerRepo *MockBorrowerRepo
	loanRepo     *MockLoanRepo
	scheduleRepo *MockScheduleRepo
	paymentRepo  *MockPaymentRepo
}

func (m *MockRepoManager) Borrowers() repositories.BorrowerRepository {
	return m.borrowerRepo
}

func (m *MockRepoManager) Loans() repositories.LoanRepository {
	return m.loanRepo
}

func (m *MockRepoManager) Schedules() repositories.ScheduleRepository {
	return m.scheduleRepo
}

func (m *MockRepoManager) Payments() repositories.PaymentRepository {
	return m.paymentRepo
}

func (m *MockRepoManager) WithTransaction(fn func(repo repositories.RepositoryManager) error) error {
	args := m.Called(fn)
	if args.Get(0) == nil {
		err := fn(m)
		return err
	}
	return args.Error(0)
}

type MockBorrowerRepo struct {
	mock.Mock
}

func (m *MockBorrowerRepo) GetByID(id uuid.UUID) (*models.Borrower, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Borrower), args.Error(1)
}

func (m *MockBorrowerRepo) GetAll() ([]models.Borrower, error) {
	args := m.Called()
	return args.Get(0).([]models.Borrower), args.Error(1)
}

func (m *MockBorrowerRepo) GetDelinquent() ([]models.Borrower, error) {
	args := m.Called()
	return args.Get(0).([]models.Borrower), args.Error(1)
}

func (m *MockBorrowerRepo) Create(name, contactInfo string) (*models.Borrower, error) {
	args := m.Called(name, contactInfo)
	return args.Get(0).(*models.Borrower), args.Error(1)
}

func (m *MockBorrowerRepo) Update(borrower *models.Borrower) error {
	args := m.Called(borrower)
	return args.Error(0)
}

func (m *MockBorrowerRepo) UpdateDelinquencyStatus(id uuid.UUID, isDelinquent bool) error {
	args := m.Called(id, isDelinquent)
	return args.Error(0)
}

type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) GetByID(id uuid.UUID) (*models.Loan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepo) GetByBorrowerID(borrowerID uuid.UUID) ([]models.Loan, error) {
	args := m.Called(borrowerID)
	return args.Get(0).([]models.Loan), args.Error(1)
}

func (m *MockLoanRepo) GetAllActive() ([]models.Loan, error) {
	args := m.Called()
	return args.Get(0).([]models.Loan), args.Error(1)
}

func (m *MockLoanRepo) Create(loan *models.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

func (m *MockLoanRepo) Update(loan *models.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

func (m *MockLoanRepo) UpdateStatus(id uuid.UUID, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockLoanRepo) UpdateBalance(id uuid.UUID, balance int64) error {
	args := m.Called(id, balance)
	return args.Error(0)
}

func (m *MockLoanRepo) UpdateLastPaymentDate(id uuid.UUID, date time.Time) error {
	args := m.Called(id, date)
	return args.Error(0)
}

func (m *MockLoanRepo) GetPotentialDelinquent() ([]models.Loan, error) {
	args := m.Called()
	return args.Get(0).([]models.Loan), args.Error(1)
}

type MockScheduleRepo struct {
	mock.Mock
}

func (m *MockScheduleRepo) GetByID(id uuid.UUID) (*models.Schedule, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Schedule), args.Error(1)
}

func (m *MockScheduleRepo) GetByLoanID(loanID uuid.UUID) ([]models.Schedule, error) {
	args := m.Called(loanID)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleRepo) GetUnpaidByLoanID(loanID uuid.UUID) ([]models.Schedule, error) {
	args := m.Called(loanID)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleRepo) Create(schedule *models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockScheduleRepo) CreateBatch(schedules []models.Schedule) error {
	args := m.Called(schedules)
	return args.Error(0)
}

func (m *MockScheduleRepo) UpdatePaidStatus(id uuid.UUID, paid bool) error {
	args := m.Called(id, paid)
	return args.Error(0)
}

func (m *MockScheduleRepo) CountUnpaidByLoanID(loanID uuid.UUID) (int64, error) {
	args := m.Called(loanID)
	return args.Get(0).(int64), args.Error(1)
}

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) GetByID(id uuid.UUID) (*models.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByLoanID(loanID uuid.UUID) ([]models.Payment, error) {
	args := m.Called(loanID)
	return args.Get(0).([]models.Payment), args.Error(1)
}

func (m *MockPaymentRepo) Create(payment *models.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

// LoanServiceTestSuite defines the test suite for loan service
type LoanServiceTestSuite struct {
	suite.Suite
	service      *services.LoanService
	repoManager  *MockRepoManager
	borrowerRepo *MockBorrowerRepo
	loanRepo     *MockLoanRepo
	scheduleRepo *MockScheduleRepo
	paymentRepo  *MockPaymentRepo
}

// SetupTest prepares the test suite before each test
func (s *LoanServiceTestSuite) SetupTest() {
	s.borrowerRepo = new(MockBorrowerRepo)
	s.loanRepo = new(MockLoanRepo)
	s.scheduleRepo = new(MockScheduleRepo)
	s.paymentRepo = new(MockPaymentRepo)

	s.repoManager = &MockRepoManager{
		borrowerRepo: s.borrowerRepo,
		loanRepo:     s.loanRepo,
		scheduleRepo: s.scheduleRepo,
		paymentRepo:  s.paymentRepo,
	}

	s.service = services.NewLoanService(s.repoManager)
}

// TestCreateLoan tests the loan creation functionality
func (s *LoanServiceTestSuite) TestCreateLoan() {
	// Prepare test data
	borrowerID := uuid.New()
	borrower := &models.Borrower{
		ID:           borrowerID,
		Name:         "Test Borrower",
		ContactInfo:  "test@example.com",
		IsDelinquent: false,
	}

	// Setup expectations
	s.borrowerRepo.On("GetByID", borrowerID).Return(borrower, nil)
	s.loanRepo.On("Create", mock.AnythingOfType("*models.Loan")).Run(func(args mock.Arguments) {
		loan := args.Get(0).(*models.Loan)
		loan.ID = uuid.New() // Simulate database assigning an ID
	}).Return(nil)
	s.scheduleRepo.On("CreateBatch", mock.AnythingOfType("[]models.Schedule")).Return(nil)

	// Call the service
	loan, err := s.service.CreateLoan(borrowerID, 5000000, 10.0, 50)

	// Assert results
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), loan)
	assert.Equal(s.T(), borrowerID, loan.BorrowerID)
	assert.Equal(s.T(), int64(5000000), loan.Amount)
	assert.Equal(s.T(), 10.0, loan.InterestRate)
	assert.Equal(s.T(), uint(50), loan.TermWeeks)
	assert.Equal(s.T(), "active", loan.Status)
	assert.Equal(s.T(), int64(5480769), loan.CurrentBalance)

	// Verify mock expectations
	s.borrowerRepo.AssertExpectations(s.T())
	s.loanRepo.AssertExpectations(s.T())
	s.scheduleRepo.AssertExpectations(s.T())
}

// TestCreateLoanBorrowerNotFound tests loan creation with non-existent borrower
func (s *LoanServiceTestSuite) TestCreateLoanBorrowerNotFound() {
	// Prepare test data
	borrowerID := uuid.New()

	// Setup expectations
	s.borrowerRepo.On("GetByID", borrowerID).Return(nil, errors.New("borrower not found"))

	// Call the service
	loan, err := s.service.CreateLoan(borrowerID, 5000000, 10.0, 50)

	// Assert results
	assert.Error(s.T(), err)
	assert.Nil(s.T(), loan)
	assert.Contains(s.T(), err.Error(), "borrower not found")

	// Verify mock expectations
	s.borrowerRepo.AssertExpectations(s.T())
}

// TestIsDelinquent tests the delinquency checking
func (s *LoanServiceTestSuite) TestIsDelinquent() {
	// Prepare test data
	loanID := uuid.New()
	borrowerID := uuid.New()
	loan := &models.Loan{
		ID:         loanID,
		BorrowerID: borrowerID,
		Amount:     5000000,
		Status:     "active",
	}

	// Create schedules with 2 missed payments
	now := time.Now()
	schedules := []models.Schedule{
		{ID: uuid.New(), LoanID: loanID, WeekNumber: 1, DueDate: now.AddDate(0, 0, -21), Amount: 109615, Paid: true},
		{ID: uuid.New(), LoanID: loanID, WeekNumber: 2, DueDate: now.AddDate(0, 0, -14), Amount: 109615, Paid: false}, // Unpaid and past due
		{ID: uuid.New(), LoanID: loanID, WeekNumber: 3, DueDate: now.AddDate(0, 0, -7), Amount: 109615, Paid: false},  // Unpaid and past due
		{ID: uuid.New(), LoanID: loanID, WeekNumber: 4, DueDate: now.AddDate(0, 0, 0), Amount: 109615, Paid: false},   // Due today (not late yet)
	}

	// Setup expectations
	s.loanRepo.On("GetByID", loanID).Return(loan, nil)
	s.scheduleRepo.On("GetByLoanID", loanID).Return(schedules, nil)
	s.borrowerRepo.On("UpdateDelinquencyStatus", borrowerID, true).Return(nil)

	// Call the service
	isDelinquent, err := s.service.IsDelinquent(loanID)

	// Assert results
	assert.NoError(s.T(), err)
	assert.True(s.T(), isDelinquent)

	// Verify mock expectations
	s.loanRepo.AssertExpectations(s.T())
	s.scheduleRepo.AssertExpectations(s.T())
	s.borrowerRepo.AssertExpectations(s.T())
}

// TestMakePayment tests the payment processing
func (s *LoanServiceTestSuite) TestMakePayment() {
	// Prepare test data
	loanID := uuid.New()
	borrowerID := uuid.New()
	scheduleID := uuid.New()

	loan := &models.Loan{
		ID:             loanID,
		BorrowerID:     borrowerID,
		Amount:         5000000,
		Status:         "active",
		CurrentBalance: 5480769,
	}

	unpaidSchedule := &models.Schedule{
		ID:         scheduleID,
		LoanID:     loanID,
		WeekNumber: 1,
		DueDate:    time.Now().AddDate(0, 0, -7),
		Amount:     109615,
		Paid:       false,
	}

	// Setup expectations for transaction
	s.repoManager.On("WithTransaction", mock.AnythingOfType("func(repositories.RepositoryManager) error")).Return(nil)

	// Setup expectations inside the transaction
	s.loanRepo.On("GetByID", loanID).Return(loan, nil)
	s.scheduleRepo.On("GetUnpaidByLoanID", loanID).Return([]models.Schedule{*unpaidSchedule}, nil)
	s.paymentRepo.On("Create", mock.AnythingOfType("*models.Payment")).Return(nil)
	s.scheduleRepo.On("UpdatePaidStatus", scheduleID, true).Return(nil)
	s.loanRepo.On("UpdateBalance", loanID, int64(5371154)).Return(nil)
	s.loanRepo.On("UpdateLastPaymentDate", loanID, mock.AnythingOfType("time.Time")).Return(nil)
	s.scheduleRepo.On("CountUnpaidByLoanID", loanID).Return(int64(5), nil)
	s.scheduleRepo.On("GetByLoanID", loanID).Return([]models.Schedule{*unpaidSchedule}, nil)
	s.borrowerRepo.On("UpdateDelinquencyStatus", borrowerID, false).Return(nil)

	// Call the service
	err := s.service.MakePayment(loanID, 109615)

	// Assert results
	assert.NoError(s.T(), err)

	// Verify mock expectations
	s.repoManager.AssertExpectations(s.T())
	s.loanRepo.AssertExpectations(s.T())
	s.scheduleRepo.AssertExpectations(s.T())
	s.paymentRepo.AssertExpectations(s.T())
	s.borrowerRepo.AssertExpectations(s.T())
}

func TestLoanServiceSuite(t *testing.T) {
	suite.Run(t, new(LoanServiceTestSuite))
}

// TestGetPotentialDelinquentLoans tests fetching potentially delinquent loans
func (s *LoanServiceTestSuite) TestGetPotentialDelinquentLoans() {
	// Prepare test data
	potentialDelinquentLoans := []models.Loan{
		{
			ID:              uuid.New(),
			BorrowerID:      uuid.New(),
			Amount:          5000000,
			Status:          "active",
			CurrentBalance:  4500000,
			LastPaymentDate: nil, // No payment yet
		},
		{
			ID:              uuid.New(),
			BorrowerID:      uuid.New(),
			Amount:          3000000,
			Status:          "active",
			CurrentBalance:  2500000,
			LastPaymentDate: ptr(time.Now().AddDate(0, 0, -20)), // Payment 20 days ago
		},
	}

	// Setup expectations
	s.loanRepo.On("GetPotentialDelinquent").Return(potentialDelinquentLoans, nil)

	// Call the service
	loans, err := s.service.GetPotentialDelinquentLoans()

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(loans))
	assert.Equal(s.T(), potentialDelinquentLoans, loans)

	// Verify mock expectations
	s.loanRepo.AssertExpectations(s.T())
}

// Helper function to get pointer to time.Time
func ptr(t time.Time) *time.Time {
	return &t
}
