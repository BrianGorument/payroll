package overtimes

import "time"

type CreateOvertimeRequest struct {
	OvertimeDate string `json:"overtime_date" binding:"required"`
	Hours        int    `json:"hours" binding:"required,min=1,max=3"`
}

type OvertimeResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	PayrollPeriodID string `json:"payroll_period_id"`
	OvertimeDate    string `json:"overtime_date"`
	Hours           int    `json:"hours"`
}

type PayrollPeriod struct {
	ID          string    `gorm:"primaryKey;autoIncrement"`
	StartDate   time.Time `gorm:"type:date;not null"`
	EndDate     time.Time `gorm:"type:date;not null"`
	IsProcessed bool      `gorm:"not null;default:false"`
}