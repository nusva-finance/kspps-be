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

func (r *PembiayaanRepository) Create(pembiayaan *models.Pembiayaan) error {
	return r.db.Create(pembiayaan).Error
}

func (r *PembiayaanRepository) FindByID(id int) (*models.Pembiayaan, error) {
	var pembiayaan models.Pembiayaan
	err := r.db.First(&pembiayaan, id).Error
	if err != nil {
		return nil, err
	}
	return &pembiayaan, nil
}

func (r *PembiayaanRepository) List(offset, limit int) ([]models.PembiayaanWithMemberName, int64, error) {
	var pembiayaans []models.PembiayaanWithMemberName
	var total int64

	// Gunakan nama tabel yang konsisten (asumsi: member sesuai DDL sebelumnya)
	dbQuery := r.db.Table("pembiayaan").
		Joins("LEFT JOIN members ON pembiayaan.idmember = members.id") // Ubah members ke member jika perlu

	// Hitung Total
	err := dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil Data dengan Select yang Lengkap
	err = dbQuery.Select("pembiayaan.*, members.full_name as nama_anggota, members.member_no as member_no").
		Order("pembiayaan.tanggalpinjaman DESC").
		Offset(offset).
		Limit(limit).
		Scan(&pembiayaans).Error // Gunakan Scan untuk struct dengan join

	if err != nil {
		return nil, 0, err
	}

	return pembiayaans, total, nil
}

func (r *PembiayaanRepository) ListByMemberID(memberID int, offset, limit int) ([]models.PembiayaanWithMemberName, int64, error) {
	var pembiayaans []models.PembiayaanWithMemberName
	var total int64

	// 1. Buat base query agar filter Where dan Join konsisten antara Count dan Select
	query := r.db.Table("pembiayaan").
		Joins("LEFT JOIN members ON pembiayaan.idmember = members.id").
		Where("pembiayaan.idmember = ?", memberID)

	// 2. Hitung total data untuk member tersebut
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 3. Ambil data dengan Select yang mencakup full_name dan member_no
	// Gunakan Scan agar field di luar struct utama (NamaAnggota & MemberNo) terisi
	err = query.Select("pembiayaan.*, members.full_name as nama_anggota, members.member_no as member_no").
		Order("pembiayaan.tanggalpinjaman DESC").
		Offset(offset).
		Limit(limit).
		Scan(&pembiayaans).Error

	if err != nil {
		return nil, 0, err
	}

	return pembiayaans, total, nil
}

func (r *PembiayaanRepository) Update(pembiayaan *models.Pembiayaan) error {
	return r.db.Save(pembiayaan).Error
}

func (r *PembiayaanRepository) Delete(id int) error {
	return r.db.Delete(&models.Pembiayaan{}, id).Error
}

func (r *PembiayaanRepository) FindByMemberID(memberID int) (*models.Pembiayaan, error) {
	var pembiayaan models.Pembiayaan
	err := r.db.Where("idmember = ?", memberID).First(&pembiayaan).Error
	if err != nil {
		return nil, err
	}
	return &pembiayaan, nil
}
