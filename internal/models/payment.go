package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment represents an actual payment made by a borrower
type Payment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LoanID      uuid.UUID      `gorm:"type:uuid;not null" json:"loan_id"`
	Loan        Loan           `gorm:"foreignKey:LoanID" json:"-"`
	ScheduleID  uuid.UUID      `gorm:"type:uuid;not null" json:"schedule_id"`
	Schedule    Schedule       `gorm:"foreignKey:ScheduleID" json:"-"`
	Amount      int64          `gorm:"not null" json:"amount"`
	PaymentDate time.Time      `gorm:"not null" json:"payment_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
