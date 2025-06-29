package payslips

import "time"

type RunPayrollRequest struct {
	PayrollPeriodID string `json:"payroll_period_id" binding:"required"`
}

type PayslipResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	PayrollPeriodID  string  `json:"payroll_period_id"`
	BaseSalary       float64 `json:"base_salary"`
	OvertimePay      float64 `json:"overtime_pay"`
	ReimbursementPay float64 `json:"reimbursement_pay"`
	TotalPay         float64 `json:"total_pay"`
}


type PayrollPeriod struct {
    ID          string    `gorm:"primaryKey;autoIncrement"`
    StartDate   time.Time `gorm:"type:date;not null"`
    EndDate     time.Time `gorm:"type:date;not null"`
    IsProcessed bool      `gorm:"not null;default:false"`
}

type User struct {
	ID         string  `gorm:"type:integer;primaryKey"`
	BaseSalary float64 `gorm:"type:decimal(15,2);not null"`
}

type Overtime struct {
	UserID       string    `gorm:"type:integer;not null"`
	Hours        int       `gorm:"type:integer;not null"`
	OvertimeDate time.Time `gorm:"type:date;not null"`
}

type Reimbursement struct {
	UserID string  `gorm:"type:integer;not null"`
	Amount float64 `gorm:"type:decimal(15,2);not null"`
}
