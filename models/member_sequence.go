package models

import "time"

type MemberSequenceCounter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	YearCode    string    `gorm:"not null" json:"year_code"` // YY
	MonthCode   string    `gorm:"not null" json:"month_code"` // MM
	GenderCode  string    `gorm:"not null" json:"gender_code"` // 01/02
	LastSeq     int       `gorm:"default:0" json:"last_seq"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MemberSequenceCounter) TableName() string {
	return "member_sequence_counter"
}
