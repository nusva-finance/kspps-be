package models

import (
	"time"
)

type QardHassan struct {
	IDQardHassan   int       `gorm:"column:idqardhassan;primaryKey" json:"idqardhassan"`
	IDMember       int       `gorm:"column:idmember;not null" json:"idmember"`
	TanggalPinjaman time.Time `gorm:"column:tanggalpinjaman;not null" json:"tanggalpinjaman"`
	BiayaAdmin     float64   `gorm:"column:biayaadmin;not null" json:"biayaadmin"`
	NominalPinjaman float64  `gorm:"column:nominalpinjaman;not null" json:"nominalpinjaman"`
	TglJtTempo     time.Time `gorm:"column:tgljttempo;not null" json:"tgljttempo"`
	Keterangan     string    `gorm:"column:keterangan" json:"keterangan"`
	NominalPembayaran float64    `gorm:"column:nominalpembayaran;default:0" json:"nominalpembayaran"`
	TanggalPembayaran *time.Time `gorm:"column:tanggalpembayaran" json:"tanggalpembayaran"` // Pakai pointer (*) agar bisa dirender sebagai null di JSON saat belum lunas	
	SysRevID       int       `gorm:"column:sysrevid" json:"sysrevid"`
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy      string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (QardHassan) TableName() string {
	return "qardhassan"
}

type QardHassanWithMemberName struct {
	QardHassan
	MemberNo      string  `gorm:"column:member_no" json:"member_no"`
	NamaAnggota   string  `gorm:"column:nama_anggota" json:"namaanggota"`
}

func (QardHassanWithMemberName) TableName() string {
	return "qardhassan"
}
