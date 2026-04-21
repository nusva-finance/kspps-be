package models

import (
	"time"
)

type SavingTransaction struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	SavingAccountID  uint          `gorm:"not null" json:"saving_account_id"`
	SavingAccount    SavingAccount `gorm:"foreignKey:SavingAccountID" json:"saving_account,omitempty"`
	TransactionType  string        `gorm:"not null" json:"transaction_type"` // credit/debit
	Amount           float64       `gorm:"not null" json:"amount"`
	Description      string        `json:"description"`
	BalanceBefore    float64       `json:"balance_before"`
	BalanceAfter     float64       `json:"balance_after"`
	TransactionDate  time.Time     `gorm:"not null" json:"transaction_date"`
	CreatedBy        string        `json:"created_by"`
	UpdatedBy        string        `json:"updated_by"`
	CreatedAt        time.Time     `json:"created_at"`

	// --- PERUBAHAN DI SINI ---
	// Gunakan gorm:"-" agar GORM mengabaikan field ini saat Insert/Update 
	// karena kolom-kolom ini tidak ada di tabel saving_transactions yang asli
	
	MemberID        uint   `gorm:"-" json:"member_id,omitempty"`
	MemberNo        string `gorm:"-" json:"member_no,omitempty"`
	MemberName      string `gorm:"-" json:"member_name,omitempty"`
	AccountTypeID   uint   `gorm:"-" json:"account_type_id,omitempty"`
	AccountType     string `gorm:"-" json:"account_type,omitempty"`
	AccountTypeName string `gorm:"-" json:"account_type_name,omitempty"`
	RekeningID      uint   `gorm:"-" json:"rekening_id,omitempty"`
	RekeningName    string `gorm:"-" json:"rekening_name,omitempty"`
	RekeningNo      string `gorm:"-" json:"rekening_no,omitempty"`
}

// TableName specifies the table name for GORM
func (SavingTransaction) TableName() string {
	return "saving_transactions"
}

const (
	TransactionTypeCredit = "credit" // deposit
	TransactionTypeDebit  = "debit"  // withdraw
)