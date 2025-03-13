package api

import (
	"loan-billing-system/internal/api/handlers"
	"loan-billing-system/internal/api/middleware"
	"loan-billing-system/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CustomValidator is the request validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, db *gorm.DB, borrowerService *services.BorrowerService, loanService *services.LoanService) {
	// Setup validator and custom binder
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Binder = &middleware.UUIDBinder{DefaultBinder: echo.DefaultBinder{}}

	// Initialize handlers
	borrowerHandler := handlers.NewBorrowerHandler(borrowerService)
	loanHandler := handlers.NewLoanHandler(loanService)

	// API group
	api := e.Group("/api")

	// Borrower routes
	borrowers := api.Group("/borrowers")
	borrowers.POST("", borrowerHandler.CreateBorrower)
	borrowers.GET("", borrowerHandler.ListBorrowers)
	borrowers.GET("/:id", borrowerHandler.GetBorrower) //use this to check borrower delinquency status
	borrowers.GET("/delinquent", borrowerHandler.ListDelinquentBorrowers)

	// Loan routes
	loans := api.Group("/loans")
	loans.POST("", loanHandler.CreateLoan)
	loans.GET("/:id", loanHandler.GetLoan)
	loans.GET("/:id/outstanding", loanHandler.GetOutstanding)
	loans.GET("/:id/delinquent", loanHandler.IsDelinquent) //check delinquency in loan level
	loans.POST("/:id/payment", loanHandler.MakePayment)
}
