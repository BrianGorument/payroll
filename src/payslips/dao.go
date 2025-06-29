package payslips

import (
    "time"
)

type Payslip struct {
    ID              string    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID          string    `gorm:"type:integer;not null" json:"user_id"`
    PayrollPeriodID string    `gorm:"type:integer;not null" json:"payroll_period_id"`
    BaseSalary      float64   `gorm:"type:decimal(15,2);not null" json:"base_salary"`
    OvertimePay     float64   `gorm:"type:decimal(15,2);not null;default:0" json:"overtime_pay"`
    ReimbursementPay float64  `gorm:"type:decimal(15,2);not null;default:0" json:"reimbursement_pay"`
    TotalPay        float64   `gorm:"type:decimal(15,2);not null" json:"total_pay"`
    CreatedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
    CreatedBy       string    `gorm:"type:integer;not null" json:"created_by"`
    UpdatedBy       string    `gorm:"type:integer;not null" json:"updated_by"`
}

func (Payslip) TableName() string {
    return "payslips"
}