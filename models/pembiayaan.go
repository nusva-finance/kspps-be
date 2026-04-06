package models

import (
	"time"


)

type Pembiayaan struct {
	IDPinjaman       int        `gorm:"column:idpinjaman;primaryKey" json:"idpinjaman"`
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


}

func (Pembiayaan) TableName() string {
	return "pembiayaan"
}

type PembiayaanWithMemberName struct {
	Pembiayaan
	// Tag column harus sama dengan alias di SELECT repository (as nama_anggota)
	NamaAnggota     string  `gorm:"column:nama_anggota" json:"namaanggota"`
	MemberNo        string  `gorm:"column:member_no" json:"member_no"`
	TotalPembayaran float64 `gorm:"column:total_pembayaran" json:"totalpembayaran"`
}

// TableName for GORM
func (PembiayaanWithMemberName) TableName() string {
	return "pembiayaan"
}
