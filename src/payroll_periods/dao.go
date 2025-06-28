package payroll_periods

import (
	"time"
)

type PayrollPeriod struct {
	ID           string           `gorm:"primaryKey;autoIncrement" json:"id"`
	StartDate    time.Time       `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time       `gorm:"type:date;not null" json:"end_date"`
	IsProcessed  bool            `gorm:"not null;default:false" json:"is_processed"`
	CreatedAt    time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy    string          `gorm:"not null" json:"created_by"`
	UpdatedBy    string          `gorm:"not null" json:"updated_by"`
}

func (PayrollPeriod) TableName() string {
	return "payroll_periods"
}