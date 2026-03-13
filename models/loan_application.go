package models

import "time"

type LoanApplication struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ApplicationNo   string    `gorm:"uniqueIndex;not null" json:"application_no"`
	MemberID        uint      `gorm:"not null" json:"member_id"`
	Member          Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	PrincipalAmount float64   `gorm:"not null" json:"principal_amount"`
	MarginRate      float64   `gorm:"not null" json:"margin_rate"` // in percentage
	MarginAmount    float64   `gorm:"not null" json:"margin_amount"`
	TermMonths      int       `gorm:"not null" json:"term_months"`
	MonthlyInstallment float64 `gorm:"not null" json:"monthly_installment"`
	StartDate       time.Time `gorm:"not null" json:"start_date"`
	EndDate         time.Time `gorm:"not null" json:"end_date"`
	Purpose         string    `json:"purpose"`
	ContractType    string    `gorm:"not null" json:"contract_type"` // murabahah, musyarakah
	Status          string    `gorm:"default:pending;not null" json:"status"` // pending, approved, rejected, active, paid_off
	ApprovedBy      string    `json:"approved_by"`
	ApprovedDate    *time.Time `json:"approved_date"`
	Notes           string    `json:"notes"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Schedules       []LoanSchedule `gorm:"foreignKey:ApplicationID" json:"schedules,omitempty"`
	Transactions    []LoanTransaction `gorm:"foreignKey:ApplicationID" json:"transactions,omitempty"`
}

const (
	ContractTypeMurabahah = "murabahah"
	ContractTypeMusyarakah = "musyarakah"
)

const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusRejected  = "rejected"
	StatusActive    = "active"
	StatusPaidOff   = "paid_off"
)
