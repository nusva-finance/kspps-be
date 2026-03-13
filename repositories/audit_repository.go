package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		db: config.GetDB(),
	}
}

func (r *AuditRepository) Create(audit *models.AuditLog) error {
	return r.db.Create(audit).Error
}

func (r *AuditRepository) List(offset, limit int) ([]models.AuditLog, int64, error) {
	var audits []models.AuditLog
	var total int64

	err := r.db.Model(&models.AuditLog{}).
		Preload("User").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&audits).Error
	if err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}

func (r *AuditRepository) ListByModule(module string, offset, limit int) ([]models.AuditLog, int64, error) {
	var audits []models.AuditLog
	var total int64

	err := r.db.Model(&models.AuditLog{}).
		Where("module = ?", module).
		Preload("User").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").
		Where("module = ?", module).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&audits).Error
	if err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}

func (r *AuditRepository) ListByUser(userID uint, offset, limit int) ([]models.AuditLog, int64, error) {
	var audits []models.AuditLog
	var total int64

	err := r.db.Model(&models.AuditLog{}).
		Where("user_id = ?", userID).
		Preload("User").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&audits).Error
	if err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}
