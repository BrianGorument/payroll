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
        StartDate: "2025-07-01",
        EndDate:   "2025-07-31",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{}, nil)
    repo.On("FindByDateRange", mock.Anything, mock.Anything).Return([]PayrollPeriod{}, nil)
    repo.On("Create", mock.Anything).Return(nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, period.ID)
    assert.Equal(t, "2025-07-01", period.StartDate)
    assert.Equal(t, "2025-07-31", period.EndDate)
    assert.False(t, period.IsProcessed)

    repo.AssertExpectations(t)
}

func TestCreatePayrollPeriod_InvalidDateFormat(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1" 
    req := CreatePayrollPeriodRequest{
        StartDate: "2025-07-01T00:00:00Z",
        EndDate:   "2025-07-31",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "invalid start_date format, use YYYY-MM-DD", err.Error())
}

func TestCreatePayrollPeriod_StartDateNotFirstDay(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1" 
    req := CreatePayrollPeriodRequest{
        StartDate: "2025-07-02",
        EndDate:   "2025-07-31",
    }

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "start_date must be the first day of the month", err.Error())
}

func TestCreatePayrollPeriod_ActivePeriodExists(t *testing.T) {
    repo := new(mockPayrollPeriodRepository)
    service := NewPayrollPeriodService(repo)
    userID := "1" 
    req := CreatePayrollPeriodRequest{
        StartDate: "2025-07-01",
        EndDate:   "2025-07-31",
    }

    repo.On("FindActivePeriods").Return([]PayrollPeriod{{ID: "1"}}, nil)

    period, err := service.CreatePayrollPeriod(req, userID)
    assert.Error(t, err)
    assert.Nil(t, period)
    assert.Equal(t, "an active payroll period already exists", err.Error())
}