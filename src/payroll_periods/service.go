package payroll_periods

import (
    "errors"
    "strconv"
    "time"
)

type payrollPeriodService struct {
    repo IPayrollPeriodRepository
}

func NewPayrollPeriodService(repo IPayrollPeriodRepository) IPayrollPeriodService {
    return &payrollPeriodService{repo}
}

func (s *payrollPeriodService) CreatePayrollPeriod(req CreatePayrollPeriodRequest, userID string) (*PayrollPeriodResponse, error) {
    // Validasi userID
    if userID == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    _, err := strconv.Atoi(userID)
    if err != nil {
        return nil, errors.New("invalid user ID format, must be a valid integer")
    }
    
    startDate, err := time.Parse("02/01/2006", req.StartDate)
    if err != nil {
        return nil, errors.New("invalid start_date format, use DD/MM/YYYY")
    }
    endDate, err := time.Parse("02/01/2006", req.EndDate)
    if err != nil {
        return nil, errors.New("invalid end_date format, use DD/MM/YYYY")
    }

    if startDate.Day() != 1 {
        return nil, errors.New("start_date must be the first day of the month")
    }

    expectedEndDate := startDate.AddDate(0, 1, -1)
    if !endDate.Equal(expectedEndDate) {
        return nil, errors.New("end_date must be the last day of the month")
    }
    
    if !startDate.Before(endDate) {
        return nil, errors.New("start_date must be before end_date")
    }

    activePeriods, err := s.repo.FindActivePeriods()
    if err != nil {
        return nil, err
    }
    if len(activePeriods) > 0 {
        return nil, errors.New("an active payroll period already exists")
    }

  //if existing period, return error
    existingPeriods, err := s.repo.FindByDateRange(startDate, endDate)
    if err != nil {
        return nil, err
    }
    if len(existingPeriods) > 0 {
        return nil, errors.New("overlapping payroll period exists")
    }

    period := PayrollPeriod{
        StartDate:   startDate,
        EndDate:     endDate,
        IsProcessed: false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        CreatedBy:   userID,
        UpdatedBy:   userID,
    }

    if err := s.repo.Create(&period); err != nil {
        return nil, err
    }

    return &PayrollPeriodResponse{
        ID:          period.ID,
        StartDate:   period.StartDate.Format("02/01/2006"),
        EndDate:     period.EndDate.Format("02/01/2006"),
        IsProcessed: period.IsProcessed,
    }, nil
}