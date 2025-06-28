package overtimes

import (
    "errors"
    "time"

    "gorm.io/gorm"
)

type overtimeRepository struct {
    db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) IOvertimeRepository {
    return &overtimeRepository{db}
}

func (r *overtimeRepository) Create(overtime *Overtime) error {
    return r.db.Create(overtime).Error
}

func (r *overtimeRepository) FindByUserAndDate(userID string, date time.Time) (*Overtime, error) {
    var overtime Overtime
    err := r.db.Where("user_id = ? AND overtime_date = ?", userID, date.Format("2006-01-02")).First(&overtime).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &overtime, nil
}

func (r *overtimeRepository) FindActivePayrollPeriod() (*PayrollPeriod, error) {
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
