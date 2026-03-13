package models

import "time"

type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MenuID    uint      `gorm:"not null" json:"menu_id"`
	Menu      Menu      `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	Action    string    `gorm:"not null" json:"action"` // view, create, edit, delete, approve, export
	Code      string    `gorm:"uniqueIndex;not null" json:"code"` // e.g., member_mgmt.view
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
