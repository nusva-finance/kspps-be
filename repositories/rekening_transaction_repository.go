package repositories

import (
	"fmt"
	"time"

	"nusvakspps/models"
	"gorm.io/gorm"
)

// Note: GetDB() is defined in base_repository.go or similar file in this package

type RekeningTransactionRepository struct {
	db *gorm.DB
}

func NewRekeningTransactionRepository() *RekeningTransactionRepository {
	return &RekeningTransactionRepository{
		db: GetDB(),
	}
}

func (r *RekeningTransactionRepository) Create(transaction *models.RekeningTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *RekeningTransactionRepository) FindByID(id uint) (*models.RekeningTransaction, error) {
	var transaction models.RekeningTransaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *RekeningTransactionRepository) Update(transaction *models.RekeningTransaction) error {
	return r.db.Save(transaction).Error
}

func (r *RekeningTransactionRepository) Delete(id uint) error {
	// Direct deletion of rekening_transaction is not allowed
	// Transactions are immutable and should only be created, never deleted
	return fmt.Errorf("penghapusan langsung rekening_transaction tidak diizinkan - transaksi bersifat immutable")
}

func (r *RekeningTransactionRepository) List(limit, offset int) ([]models.RekeningTransaction, int64, error) {
	var transactions []models.RekeningTransaction
	var total int64

	if err := r.db.Model(&models.RekeningTransaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// MutasiRekeningRow represents a row in mutasi rekening response
type MutasiRekeningRow struct {
	IDRekeningTransaction uint      `gorm:"column:idrekeningtransaction" json:"idrekeningtransaction"`
	TanggalTransaksi      time.Time `gorm:"column:tanggaltransaksi" json:"tanggaltransaksi"`
	NominalTransaksi      float64   `gorm:"column:nominaltransaction" json:"nominaltransaction"`
	TransactionType       string    `gorm:"column:transactiontype" json:"transactiontype"`
	TableTransaction      string    `gorm:"column:tabletransaction" json:"tabletransaction"`
	IDTableTransaction    uint      `gorm:"column:idtabletransaction" json:"idtabletransaction"`
}

// GetMutasiByRekeningID gets transactions for a rekening with date range filter
func (r *RekeningTransactionRepository) GetMutasiByRekeningID(idNusvaRekening uint, dateFrom, dateTo time.Time) ([]MutasiRekeningRow, error) {
	var transactions []MutasiRekeningRow

	query := `
		SELECT
			rt.idrekeningtransaction,
			rt.tanggaltransaksi,
			rt.nominaltransaction,
			rt.transactiontype,
			rt.tabletransaction,
			rt.idtabletransaction
		FROM rekening_transaction rt
		WHERE rt.idnusvarekening = ?
		AND rt.tanggaltransaksi >= ?
		AND rt.tanggaltransaksi <= ?
		ORDER BY rt.tanggaltransaksi ASC, rt.idrekeningtransaction ASC
	`

	err := r.db.Raw(query, idNusvaRekening, dateFrom, dateTo).Scan(&transactions).Error
	return transactions, err
}

// GetOpeningBalance calculates opening balance before the date from
func (r *RekeningTransactionRepository) GetOpeningBalance(idNusvaRekening uint, dateFrom time.Time) (float64, error) {
	var result struct {
		Total float64 `json:"total"`
	}

	query := `
		SELECT COALESCE(SUM(nominaltransaction), 0) as total
		FROM rekening_transaction
		WHERE idnusvarekening = ?
		AND tanggaltransaksi < ?
	`

	err := r.db.Raw(query, idNusvaRekening, dateFrom).Scan(&result).Error
	return result.Total, err
}

// GetIDNusvaRekeningByPembiayaanID gets the idnusvarekening from the original Insert transaction
func (r *RekeningTransactionRepository) GetIDNusvaRekeningByPembiayaanID(pembiayaanID uint) (uint, error) {
    var transaction models.RekeningTransaction
    
    // Logika: Cari baris terakhir (Top 1) berdasarkan waktu dibuat
    err := r.db.Where("tabletransaction = ? AND idtabletransaction = ?", "pembiayaan", pembiayaanID).
        Order("created_at DESC"). // Urutkan dari yang terbaru (Descending)
        First(&transaction).      // Ambil baris pertama hasil urutan tersebut (Limit 1)
        Error

    if err != nil {
        return 0, err
    }
    return transaction.IDNusvaRekening, nil
}
