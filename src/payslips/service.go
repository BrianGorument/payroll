package payslips

import (
    "errors"
    "strconv"
    "time"
)

type payslipService struct {
    repo IPayslipRepository
}

func NewPayslipService(repo IPayslipRepository) IPayslipService {
    return &payslipService{repo}
}

func (s *payslipService) RunPayroll(req RunPayrollRequest, adminID string) ([]PayslipResponse, error) {
    // Validasi adminID
    if adminID == "" {
        return nil, errors.New("admin ID cannot be empty")
    }
    _, err := strconv.Atoi(adminID)
    if err != nil {
        return nil, errors.New("invalid admin ID format, must be a valid integer")
    }

    // Validasi payroll_period_id
    _, err = strconv.Atoi(req.PayrollPeriodID)
    if err != nil {
        return nil, errors.New("invalid payroll_period_id format, must be a valid integer")
    }

    // Cek apakah periode ada dan belum diproses
    period, err := s.repo.FindPayrollPeriodByID(req.PayrollPeriodID)
    if err != nil {
        return nil, err
    }
    if period.IsProcessed {
        return nil, errors.New("payroll period already processed")
    }

    // Ambil waktu saat ini (WIB)
    wib, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return nil, errors.New("failed to load WIB timezone")
    }
    now := time.Now().In(wib)

    // Ambil semua pengguna dengan aktivitas
    users, err := s.repo.FindUsersWithActivity(req.PayrollPeriodID)
    if err != nil {
        return nil, err
    }
    if len(users) == 0 {
        return nil, errors.New("no users with activity found for this period")
    }

    var payslips []PayslipResponse
    const overtimeRate = 50000.0 // Rp50,000 per jam (asumsikan)

    for _, user := range users {
        // Hitung overtime_pay
        overtimes, err := s.repo.FindOvertimeByUserAndPeriod(user.ID, req.PayrollPeriodID)
        if err != nil {
            return nil, err
        }
        var totalOvertimeHours int
        for _, ot := range overtimes {
            totalOvertimeHours += ot.Hours
        }
        overtimePay := float64(totalOvertimeHours) * overtimeRate

        // Hitung reimbursement_pay
        reimbursements, err := s.repo.FindReimbursementByUserAndPeriod(user.ID, req.PayrollPeriodID)
        if err != nil {
            return nil, err
        }
        var reimbursementPay float64
        for _, r := range reimbursements {
            reimbursementPay += r.Amount
        }

        // Hitung total_pay
        totalPay := user.BaseSalary + overtimePay + reimbursementPay

        // Buat payslip
        payslip := Payslip{
            ID:              "", // Auto-increment by database
            UserID:          user.ID,
            PayrollPeriodID: req.PayrollPeriodID,
            BaseSalary:      user.BaseSalary,
            OvertimePay:     overtimePay,
            ReimbursementPay: reimbursementPay,
            TotalPay:        totalPay,
            CreatedAt:       now,
            UpdatedAt:       now,
            CreatedBy:       adminID,
            UpdatedBy:       adminID,
        }

        if err := s.repo.Create(&payslip); err != nil {
            return nil, err
        }

        payslips = append(payslips, PayslipResponse{
            ID:              payslip.ID,
            UserID:          payslip.UserID,
            PayrollPeriodID: payslip.PayrollPeriodID,
            BaseSalary:      payslip.BaseSalary,
            OvertimePay:     payslip.OvertimePay,
            ReimbursementPay: payslip.ReimbursementPay,
            TotalPay:        payslip.TotalPay,
        })
    }

    // Tandai periode sebagai diproses
    period.IsProcessed = true
    if err := s.repo.UpdatePayrollPeriod(period); err != nil {
        return nil, err
    }

    return payslips, nil
}