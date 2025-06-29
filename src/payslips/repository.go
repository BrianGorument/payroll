package payslips

import (
    "errors"

    "gorm.io/gorm"
)

type payslipRepository struct {
    db *gorm.DB
}

func NewPayslipRepository(db *gorm.DB) IPayslipRepository {
    return &payslipRepository{db}
}

func (r *payslipRepository) Create(payslip *Payslip) error {
    return r.db.Create(payslip).Error
}

func (r *payslipRepository) FindPayrollPeriodByID(periodID string) (*PayrollPeriod, error) {
    var period PayrollPeriod
    err := r.db.Table("payroll_periods").Where("id = ?", periodID).First(&period).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("payroll period not found")
        }
        return nil, err
    }
    return &period, nil
}

func (r *payslipRepository) FindUsersWithActivity(periodID string) ([]User, error) {
    var users []User
    err := r.db.Table("users").
        Where("id IN (?) OR id IN (?) OR id IN (?)",
            r.db.Table("attendances").Select("user_id").Where("payroll_period_id = ?", periodID),
            r.db.Table("overtimes").Select("user_id").Where("payroll_period_id = ?", periodID),
            r.db.Table("reimbursements").Select("user_id").Where("payroll_period_id = ?", periodID)).
        Find(&users).Error
    return users, err
}

func (r *payslipRepository) FindOvertimeByUserAndPeriod(userID, periodID string) ([]Overtime, error) {
    var overtimes []Overtime
    err := r.db.Table("overtimes").
        Where("user_id = ? AND payroll_period_id = ?", userID, periodID).
        Find(&overtimes).Error
    return overtimes, err
}

func (r *payslipRepository) FindReimbursementByUserAndPeriod(userID, periodID string) ([]Reimbursement, error) {
    var reimbursements []Reimbursement
    err := r.db.Table("reimbursements").
        Where("user_id = ? AND payroll_period_id = ?", userID, periodID).
        Find(&reimbursements).Error
    return reimbursements, err
}

func (r *payslipRepository) UpdatePayrollPeriod(period *PayrollPeriod) error {
    result := r.db.Table("payroll_periods").Where("id = ?", period.ID).Update("is_processed", true)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("no payroll period found to update")
    }
    return nil
}