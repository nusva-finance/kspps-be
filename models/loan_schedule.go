package models

import "time"

type LoanSchedule struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ApplicationID  uint      `gorm:"not null" json:"application_id"`
	Application    LoanApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	SequenceNumber int       `gorm:"not null" json:"sequence_number"`
	DueDate        time.Time `gorm:"not null" json:"due_date"`
	Principal      float64   `gorm:"not null" json:"principal"`
	Margin         float64   `gorm:"not null" json:"margin"`
	TotalAmount    float64   `gorm:"not null" json:"total_amount"`
	IsPaid         bool      `gorm:"default:false" json:"is_paid"`
	PaidDate       *time.Time `json:"paid_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
