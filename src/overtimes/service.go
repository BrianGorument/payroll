package overtimes

import (
    "errors"
    "strconv"
    "time"
)

type overtimeService struct {
    repo IOvertimeRepository
}

func NewOvertimeService(repo IOvertimeRepository) IOvertimeService {
    return &overtimeService{repo}
}

func (s *overtimeService) CreateOvertime(req CreateOvertimeRequest, userID string) (*OvertimeResponse, error) {
    // Validasi userID
    if userID == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    _, err := strconv.Atoi(userID)
    if err != nil {
        return nil, errors.New("invalid user ID format, must be a valid integer")
    }

    // Parse overtime_date
    overtimeDate, err := time.Parse("02/01/2006", req.OvertimeDate)
    if err != nil {
        return nil, errors.New("invalid overtime_date format, use DD/MM/YYYY")
    }

    // Validasi: Harus ada periode aktif
    period, err := s.repo.FindActivePayrollPeriod()
    if err != nil {
        return nil, err
    }

    // Validasi: Lembur harus diajukan setelah jam kerja (17:00 WIB)
    wib, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return nil, errors.New("failed to load WIB timezone")
    }
    now := time.Now().In(wib)
    workEndTime, _ := time.ParseInLocation("15:04:05", "17:00:00", wib)
    workEndToday := time.Date(now.Year(), now.Month(), now.Day(), workEndTime.Hour(), workEndTime.Minute(), workEndTime.Second(), 0, wib)
	
    if now.Before(workEndToday) && (time.Weekday(now.Weekday()) != time.Saturday) && (time.Weekday(now.Weekday()) != time.Sunday) {
        return nil, errors.New("overtime must be submitted after working hours (17:00 WIB)")
    }

	
    existingOvertime, err := s.repo.FindByUserAndDate(userID, overtimeDate)
    if err != nil {
        return nil, err
    }
    if existingOvertime != nil {
        return nil, errors.New("overtime already submitted for this date")
    }

    if req.Hours < 1 || req.Hours > 3 {
        return nil, errors.New("hours must be between 1 and 3")
    }

    // Buat lembur baru
    overtime := Overtime{
        UserID:         userID,
        PayrollPeriodID: period.ID,
        OvertimeDate:   overtimeDate,
        Hours:          req.Hours,
        CreatedAt:      now,
        UpdatedAt:      now,
        CreatedBy:      userID,
        UpdatedBy:      userID,
    }

    if err := s.repo.Create(&overtime); err != nil {
        return nil, err
    }

    return &OvertimeResponse{
        ID:             overtime.ID,
        UserID:         userID,
        PayrollPeriodID: overtime.PayrollPeriodID,
        OvertimeDate:   overtime.OvertimeDate.Format("02/01/2006"),
        Hours:          overtime.Hours,
    }, nil
}