package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`	
	FullName     string    `gorm:"not null" json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	IsLocked     bool      `gorm:"default:false" json:"is_locked"`
	FailedLogin  int       `gorm:"default:0" json:"failed_login"`
	LastLogin    *time.Time `json:"last_login"`
	LastIP       string    `json:"last_ip"`
	ForceChange  bool      `gorm:"default:false" json:"force_change"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Roles        []Role   `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
