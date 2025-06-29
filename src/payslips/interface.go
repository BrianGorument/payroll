package payslips



type IPayslipService interface {
    RunPayroll(req RunPayrollRequest, adminID string) ([]PayslipResponse, error)
}

type IPayslipRepository interface {
    Create(payslip *Payslip) error
    FindPayrollPeriodByID(periodID string) (*PayrollPeriod, error)
    UpdatePayrollPeriod(period *PayrollPeriod) error
    FindUsersWithActivity(periodID string) ([]User, error)
    FindOvertimeByUserAndPeriod(userID, periodID string) ([]Overtime, error)
    FindReimbursementByUserAndPeriod(userID, periodID string) ([]Reimbursement, error)
}

