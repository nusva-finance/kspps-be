package models

import (
	"time"

	"gorm.io/gorm"
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

	// Legacy fields for backward compatibility (deprecated)
	AccountType string `gorm:"-" json:"account_type,omitempty"` // derived from saving_account.account_type.code
	MemberID   uint   `gorm:"-" json:"member_id,omitempty"`   // derived from saving_account.member_id
	MemberName string `gorm:"-" json:"member_name,omitempty"` // derived from saving_account.member.full_name
}

// TableName specifies the table name for GORM
func (SavingTransaction) TableName() string {
	return "saving_transactions"
}

// GetAccountTypeCode returns the account type code (backward compatibility helper)
func (st *SavingTransaction) GetAccountTypeCode() string {
	if st.AccountType != "" {
		return st.AccountType // Use legacy field if available
	}
	return st.SavingAccount.GetAccountTypeCode()
}

// AfterFind GORM hook to populate legacy fields for backward compatibility
func (st *SavingTransaction) AfterFind(tx *gorm.DB) error {
	if st.SavingAccount.ID > 0 {
		st.AccountType = st.SavingAccount.GetAccountTypeCode()
		st.MemberID = st.SavingAccount.MemberID
		st.MemberName = st.SavingAccount.Member.FullName
	}
	return nil
}

const (
	TransactionTypeCredit = "credit" // deposit
	TransactionTypeDebit  = "debit"  // withdraw
)
