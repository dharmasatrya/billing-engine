package handlers

import (
	"net/http"
	"time"

	"loan-billing-system/internal/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// LoanHandler handles HTTP requests related to loans
type LoanHandler struct {
	loanService *services.LoanService
}

// NewLoanHandler creates a new loan handler
func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

// CreateLoanRequest represents the request body for creating a loan
// @Description Request body for creating a new loan
type CreateLoanRequest struct {
	BorrowerID   uuid.UUID `json:"borrower_id" validate:"required"`
	Amount       int64     `json:"amount" validate:"required,min=1"`
	InterestRate float64   `json:"interest_rate" validate:"required,min=0"`
	TermWeeks    uint      `json:"term_weeks" validate:"required,min=1"`
}

// LoanResponse represents the loan data in responses
// @Description Response containing loan data
type LoanResponse struct {
	ID           uuid.UUID `json:"id"`
	BorrowerID   uuid.UUID `json:"borrower_id"`
	Amount       int64     `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	TermWeeks    uint      `json:"term_weeks"`
	StartDate    time.Time `json:"start_date"`
	Status       string    `json:"status"`
}

// PaymentRequest represents the request body for making a payment
// @Description Request body for making a payment
type PaymentRequest struct {
	Amount int64 `json:"amount" validate:"required,min=1"`
}

// CreateLoan godoc
// @Summary Create a new loan
// @Description Creates a new loan for a borrower
// @Tags Loans
// @Accept json
// @Produce json
// @Param request body handlers.CreateLoanRequest true "Loan details"
// @Success 201 {object} handlers.LoanResponse
// @Failure 400 {object} map[string]string "Error response"
// @Failure 500 {object} map[string]string "Error response"
// @Router /api/loans [post]
func (h *LoanHandler) CreateLoan(c echo.Context) error {
	var req CreateLoanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	loan, err := h.loanService.CreateLoan(req.BorrowerID, req.Amount, req.InterestRate, req.TermWeeks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		TermWeeks:    loan.TermWeeks,
		StartDate:    loan.StartDate,
		Status:       loan.Status,
	})
}

// GetLoan godoc
// @Summary Get loan details
// @Description Retrieves details for a specific loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID" format(uuid)
// @Success 200 {object} handlers.LoanResponse
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /api/loans/{id} [get]
func (h *LoanHandler) GetLoan(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID format"})
	}

	loan, err := h.loanService.GetLoan(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Loan not found"})
	}

	return c.JSON(http.StatusOK, LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		TermWeeks:    loan.TermWeeks,
		StartDate:    loan.StartDate,
		Status:       loan.Status,
	})
}

// GetOutstanding godoc
// @Summary Get outstanding balance
// @Description Retrieves the current outstanding balance for a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID" format(uuid)
// @Success 200 {object} map[string]int64 "Outstanding amount"
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /api/loans/{id}/outstanding [get]
func (h *LoanHandler) GetOutstanding(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID format"})
	}

	outstanding, err := h.loanService.GetOutstanding(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Loan not found"})
	}

	return c.JSON(http.StatusOK, map[string]int64{"outstanding": outstanding})
}

// IsDelinquent godoc
// @Summary Check if loan is delinquent
// @Description Checks if a loan is currently delinquent (2+ missed payments)
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID" format(uuid)
// @Success 200 {object} map[string]bool "Delinquency status"
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /api/loans/{id}/delinquent [get]
func (h *LoanHandler) IsDelinquent(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID format"})
	}

	isDelinquent, err := h.loanService.IsDelinquent(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Loan not found"})
	}

	return c.JSON(http.StatusOK, map[string]bool{"is_delinquent": isDelinquent})
}

// MakePayment godoc
// @Summary Make a payment
// @Description Makes a payment for a loan
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Loan ID" format(uuid)
// @Param request body handlers.PaymentRequest true "Payment details"
// @Success 200 {object} map[string]string "Success response"
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /api/loans/{id}/payment [post]
func (h *LoanHandler) MakePayment(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID format"})
	}

	var req PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.loanService.MakePayment(id, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Payment successful"})
}
