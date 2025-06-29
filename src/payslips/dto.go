package payslips


import (
	"time"
)

type RunPayrollRequest struct {
	PayrollPeriodID string `json:"payroll_period_id" binding:"required"`
}

type PayslipResponse struct {
    ID                    string `json:"id"`
    UserID                string `json:"user_id"`
    PayrollPeriodID       string `json:"payroll_period_id"`
    BaseSalary            string `json:"base_salary"`             
    SalaryBaseOnAttended  string `json:"salary_base_on_attended"` 
    OvertimePay           string `json:"overtime_pay"`           
    ReimbursementPay      string `json:"reimbursement_pay"`      
	TotalPay              string `json:"total_pay"`             
}


type PayrollPeriod struct {
    ID          string    `gorm:"primaryKey;autoIncrement"`
    StartDate   time.Time `gorm:"type:date;not null"`
    EndDate     time.Time `gorm:"type:date;not null"`
    IsProcessed bool      `gorm:"not null;default:false"`
}

type User struct {
	ID         string  `gorm:"type:integer;primaryKey"`
	Salary     float64 `gorm:"column:salary,type:decimal(15,2); not null" json:"salary"`
}

type Overtime struct {
	UserID       string    `gorm:"type:integer;not null"`
	Hours        int       `gorm:"type:integer;not null"`
	OvertimeDate time.Time `gorm:"type:date;not null"`
}

type Reimbursement struct {
	UserID string  `gorm:"type:integer;not null"`
	Amount float64 `gorm:"type:decimal(15,2);not null"`
	Description *string `gorm:"type:text column:description" json:"description"`
}

type GeneratePayslipRequest struct {
    PayrollPeriodID string `json:"payroll_period_id" binding:"required"`
    UserID          string `json:"user_id" binding:"required_if=Role admin"`
}

type Attendance struct {
    UserID    string    `gorm:"type:integer;not null"`
    CheckIn   time.Time `gorm:"type:timestamp;not null"`
    CheckOut  time.Time `gorm:"type:timestamp"`
}

// type User struct {
//    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`   
// 	Username  string    `gorm:"not null" json:"username"`
// 	Password  string    `gorm:"not null" json:"-"`
// 	Salary    float64   `gorm:"type:decimal(15,2)" json:"salary"`
// 	Role      string    `gorm:"type:user_role;not null" json:"role"`
// 	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
// 	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
// 	CreatedBy string    `gorm:"not null" json:"created_by"`
// 	UpdatedBy string    `gorm:"not null" json:"updated_by"`
// }