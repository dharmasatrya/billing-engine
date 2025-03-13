package handlers

import (
	"net/http"

	"loan-billing-system/internal/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// BorrowerHandler handles HTTP requests related to borrowers
type BorrowerHandler struct {
	borrowerService *services.BorrowerService
}

// NewBorrowerHandler creates a new borrower handler
func NewBorrowerHandler(borrowerService *services.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{
		borrowerService: borrowerService,
	}
}

// CreateBorrowerRequest represents the request body for creating a borrower
// @Description Request body for creating a new borrower
type CreateBorrowerRequest struct {
	Name        string `json:"name" validate:"required"`
	ContactInfo string `json:"contact_info" validate:"required"`
}

// BorrowerResponse represents the borrower data in responses
// @Description Response containing borrower data
type BorrowerResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ContactInfo  string    `json:"contact_info"`
	IsDelinquent bool      `json:"is_delinquent"`
}

// CreateBorrower godoc
// @Summary Create a new borrower
// @Description Creates a new borrower in the system
// @Tags Borrowers
// @Accept json
// @Produce json
// @Param request body handlers.CreateBorrowerRequest true "Borrower details"
// @Success 201 {object} handlers.BorrowerResponse
// @Failure 400 {object} map[string]string "Error response"
// @Failure 500 {object} map[string]string "Error response"
// @Router /api/borrowers [post]
func (h *BorrowerHandler) CreateBorrower(c echo.Context) error {
	var req CreateBorrowerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	borrower, err := h.borrowerService.CreateBorrower(req.Name, req.ContactInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, BorrowerResponse{
		ID:           borrower.ID,
		Name:         borrower.Name,
		ContactInfo:  borrower.ContactInfo,
		IsDelinquent: borrower.IsDelinquent,
	})
}

// GetBorrower godoc
// @Summary Get borrower details
// @Description Retrieves details for a specific borrower
// @Tags Borrowers
// @Accept json
// @Produce json
// @Param id path string true "Borrower ID" format(uuid)
// @Success 200 {object} handlers.BorrowerResponse
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /api/borrowers/{id} [get]
func (h *BorrowerHandler) GetBorrower(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid borrower ID format"})
	}

	borrower, err := h.borrowerService.GetBorrower(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Borrower not found"})
	}

	return c.JSON(http.StatusOK, BorrowerResponse{
		ID:           borrower.ID,
		Name:         borrower.Name,
		ContactInfo:  borrower.ContactInfo,
		IsDelinquent: borrower.IsDelinquent,
	})
}

// ListBorrowers godoc
// @Summary List all borrowers
// @Description Retrieves a list of all borrowers in the system
// @Tags Borrowers
// @Accept json
// @Produce json
// @Success 200 {array} handlers.BorrowerResponse
// @Failure 500 {object} map[string]string "Error response"
// @Router /api/borrowers [get]
func (h *BorrowerHandler) ListBorrowers(c echo.Context) error {
	borrowers, err := h.borrowerService.ListBorrowers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var response []BorrowerResponse
	for _, borrower := range borrowers {
		response = append(response, BorrowerResponse{
			ID:           borrower.ID,
			Name:         borrower.Name,
			ContactInfo:  borrower.ContactInfo,
			IsDelinquent: borrower.IsDelinquent,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// ListDelinquentBorrowers godoc
// @Summary List delinquent borrowers
// @Description Retrieves a list of all borrowers who are currently delinquent
// @Tags Borrowers
// @Accept json
// @Produce json
// @Success 200 {array} handlers.BorrowerResponse
// @Failure 500 {object} map[string]string "Error response"
// @Router /api/borrowers/delinquent [get]
func (h *BorrowerHandler) ListDelinquentBorrowers(c echo.Context) error {
	borrowers, err := h.borrowerService.GetDelinquentBorrowers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var response []BorrowerResponse
	for _, borrower := range borrowers {
		response = append(response, BorrowerResponse{
			ID:           borrower.ID,
			Name:         borrower.Name,
			ContactInfo:  borrower.ContactInfo,
			IsDelinquent: borrower.IsDelinquent,
		})
	}

	return c.JSON(http.StatusOK, response)
}
