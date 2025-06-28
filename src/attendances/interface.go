package attendances

import (
	"time"
)

type IAttendancesService interface {
    CreateAttendance(req CreateAttendanceRequest, userID string) (*AttendanceResponse, error)
    
}

type IAttendanceRepository interface {
    Create(attendance *Attendance) error
    FindByUserAndDate(userID string, date time.Time) (*Attendance, error)
    FindActivePayrollPeriod() (*PayrollPeriod, error)
}
