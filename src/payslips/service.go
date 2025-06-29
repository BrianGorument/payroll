package payslips

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
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

    for _, user := range users {
        // Hitung gaji berdasarkan absensi
        attendances, err := s.repo.FindAttendanceByUserAndPeriod(strconv.Itoa(user.ID), req.PayrollPeriodID)
        if err != nil {
            return nil, err
        }
        var validAttendanceDays int
        for _, att := range attendances {
            if !att.CheckOut.IsZero() {
                checkInDate := att.CheckIn.Truncate(24 * time.Hour)
                checkOutDate := att.CheckOut.Truncate(24 * time.Hour)
                if checkInDate.Equal(checkOutDate) {
                    checkIn := att.CheckIn.Truncate(time.Minute)
                    checkOut := att.CheckOut.Truncate(time.Minute)
                    durationHours := checkOut.Sub(checkIn).Hours()
                    if durationHours >= 8 {
                        validAttendanceDays++
                    }
                }
            }
        }
        dailySalary := user.Salary / 20.0
        salaryBaseOnAttended := float64(validAttendanceDays) * dailySalary

        // Hitung overtime_pay: 2 x gaji prorata per jam
        overtimes, err := s.repo.FindOvertimeByUserAndPeriod(strconv.Itoa(user.ID), req.PayrollPeriodID)
        if err != nil {
            return nil, err
        }
        var totalOvertimeHours int
        for _, ot := range overtimes {
            totalOvertimeHours += ot.Hours
        }
        overtimeRate := (user.Salary / 160.0) * 2
        overtimePay := float64(totalOvertimeHours) * overtimeRate

        // Hitung reimbursement_pay
        reimbursements, err := s.repo.FindReimbursementByUserAndPeriod(strconv.Itoa(user.ID), req.PayrollPeriodID)
        if err != nil {
            return nil, err
        }
        var reimbursementPay float64
        for _, r := range reimbursements {
            reimbursementPay += r.Amount
        }

        totalPay := salaryBaseOnAttended + overtimePay + reimbursementPay


        payslip := Payslip{
            UserID:          strconv.Itoa(user.ID),
            PayrollPeriodID: req.PayrollPeriodID,
            BaseSalary:      user.Salary,
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
            ID:                   payslip.ID,
            UserID:               payslip.UserID,
            PayrollPeriodID:      payslip.PayrollPeriodID,
            BaseSalary:           FormatRupiah(payslip.BaseSalary),
            SalaryBaseOnAttended: FormatRupiah(salaryBaseOnAttended),
            OvertimePay:          FormatRupiah(payslip.OvertimePay),
            ReimbursementPay:     FormatRupiah(payslip.ReimbursementPay),
            TotalPay:             FormatRupiah(payslip.TotalPay),
        })
    }

    period.IsProcessed = true
    if err := s.repo.UpdatePayrollPeriod(period); err != nil {
        return nil, err
    }

    return payslips, nil
}

func (s *payslipService) GeneratePayslip(req GeneratePayslipRequest, userID, role string) (string, error) {
    // Validasi payroll_period_id
    _, err := strconv.Atoi(req.PayrollPeriodID)
    if err != nil {
        return "", errors.New("invalid payroll_period_id format, must be a valid integer")
    }

    targetUserID := userID
    if role == "admin" {
        if req.UserID == "" {
            return "", errors.New("user_id is required for admin role")
        }
        _, err := strconv.Atoi(req.UserID)
        if err != nil {
            return "", errors.New("invalid user_id format, must be a valid integer")
        }
        targetUserID = req.UserID
    }

    // Cek apakah periode ada dan sudah diproses
    period, err := s.repo.FindPayrollPeriodByID(req.PayrollPeriodID)
    if err != nil {
        return "", err
    }
    if !period.IsProcessed {
        return "", errors.New("payroll period not processed")
    }

    // Ambil payslip
    payslip, err := s.repo.FindPayslipByUserAndPeriod(targetUserID, req.PayrollPeriodID)
    if err != nil {
        return "", err
    }

    // Ambil data pengguna
    user, err := s.repo.FindUserByID(targetUserID)
    if err != nil {
        return "", err
    }

    attendances, err := s.repo.FindAttendanceByUserAndPeriod(targetUserID, req.PayrollPeriodID)
    if err != nil {
        return "", err
    }

    overtimes, err := s.repo.FindOvertimeByUserAndPeriod(targetUserID, req.PayrollPeriodID)
    if err != nil {
        return "", err
    }

    reimbursements, err := s.repo.FindReimbursementByUserAndPeriod(targetUserID, req.PayrollPeriodID)
    if err != nil {
        return "", err
    }

    // Buat folder tmp/payslips jika belum ada
    outputDir := "tmp/payslips"
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return "", errors.New("failed to create output directory")
    }

    // Generate nama file PDF
    pdfFileName := fmt.Sprintf("payslip_%s_%s.pdf", targetUserID, req.PayrollPeriodID)
    pdfFilePath := filepath.Join(outputDir, pdfFileName)

    // Generate HTML template untuk PDF
    htmlContent, err := GeneratePayslipHTML(user, period, payslip, attendances, overtimes, reimbursements)
    if err != nil {
        return "", err
    }

    // Generate PDF dari HTML
    err = GeneratePDF(htmlContent, pdfFilePath)
    if err != nil {
        return "", err
    }

    return pdfFilePath, nil
}