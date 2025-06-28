package overtimes

import (
    "time"
)

type Overtime struct {
    ID             string    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID         string    `gorm:"type:integer;not null" json:"user_id"`
    PayrollPeriodID string   `gorm:"type:integer;not null" json:"payroll_period_id"`
    OvertimeDate   time.Time `gorm:"type:date;not null" json:"overtime_date"`
    Hours          int       `gorm:"type:integer;not null" json:"hours"`
    CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
    CreatedBy      string    `gorm:"type:integer;not null" json:"created_by"`
    UpdatedBy      string    `gorm:"type:integer;not null" json:"updated_by"`
}

func (Overtime) TableName() string {
    return "overtimes"
}