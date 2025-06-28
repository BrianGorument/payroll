package attendances

import "time"

type CreateAttendanceRequest struct {
	Action string `json:"action" binding:"required"`
}

type AttendanceResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	PayrollPeriodID string  `json:"payroll_period_id"`
	CheckIn         string  `json:"check_in"`
	CheckOut        *string `json:"check_out"`
}

type PayrollPeriod struct {
	ID          string    `gorm:"primaryKey;autoIncrement"`
	StartDate   time.Time `gorm:"type:date;not null"`
	EndDate     time.Time `gorm:"type:date;not null"`
	IsProcessed bool      `gorm:"not null;default:false"`
}