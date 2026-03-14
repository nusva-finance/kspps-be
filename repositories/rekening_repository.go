package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type NusvaRekeningRepository struct {
	db *gorm.DB
}

func NewNusvaRekeningRepository() *NusvaRekeningRepository {
	return &NusvaRekeningRepository{
		db: config.GetDB(),
	}
}

func (r *NusvaRekeningRepository) Create(rekening *models.NusvaRekening) error {
	return r.db.Create(rekening).Error
}

func (r *NusvaRekeningRepository) FindByID(id uint) (*models.NusvaRekening, error) {
	var rekening models.NusvaRekening
	err := r.db.First(&rekening, id).Error
	if err != nil {
		return nil, err
	}
	return &rekening, nil
}

func (r *NusvaRekeningRepository) List(offset, limit int) ([]models.NusvaRekening, int64, error) {
	var rekenings []models.NusvaRekening
	var total int64

	err := r.db.Model(&models.NusvaRekening{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("namarekening").Offset(offset).Limit(limit).Find(&rekenings).Error
	if err != nil {
		return nil, 0, err
	}

	return rekenings, total, nil
}

func (r *NusvaRekeningRepository) ListActive(offset, limit int) ([]models.NusvaRekening, int64, error) {
	var rekenings []models.NusvaRekening
	var total int64

	err := r.db.Model(&models.NusvaRekening{}).Where("aktif = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("aktif = ?", true).Order("namarekening").Offset(offset).Limit(limit).Find(&rekenings).Error
	if err != nil {
		return nil, 0, err
	}

	return rekenings, total, nil
}

func (r *NusvaRekeningRepository) Update(rekening *models.NusvaRekening) error {
	return r.db.Save(rekening).Error
}

func (r *NusvaRekeningRepository) Delete(id uint) error {
	return r.db.Delete(&models.NusvaRekening{}, id).Error
}

func (r *NusvaRekeningRepository) FindByNoRekening(noRekening string) (*models.NusvaRekening, error) {
	var rekening models.NusvaRekening
	err := r.db.Where("norekening = ?", noRekening).First(&rekening).Error
	if err != nil {
		return nil, err
	}
	return &rekening, nil
}
