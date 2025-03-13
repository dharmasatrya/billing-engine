package scheduler

import (
	"loan-billing-system/internal/services"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Scheduler struct {
	cron        *cron.Cron
	db          *gorm.DB
	loanService *services.LoanService
}

func NewScheduler(db *gorm.DB, loanService *services.LoanService) *Scheduler {
	return &Scheduler{
		cron:        cron.New(),
		db:          db,
		loanService: loanService,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	// Run delinquency check daily at midnight
	s.cron.AddFunc("0 0 * * *", s.checkDelinquency)
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// checkDelinquency checks all active loans for delinquency
func (s *Scheduler) checkDelinquency() {
	log.Println("Running delinquency check...")
	startTime := time.Now()

	// Get potentially delinquent loans (haven't been paid in the last 2 weeks)
	potentialDelinquentLoans, err := s.loanService.GetPotentialDelinquentLoans()
	if err != nil {
		log.Printf("Error fetching potential delinquent loans: %v", err)
		return
	}

	log.Printf("Found %d loans that may be delinquent", len(potentialDelinquentLoans))

	// Process loans in batches
	batchSize := 500
	totalLoans := len(potentialDelinquentLoans)
	delinquentCount := 0

	for i := 0; i < totalLoans; i += batchSize {
		end := i + batchSize
		if end > totalLoans {
			end = totalLoans
		}

		batch := potentialDelinquentLoans[i:end]
		log.Printf("Processing batch %d to %d of %d", i, end-1, totalLoans)

		for _, loan := range batch {
			isDelinquent, err := s.loanService.IsDelinquent(loan.ID)
			if err != nil {
				log.Printf("Error checking delinquency for loan %s: %v", loan.ID, err)
				continue
			}

			if isDelinquent {
				delinquentCount++
				log.Printf("Loan %s marked as delinquent", loan.ID)
			}
		}

		// Small delay between batches to reduce database load
		if end < totalLoans {
			time.Sleep(100 * time.Millisecond)
		}
	}

	duration := time.Since(startTime)
	log.Printf("Delinquency check completed in %v", duration)
	log.Printf("Processed %d loans, found %d delinquent", totalLoans, delinquentCount)
}

// RunNow runs the delinquency check immediately
func (s *Scheduler) RunNow() {
	s.checkDelinquency()
}
