package models

import "time"

type SavingAccount struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	MemberID       uint        `gorm:"not null" json:"member_id"`
	Member         Member      `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	AccountTypeID  uint        `gorm:"not null" json:"account_type_id"` // foreign key to saving_types
	AccountType    SavingType  `gorm:"foreignKey:AccountTypeID" json:"account_type,omitempty"`
	AccountNumber  string      `gorm:"uniqueIndex;not null" json:"account_number"`
	Balance        float64     `gorm:"default:0" json:"balance"`
	IsActive       bool        `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (SavingAccount) TableName() string {
	return "saving_accounts"
}

// GetAccountTypeCode returns the account type code (backward compatibility helper)
func (sa *SavingAccount) GetAccountTypeCode() string {
	return sa.AccountType.Code
}

// GetAccountTypeName returns the account type name
func (sa *SavingAccount) GetAccountTypeName() string {
	return sa.AccountType.Name
}

// Backward compatibility constants (deprecated, use SavingType constants)
const (
	AccountTypePokok = "pokok"
	AccountTypeWajib = "wajib"
	AccountTypeModal = "modal"
)
