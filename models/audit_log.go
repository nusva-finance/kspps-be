package models

import "time"

type AuditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      *uint     `json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Username    string    `json:"username"`
	Action      string    `gorm:"not null" json:"action"` // login, create, update, delete, approve
	Module      string    `gorm:"not null" json:"module"`
	RecordID    *uint     `json:"record_id"`
	OldData     string    `gorm:"type:text" json:"old_data"`
	NewData     string    `gorm:"type:text" json:"new_data"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	Status      string    `gorm:"default:success" json:"status"` // success, failed
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	ActionLogin   = "login"
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionApprove = "approve"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)
