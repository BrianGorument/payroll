package reimbursements

import (
    "errors"

    "gorm.io/gorm"
)

type reimbursementRepository struct {
    db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) IReimbursementRepository {
    return &reimbursementRepository{db}
}

func (r *reimbursementRepository) Create(reimbursement *Reimbursement) error {
    return r.db.Create(reimbursement).Error
}

func (r *reimbursementRepository) FindActivePayrollPeriod() (*PayrollPeriod, error) {
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