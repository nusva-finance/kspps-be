package models

import (
	"time"
)

type RekeningTransaction struct {
	IDRekeningTransaction uint      `gorm:"column:idrekeningtransaction;primaryKey" json:"id"`
	TransactionType      string    `gorm:"column:transactiontype;not null" json:"transactiontype"`
	IDNusvaRekening      uint      `gorm:"column:idnusvarekening;not null" json:"idnusvarekening"`
	TableTransaction     string    `gorm:"column:tabletransaction;not null" json:"tabletransaction"`
	IDTableTransaction   uint      `gorm:"column:idtabletransaction;not null" json:"idtabletransaction"`
	TanggalTransaksi     time.Time `gorm:"column:tanggaltransaksi;not null" json:"tanggaltransaksi"`
	NominalTransaksi     float64   `gorm:"column:nominaltransaction;not null" json:"nominaltransaction"`
	CreatedBy            string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (RekeningTransaction) TableName() string {
	return "rekening_transaction"
}
