package reimbursements


type IReimbursementService interface {
    CreateReimbursement(req CreateReimbursementRequest, userID string) (*ReimbursementResponse, error)
}

type IReimbursementRepository interface {
    Create(reimbursement *Reimbursement) error
    FindActivePayrollPeriod() (*PayrollPeriod, error)
}

