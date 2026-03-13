package models

import "time"

type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Icon       string    `json:"icon"`
	Path       string    `json:"path"`
	ParentID   *uint     `json:"parent_id"`
	Order      int       `gorm:"default:0" json:"order"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Parent     *Menu     `json:"parent,omitempty"`
	Children   []Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
