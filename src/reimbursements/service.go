package reimbursements

import (
    "errors"
    "strconv"
    "time"
)

type reimbursementService struct {
    repo IReimbursementRepository
}

func NewReimbursementService(repo IReimbursementRepository) IReimbursementService {
    return &reimbursementService{repo}
}

func (s *reimbursementService) CreateReimbursement(req CreateReimbursementRequest, userID string) (*ReimbursementResponse, error) {
    // Validasi userID
    if userID == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    _, err := strconv.Atoi(userID)
    if err != nil {
        return nil, errors.New("invalid user ID format, must be a valid integer")
    }

    // Validasi amount
    if req.Amount == 0 || req.Amount < 0 {
        return nil, errors.New("amount must be greater than 0")
    }

    period, err := s.repo.FindActivePayrollPeriod()
    if err != nil {
        return nil, err
    }
	
    wib, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return nil, errors.New("failed to load WIB timezone")
    }
    now := time.Now().In(wib)

    // Buat reimbursement baru
    reimbursement := Reimbursement{
        UserID:         userID,
        PayrollPeriodID: period.ID,
        Amount:         req.Amount,
        Description:    req.Description,
        CreatedAt:      now,
        UpdatedAt:      now,
        CreatedBy:      userID,
        UpdatedBy:      userID,
    }

    if err := s.repo.Create(&reimbursement); err != nil {
        return nil, err
    }

    return &ReimbursementResponse{
        ID:             reimbursement.ID,
        UserID:         userID,
        PayrollPeriodID: reimbursement.PayrollPeriodID,
        Amount:         reimbursement.Amount,
        Description:    reimbursement.Description,
    }, nil
}