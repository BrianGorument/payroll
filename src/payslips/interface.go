package payslips

import "payroll/src/users"

type IPayslipService interface {
	RunPayroll(req RunPayrollRequest, adminID string) ([]PayslipResponse, error)
	GeneratePayslip(req GeneratePayslipRequest, userID, role string) (string, error)
}

type IPayslipRepository interface {
	Create(payslip *Payslip) error
	FindPayrollPeriodByID(periodID string) (*PayrollPeriod, error)
	UpdatePayrollPeriod(period *PayrollPeriod) error
	FindUsersWithActivity(periodID string) ([]users.User, error)
	FindOvertimeByUserAndPeriod(userID, periodID string) ([]Overtime, error)
	FindReimbursementByUserAndPeriod(userID, periodID string) ([]Reimbursement, error)
	FindPayslipByUserAndPeriod(userID, periodID string) (*Payslip, error)
	FindAttendanceByUserAndPeriod(userID, periodID string) ([]Attendance, error)
	FindUserByID(userID string) (*User, error)
}
