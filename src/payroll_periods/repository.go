package payroll_periods

import (
    "time"

    "gorm.io/gorm"
)

type payrollPeriodRepository struct {
    db *gorm.DB
}

func NewPayrollPeriodRepository(db *gorm.DB) IPayrollPeriodRepository {
    return &payrollPeriodRepository{db: db}
}

func (r *payrollPeriodRepository) Create(period *PayrollPeriod) error {
    return r.db.Create(period).Error
}

func (r *payrollPeriodRepository) FindByDateRange(startDate, endDate time.Time) ([]PayrollPeriod, error) {
    var periods []PayrollPeriod
    err := r.db.Where("start_date <= ? AND end_date >= ?", endDate, startDate).Find(&periods).Error
    return periods, err
}

func (r *payrollPeriodRepository) FindActivePeriods() ([]PayrollPeriod, error) {
    var periods []PayrollPeriod
    err := r.db.Where("is_processed = ?", false).Find(&periods).Error
    return periods, err
}