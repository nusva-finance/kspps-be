package models

import (
	"time"
)

type SavingTransaction struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SavingAccountID  uint      `gorm:"not null" json:"saving_account_id"`
	SavingAccount    SavingAccount `gorm:"foreignKey:SavingAccountID" json:"saving_account,omitempty"`
	TransactionType  string    `gorm:"not null" json:"transaction_type"` // debit/credit
	Amount           float64   `gorm:"not null" json:"amount"`
	Description      string    `json:"description"`
	BalanceBefore    float64   `json:"balance_before"`
	BalanceAfter     float64   `json:"balance_after"`
	TransactionDate  time.Time `gorm:"not null" json:"transaction_date"`
	CreatedBy        string    `json:"created_by"`
	UpdatedBy        string    `json:"updated_by"`
	CreatedAt        time.Time `json:"created_at"`

	// Fields for list queries (from joins) - use column tag for raw SQL mapping
	MemberID        uint   `gorm:"column:member_id" json:"member_id,omitempty"`
	MemberNo        string `gorm:"column:member_no" json:"member_no,omitempty"`
	MemberName      string `gorm:"column:member_name" json:"member_name,omitempty"`
	AccountTypeID   uint   `gorm:"column:account_type_id" json:"account_type_id,omitempty"`
	AccountType     string `gorm:"column:account_type" json:"account_type,omitempty"`
	AccountTypeName string `gorm:"column:account_type_name" json:"account_type_name,omitempty"`
	RekeningID      uint   `gorm:"column:rekening_id" json:"rekening_id,omitempty"`
	RekeningName    string `gorm:"column:rekening_name" json:"rekening_name,omitempty"`
	RekeningNo      string `gorm:"column:rekening_no" json:"rekening_no,omitempty"`
}

// TableName specifies the table name for GORM
func (SavingTransaction) TableName() string {
	return "saving_transactions"
}

const (
	TransactionTypeCredit = "credit" // deposit
	TransactionTypeDebit  = "debit"  // withdraw
)
