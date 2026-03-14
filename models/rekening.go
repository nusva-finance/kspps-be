package models

import (
	"time"
)

type NusvaRekening struct {
	IDRekening    uint      `gorm:"column:idnusvarekening;primaryKey" json:"id"`
	NamaRekening  string    `gorm:"column:namarekening;not null" json:"namarekening"`
	NoRekening    string    `gorm:"column:norekening;not null;uniqueIndex" json:"norekening"`
	Aktif         bool      `gorm:"column:is_active;default:true" json:"aktif"`
	Deskripsi     string    `gorm:"column:description" json:"deskripsi"`
	CreatedBy     string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (NusvaRekening) TableName() string {
	return "nusva_rekening"
}
