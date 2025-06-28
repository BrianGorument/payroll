package payroll_periods


type CreatePayrollPeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type PayrollPeriodResponse struct {
	ID          string    `json:"id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	IsProcessed bool      `json:"is_processed"`
}