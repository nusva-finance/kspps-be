package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type KategoriBarangRepository struct {
	db *gorm.DB
}

func NewKategoriBarangRepository() *KategoriBarangRepository {
	return &KategoriBarangRepository{
		db: config.GetDB(),
	}
}

func (r *KategoriBarangRepository) Create(kategoriBarang *models.KategoriBarang) error {
	return r.db.Create(kategoriBarang).Error
}

func (r *KategoriBarangRepository) FindByID(id uint) (*models.KategoriBarang, error) {
	var kategoriBarang models.KategoriBarang
	err := r.db.First(&kategoriBarang, id).Error
	if err != nil {
		return nil, err
	}
	return &kategoriBarang, nil
}

func (r *KategoriBarangRepository) List(offset, limit int) ([]models.KategoriBarang, int64, error) {
	var kategoriBarangs []models.KategoriBarang
	var total int64

	err := r.db.Model(&models.KategoriBarang{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("namakategoribarang").Offset(offset).Limit(limit).Find(&kategoriBarangs).Error
	if err != nil {
		return nil, 0, err
	}

	return kategoriBarangs, total, nil
}

func (r *KategoriBarangRepository) ListActive(offset, limit int) ([]models.KategoriBarang, int64, error) {
	var kategoriBarangs []models.KategoriBarang
	var total int64

	err := r.db.Model(&models.KategoriBarang{}).Where("is_active = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("is_active = ?", true).Order("namakategoribarang").Offset(offset).Limit(limit).Find(&kategoriBarangs).Error
	if err != nil {
		return nil, 0, err
	}

	return kategoriBarangs, total, nil
}

func (r *KategoriBarangRepository) Update(kategoriBarang *models.KategoriBarang) error {
	return r.db.Save(kategoriBarang).Error
}

func (r *KategoriBarangRepository) Delete(id uint) error {
	return r.db.Delete(&models.KategoriBarang{}, id).Error
}

func (r *KategoriBarangRepository) FindByKategori(kategori string) (*models.KategoriBarang, error) {
	var kategoriBarang models.KategoriBarang
	err := r.db.Where("namakategoribarang = ?", kategori).First(&kategoriBarang).Error
	if err != nil {
		return nil, err
	}
	return &kategoriBarang, nil
}
