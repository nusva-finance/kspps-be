package models

import (
	"time"
)

type PembayaranPembiayaan struct {
	IDPembayaranPembiayaan int       `gorm:"column:idpembayaranpembiayaan;primaryKey;autoIncrement" json:"idpembayaranpembiayaan"`
	IDPinjaman             int       `gorm:"column:idpinjaman;not null" json:"idpinjaman"`
	NominalPembayaran      float64   `gorm:"column:nominalpembayaran;not null;default:0" json:"nominalpembayaran"`
	NominalAngsuran        float64   `gorm:"column:nominalangsuran;not null;default:0" json:"nominalangsuran"`
	NominalPendapatanLain  float64   `gorm:"column:nominalpendapatanlainlain;not null;default:0" json:"nominalpendapatanlainlain"`
	AngsuranKe             int       `gorm:"column:angsuranke;not null" json:"angsuranke"`
	TglJtAngsuran          time.Time `gorm:"column:tgljtangsuran;not null" json:"tgljtangsuran"`
	TglPembayaran          time.Time `gorm:"column:tglpembayaran;not null" json:"tglpembayaran"`
	Keterangan             string    `gorm:"column:keterangan" json:"keterangan"`
	SysRevID               int       `gorm:"column:sysrevid;default:1" json:"sysrevid"`
	SysRowID               string    `gorm:"column:sysrowid;default:gen_random_uuid()" json:"sysrowid"`

	CreatedBy              string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy              string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt              time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (PembayaranPembiayaan) TableName() string {
	return "pembayaran_pembiayaan"
}

// PembayaranPembiayaanWithDetails includes additional info from joins
type PembayaranPembiayaanWithDetails struct {
	PembayaranPembiayaan
	NamaAnggota string `gorm:"column:nama_anggota" json:"nama_anggota"`
	MemberNo    string `gorm:"column:member_no" json:"member_no"`
	IDMember    int    `gorm:"column:idmember" json:"idmember"`
}
