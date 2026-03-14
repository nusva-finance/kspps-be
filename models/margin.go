package models

import (
	"time"
)

type MarginSetup struct {
	ID        uint      `gorm:"column:idmargin;primaryKey" json:"id"`
	Category  string    `gorm:"column:category;not null" json:"category"`
	Tenor     int       `gorm:"column:tenor;not null" json:"tenor"`
	Margin    float64   `gorm:"column:margin;not null" json:"margin"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
