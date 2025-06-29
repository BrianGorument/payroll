package reimbursements

import "time"

type CreateReimbursementRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Description *string `json:"description"`
}

type ReimbursementResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	Amount          float64 `json:"amount"`
	Description     *string `json:"description"`
}

type PayrollPeriod struct {
	ID          string    `gorm:"primaryKey;autoIncrement"`
	StartDate   time.Time `gorm:"type:date;not null"`
	EndDate     time.Time `gorm:"type:date;not null"`
	IsProcessed bool      `gorm:"not null;default:false"`
}