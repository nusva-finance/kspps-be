package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type MarginRepository struct {
	db *gorm.DB
}

func NewMarginRepository() *MarginRepository {
	return &MarginRepository{
		db: config.GetDB(),
	}
}

func (r *MarginRepository) Create(margin *models.MarginSetup) error {
	return r.db.Create(margin).Error
}

func (r *MarginRepository) FindByID(id uint) (*models.MarginSetup, error) {
	var margin models.MarginSetup
	err := r.db.First(&margin, id).Error
	if err != nil {
		return nil, err
	}
	return &margin, nil
}

func (r *MarginRepository) List(offset, limit int) ([]models.MarginSetup, int64, error) {
	var margins []models.MarginSetup
	var total int64

	err := r.db.Model(&models.MarginSetup{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("category, tenor").Offset(offset).Limit(limit).Find(&margins).Error
	if err != nil {
		return nil, 0, err
	}

	return margins, total, nil
}

func (r *MarginRepository) Update(margin *models.MarginSetup) error {
	return r.db.Save(margin).Error
}

func (r *MarginRepository) Delete(id uint) error {
	return r.db.Delete(&models.MarginSetup{}, id).Error
}

func (r *MarginRepository) FindByCategoryAndTenor(category string, tenor int) (*models.MarginSetup, error) {
	var margin models.MarginSetup
	err := r.db.Where("category = ? AND tenor = ?", category, tenor).First(&margin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &margin, nil
}
