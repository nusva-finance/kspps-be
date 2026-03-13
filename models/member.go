package models

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	MemberNo       string    `gorm:"uniqueIndex;not null" json:"member_no"` // YYMMXXZZZZZZ
	FullName       string    `gorm:"not null" json:"full_name"`
	Gender         string    `gorm:"not null" json:"gender"` // Laki-laki/Perempuan
	JoinDate       time.Time `gorm:"not null" json:"join_date"`
	JoinYear       string    `gorm:"type:char(2)" json:"join_year"` // YY - tahun masuk
	JoinMonth      string    `gorm:"type:char(2)" json:"join_month"` // MM - bulan masuk
	BirthDate      time.Time `json:"birth_date"`
	BirthPlace     string    `json:"birth_place"`
	KtpNo          string    `gorm:"uniqueIndex" json:"ktp_no"` // NIK/KTP (renamed from nik)
	NpwpNo         string    `json:"npwp_no"` // NPWP (renamed from npwp)
	AddressKtp     string    `gorm:"not null" json:"address_ktp"` // Alamat sesuai KTP (renamed from address)
	City           string    `json:"city"`
	Province       string    `json:"province"`
	PostalCode     string    `json:"postal_code"`
	PhoneNumber    string    `gorm:"not null" json:"phone_number"`
	Email          string    `json:"email"`
	KTPPhoto       string    `json:"ktp_photo"`
	NPWPPhoto      string    `json:"npwp_photo"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedBy      string    `json:"created_by"`
	UpdatedBy      string    `json:"updated_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Emergency Contact Information
	EmergencyName    string `json:"emergency_name"`
	EmergencyRelation string `json:"emergency_relation"`
	EmergencyPhone   string `json:"emergency_phone"`
	EmergencyAddress string `json:"emergency_address"`

	// Work Information
	CompanyName string `json:"company_name"`
	JobTitle   string `json:"job_title"`

	// Bank Information
	BankAccountNo string `json:"bank_account_no"`
	BankName      string `json:"bank_name"`

	SavingAccounts []SavingAccount `gorm:"foreignKey:MemberID" json:"saving_accounts,omitempty"`
}

const (
	GenderMale   = "Laki-laki"
	GenderFemale = "Perempuan"
)
