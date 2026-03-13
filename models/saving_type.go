package models

import "time"

// SavingType represents a type of savings account (pokok, wajib, manasuka, etc.)
type SavingType struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"uniqueIndex;not null" json:"code"`            // pokok, wajib, manasuka
	Name          string    `gorm:"not null" json:"name"`                        // Simpanan Pokok, etc
	Description   string    `json:"description"`
	IsRequired    bool      `gorm:"default:false" json:"is_required"`             // mandatory or optional
	MinBalance    float64   `gorm:"default:0" json:"min_balance"`               // minimum balance required
	IsActive      bool      `gorm:"default:true" json:"is_active"`               // active or inactive
	DisplayOrder  int       `gorm:"default:0" json:"display_order"`             // order for display
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	SavingAccounts []SavingAccount `gorm:"foreignKey:AccountTypeID" json:"saving_accounts,omitempty"`
}

// TableName specifies the table name for GORM
func (SavingType) TableName() string {
	return "saving_types"
}

// Common Saving Types (for backward compatibility)
const (
	SavingTypeCodePokok = "pokok"
	SavingTypeCodeWajib = "wajib"
	SavingTypeCodeModal = "modal"
)

// DefaultSavingTypes returns the default saving types that should be initialized
func DefaultSavingTypes() []SavingType {
	return []SavingType{
		{
			Code:         SavingTypeCodePokok,
			Name:         "Simpanan Pokok",
			Description:  "Simpanan wajib pertama kali menjadi anggota koperasi",
			IsRequired:   true,
			MinBalance:   0,
			IsActive:     true,
			DisplayOrder: 1,
		},
		{
			Code:         SavingTypeCodeWajib,
			Name:         "Simpanan Wajib",
			Description:  "Simpanan rutin yang wajib dibayar oleh anggota",
			IsRequired:   true,
			MinBalance:   0,
			IsActive:     true,
			DisplayOrder: 2,
		},
		{
			Code:         SavingTypeCodeModal,
			Name:         "Simpanan Modal",
			Description:  "Simpanan modal untuk keperluan investasi dan pengembangan usaha",
			IsRequired:   false,
			MinBalance:   0,
			IsActive:     true,
			DisplayOrder: 3,
		},
	}
}
