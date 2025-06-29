package attendances

import (
    "errors"
    "strconv"
    "time"
)

type attendanceService struct {
    repo IAttendanceRepository
	nowFunc func() time.Time
}


func NewAttendanceService(repo IAttendanceRepository) IAttendancesService {
    return &attendanceService{
		repo: repo ,
		nowFunc: time.Now,
	}
}

func (s *attendanceService) CreateAttendance(req CreateAttendanceRequest, userID string) (*AttendanceResponse, error) {
    // Validasi userID
    if userID == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    _, err := strconv.Atoi(userID)
    if err != nil {
        return nil, errors.New("invalid user ID format, must be a valid integer")
    }

    wib, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return nil, errors.New("failed to load WIB timezone")
    }
    now := time.Now().In(wib)


    // if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
    //     return nil, errors.New("attendance cannot be submitted on weekends")
    // }

    period, err := s.repo.FindActivePayrollPeriod()
    if err != nil {
        return nil, err
    }

    existingAttendance, err := s.repo.FindByUserAndDate(userID, now)
    if err != nil {
        return nil, err
    }

    var attendance *Attendance
    if req.Action == "check-in" {
        if existingAttendance == nil {
            attendance = &Attendance{
                ID:             0,
                UserID:         userID,
                PayrollPeriodID: period.ID,
                CheckIn:        now,
                CreatedAt:      now,
                UpdatedAt:      now,
                CreatedBy:      userID,
                UpdatedBy:      userID,
            }
        } else if now.Before(existingAttendance.CheckIn) {
            existingAttendance.CheckIn = now
            existingAttendance.UpdatedAt = now
            existingAttendance.UpdatedBy = userID
            attendance = existingAttendance
        } else {
            return nil, errors.New("check-in already submitted for this date")
        }
        if err := s.repo.Create(attendance); err != nil {
            return nil, err
        }
    } else if req.Action == "check-out" {
        if existingAttendance == nil {
            return nil, errors.New("no check-in found for this date")
        }
        if existingAttendance.CheckOut != nil {
            return nil, errors.New("check-out already submitted for this date")
        }
        attendance = existingAttendance
        attendance.CheckOut = &now
        attendance.UpdatedAt = now
        attendance.UpdatedBy = userID
        if err := s.repo.Create(attendance); err != nil {
            return nil, err
        }
    } else {
        return nil, errors.New("invalid action, must be check-in or check-out")
    }

    // Format respons
    var checkOutStr *string
    if attendance.CheckOut != nil {
        str := attendance.CheckOut.In(wib).Format("02/01/2006 : 15:04:05")
        checkOutStr = &str
    }
    return &AttendanceResponse{
        ID:             strconv.Itoa(attendance.ID),
        UserID:         userID,
        PayrollPeriodID: attendance.PayrollPeriodID,
        CheckIn:        attendance.CheckIn.In(wib).Format("02/01/2006 : 15:04:05"),
        CheckOut:       checkOutStr,
    }, nil
}