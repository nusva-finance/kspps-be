package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type PembiayaanRepository struct {
	db *gorm.DB
}

func NewPembiayaanRepository() *PembiayaanRepository {
	return &PembiayaanRepository{
		db: config.GetDB(),
	}
}

// Create - Membuat data pembiayaan baru
func (r *PembiayaanRepository) Create(pembiayaan *models.Pembiayaan) error {
	return r.db.Create(pembiayaan).Error
}

// FindByID - Mencari pembiayaan berdasarkan ID tunggal
func (r *PembiayaanRepository) FindByID(id int) (*models.Pembiayaan, error) {
	var pembiayaan models.Pembiayaan
	// Pastikan kolom primary key adalah idpinjaman
	err := r.db.Where("idpinjaman = ?", id).First(&pembiayaan).Error
	if err != nil {
		return nil, err
	}
	return &pembiayaan, nil
}

// List - Mengambil daftar semua pembiayaan dengan data Anggota (JOIN)
func (r *PembiayaanRepository) List(offset, limit int) ([]models.PembiayaanWithMemberName, int64, error) {
	var result []models.PembiayaanWithMemberName
	var total int64

	// Hitung Total Baris
	err := r.db.Table("pembiayaan").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil Data dengan JOIN ke tabel members dan subquery untuk total pembayaran
	query := `
		SELECT
			p.idpinjaman, p.idmember, p.tipepinjaman, p.tanggalpinjaman, p.kategoribarang,
			p.tenor, p.margin, p.nominalpinjaman, p.nominalpembelian, p.tgljtangsuran1,
			p.sysrevid, p.created_by, p.updated_by, p.created_at, p.updated_at,
			m.full_name as nama_anggota,
			m.member_no as member_no,
			COALESCE((SELECT SUM(nominalpembayaran) FROM pembayaran_pembiayaan WHERE idpinjaman = p.idpinjaman), 0) as total_pembayaran
		FROM pembiayaan p
		LEFT JOIN members m ON p.idmember = m.id
		ORDER BY p.tanggalpinjaman DESC
		LIMIT ? OFFSET ?
	`

	err = r.db.Raw(query, limit, offset).Scan(&result).Error

	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ListByMemberID - Mengambil daftar pembiayaan khusus untuk satu anggota tertentu
func (r *PembiayaanRepository) ListByMemberID(memberID int, offset, limit int) ([]models.PembiayaanWithMemberName, int64, error) {
	var result []models.PembiayaanWithMemberName
	var total int64

	// Hitung total data untuk member tersebut
	err := r.db.Table("pembiayaan").Where("idmember = ?", memberID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil data dengan JOIN dan subquery untuk total pembayaran
	query := `
		SELECT
			p.idpinjaman, p.idmember, p.tipepinjaman, p.tanggalpinjaman, p.kategoribarang,
			p.tenor, p.margin, p.nominalpinjaman, p.nominalpembelian, p.tgljtangsuran1,
			p.sysrevid, p.created_by, p.updated_by, p.created_at, p.updated_at,
			m.full_name as nama_anggota,
			m.member_no as member_no,
			COALESCE((SELECT SUM(nominalpembayaran) FROM pembayaran_pembiayaan WHERE idpinjaman = p.idpinjaman), 0) as total_pembayaran
		FROM pembiayaan p
		LEFT JOIN members m ON p.idmember = m.id
		WHERE p.idmember = ?
		ORDER BY p.tanggalpinjaman DESC
		LIMIT ? OFFSET ?
	`

	err = r.db.Raw(query, memberID, limit, offset).Scan(&result).Error

	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// Update - Memperbarui data pembiayaan yang sudah ada
func (r *PembiayaanRepository) Update(pembiayaan *models.Pembiayaan) error {
	return r.db.Save(pembiayaan).Error
}

// Delete - Menghapus data pembiayaan berdasarkan ID
func (r *PembiayaanRepository) Delete(id int) error {
	// Pastikan menghapus berdasarkan kolom idpinjaman
	return r.db.Where("idpinjaman = ?", id).Delete(&models.Pembiayaan{}).Error
}

// FindByMemberID - Mencari data pembiayaan pertama yang ditemukan untuk seorang anggota
func (r *PembiayaanRepository) FindByMemberID(memberID int) (*models.Pembiayaan, error) {
	var pembiayaan models.Pembiayaan
	err := r.db.Where("idmember = ?", memberID).First(&pembiayaan).Error
	if err != nil {
		return nil, err
	}
	return &pembiayaan, nil
}