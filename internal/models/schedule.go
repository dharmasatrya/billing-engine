package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule represents a weekly payment schedule for a loan
type Schedule struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LoanID     uuid.UUID      `gorm:"type:uuid;not null" json:"loan_id"`
	Loan       Loan           `gorm:"foreignKey:LoanID" json:"-"`
	WeekNumber uint           `gorm:"not null" json:"week_number"`
	DueDate    time.Time      `gorm:"not null" json:"due_date"`
	Amount     int64          `gorm:"not null" json:"amount"`
	Paid       bool           `gorm:"default:false" json:"paid"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
