package payroll_periods

import "time"

type IPayrollPeriodService interface {
    CreatePayrollPeriod(req CreatePayrollPeriodRequest, userID string) (*PayrollPeriodResponse, error)
}

type IPayrollPeriodRepository interface {
    Create(period *PayrollPeriod) error
    FindByDateRange(startDate, endDate time.Time) ([]PayrollPeriod, error)
    FindActivePeriods() ([]PayrollPeriod, error)
}