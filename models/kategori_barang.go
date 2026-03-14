package models

import (
	"time"
)

type KategoriBarang struct {
	ID        uint      `gorm:"column:idkategoribarang;primaryKey" json:"id"`
	Kategori  string    `gorm:"column:namakategoribarang;not null;uniqueIndex" json:"kategori"`
	Aktif     bool      `gorm:"column:is_active;default:true" json:"aktif"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (KategoriBarang) TableName() string {
    return "kategori_barang"
}