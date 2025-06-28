package attendances

import (
    "time"
)

type Attendance struct {
    ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID         string     `gorm:"type:integer;not null" json:"user_id"`
    PayrollPeriodID string    `gorm:"type:integer;not null" json:"payroll_period_id"`
    CheckIn        time.Time  `gorm:"type:timestamp;not null" json:"check_in"`
    CheckOut       *time.Time `gorm:"type:timestamp" json:"check_out"`
    CreatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
    CreatedBy      string     `gorm:"type:integer;not null" json:"created_by"`
    UpdatedBy      string     `gorm:"type:integer;not null" json:"updated_by"`
}


func (Attendance) TableName() string {
	return "attendances"
}
