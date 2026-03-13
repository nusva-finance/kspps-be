package models

import "time"

type LoanTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ApplicationID   uint      `gorm:"not null" json:"application_id"`
	Application     LoanApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	ScheduleID      *uint     `json:"schedule_id"`
	Schedule        *LoanSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Amount          float64   `gorm:"not null" json:"amount"`
	TransactionDate time.Time `gorm:"not null" json:"transaction_date"`
	Description     string    `json:"description"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}
