package repositories

import (
	"loan-billing-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormScheduleRepository struct {
	db *gorm.DB
}

func NewGormScheduleRepository(db *gorm.DB) *GormScheduleRepository {
	return &GormScheduleRepository{db: db}
}

// GetByID retrieves a schedule by ID
func (r *GormScheduleRepository) GetByID(id uuid.UUID) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// GetByLoanID retrieves schedules by loan ID
func (r *GormScheduleRepository) GetByLoanID(loanID uuid.UUID) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.Where("loan_id = ?", loanID).Order("week_number").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetUnpaidByLoanID retrieves unpaid schedules by loan ID
func (r *GormScheduleRepository) GetUnpaidByLoanID(loanID uuid.UUID) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.Where("loan_id = ? AND paid = ?", loanID, false).Order("week_number").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// Create creates a new schedule
func (r *GormScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

// CreateBatch creates multiple schedules in one go
func (r *GormScheduleRepository) CreateBatch(schedules []models.Schedule) error {
	return r.db.Create(&schedules).Error
}

// UpdatePaidStatus updates a schedule's paid status
func (r *GormScheduleRepository) UpdatePaidStatus(id uuid.UUID, paid bool) error {
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).
		Update("paid", paid).Error
}

// CountUnpaidByLoanID counts the number of unpaid schedules for a loan
func (r *GormScheduleRepository) CountUnpaidByLoanID(loanID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Schedule{}).Where("loan_id = ? AND paid = ?", loanID, false).Count(&count).Error
	return count, err
}
