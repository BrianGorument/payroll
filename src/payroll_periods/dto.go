package payroll_periods

import "time"

type CreatePayrollPeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type PayrollPeriodResponse struct {
	ID          string    `json:"id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsProcessed bool      `json:"is_processed"`
}