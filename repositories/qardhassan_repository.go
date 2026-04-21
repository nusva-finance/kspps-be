package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type QardHassanRepository struct {
	db *gorm.DB
}

func NewQardHassanRepository() *QardHassanRepository {
	return &QardHassanRepository{
		db: config.GetDB(),
	}
}

// Create - Membuat data qardhassan baru
func (r *QardHassanRepository) Create(qardHassan *models.QardHassan) error {
	return r.db.Create(qardHassan).Error
}

// FindByID - Mencari qardhassan berdasarkan ID
func (r *QardHassanRepository) FindByID(id int) (*models.QardHassan, error) {
	var qardHassan models.QardHassan
	err := r.db.Where("idqardhassan = ?", id).First(&qardHassan).Error
	if err != nil {
		return nil, err
	}
	return &qardHassan, nil
}

// List - Mengambil daftar qardhassan dengan info member
func (r *QardHassanRepository) List(offset, limit int) ([]models.QardHassanWithMemberName, int64, error) {
	var result []models.QardHassanWithMemberName
	var total int64

	// 1. Hitung total data terlebih dahulu
	err := r.db.Table("qardhassan").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. Ambil data dengan LIMIT dan OFFSET
	query := `
		SELECT
			q.*,
			m.member_no,
			m.full_name as nama_anggota
		FROM qardhassan q
		LEFT JOIN members m ON q.idmember = m.id
		ORDER BY q.idqardhassan DESC
		LIMIT ? OFFSET ?
	`

	err = r.db.Raw(query, limit, offset).Scan(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}


// Update - Memperbarui data qardhassan
func (r *QardHassanRepository) Update(qardHassan *models.QardHassan) error {
	return r.db.Model(&models.QardHassan{}).
		Where("idqardhassan = ?", qardHassan.IDQardHassan).
		Updates(map[string]interface{}{
			"idmember":        qardHassan.IDMember,
			"tanggalpinjaman": qardHassan.TanggalPinjaman,
			"biayaadmin":      qardHassan.BiayaAdmin,
			"nominalpinjaman": qardHassan.NominalPinjaman,
			"tgljttempo":      qardHassan.TglJtTempo,
			"keterangan":      qardHassan.Keterangan,
			"nominalpembayaran": qardHassan.NominalPembayaran,
			"tanggalpembayaran": qardHassan.TanggalPembayaran,
			"updated_by":      qardHassan.UpdatedBy,
			"updated_at":      qardHassan.UpdatedAt,
		}).Error
}

// Delete - Menghapus data qardhassan
func (r *QardHassanRepository) Delete(id int) error {
	return r.db.Where("idqardhassan = ?", id).Delete(&models.QardHassan{}).Error
}

// CreateRekeningTransaction - Membuat record transaksi rekening
func (r *QardHassanRepository) CreateRekeningTransaction(transaction *models.RekeningTransaction) error {
	return r.db.Create(transaction).Error
}

// GetOutstandingByMemberID menghitung total sisa pinjaman Qard Hassan milik seorang anggota
func (r *QardHassanRepository) GetOutstandingByMemberID(memberID int) (float64, error) {
	var outstanding float64
	
	// COALESCE memastikan jika member belum pernah pinjam, hasilnya 0 (bukan error/null)
	query := `
		SELECT COALESCE(SUM(nominalpinjaman - nominalpembayaran), 0) 
		FROM qardhassan 
		WHERE idmember = ?
	`
	
	err := r.db.Raw(query, memberID).Scan(&outstanding).Error
	if err != nil {
		return 0, err
	}
	
	return outstanding, nil
}
