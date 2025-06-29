package reimbursements

import (
    "time"
)

type Reimbursement struct {
    ID             string    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID         string    `gorm:"type:integer;not null" json:"user_id"`
    PayrollPeriodID string   `gorm:"type:integer;not null" json:"payroll_period_id"`
    Amount         float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
    Description    *string   `gorm:"type:text" json:"description"`
    CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
    CreatedBy      string    `gorm:"type:integer;not null" json:"created_by"`
    UpdatedBy      string    `gorm:"type:integer;not null" json:"updated_by"`
}

func (Reimbursement) TableName() string {
    return "reimbursements"
}