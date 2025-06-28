package payroll_periods

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type mockPayrollPeriodRepository struct {
    mock.Mock
}

func (m *mockPayrollPeriodRepository) Create(period *PayrollPeriod) error {
    args := m.Called(period)
    period.ID = "1"
    return args.Error(0)
}

func (m *mockPayrollPeriodRepository) FindByDateRange(startDate, endDate time.Time) ([]PayrollPeriod, error) {
    args := m.Called(startDate, endDate)
    return args.Get(0).([]PayrollPeriod), args.Error(1)
}

func (m *mockPayrollPeriodRepository) FindActivePeriods() ([]PayrollPeriod, error) {
    args := m.Called()
    return args.Get(0).([]PayrollPeriod), args.Error(1)
}

func TestCreatePayrollPeriod_Success(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "31/07/2025",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{}, nil)
    repo.On("FindByDateRange", mock.Anything, mock.Anything).Return([]PayrollPeriod{}, nil)
    repo.On("Create", mock.Anything).Return(nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.NoError(t, err)
    assert.NotNil(t, period)
    assert.Equal(t, "1", period.ID)
    assert.Equal(t, "01/07/2025", period.StartDate)
    assert.Equal(t, "31/07/2025", period.EndDate)
    assert.False(t, period.IsProcessed)

    repo.AssertExpectations(t)
}

func TestCreatePayrollPeriod_Success_PastDate(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/06/2025",
        EndDate:   "30/06/2025",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{}, nil)
    repo.On("FindByDateRange", mock.Anything, mock.Anything).Return([]PayrollPeriod{}, nil)
    repo.On("Create", mock.Anything).Return(nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.NoError(t, err)
    assert.NotNil(t, period)
    assert.Equal(t, "1", period.ID)
    assert.Equal(t, "01/06/2025", period.StartDate)
    assert.Equal(t, "30/06/2025", period.EndDate)
    assert.False(t, period.IsProcessed)

    repo.AssertExpectations(t)
}

func TestCreatePayrollPeriod_EmptyUserID(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := ""
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "31/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "user ID cannot be empty", err.Error())
}

func TestCreatePayrollPeriod_InvalidUserID(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "invalid"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "31/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "invalid user ID format, must be a valid integer", err.Error())
}

func TestCreatePayrollPeriod_InvalidStartDateFormat(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "2025-07-01",
        EndDate:   "31/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "invalid start_date format, use DD/MM/YYYY", err.Error())
}

func TestCreatePayrollPeriod_InvalidEndDateFormat(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "2025-07-31",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "invalid end_date format, use DD/MM/YYYY", err.Error())
}

func TestCreatePayrollPeriod_StartDateNotFirstDay(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "02/07/2025",
        EndDate:   "31/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "start_date must be the first day of the month", err.Error())
}

func TestCreatePayrollPeriod_EndDateNotLastDay(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "30/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "end_date must be the last day of the month", err.Error())
}

func TestCreatePayrollPeriod_StartDateNotBeforeEndDate(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/08/2025",
        EndDate:   "31/07/2025",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "start_date must be before end_date", err.Error())
}

func TestCreatePayrollPeriod_ActivePeriodExists(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "31/07/2025",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{{ID: "1"}}, nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "an active payroll period already exists", err.Error())
}

func TestCreatePayrollPeriod_OverlappingPeriod(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1"
    req := CreatePayrollPeriodRequest{
        StartDate: "01/07/2025",
        EndDate:   "31/07/2025",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{}, nil)
    repo.On("FindByDateRange", mock.Anything, mock.Anything).Return([]PayrollPeriod{{ID: "1"}}, nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "overlapping payroll period exists", err.Error())
}