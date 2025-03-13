package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Borrower represents a person who borrows money
type Borrower struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	ContactInfo  string         `gorm:"size:255" json:"contact_info"`
	IsDelinquent bool           `gorm:"default:false" json:"is_delinquent"`
	Loans        []Loan         `gorm:"foreignKey:BorrowerID" json:"loans,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
