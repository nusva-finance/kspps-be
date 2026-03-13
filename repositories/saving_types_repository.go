package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type SavingTypesRepository struct {
	db *gorm.DB
}

func NewSavingTypesRepository() *SavingTypesRepository {
	return &SavingTypesRepository{
		db: config.GetDB(),
	}
}

// Create creates a new saving type
func (r *SavingTypesRepository) Create(savingType *models.SavingType) error {
	return r.db.Create(savingType).Error
}

// FindByID finds a saving type by ID
func (r *SavingTypesRepository) FindByID(id uint) (*models.SavingType, error) {
	var savingType models.SavingType
	err := r.db.First(&savingType, id).Error
	if err != nil {
		return nil, err
	}
	return &savingType, nil
}

// FindByCode finds a saving type by code
func (r *SavingTypesRepository) FindByCode(code string) (*models.SavingType, error) {
	var savingType models.SavingType
	err := r.db.Where("code = ?", code).First(&savingType).Error
	if err != nil {
		return nil, err
	}
	return &savingType, nil
}

// List returns all active saving types
func (r *SavingTypesRepository) List() ([]models.SavingType, error) {
	var savingTypes []models.SavingType
	err := r.db.Where("is_active = ?", true).Order("display_order ASC").Find(&savingTypes).Error
	if err != nil {
		return nil, err
	}
	return savingTypes, nil
}

// ListAll returns all saving types including inactive ones
func (r *SavingTypesRepository) ListAll() ([]models.SavingType, error) {
	var savingTypes []models.SavingType
	err := r.db.Order("display_order ASC").Find(&savingTypes).Error
	if err != nil {
		return nil, err
	}
	return savingTypes, nil
}

// Update updates an existing saving type
func (r *SavingTypesRepository) Update(savingType *models.SavingType) error {
	return r.db.Save(savingType).Error
}

// Delete soft deletes a saving type
func (r *SavingTypesRepository) Delete(id uint) error {
	return r.db.Delete(&models.SavingType{}, id).Error
}

// InitializeDefaultTypes creates default saving types if they don't exist
func (r *SavingTypesRepository) InitializeDefaultTypes() error {
	for _, defaultType := range models.DefaultSavingTypes() {
		// Check if type with this code already exists
		existingType, err := r.FindByCode(defaultType.Code)
		if err == gorm.ErrRecordNotFound {
			// Type doesn't exist, create it
			if err := r.Create(&defaultType); err != nil {
				return err
			}
		} else if err != nil {
			// Other error
			return err
		} else {
			// Type exists, update if needed (optional)
			existingType.Name = defaultType.Name
			existingType.Description = defaultType.Description
			existingType.IsRequired = defaultType.IsRequired
			existingType.MinBalance = defaultType.MinBalance
			existingType.DisplayOrder = defaultType.DisplayOrder
			if err := r.Update(existingType); err != nil {
				return err
			}
		}
	}
	return nil
}
