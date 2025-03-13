package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Loan represents a loan issued to a borrower
type Loan struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BorrowerID      uuid.UUID      `gorm:"type:uuid;not null" json:"borrower_id"`
	Borrower        Borrower       `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
	Amount          int64          `gorm:"not null" json:"amount"`
	InterestRate    float64        `gorm:"not null" json:"interest_rate"`
	TermWeeks       uint           `gorm:"not null" json:"term_weeks"`
	StartDate       time.Time      `gorm:"not null" json:"start_date"`
	Status          string         `gorm:"size:20;not null;default:'active'" json:"status"`
	CurrentBalance  int64          `gorm:"not null" json:"current_balance"`
	LastPaymentDate *time.Time     `json:"last_payment_date"` // Date of last payment for query optimization
	Schedules       []Schedule     `gorm:"foreignKey:LoanID" json:"schedules,omitempty"`
	Payments        []Payment      `gorm:"foreignKey:LoanID" json:"payments,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// CalculateTotalDue returns the total amount due including interest
func (l *Loan) CalculateTotalDue() int64 {
	// Calculate total with 10% annual interest for the loan period
	interestFactor := l.InterestRate / 100 * float64(l.TermWeeks) / 52
	return int64(float64(l.Amount) * (1 + interestFactor))
}

// CalculateWeeklyPayment returns the weekly payment amount
func (l *Loan) CalculateWeeklyPayment() int64 {
	totalDue := l.CalculateTotalDue()
	return totalDue / int64(l.TermWeeks)
}
