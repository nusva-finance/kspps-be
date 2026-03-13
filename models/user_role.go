package models

import "time"

type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
