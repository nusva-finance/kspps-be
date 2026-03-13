package models

import "time"

type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	PermissionID uint      `gorm:"not null" json:"permission_id"`
	IsAllowed    bool      `gorm:"default:true" json:"is_allowed"`
	CreatedAt    time.Time `json:"created_at"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
