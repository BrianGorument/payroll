package payroll_periods

import (
	"errors"
	"time"


)

type payrollPeriodService struct {
	repo IPayrollPeriodRepository
}

func NewPayrollPeriodService(repo IPayrollPeriodRepository) IPayrollPeriodService {
	return &payrollPeriodService{repo}
}

func (s *payrollPeriodService) CreatePayrollPeriod(req CreatePayrollPeriodRequest, userID string) (*PayrollPeriodResponse, error) {
	// Parse string tanggal ke time.Time
    startDate, err := time.Parse("01/02/2006", req.StartDate)
    if err != nil {
        return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
    }
    endDate, err := time.Parse("01/02/2006", req.EndDate)
	
    if startDate.After(endDate) {
        return nil, errors.New("start_date must be before end_date")
    }

	
    if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
        return nil, errors.New("start_date cannot be in the past")
    }

    if endDate.Sub(startDate).Hours() > 31*24 {
        return nil, errors.New("payroll period cannot exceed 31 days")
    }

    existingPeriods, err := s.repo.FindByDateRange(startDate, endDate)
    if err != nil {
        return nil, err
    }
    if len(existingPeriods) > 0 {
        return nil, errors.New("overlapping payroll period exists")
    }

    period := PayrollPeriod{
        StartDate:  startDate,
        EndDate:    endDate,
        IsProcessed: false,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        CreatedBy:  userID,
        UpdatedBy:  userID,
    }

    if err := s.repo.Create(&period); err != nil {
        return nil, err
    }

    return &PayrollPeriodResponse{
        ID:          period.ID,
        StartDate:   period.StartDate,
        EndDate:     period.EndDate,
        IsProcessed: period.IsProcessed,
    }, nil
}