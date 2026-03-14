package models

import (
	"time"


)

type Pembiayaan struct {
	IDPinjaman       int        `gorm:"column:idpinjaman;primaryKey" json:"id"`
	IDMember         int        `gorm:"column:idmember;not null" json:"idmember"`
	TipePinjaman     string     `gorm:"column:tipepinjaman;not null" json:"tipepinjaman"`
	TanggalPinjaman  time.Time  `gorm:"column:tanggalpinjaman;not null" json:"tanggalpinjaman"`
	KategoriBarang   string     `gorm:"column:kategoribarang" json:"kategoribarang"`
	Tenor            int        `gorm:"column:tenor;not null" json:"tenor"`
	Margin           float64    `gorm:"column:margin;not null" json:"margin"`
	NominalPinjaman  float64    `gorm:"column:nominalpinjaman;not null" json:"nominalpinjaman"`
	NominalPembelian float64    `gorm:"column:nominalpembelian;not null" json:"nominalpembelian"`
	TglJtAngsuran1  time.Time  `gorm:"column:tgljtangsuran1;not null" json:"tgljtangsuran1"`
	SysRevID        int        `gorm:"column:sysrevid" json:"sysrevid"`
	
	CreatedBy        string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy        string     `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	

	// Joined fields (not persisted)
	NamaAnggota string `gorm:"-" json:"namaanggota,omitempty"`
	MemberNo     string `gorm:"-" json:"member_no,omitempty"`
}

func (Pembiayaan) TableName() string {
	return "pembiayaan"
}

type PembiayaanWithMemberName struct {
	Pembiayaan
	NamaAnggota string `json:"namaanggota"` 
	MemberNo    string `json:"member_no"`
}
