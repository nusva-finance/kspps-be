package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type PembayaranPembiayaanRepository struct {
	db *gorm.DB
}

func NewPembayaranPembiayaanRepository() *PembayaranPembiayaanRepository {
	return &PembayaranPembiayaanRepository{
		db: config.GetDB(),
	}
}

// Create - Membuat data pembayaran baru
func (r *PembayaranPembiayaanRepository) Create(pembayaran *models.PembayaranPembiayaan) error {
	return r.db.Create(pembayaran).Error
}

// FindByID - Mencari pembayaran berdasarkan ID
func (r *PembayaranPembiayaanRepository) FindByID(id int) (*models.PembayaranPembiayaan, error) {
	var pembayaran models.PembayaranPembiayaan
	err := r.db.Where("idpembayaranpembiayaan = ?", id).First(&pembayaran).Error
	if err != nil {
		return nil, err
	}
	return &pembayaran, nil
}

// ListByPinjamanID - Mengambil daftar pembayaran berdasarkan ID Pinjaman
// ListByPinjamanID - Mengambil daftar pembayaran berdasarkan ID Pinjaman
func (r *PembayaranPembiayaanRepository) ListByPinjamanID(idPinjaman int) ([]models.PembayaranPembiayaanWithDetails, error) {
	var result []models.PembayaranPembiayaanWithDetails

	query := `
		SELECT
			pp.*,
			m.full_name as nama_anggota,
			m.member_no,
			pem.idmember,
			COALESCE(rt.idnusvarekening, 0) as idnusvarekening,
			COALESCE(nr.namarekening, '') as namarekening,
			COALESCE(nr.norekening, '') as norekening
		FROM pembayaran_pembiayaan pp
		LEFT JOIN pembiayaan pem ON pp.idpinjaman = pem.idpinjaman
		LEFT JOIN members m ON pem.idmember = m.id
		-- Mengambil TOP 1 (LIMIT 1) transaksi rekening terbaru (Insert atau Update)
		LEFT JOIN LATERAL (
			SELECT idnusvarekening 
			FROM rekening_transaction 
			WHERE tabletransaction = 'pembayaran_pembiayaan' 
			  AND idtabletransaction = pp.idpembayaranpembiayaan 
			  AND UPPER(transactiontype) IN ('INSERT', 'UPDATE')
			ORDER BY created_at DESC 
			LIMIT 1
		) rt ON true
		-- Relasi ke tabel master rekening
		LEFT JOIN nusva_rekening nr ON rt.idnusvarekening = nr.idnusvarekening
		WHERE pp.idpinjaman = ?
		ORDER BY pp.angsuranke ASC, pp.tglpembayaran ASC
	`

	err := r.db.Raw(query, idPinjaman).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CountByPinjamanID - Menghitung jumlah pembayaran berdasarkan ID Pinjaman
func (r *PembayaranPembiayaanRepository) CountByPinjamanID(idPinjaman int) (int64, error) {
	var count int64
	err := r.db.Model(&models.PembayaranPembiayaan{}).
		Where("idpinjaman = ?", idPinjaman).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
// Update - Memperbarui data pembayaran
func (r *PembayaranPembiayaanRepository) Update(pembayaran *models.PembayaranPembiayaan) error {
	return r.db.Model(&models.PembayaranPembiayaan{}).
		Where("idpembayaranpembiayaan = ?", pembayaran.IDPembayaranPembiayaan).
		Updates(map[string]interface{}{
			"tglpembayaran":             pembayaran.TglPembayaran,
			"nominalpembayaran":         pembayaran.NominalPembayaran,
			"nominalangsuran":           pembayaran.NominalAngsuran,
			"nominalpendapatanlainlain": pembayaran.NominalPendapatanLain,
			"angsuranke":                pembayaran.AngsuranKe,
			"tgljtangsuran":             pembayaran.TglJtAngsuran,
			"keterangan":                pembayaran.Keterangan,
			"updated_by":                pembayaran.UpdatedBy,
			"updated_at":                pembayaran.UpdatedAt,
		}).Error
}

// Delete - Menghapus data pembayaran
func (r *PembayaranPembiayaanRepository) Delete(id int) error {
	return r.db.Where("idpembayaranpembiayaan = ?", id).Delete(&models.PembayaranPembiayaan{}).Error
}

// CreateRekeningTransaction - Membuat record transaksi rekening
func (r *PembayaranPembiayaanRepository) CreateRekeningTransaction(transaction *models.RekeningTransaction) error {
	return r.db.Create(transaction).Error
}

// FindRekeningTransactionByReference - Mencari transaksi rekening berdasarkan table dan id
func (r *PembayaranPembiayaanRepository) FindRekeningTransactionByReference(tableName string, idTableTransaction uint) (*models.RekeningTransaction, error) {
	var transaction models.RekeningTransaction
	err := r.db.Where("tabletransaction = ? AND idtabletransaction = ?", tableName, idTableTransaction).
		Order("idrekeningtransaction DESC").
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
