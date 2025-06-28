package attendances

import (
    "errors"
    "time"

    "gorm.io/gorm"
)

type attendanceRepository struct {
    db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) IAttendanceRepository {
    return &attendanceRepository{db}
}

func (r *attendanceRepository) Create(attendance *Attendance) error {
    return r.db.Save(attendance).Error
}

func (r *attendanceRepository) FindByUserAndDate(userID string, date time.Time) (*Attendance, error) {
    var attendance Attendance
    dateStart := date.Truncate(24 * time.Hour)
    err := r.db.Where("user_id = ? AND DATE(check_in) = ?", userID, dateStart).First(&attendance).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &attendance, nil
}

func (r *attendanceRepository) FindActivePayrollPeriod() (*PayrollPeriod, error) {
    var period PayrollPeriod
    err := r.db.Table("payroll_periods").Where("is_processed = ?", false).First(&period).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("no active payroll period found")
        }
        return nil, err
    }
    return &period, nil
}