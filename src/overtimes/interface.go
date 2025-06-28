package overtimes

import (
    "time"
)

type IOvertimeService interface {
    CreateOvertime(req CreateOvertimeRequest, userID string) (*OvertimeResponse, error)
}

type IOvertimeRepository interface {
    Create(overtime *Overtime) error
    FindByUserAndDate(userID string, date time.Time) (*Overtime, error)
    FindActivePayrollPeriod() (*PayrollPeriod, error)
}
